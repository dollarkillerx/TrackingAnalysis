package handler

import (
	"crypto/rsa"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
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
	PubKey     *rsa.PublicKey
	PrivKey    *rsa.PrivateKey
}

// GET /t/:token — JS-based click tracking page
func (h *TrackingHandler) HandleJSTrack(c *gin.Context) {
	token := c.Param("token")
	payload, err := security.VerifyToken(token, h.Config.SecurityConfiguration.TokenSecret)
	if err != nil {
		c.String(http.StatusBadRequest, "invalid or expired token")
		return
	}

	pubPEM, err := security.PublicKeyPEM(h.PubKey)
	if err != nil {
		c.String(http.StatusInternalServerError, "server error")
		return
	}

	// Look up target URL
	target, err := h.TargetRepo.GetByID(payload.TargetID)
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
	payload, err := security.VerifyToken(token, h.Config.SecurityConfiguration.TokenSecret)
	if err != nil {
		c.String(http.StatusBadRequest, "invalid or expired token")
		return
	}

	// Look up target URL from DB (never from token)
	target, err := h.TargetRepo.GetByID(payload.TargetID)
	if err != nil {
		c.String(http.StatusNotFound, "target not found")
		return
	}

	// Record click
	click := &models.Click{
		TS:         time.Now(),
		TrackerID:  payload.TrackerID,
		CampaignID: payload.CampaignID,
		ChannelID:  payload.ChannelID,
		TargetID:   payload.TargetID,
		IP:         c.ClientIP(),
		UA:         c.GetHeader("User-Agent"),
		Lang:       c.GetHeader("Accept-Language"),
		Referer:    c.GetHeader("Referer"),
	}
	h.ClickRepo.Create(click)

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
