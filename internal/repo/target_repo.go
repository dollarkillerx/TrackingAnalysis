package repo

import (
	"github.com/tracking/analysis/internal/models"
	"gorm.io/gorm"
)

type TargetRepo struct {
	DB *gorm.DB
}

func NewTargetRepo(db *gorm.DB) *TargetRepo {
	return &TargetRepo{DB: db}
}

func (r *TargetRepo) Create(t *models.Target) error {
	return r.DB.Create(t).Error
}

func (r *TargetRepo) List(trackerID string) ([]models.Target, error) {
	var targets []models.Target
	q := r.DB
	if trackerID != "" {
		q = q.Where("tracker_id = ?", trackerID)
	}
	err := q.Order("created_at DESC").Find(&targets).Error
	return targets, err
}

func (r *TargetRepo) GetByID(id string) (*models.Target, error) {
	var t models.Target
	err := r.DB.First(&t, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &t, nil
}
