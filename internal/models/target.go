package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Target struct {
	ID        string    `gorm:"type:uuid;primaryKey" json:"id"`
	TrackerID string    `gorm:"type:uuid;not null;index" json:"tracker_id"`
	URL       string    `gorm:"type:text;not null" json:"url"`
	CreatedAt time.Time `json:"created_at"`
	Tracker   Tracker   `gorm:"foreignKey:TrackerID" json:"-"`
}

func (t *Target) BeforeCreate(tx *gorm.DB) error {
	if t.ID == "" {
		t.ID = uuid.New().String()
	}
	return nil
}
