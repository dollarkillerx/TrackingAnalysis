package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Site struct {
	ID        string    `gorm:"type:uuid;primaryKey" json:"id"`
	Name      string    `gorm:"type:varchar(255);not null" json:"name"`
	Domain    string    `gorm:"type:varchar(255);not null" json:"domain"`
	SiteKey   string    `gorm:"type:varchar(64);uniqueIndex;not null" json:"site_key"`
	Status    string    `gorm:"type:varchar(20);not null;default:'active'" json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (s *Site) BeforeCreate(tx *gorm.DB) error {
	if s.ID == "" {
		s.ID = uuid.New().String()
	}
	if s.SiteKey == "" {
		s.SiteKey = uuid.New().String()
	}
	return nil
}
