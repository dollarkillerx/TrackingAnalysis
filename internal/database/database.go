package database

import (
	"github.com/tracking/analysis/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Init(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(
		&models.Tracker{},
		&models.Campaign{},
		&models.Channel{},
		&models.Target{},
		&models.Site{},
		&models.Click{},
		&models.Event{},
	)
	if err != nil {
		return nil, err
	}
	return db, nil
}
