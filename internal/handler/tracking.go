package handler

import (
	"crypto/rsa"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/tracking/analysis/internal/bot"
	"github.com/tracking/analysis/internal/config"
	"github.com/tracking/analysis/internal/models"
	"github.com/tracking/analysis/internal/repo"
	"github.com/tracking/analysis/internal/sdk"
	"github.com/tracking/analysis/internal/security"
)

type TrackingHandler struct {
	Config     *config.Config
	ClickRepo  *repo.ClickRepo
	TargetRepo *repo.TargetRepo
	TokenRepo  *repo.TokenRepo
	PubKey     *rsa.PublicKey
	PrivKey    *rsa.PrivateKey
	Redis      *redis.Client
	BotCfg     *config.BotConfiguration
}

// GET /t/:token — JS-based click tracking page
func (h *TrackingHandler) HandleJSTrack(c *gin.Context) {
	token := c.Param("token")
	tkn, err := h.TokenRepo.GetByShortCode(token)
	if err != nil {
		c.String(http.StatusBadRequest, "invalid token")
		return
	}

	pubPEM, err := security.PublicKeyPEM(h.PubKey)
	if err != nil {
		c.String(http.StatusInternalServerError, "server error")
		return
	}

	// Look up target URL
	target, err := h.TargetRepo.GetByID(tkn.TargetID)
	if err != nil {
		c.String(http.StatusNotFound, "target not found")
		return
	}

	html := sdk.GenerateClickPage(token, pubPEM, h.Config.SecurityConfiguration.KID, h.Config.ServiceConfiguration.ExportURL, target.URL)
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, html)
}

// GET /r/:token — 302 redirect click tracking
func (h *TrackingHandler) HandleRedirectTrack(c *gin.Context) {
	token := c.Param("token")
	tkn, err := h.TokenRepo.GetByShortCode(token)
	if err != nil {
		c.String(http.StatusBadRequest, "invalid token")
		return
	}

	// Look up target URL from DB (never from token)
	target, err := h.TargetRepo.GetByID(tkn.TargetID)
	if err != nil {
		c.String(http.StatusNotFound, "target not found")
		return
	}

	// Bot detection (mark only, never block on 302)
	ua := c.GetHeader("User-Agent")
	lang := c.GetHeader("Accept-Language")
	secFetch := c.GetHeader("Sec-Fetch-Mode")
	referer := c.GetHeader("Referer")
	recentHits := bot.CountRecentHits(c.Request.Context(), h.Redis, c.ClientIP())
	botScore := bot.Score(ua, lang, secFetch, referer, recentHits)
	_, suspected := bot.IsBot(botScore, h.BotCfg)

	// Record click
	click := &models.Click{
		TS:           time.Now(),
		TrackerID:    tkn.TrackerID,
		TargetID:     tkn.TargetID,
		IP:           c.ClientIP(),
		UA:           ua,
		Lang:         lang,
		Referer:      referer,
		SuspectedBot: suspected,
		IsBot:        suspected,
	}
	if tkn.CampaignID != "" {
		click.CampaignID = tkn.CampaignID
	}
	if tkn.ChannelID != "" {
		click.ChannelID = tkn.ChannelID
	}
	if err := h.ClickRepo.Create(click); err != nil {
		slog.Error("failed to record click", "error", err, "token", token)
	}

	c.Redirect(http.StatusFound, target.URL)
}

// GET /sdk/track.js — serve the JS SDK
func (h *TrackingHandler) HandleSDK(c *gin.Context) {
	pubPEM, err := security.PublicKeyPEM(h.PubKey)
	if err != nil {
		c.String(http.StatusInternalServerError, "server error")
		return
	}

	js := sdk.GenerateSDK(pubPEM, h.Config.SecurityConfiguration.KID, h.Config.ServiceConfiguration.ExportURL)
	c.Header("Content-Type", "application/javascript; charset=utf-8")
	c.Header("Cache-Control", "public, max-age=3600")
	c.String(http.StatusOK, js)
}

// GET /public-keys.json — serve the public key for client-side encryption
func (h *TrackingHandler) HandlePublicKeys(c *gin.Context) {
	pubPEM, err := security.PublicKeyPEM(h.PubKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"kid":        h.Config.SecurityConfiguration.KID,
		"public_key": pubPEM,
	})
}
