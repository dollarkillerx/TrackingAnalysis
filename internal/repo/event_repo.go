package repo

import (
	"github.com/tracking/analysis/internal/models"
	"gorm.io/gorm"
)

type EventRepo struct {
	DB *gorm.DB
}

func NewEventRepo(db *gorm.DB) *EventRepo {
	return &EventRepo{DB: db}
}

func (r *EventRepo) BatchCreate(events []models.Event) error {
	return r.DB.CreateInBatches(events, 100).Error
}
