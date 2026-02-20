package rpc

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"time"

	"github.com/tracking/analysis/internal/config"
	"github.com/tracking/analysis/internal/models"
	"github.com/tracking/analysis/internal/repo"
	"github.com/tracking/analysis/internal/security"
	"gorm.io/gorm"
)

type AdminHandlers struct {
	Config       *config.Config
	TrackerRepo  *repo.TrackerRepo
	CampaignRepo *repo.CampaignRepo
	ChannelRepo  *repo.ChannelRepo
	TargetRepo   *repo.TargetRepo
	SiteRepo     *repo.SiteRepo
	TokenRepo    *repo.TokenRepo
	ClickRepo    *repo.ClickRepo
	EventRepo    *repo.EventRepo
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
	d.Register("admin.token.list", h.TokenList)
	d.Register("admin.token.delete", h.TokenDelete)
	d.Register("admin.stats.clicks", h.StatsClicks)
	d.Register("admin.stats.events", h.StatsEvents)
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
	}
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, NewRPCError(ErrCodeInvalidParams, err.Error())
	}
	if p.Type != "ad" && p.Type != "web" {
		return nil, NewRPCErrorWithMessage(ErrCodeInvalidParams, "type must be 'ad' or 'web'")
	}
	tracker := &models.Tracker{
		Name:   p.Name,
		Type:   p.Type,
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

// admin.token.generate — generates a short-code tracking token stored in DB
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
	}
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, NewRPCError(ErrCodeInvalidParams, err.Error())
	}
	if p.Mode == "" {
		p.Mode = "302"
	}

	// Generate unique short code with collision retry
	var shortCode string
	for i := 0; i < 10; i++ {
		candidate := security.GenerateShortCode()
		exists, err := h.TokenRepo.ExistsByShortCode(candidate)
		if err != nil {
			return nil, NewRPCError(ErrCodeDBError, err.Error())
		}
		if !exists {
			shortCode = candidate
			break
		}
	}
	if shortCode == "" {
		return nil, NewRPCErrorWithMessage(ErrCodeInternalError, "failed to generate unique short code")
	}

	token := &models.Token{
		ShortCode:  shortCode,
		TrackerID:  p.TrackerID,
		CampaignID: p.CampaignID,
		ChannelID:  p.ChannelID,
		TargetID:   p.TargetID,
		Mode:       p.Mode,
	}
	if err := h.TokenRepo.Create(token); err != nil {
		return nil, NewRPCError(ErrCodeDBError, err.Error())
	}

	prefix := "r"
	if token.Mode == "js" {
		prefix = "t"
	}
	trackingURL := fmt.Sprintf("%s/%s/%s", h.Config.ServiceConfiguration.ExportURL, prefix, token.ShortCode)

	return map[string]any{
		"id":           token.ID,
		"short_code":   token.ShortCode,
		"tracker_id":   token.TrackerID,
		"campaign_id":  token.CampaignID,
		"channel_id":   token.ChannelID,
		"target_id":    token.TargetID,
		"mode":         token.Mode,
		"created_at":   token.CreatedAt,
		"tracking_url": trackingURL,
	}, nil
}

// admin.token.list
func (h *AdminHandlers) TokenList(_ context.Context, params json.RawMessage) (any, *RPCError) {
	if err := h.requireAuth(params); err != nil {
		return nil, err
	}
	var p struct {
		TrackerID string `json:"tracker_id"`
	}
	json.Unmarshal(params, &p)
	tokens, err := h.TokenRepo.List(p.TrackerID)
	if err != nil {
		return nil, NewRPCError(ErrCodeDBError, err.Error())
	}

	results := make([]map[string]any, len(tokens))
	for i, t := range tokens {
		prefix := "r"
		if t.Mode == "js" {
			prefix = "t"
		}
		results[i] = map[string]any{
			"id":           t.ID,
			"short_code":   t.ShortCode,
			"tracker_id":   t.TrackerID,
			"campaign_id":  t.CampaignID,
			"channel_id":   t.ChannelID,
			"target_id":    t.TargetID,
			"mode":         t.Mode,
			"created_at":   t.CreatedAt,
			"tracking_url": fmt.Sprintf("%s/%s/%s", h.Config.ServiceConfiguration.ExportURL, prefix, t.ShortCode),
		}
	}
	return results, nil
}

// admin.token.delete
func (h *AdminHandlers) TokenDelete(_ context.Context, params json.RawMessage) (any, *RPCError) {
	if err := h.requireAuth(params); err != nil {
		return nil, err
	}
	var p struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, NewRPCError(ErrCodeInvalidParams, err.Error())
	}
	if err := h.TokenRepo.Delete(p.ID); err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, NewRPCErrorWithMessage(ErrCodeInvalidParams, "token not found")
		}
		return nil, NewRPCError(ErrCodeDBError, err.Error())
	}
	return map[string]bool{"ok": true}, nil
}

func safeDivide(numerator, denominator int64) float64 {
	if denominator == 0 {
		return 0
	}
	return math.Round(float64(numerator)/float64(denominator)*10000) / 100
}

// admin.stats.clicks
func (h *AdminHandlers) StatsClicks(_ context.Context, params json.RawMessage) (any, *RPCError) {
	if err := h.requireAuth(params); err != nil {
		return nil, err
	}
	var p struct {
		StartDate  string `json:"start_date"`
		EndDate    string `json:"end_date"`
		TrackerID  string `json:"tracker_id"`
		CampaignID string `json:"campaign_id"`
		ChannelID  string `json:"channel_id"`
	}
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, NewRPCError(ErrCodeInvalidParams, err.Error())
	}
	start, err := time.Parse("2006-01-02", p.StartDate)
	if err != nil {
		return nil, NewRPCErrorWithMessage(ErrCodeInvalidParams, "invalid start_date, expected YYYY-MM-DD")
	}
	end, err := time.Parse("2006-01-02", p.EndDate)
	if err != nil {
		return nil, NewRPCErrorWithMessage(ErrCodeInvalidParams, "invalid end_date, expected YYYY-MM-DD")
	}
	end = end.Add(24*time.Hour - time.Nanosecond)

	daily, err := h.ClickRepo.CountByDay(start, end, p.TrackerID, p.CampaignID, p.ChannelID)
	if err != nil {
		return nil, NewRPCError(ErrCodeDBError, err.Error())
	}
	total, uniqueVisitors, bots, err := h.ClickRepo.Summary(start, end, p.TrackerID, p.CampaignID, p.ChannelID)
	if err != nil {
		return nil, NewRPCError(ErrCodeDBError, err.Error())
	}
	topTrackers, err := h.ClickRepo.TopByGroup(start, end, "tracker_id", 10)
	if err != nil {
		return nil, NewRPCError(ErrCodeDBError, err.Error())
	}
	topChannels, err := h.ClickRepo.TopByGroup(start, end, "channel_id", 10)
	if err != nil {
		return nil, NewRPCError(ErrCodeDBError, err.Error())
	}
	topCampaigns, err := h.ClickRepo.TopByGroup(start, end, "campaign_id", 10)
	if err != nil {
		return nil, NewRPCError(ErrCodeDBError, err.Error())
	}

	// New aggregations — graceful on failure
	topReferrers, err := h.ClickRepo.TopReferrers(start, end, p.TrackerID, p.CampaignID, p.ChannelID, 10)
	if err != nil {
		log.Printf("StatsClicks: TopReferrers error: %v", err)
		topReferrers = []repo.NameCount{}
	}
	rawUA, err := h.ClickRepo.RawUACounts(start, end, p.TrackerID, p.CampaignID, p.ChannelID)
	var browsers, oses []repo.NameCount
	if err != nil {
		log.Printf("StatsClicks: RawUACounts error: %v", err)
		browsers, oses = []repo.NameCount{}, []repo.NameCount{}
	} else {
		browsers, oses = repo.ParseUADistribution(rawUA, 10)
	}
	languages, err := h.ClickRepo.LanguageDistribution(start, end, p.TrackerID, p.CampaignID, p.ChannelID, 10)
	if err != nil {
		log.Printf("StatsClicks: LanguageDistribution error: %v", err)
		languages = []repo.NameCount{}
	}
	countries, err := h.ClickRepo.CountryDistribution(start, end, p.TrackerID, p.CampaignID, p.ChannelID, 10)
	if err != nil {
		log.Printf("StatsClicks: CountryDistribution error: %v", err)
		countries = []repo.NameCount{}
	}
	botDaily, err := h.ClickRepo.BotCountByDay(start, end, p.TrackerID, p.CampaignID, p.ChannelID)
	if err != nil {
		log.Printf("StatsClicks: BotCountByDay error: %v", err)
		botDaily = []repo.DailyCount{}
	}
	hourly, err := h.ClickRepo.CountByHour(start, end, p.TrackerID, p.CampaignID, p.ChannelID)
	if err != nil {
		log.Printf("StatsClicks: CountByHour error: %v", err)
		hourly = []repo.HourlyCount{}
	}

	return map[string]any{
		"summary": map[string]any{
			"total":           total,
			"unique_visitors": uniqueVisitors,
			"bots":            bots,
			"bot_rate":        safeDivide(bots, total),
		},
		"daily":          daily,
		"top_trackers":   topTrackers,
		"top_channels":   topChannels,
		"top_campaigns":  topCampaigns,
		"top_referrers":  topReferrers,
		"browsers":       browsers,
		"oses":           oses,
		"languages":      languages,
		"countries":      countries,
		"bot_daily":      botDaily,
		"hourly":         hourly,
	}, nil
}

// admin.stats.events
func (h *AdminHandlers) StatsEvents(_ context.Context, params json.RawMessage) (any, *RPCError) {
	if err := h.requireAuth(params); err != nil {
		return nil, err
	}
	var p struct {
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
		SiteID    string `json:"site_id"`
	}
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, NewRPCError(ErrCodeInvalidParams, err.Error())
	}
	start, err := time.Parse("2006-01-02", p.StartDate)
	if err != nil {
		return nil, NewRPCErrorWithMessage(ErrCodeInvalidParams, "invalid start_date, expected YYYY-MM-DD")
	}
	end, err := time.Parse("2006-01-02", p.EndDate)
	if err != nil {
		return nil, NewRPCErrorWithMessage(ErrCodeInvalidParams, "invalid end_date, expected YYYY-MM-DD")
	}
	end = end.Add(24*time.Hour - time.Nanosecond)

	daily, err := h.EventRepo.CountByDay(start, end, p.SiteID)
	if err != nil {
		return nil, NewRPCError(ErrCodeDBError, err.Error())
	}
	total, uniqueVisitors, uniqueSessions, bots, err := h.EventRepo.Summary(start, end, p.SiteID)
	if err != nil {
		return nil, NewRPCError(ErrCodeDBError, err.Error())
	}
	topSites, err := h.EventRepo.TopByGroup(start, end, "site_id", 10)
	if err != nil {
		return nil, NewRPCError(ErrCodeDBError, err.Error())
	}
	topTypes, err := h.EventRepo.TopByGroup(start, end, "type", 10)
	if err != nil {
		return nil, NewRPCError(ErrCodeDBError, err.Error())
	}

	// New aggregations — graceful on failure
	topReferrers, err := h.EventRepo.TopReferrers(start, end, p.SiteID, 10)
	if err != nil {
		log.Printf("StatsEvents: TopReferrers error: %v", err)
		topReferrers = []repo.NameCount{}
	}
	topPages, err := h.EventRepo.TopPages(start, end, p.SiteID, 10)
	if err != nil {
		log.Printf("StatsEvents: TopPages error: %v", err)
		topPages = []repo.NameCount{}
	}
	rawUA, err := h.EventRepo.RawUACounts(start, end, p.SiteID)
	var browsers, oses []repo.NameCount
	if err != nil {
		log.Printf("StatsEvents: RawUACounts error: %v", err)
		browsers, oses = []repo.NameCount{}, []repo.NameCount{}
	} else {
		browsers, oses = repo.ParseUADistribution(rawUA, 10)
	}
	languages, err := h.EventRepo.LanguageDistribution(start, end, p.SiteID, 10)
	if err != nil {
		log.Printf("StatsEvents: LanguageDistribution error: %v", err)
		languages = []repo.NameCount{}
	}
	countries, err := h.EventRepo.CountryDistribution(start, end, p.SiteID, 10)
	if err != nil {
		log.Printf("StatsEvents: CountryDistribution error: %v", err)
		countries = []repo.NameCount{}
	}
	botDaily, err := h.EventRepo.BotCountByDay(start, end, p.SiteID)
	if err != nil {
		log.Printf("StatsEvents: BotCountByDay error: %v", err)
		botDaily = []repo.DailyCount{}
	}
	hourly, err := h.EventRepo.CountByHour(start, end, p.SiteID)
	if err != nil {
		log.Printf("StatsEvents: CountByHour error: %v", err)
		hourly = []repo.HourlyCount{}
	}

	return map[string]any{
		"summary": map[string]any{
			"total":           total,
			"unique_visitors": uniqueVisitors,
			"unique_sessions": uniqueSessions,
			"bots":            bots,
			"bot_rate":        safeDivide(bots, total),
		},
		"daily":          daily,
		"top_sites":      topSites,
		"top_types":      topTypes,
		"top_referrers":  topReferrers,
		"top_pages":      topPages,
		"browsers":       browsers,
		"oses":           oses,
		"languages":      languages,
		"countries":      countries,
		"bot_daily":      botDaily,
		"hourly":         hourly,
	}, nil
}
