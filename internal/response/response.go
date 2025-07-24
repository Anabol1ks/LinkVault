package response

import (
	"time"

	"github.com/google/uuid"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type UserRegisterResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type SuccessShortLinkResponse struct {
	ID          uuid.UUID  `json:"id"`
	ShortURL    string     `json:"short_url"`
	OriginalURL string     `json:"original_url"`
	ExpireAt    *time.Time `json:"expire_at,omitempty"`
}

// Подробная структура для статистики по ссылке
// Используется для Swagger и JSON-ответа
// total — общее количество кликов
// unique_ip_count — количество уникальных IP
// unique_ips — список уникальных IP
// countries_count — количество стран
// countries — список стран
// countries_stats — карта страна: количество
// daily_stats — карта дата: количество

// swagger:model DetailedLinkStats
// swagger:response LinkStatsResponse
// @name DetailedLinkStats
// @description Подробная статистика по короткой ссылке
// @property total int64 "Общее количество кликов"
// @property unique_ip_count int64 "Количество уникальных IP"
// @property unique_ips []string "Список уникальных IP"
// @property countries_count int "Количество стран"
// @property countries []string "Список стран"
// @property countries_stats map[string]int64 "Статистика по странам"
// @property daily_stats map[string]int64 "Статистика по дням"
type DetailedLinkStats struct {
	Total          int64            `json:"total"`
	UniqueIPCount  int64            `json:"unique_ip_count"`
	UniqueIPs      []string         `json:"unique_ips"`
	CountriesCount int              `json:"countries_count"`
	Countries      []string         `json:"countries"`
	CountriesStats map[string]int64 `json:"countries_stats"`
	DailyStats     map[string]int64 `json:"daily_stats"`
}

type LinkStatsResponse struct {
	Stats DetailedLinkStats `json:"stats"`
}

type ShortLinkListResponse struct {
	Links []*SuccessShortLinkResponse `json:"links"`
}

// DTO для ответа по кликам
type ClickResponse struct {
	ID        string    `json:"id"`
	IP        string    `json:"ip"`
	UserAgent string    `json:"user_agent"`
	ClickedAt time.Time `json:"clicked_at"`
	Country   string    `json:"country"`
	Region    string    `json:"region"`
}
