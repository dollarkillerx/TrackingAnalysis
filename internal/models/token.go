package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Token struct {
	ID         string    `gorm:"type:uuid;primaryKey" json:"id"`
	ShortCode  string    `gorm:"type:varchar(8);uniqueIndex;not null" json:"short_code"`
	TrackerID  string    `gorm:"type:uuid;not null;index" json:"tracker_id"`
	CampaignID string    `gorm:"type:uuid" json:"campaign_id"`
	ChannelID  string    `gorm:"type:uuid" json:"channel_id"`
	TargetID   string    `gorm:"type:uuid;not null" json:"target_id"`
	Mode       string    `gorm:"type:varchar(10);not null;default:'302'" json:"mode"`
	CreatedAt  time.Time `json:"created_at"`
}

func (t *Token) BeforeCreate(tx *gorm.DB) error {
	if t.ID == "" {
		t.ID = uuid.New().String()
	}
	return nil
}
