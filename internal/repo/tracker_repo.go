package repo

import (
	"github.com/tracking/analysis/internal/models"
	"gorm.io/gorm"
)

type TrackerRepo struct {
	DB *gorm.DB
}

func NewTrackerRepo(db *gorm.DB) *TrackerRepo {
	return &TrackerRepo{DB: db}
}

func (r *TrackerRepo) Create(t *models.Tracker) error {
	return r.DB.Create(t).Error
}

func (r *TrackerRepo) List(trackerType string) ([]models.Tracker, error) {
	var trackers []models.Tracker
	q := r.DB
	if trackerType != "" {
		q = q.Where("type = ?", trackerType)
	}
	err := q.Order("created_at DESC").Find(&trackers).Error
	return trackers, err
}

func (r *TrackerRepo) GetByID(id string) (*models.Tracker, error) {
	var t models.Tracker
	err := r.DB.First(&t, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *TrackerRepo) Update(t *models.Tracker) error {
	return r.DB.Save(t).Error
}

func (r *TrackerRepo) Delete(id string) error {
	return r.DB.Delete(&models.Tracker{}, "id = ?", id).Error
}
