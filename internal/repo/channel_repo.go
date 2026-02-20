package repo

import (
	"github.com/tracking/analysis/internal/models"
	"gorm.io/gorm"
)

type ChannelRepo struct {
	DB *gorm.DB
}

func NewChannelRepo(db *gorm.DB) *ChannelRepo {
	return &ChannelRepo{DB: db}
}

func (r *ChannelRepo) Create(ch *models.Channel) error {
	return r.DB.Create(ch).Error
}

func (r *ChannelRepo) BatchImport(channels []models.Channel) error {
	return r.DB.CreateInBatches(channels, 100).Error
}

func (r *ChannelRepo) List(trackerID, campaignID string) ([]models.Channel, error) {
	var channels []models.Channel
	q := r.DB
	if trackerID != "" {
		q = q.Where("tracker_id = ?", trackerID)
	}
	if campaignID != "" {
		q = q.Where("campaign_id = ?", campaignID)
	}
	err := q.Order("created_at DESC").Find(&channels).Error
	return channels, err
}
