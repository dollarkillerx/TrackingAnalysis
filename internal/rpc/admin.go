package rpc

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/tracking/analysis/internal/config"
	"github.com/tracking/analysis/internal/models"
	"github.com/tracking/analysis/internal/repo"
	"github.com/tracking/analysis/internal/security"
)

type AdminHandlers struct {
	Config      *config.Config
	TrackerRepo *repo.TrackerRepo
	CampaignRepo *repo.CampaignRepo
	ChannelRepo *repo.ChannelRepo
	TargetRepo  *repo.TargetRepo
	SiteRepo    *repo.SiteRepo
}

// Session token generation using HMAC
func (h *AdminHandlers) generateSessionToken(username string) string {
	mac := hmac.New(sha256.New, []byte(h.Config.SecurityConfiguration.TokenSecret))
	mac.Write([]byte(fmt.Sprintf("%s:%d", username, time.Now().Unix()/3600)))
	return hex.EncodeToString(mac.Sum(nil))
}

func (h *AdminHandlers) verifySessionToken(token string) bool {
	expected := h.generateSessionToken(h.Config.AdminConfiguration.Username)
	return hmac.Equal([]byte(token), []byte(expected))
}

func (h *AdminHandlers) requireAuth(params json.RawMessage) *RPCError {
	var p struct {
		AdminToken string `json:"admin_token"`
	}
	if err := json.Unmarshal(params, &p); err != nil || p.AdminToken == "" {
		return NewRPCErrorWithMessage(ErrCodeInvalidToken, "admin_token required")
	}
	if !h.verifySessionToken(p.AdminToken) {
		return NewRPCErrorWithMessage(ErrCodeInvalidToken, "invalid admin_token")
	}
	return nil
}

func (h *AdminHandlers) Register(d *Dispatcher) {
	d.Register("admin.login", h.Login)
	d.Register("admin.tracker.create", h.TrackerCreate)
	d.Register("admin.tracker.list", h.TrackerList)
	d.Register("admin.tracker.update", h.TrackerUpdate)
	d.Register("admin.tracker.delete", h.TrackerDelete)
	d.Register("admin.campaign.create", h.CampaignCreate)
	d.Register("admin.campaign.list", h.CampaignList)
	d.Register("admin.channel.create", h.ChannelCreate)
	d.Register("admin.channel.batchImport", h.ChannelBatchImport)
	d.Register("admin.channel.list", h.ChannelList)
	d.Register("admin.target.create", h.TargetCreate)
	d.Register("admin.target.list", h.TargetList)
	d.Register("admin.site.create", h.SiteCreate)
	d.Register("admin.site.list", h.SiteList)
	d.Register("admin.token.generate", h.TokenGenerate)
}

// admin.login
func (h *AdminHandlers) Login(_ context.Context, params json.RawMessage) (any, *RPCError) {
	var p struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, NewRPCError(ErrCodeInvalidParams, err.Error())
	}
	if p.Username != h.Config.AdminConfiguration.Username || p.Password != h.Config.AdminConfiguration.Password {
		return nil, NewRPCErrorWithMessage(ErrCodeInvalidToken, "invalid credentials")
	}
	token := h.generateSessionToken(p.Username)
	return map[string]string{"admin_token": token}, nil
}

// admin.tracker.create
func (h *AdminHandlers) TrackerCreate(_ context.Context, params json.RawMessage) (any, *RPCError) {
	if err := h.requireAuth(params); err != nil {
		return nil, err
	}
	var p struct {
		AdminToken string `json:"admin_token"`
		Name       string `json:"name"`
		Type       string `json:"type"`
		Mode       string `json:"mode"`
	}
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, NewRPCError(ErrCodeInvalidParams, err.Error())
	}
	if p.Type != "ad" && p.Type != "web" {
		return nil, NewRPCErrorWithMessage(ErrCodeInvalidParams, "type must be 'ad' or 'web'")
	}
	if p.Mode == "" {
		p.Mode = "302"
	}
	tracker := &models.Tracker{
		Name:   p.Name,
		Type:   p.Type,
		Mode:   p.Mode,
		Status: "active",
	}
	if err := h.TrackerRepo.Create(tracker); err != nil {
		return nil, NewRPCError(ErrCodeDBError, err.Error())
	}
	return tracker, nil
}

// admin.tracker.list
func (h *AdminHandlers) TrackerList(_ context.Context, params json.RawMessage) (any, *RPCError) {
	if err := h.requireAuth(params); err != nil {
		return nil, err
	}
	var p struct {
		Type string `json:"type"`
	}
	json.Unmarshal(params, &p)
	trackers, err := h.TrackerRepo.List(p.Type)
	if err != nil {
		return nil, NewRPCError(ErrCodeDBError, err.Error())
	}
	return trackers, nil
}

// admin.tracker.update
func (h *AdminHandlers) TrackerUpdate(_ context.Context, params json.RawMessage) (any, *RPCError) {
	if err := h.requireAuth(params); err != nil {
		return nil, err
	}
	var p struct {
		AdminToken string `json:"admin_token"`
		ID         string `json:"id"`
		Name       string `json:"name"`
		Status     string `json:"status"`
		Mode       string `json:"mode"`
	}
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, NewRPCError(ErrCodeInvalidParams, err.Error())
	}
	tracker, err := h.TrackerRepo.GetByID(p.ID)
	if err != nil {
		return nil, NewRPCError(ErrCodeDBError, "tracker not found")
	}
	if p.Name != "" {
		tracker.Name = p.Name
	}
	if p.Status != "" {
		tracker.Status = p.Status
	}
	if p.Mode != "" {
		tracker.Mode = p.Mode
	}
	if err := h.TrackerRepo.Update(tracker); err != nil {
		return nil, NewRPCError(ErrCodeDBError, err.Error())
	}
	return tracker, nil
}

// admin.tracker.delete
func (h *AdminHandlers) TrackerDelete(_ context.Context, params json.RawMessage) (any, *RPCError) {
	if err := h.requireAuth(params); err != nil {
		return nil, err
	}
	var p struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, NewRPCError(ErrCodeInvalidParams, err.Error())
	}
	if err := h.TrackerRepo.Delete(p.ID); err != nil {
		return nil, NewRPCError(ErrCodeDBError, err.Error())
	}
	return map[string]bool{"ok": true}, nil
}

// admin.campaign.create
func (h *AdminHandlers) CampaignCreate(_ context.Context, params json.RawMessage) (any, *RPCError) {
	if err := h.requireAuth(params); err != nil {
		return nil, err
	}
	var p struct {
		AdminToken string `json:"admin_token"`
		TrackerID  string `json:"tracker_id"`
		Name       string `json:"name"`
	}
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, NewRPCError(ErrCodeInvalidParams, err.Error())
	}
	// Validate tracker exists and is type=ad
	tracker, err := h.TrackerRepo.GetByID(p.TrackerID)
	if err != nil {
		return nil, NewRPCErrorWithMessage(ErrCodeInvalidParams, "tracker not found")
	}
	if tracker.Type != "ad" {
		return nil, NewRPCErrorWithMessage(ErrCodeInvalidParams, "campaigns require tracker type 'ad'")
	}
	campaign := &models.Campaign{
		TrackerID: p.TrackerID,
		Name:      p.Name,
		Status:    "active",
	}
	if err := h.CampaignRepo.Create(campaign); err != nil {
		return nil, NewRPCError(ErrCodeDBError, err.Error())
	}
	return campaign, nil
}

// admin.campaign.list
func (h *AdminHandlers) CampaignList(_ context.Context, params json.RawMessage) (any, *RPCError) {
	if err := h.requireAuth(params); err != nil {
		return nil, err
	}
	var p struct {
		TrackerID string `json:"tracker_id"`
	}
	json.Unmarshal(params, &p)
	campaigns, err := h.CampaignRepo.List(p.TrackerID)
	if err != nil {
		return nil, NewRPCError(ErrCodeDBError, err.Error())
	}
	return campaigns, nil
}

// admin.channel.create
func (h *AdminHandlers) ChannelCreate(_ context.Context, params json.RawMessage) (any, *RPCError) {
	if err := h.requireAuth(params); err != nil {
		return nil, err
	}
	var p struct {
		AdminToken string     `json:"admin_token"`
		TrackerID  string     `json:"tracker_id"`
		CampaignID string     `json:"campaign_id"`
		Name       string     `json:"name"`
		Source     string     `json:"source"`
		Medium     string     `json:"medium"`
		Tags       models.JSONMap `json:"tags"`
	}
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, NewRPCError(ErrCodeInvalidParams, err.Error())
	}
	channel := &models.Channel{
		TrackerID:  p.TrackerID,
		CampaignID: p.CampaignID,
		Name:       p.Name,
		Source:     p.Source,
		Medium:     p.Medium,
		Tags:       p.Tags,
	}
	if err := h.ChannelRepo.Create(channel); err != nil {
		return nil, NewRPCError(ErrCodeDBError, err.Error())
	}
	return channel, nil
}

// admin.channel.batchImport
func (h *AdminHandlers) ChannelBatchImport(_ context.Context, params json.RawMessage) (any, *RPCError) {
	if err := h.requireAuth(params); err != nil {
		return nil, err
	}
	var p struct {
		Channels []models.Channel `json:"channels"`
	}
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, NewRPCError(ErrCodeInvalidParams, err.Error())
	}
	if err := h.ChannelRepo.BatchImport(p.Channels); err != nil {
		return nil, NewRPCError(ErrCodeDBError, err.Error())
	}
	return map[string]int{"imported": len(p.Channels)}, nil
}

// admin.channel.list
func (h *AdminHandlers) ChannelList(_ context.Context, params json.RawMessage) (any, *RPCError) {
	if err := h.requireAuth(params); err != nil {
		return nil, err
	}
	var p struct {
		TrackerID  string `json:"tracker_id"`
		CampaignID string `json:"campaign_id"`
	}
	json.Unmarshal(params, &p)
	channels, err := h.ChannelRepo.List(p.TrackerID, p.CampaignID)
	if err != nil {
		return nil, NewRPCError(ErrCodeDBError, err.Error())
	}
	return channels, nil
}

// admin.target.create
func (h *AdminHandlers) TargetCreate(_ context.Context, params json.RawMessage) (any, *RPCError) {
	if err := h.requireAuth(params); err != nil {
		return nil, err
	}
	var p struct {
		AdminToken string `json:"admin_token"`
		TrackerID  string `json:"tracker_id"`
		URL        string `json:"url"`
	}
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, NewRPCError(ErrCodeInvalidParams, err.Error())
	}
	target := &models.Target{
		TrackerID: p.TrackerID,
		URL:       p.URL,
	}
	if err := h.TargetRepo.Create(target); err != nil {
		return nil, NewRPCError(ErrCodeDBError, err.Error())
	}
	return target, nil
}

// admin.target.list
func (h *AdminHandlers) TargetList(_ context.Context, params json.RawMessage) (any, *RPCError) {
	if err := h.requireAuth(params); err != nil {
		return nil, err
	}
	var p struct {
		TrackerID string `json:"tracker_id"`
	}
	json.Unmarshal(params, &p)
	targets, err := h.TargetRepo.List(p.TrackerID)
	if err != nil {
		return nil, NewRPCError(ErrCodeDBError, err.Error())
	}
	return targets, nil
}

// admin.site.create
func (h *AdminHandlers) SiteCreate(_ context.Context, params json.RawMessage) (any, *RPCError) {
	if err := h.requireAuth(params); err != nil {
		return nil, err
	}
	var p struct {
		AdminToken string `json:"admin_token"`
		Name       string `json:"name"`
		Domain     string `json:"domain"`
	}
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, NewRPCError(ErrCodeInvalidParams, err.Error())
	}
	site := &models.Site{
		Name:   p.Name,
		Domain: p.Domain,
		Status: "active",
	}
	if err := h.SiteRepo.Create(site); err != nil {
		return nil, NewRPCError(ErrCodeDBError, err.Error())
	}
	return site, nil
}

// admin.site.list
func (h *AdminHandlers) SiteList(_ context.Context, params json.RawMessage) (any, *RPCError) {
	if err := h.requireAuth(params); err != nil {
		return nil, err
	}
	sites, err := h.SiteRepo.List()
	if err != nil {
		return nil, NewRPCError(ErrCodeDBError, err.Error())
	}
	return sites, nil
}

// admin.token.generate — generates an HMAC tracking token
func (h *AdminHandlers) TokenGenerate(_ context.Context, params json.RawMessage) (any, *RPCError) {
	if err := h.requireAuth(params); err != nil {
		return nil, err
	}
	var p struct {
		AdminToken string `json:"admin_token"`
		TrackerID  string `json:"tracker_id"`
		CampaignID string `json:"campaign_id"`
		ChannelID  string `json:"channel_id"`
		TargetID   string `json:"target_id"`
		Mode       string `json:"mode"`
		ExpSeconds int64  `json:"exp_seconds"`
	}
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, NewRPCError(ErrCodeInvalidParams, err.Error())
	}

	var exp int64
	if p.ExpSeconds > 0 {
		exp = time.Now().Unix() + p.ExpSeconds
	}

	payload := security.TokenPayload{
		TrackerID:  p.TrackerID,
		CampaignID: p.CampaignID,
		ChannelID:  p.ChannelID,
		TargetID:   p.TargetID,
		Exp:        exp,
		Mode:       p.Mode,
	}
	token, err := security.GenerateToken(payload, h.Config.SecurityConfiguration.TokenSecret)
	if err != nil {
		return nil, NewRPCError(ErrCodeInternalError, err.Error())
	}
	return map[string]string{"token": token}, nil
}
