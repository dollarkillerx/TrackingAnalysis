package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Channel struct {
	ID         string    `gorm:"type:uuid;primaryKey" json:"id"`
	TrackerID  string    `gorm:"type:uuid;not null;index" json:"tracker_id"`
	CampaignID string    `gorm:"type:uuid;not null;index" json:"campaign_id"`
	Name       string    `gorm:"type:varchar(255);not null" json:"name"`
	Source     string    `gorm:"type:varchar(255)" json:"source"`
	Medium     string    `gorm:"type:varchar(255)" json:"medium"`
	Tags       JSONMap   `gorm:"type:jsonb" json:"tags"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Tracker    Tracker   `gorm:"foreignKey:TrackerID" json:"-"`
	Campaign   Campaign  `gorm:"foreignKey:CampaignID" json:"-"`
}

func (ch *Channel) BeforeCreate(tx *gorm.DB) error {
	if ch.ID == "" {
		ch.ID = uuid.New().String()
	}
	return nil
}
