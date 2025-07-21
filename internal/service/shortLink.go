package service

import (
	"errors"
	"linkvault/internal/models"
	"linkvault/internal/repository"
	"time"

	"github.com/google/uuid"
	"github.com/teris-io/shortid"
	"go.uber.org/zap"
)

type ShortLinkService struct {
	repo repository.ShortLinkRepository
	log  *zap.Logger
}

func NewShortLinkService(repo repository.ShortLinkRepository, log *zap.Logger) *ShortLinkService {
	return &ShortLinkService{
		repo: repo,
		log:  log,
	}
}

var ErrGenerateShortCode = errors.New("error generating short code")
var ErrCreateShortLink = errors.New("error creating short link")

func (s *ShortLinkService) CreateShortLink(originalURL string, userID *uuid.UUID, expireAt *time.Time) (*models.ShortLink, error) {
	var finalExpireAt *time.Time
	if userID == nil {
		exp := time.Now().Add(7 * 24 * time.Hour)
		finalExpireAt = &exp
	} else {
		if expireAt != nil && !expireAt.IsZero() {
			finalExpireAt = expireAt
		} else {
			finalExpireAt = nil
		}
	}

	shortCode, err := generateShortCode()
	if err != nil {
		s.log.Error("Failed to generate short code", zap.Error(err))
		return nil, ErrGenerateShortCode
	}

	shortLink := &models.ShortLink{
		OriginalURL: originalURL,
		UserID:      userID,
		ShortCode:   shortCode,
		IsActive:    true,
		ExpireAt:    finalExpireAt,
	}

	if err := s.repo.Create(shortLink); err != nil {
		s.log.Error("Failed to create short link", zap.Error(err))
		return nil, err
	}

	return shortLink, nil
}

func generateShortCode() (string, error) {
	id, err := shortid.Generate()
	if err != nil {
		return "", err
	}
	return id, nil
}
