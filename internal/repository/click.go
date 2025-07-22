package repository

import (
	"linkvault/internal/models"

	"gorm.io/gorm"
)

type ClickRepository struct {
	db *gorm.DB
}

func NewClickRepository(db *gorm.DB) *ClickRepository {
	return &ClickRepository{
		db: db,
	}
}

func (r *ClickRepository) Create(click *models.Click) error {
	return r.db.Create(click).Error
}

func (r *ClickRepository) GetByShortLinkID(shortLinkID string) ([]models.Click, error) {
	var clicks []models.Click
	err := r.db.Where("short_link_id = ?", shortLinkID).Find(&clicks).Error
	return clicks, err
}
