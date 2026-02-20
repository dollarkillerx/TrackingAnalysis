package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Campaign struct {
	ID        string    `gorm:"type:uuid;primaryKey" json:"id"`
	TrackerID string    `gorm:"type:uuid;not null;index" json:"tracker_id"`
	Name      string    `gorm:"type:varchar(255);not null" json:"name"`
	Status    string    `gorm:"type:varchar(20);not null;default:'active'" json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Tracker   Tracker   `gorm:"foreignKey:TrackerID" json:"-"`
}

func (c *Campaign) BeforeCreate(tx *gorm.DB) error {
	if c.ID == "" {
		c.ID = uuid.New().String()
	}
	return nil
}
