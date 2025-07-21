package repository

import (
	"linkvault/internal/models"

	"gorm.io/gorm"
)

type ShortLinkRepository struct {
	db *gorm.DB
}

func NewShortLinkRepository(db *gorm.DB) *ShortLinkRepository {
	return &ShortLinkRepository{
		db: db,
	}
}

func (r *ShortLinkRepository) Create(shortLink *models.ShortLink) error {
	return r.db.Create(shortLink).Error
}
