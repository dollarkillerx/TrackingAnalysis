package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Event struct {
	ID           string    `gorm:"type:uuid;primaryKey" json:"id"`
	TS           time.Time `gorm:"not null;index:idx_events_site_ts;index:idx_events_type_ts" json:"ts"`
	SiteID       string    `gorm:"type:uuid;not null;index:idx_events_site_ts" json:"site_id"`
	Type         string    `gorm:"type:varchar(50);not null;index:idx_events_type_ts" json:"type"`
	VisitorID    string    `gorm:"type:varchar(255)" json:"visitor_id"`
	SessionID    string    `gorm:"type:varchar(255)" json:"session_id"`
	URL          string    `gorm:"type:text" json:"url"`
	Title        string    `gorm:"type:text" json:"title"`
	Referrer     string    `gorm:"type:text" json:"referrer"`
	IP           string    `gorm:"type:varchar(45)" json:"ip"`
	Country      string    `gorm:"type:varchar(2)" json:"country"`
	UA           string    `gorm:"type:text" json:"ua"`
	Lang         string    `gorm:"type:varchar(50)" json:"lang"`
	Props        JSONMap   `gorm:"type:jsonb" json:"props"`
	SuspectedBot bool      `gorm:"default:false" json:"suspected_bot"`
	IsBot        bool      `gorm:"default:false" json:"is_bot"`
	CreatedAt    time.Time `json:"created_at"`
}

func (e *Event) BeforeCreate(tx *gorm.DB) error {
	if e.ID == "" {
		e.ID = uuid.New().String()
	}
	return nil
}
