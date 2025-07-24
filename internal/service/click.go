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

func (s *ClickService) GetStats(shortLinkID string) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Общее количество переходов
	total, err := s.repo.GetCount(shortLinkID)
	if err != nil {
		return nil, err
	}
	stats["total"] = total

	// Уникальные IP
	uniqueIP, err := s.repo.GetUniqueIPCount(shortLinkID)
	if err != nil {
		return nil, err
	}
	stats["unique_ip"] = uniqueIP

	// География по странам
	countryStats, err := s.repo.GetCountryStats(shortLinkID)
	if err != nil {
		return nil, err
	}
	stats["countries"] = countryStats

	// График по дням
	dailyStats, err := s.repo.GetDailyStats(shortLinkID)
	if err != nil {
		return nil, err
	}
	stats["daily"] = dailyStats

	return stats, nil
}
