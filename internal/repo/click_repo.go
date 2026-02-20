package repo

import (
	"github.com/tracking/analysis/internal/models"
	"gorm.io/gorm"
)

type ClickRepo struct {
	DB *gorm.DB
}

func NewClickRepo(db *gorm.DB) *ClickRepo {
	return &ClickRepo{DB: db}
}

func (r *ClickRepo) Create(c *models.Click) error {
	return r.DB.Create(c).Error
}

func (r *ClickRepo) GetByID(id string) (*models.Click, error) {
	var c models.Click
	err := r.DB.First(&c, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &c, nil
}
