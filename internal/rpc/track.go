package rpc

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/tracking/analysis/internal/bot"
	"github.com/tracking/analysis/internal/config"
	"github.com/tracking/analysis/internal/dedup"
	"github.com/tracking/analysis/internal/geo"
	"github.com/tracking/analysis/internal/middleware"
	"github.com/tracking/analysis/internal/models"
	"github.com/tracking/analysis/internal/repo"
)

type TrackHandlers struct {
	Config      *config.Config
	Redis       *redis.Client
	PrivKey     *rsa.PrivateKey
	ClickRepo   *repo.ClickRepo
	EventRepo   *repo.EventRepo
	SiteRepo    *repo.SiteRepo
	TokenRepo   *repo.TokenRepo
	GeoResolver *geo.Resolver
}

// track.collectClick
func (h *TrackHandlers) CollectClick(ctx context.Context, params json.RawMessage) (any, *RPCError) {
	var envelope middleware.EncryptedParams
	if err := json.Unmarshal(params, &envelope); err != nil {
		return nil, NewRPCError(ErrCodeInvalidParams, err.Error())
	}

	ip, ua, lang, referer, secFetch := extractClientInfo(ctx)

	// Rate limiting
	if err := middleware.CheckRateLimit(ctx, h.Redis, &h.Config.RateLimitConfiguration, ip, ua, ""); err != nil {
		return nil, NewRPCError(ErrCodeRateLimited, nil)
	}

	// Anti-replay
	if err := middleware.CheckAntiReplay(ctx, h.Redis, &h.Config.SecurityConfiguration, envelope.TS, envelope.Nonce2); err != nil {
		if err.Error() == "replay_detected" {
			return nil, NewRPCError(ErrCodeReplayDetected, nil)
		}
		return nil, NewRPCError(ErrCodeExpiredToken, nil)
	}

	// Decrypt
	plaintext, err := middleware.DecryptRequest(envelope, h.PrivKey)
	if err != nil {
		return nil, NewRPCError(ErrCodeDecryptFailed, nil)
	}

	// Parse decrypted payload
	var payload struct {
		Token     string         `json:"token"`
		VisitorID string         `json:"visitor_id"`
		Env       models.JSONMap `json:"env"`
	}
	if err := json.Unmarshal(plaintext, &payload); err != nil {
		return nil, NewRPCError(ErrCodeInvalidParams, "invalid decrypted payload")
	}

	// Look up token by short code
	tkn, err := h.TokenRepo.GetByShortCode(payload.Token)
	if err != nil {
		return nil, NewRPCError(ErrCodeInvalidToken, nil)
	}

	// Bot detection
	recentHits := bot.CountRecentHits(ctx, h.Redis, ip)
	botScore := bot.Score(ua, lang, secFetch, referer, recentHits)
	blocked, suspected := bot.IsBot(botScore, &h.Config.BotConfiguration)
	if blocked && h.Config.BotConfiguration.BlockMode == "reject" {
		return nil, NewRPCError(ErrCodeBotBlocked, nil)
	}

	// Dedup check
	if dedup.CheckClickDedup(ctx, h.Redis, &h.Config.SecurityConfiguration, tkn.TrackerID, tkn.ChannelID, payload.VisitorID) {
		return map[string]any{"target_url": "", "click_id": "", "dedup": true}, nil
	}

	// Write click
	click := &models.Click{
		TS:           time.Now(),
		TrackerID:    tkn.TrackerID,
		CampaignID:   tkn.CampaignID,
		ChannelID:    tkn.ChannelID,
		TargetID:     tkn.TargetID,
		VisitorID:    payload.VisitorID,
		IP:           ip,
		Country:      h.GeoResolver.Country(ip),
		UA:           ua,
		Lang:         lang,
		Referer:      referer,
		Props:        payload.Env,
		SuspectedBot: suspected,
		IsBot:        blocked,
	}
	if err := h.ClickRepo.Create(click); err != nil {
		return nil, NewRPCError(ErrCodeDBError, nil)
	}

	return map[string]any{"click_id": click.ID, "target_id": tkn.TargetID}, nil
}

// track.collectEvents
func (h *TrackHandlers) CollectEvents(ctx context.Context, params json.RawMessage) (any, *RPCError) {
	var envelope middleware.EncryptedParams
	if err := json.Unmarshal(params, &envelope); err != nil {
		return nil, NewRPCError(ErrCodeInvalidParams, err.Error())
	}

	ip, ua, lang, _, secFetch := extractClientInfo(ctx)

	// Rate limiting
	if err := middleware.CheckRateLimit(ctx, h.Redis, &h.Config.RateLimitConfiguration, ip, ua, ""); err != nil {
		return nil, NewRPCError(ErrCodeRateLimited, nil)
	}

	// Anti-replay
	if err := middleware.CheckAntiReplay(ctx, h.Redis, &h.Config.SecurityConfiguration, envelope.TS, envelope.Nonce2); err != nil {
		if err.Error() == "replay_detected" {
			return nil, NewRPCError(ErrCodeReplayDetected, nil)
		}
		return nil, NewRPCError(ErrCodeExpiredToken, nil)
	}

	// Decrypt
	plaintext, err := middleware.DecryptRequest(envelope, h.PrivKey)
	if err != nil {
		return nil, NewRPCError(ErrCodeDecryptFailed, nil)
	}

	// Parse decrypted payload
	var payload struct {
		SiteKey   string `json:"site_key"`
		VisitorID string `json:"visitor_id"`
		SessionID string `json:"session_id"`
		Events    []struct {
			Type     string         `json:"type"`
			URL      string         `json:"url"`
			Title    string         `json:"title"`
			Referrer string         `json:"referrer"`
			Props    models.JSONMap `json:"props"`
		} `json:"events"`
	}
	if err := json.Unmarshal(plaintext, &payload); err != nil {
		return nil, NewRPCError(ErrCodeInvalidParams, "invalid decrypted payload")
	}

	// Verify site key
	site, err := h.SiteRepo.GetByKey(payload.SiteKey)
	if err != nil {
		return nil, NewRPCErrorWithMessage(ErrCodeInvalidParams, "invalid site_key")
	}

	// Bot detection
	recentHits := bot.CountRecentHits(ctx, h.Redis, ip)
	botScore := bot.Score(ua, lang, secFetch, "", recentHits)
	blocked, suspected := bot.IsBot(botScore, &h.Config.BotConfiguration)
	if blocked && h.Config.BotConfiguration.BlockMode == "reject" {
		return nil, NewRPCError(ErrCodeBotBlocked, nil)
	}

	// Build events
	country := h.GeoResolver.Country(ip)
	events := make([]models.Event, 0, len(payload.Events))
	now := time.Now()
	for _, e := range payload.Events {
		events = append(events, models.Event{
			TS:           now,
			SiteID:       site.ID,
			Type:         e.Type,
			VisitorID:    payload.VisitorID,
			SessionID:    payload.SessionID,
			URL:          e.URL,
			Title:        e.Title,
			Referrer:     e.Referrer,
			IP:           ip,
			Country:      country,
			UA:           ua,
			Lang:         lang,
			Props:        e.Props,
			SuspectedBot: suspected,
			IsBot:        blocked,
		})
	}

	if err := h.EventRepo.BatchCreate(events); err != nil {
		return nil, NewRPCError(ErrCodeDBError, nil)
	}

	return map[string]any{"ok": true, "server_time": now.Unix()}, nil
}

func extractClientInfo(ctx context.Context) (ip, ua, lang, referer, secFetch string) {
	c := GinContext(ctx)
	if c == nil {
		return
	}
	ip = c.ClientIP()
	ua = c.GetHeader("User-Agent")
	lang = c.GetHeader("Accept-Language")
	referer = c.GetHeader("Referer")
	secFetch = c.GetHeader("Sec-Fetch-Mode")
	return
}
