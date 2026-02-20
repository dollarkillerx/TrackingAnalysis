package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Click struct {
	ID           string    `gorm:"type:uuid;primaryKey" json:"id"`
	TS           time.Time `gorm:"not null;index:idx_clicks_tracker_ts;index:idx_clicks_channel_ts" json:"ts"`
	TrackerID    string    `gorm:"type:uuid;not null;index:idx_clicks_tracker_ts" json:"tracker_id"`
	CampaignID   string    `gorm:"type:uuid" json:"campaign_id"`
	ChannelID    string    `gorm:"type:uuid;index:idx_clicks_channel_ts" json:"channel_id"`
	TargetID     string    `gorm:"type:uuid" json:"target_id"`
	VisitorID    string    `gorm:"type:varchar(255)" json:"visitor_id"`
	IP           string    `gorm:"type:varchar(45)" json:"ip"`
	UA           string    `gorm:"type:text" json:"ua"`
	Lang         string    `gorm:"type:varchar(50)" json:"lang"`
	Referer      string    `gorm:"type:text" json:"referer"`
	Props        JSONMap   `gorm:"type:jsonb" json:"props"`
	SuspectedBot bool      `gorm:"default:false" json:"suspected_bot"`
	IsBot        bool      `gorm:"default:false" json:"is_bot"`
	CreatedAt    time.Time `json:"created_at"`
}

func (c *Click) BeforeCreate(tx *gorm.DB) error {
	if c.ID == "" {
		c.ID = uuid.New().String()
	}
	return nil
}
