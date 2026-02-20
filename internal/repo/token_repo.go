package repo

import (
	"github.com/tracking/analysis/internal/models"
	"gorm.io/gorm"
)

type TokenRepo struct {
	DB *gorm.DB
}

func NewTokenRepo(db *gorm.DB) *TokenRepo {
	return &TokenRepo{DB: db}
}

func (r *TokenRepo) Create(t *models.Token) error {
	return r.DB.Create(t).Error
}

func (r *TokenRepo) GetByShortCode(shortCode string) (*models.Token, error) {
	var t models.Token
	err := r.DB.First(&t, "short_code = ?", shortCode).Error
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *TokenRepo) List(trackerID string) ([]models.Token, error) {
	var tokens []models.Token
	q := r.DB
	if trackerID != "" {
		q = q.Where("tracker_id = ?", trackerID)
	}
	err := q.Order("created_at DESC").Find(&tokens).Error
	return tokens, err
}

func (r *TokenRepo) Delete(id string) error {
	return r.DB.Delete(&models.Token{}, "id = ?", id).Error
}

func (r *TokenRepo) ExistsByShortCode(shortCode string) (bool, error) {
	var count int64
	err := r.DB.Model(&models.Token{}).Where("short_code = ?", shortCode).Count(&count).Error
	return count > 0, err
}
