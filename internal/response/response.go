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

type LinkStatsResponse struct {
	Stats map[string]interface{} `json:"stats"`
}

type ShortLinkListResponse struct {
	Links []*SuccessShortLinkResponse `json:"links"`
}
