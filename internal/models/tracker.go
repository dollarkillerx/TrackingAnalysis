package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Tracker struct {
	ID        string    `gorm:"type:uuid;primaryKey" json:"id"`
	Type      string    `gorm:"type:varchar(10);not null" json:"type"` // "ad" or "web"
	Name      string    `gorm:"type:varchar(255);not null" json:"name"`
	Mode      string    `gorm:"type:varchar(10);not null;default:'302'" json:"mode"` // "js" or "302"
	Status    string    `gorm:"type:varchar(20);not null;default:'active'" json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (t *Tracker) BeforeCreate(tx *gorm.DB) error {
	if t.ID == "" {
		t.ID = uuid.New().String()
	}
	return nil
}
