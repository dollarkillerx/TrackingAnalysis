package repo

import (
	"github.com/tracking/analysis/internal/models"
	"gorm.io/gorm"
)

type CampaignRepo struct {
	DB *gorm.DB
}

func NewCampaignRepo(db *gorm.DB) *CampaignRepo {
	return &CampaignRepo{DB: db}
}

func (r *CampaignRepo) Create(c *models.Campaign) error {
	return r.DB.Create(c).Error
}

func (r *CampaignRepo) List(trackerID string) ([]models.Campaign, error) {
	var campaigns []models.Campaign
	q := r.DB
	if trackerID != "" {
		q = q.Where("tracker_id = ?", trackerID)
	}
	err := q.Order("created_at DESC").Find(&campaigns).Error
	return campaigns, err
}

func (r *CampaignRepo) GetByID(id string) (*models.Campaign, error) {
	var c models.Campaign
	err := r.DB.First(&c, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &c, nil
}
