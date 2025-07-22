package service

import (
	"encoding/json"
	"linkvault/internal/models"
	"linkvault/internal/repository"
	"net/http"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ClickService struct {
	repo *repository.ClickRepository
	log  *zap.Logger
}

func NewClickService(repo *repository.ClickRepository, log *zap.Logger) *ClickService {
	return &ClickService{
		repo: repo,
		log:  log,
	}
}

func (s *ClickService) CreateClick(shortLinkID uuid.UUID, ip, userAgent string) error {
	click := &models.Click{
		ShortLinkID: shortLinkID,
		IP:          ip,
		UserAgent:   userAgent,
		ClickedAt:   time.Now(),
	}
	resp, err := http.Get("http://ip-api.com/json/" + ip)
	if err == nil {
		defer resp.Body.Close()
		var data struct {
			Country    string `json:"country"`
			RegionName string `json:"regionName"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&data); err == nil {
			click.Country = data.Country
			click.Region = data.RegionName
		}
	}

	return s.repo.Create(click)
}
