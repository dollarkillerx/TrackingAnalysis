package repo

import (
	"github.com/tracking/analysis/internal/models"
	"gorm.io/gorm"
)

type SiteRepo struct {
	DB *gorm.DB
}

func NewSiteRepo(db *gorm.DB) *SiteRepo {
	return &SiteRepo{DB: db}
}

func (r *SiteRepo) Create(s *models.Site) error {
	return r.DB.Create(s).Error
}

func (r *SiteRepo) List() ([]models.Site, error) {
	var sites []models.Site
	err := r.DB.Order("created_at DESC").Find(&sites).Error
	return sites, err
}

func (r *SiteRepo) GetByKey(siteKey string) (*models.Site, error) {
	var s models.Site
	err := r.DB.First(&s, "site_key = ?", siteKey).Error
	if err != nil {
		return nil, err
	}
	return &s, nil
}
