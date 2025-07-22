package handler

import (
	"fmt"
	"linkvault/internal/config"
	"linkvault/internal/response"
	"linkvault/internal/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ShortLinkHandler struct {
	service      *service.ShortLinkService
	clickService *service.ClickService
	cfg          *config.Config
}

func NewShortLinkHandler(service *service.ShortLinkService, clickService *service.ClickService, cfg *config.Config) *ShortLinkHandler {
	return &ShortLinkHandler{
		service:      service,
		clickService: clickService,
		cfg:          cfg,
	}
}

type CreateShortLinkRequest struct {
	OriginalURL string `json:"original_url" binding:"required,url"`
	ExpireAfter string `json:"expire_after" binding:"omitempty"` // например, "2h", "30m", "7d"
}

// CreateShortLink godoc
// @Summary Create a short link
// @Description Create a short link
// @Tags links
// @Accept json
// @Produce json
// @Param request body CreateShortLinkRequest true "CreateShortLinkRequest"
// @Success 200 {object} response.SuccessShortLinkResponse "Успешное создание короткой ссылки"
// @Failure 400 {object} response.ErrorResponse "Ошибка валидации"
// @Failure 500 {object} response.ErrorResponse "Ошибка создания короткой ссылки"
// @Router /links [post]
func (h *ShortLinkHandler) CreateShortLink(c *gin.Context) {
	var req CreateShortLinkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Ошибка валидации"})
		return
	}

	userIDStr, _ := c.GetQuery("user_id")
	var userID *uuid.UUID
	if userIDStr != "" {
		parsed, err := uuid.Parse(userIDStr)
		if err == nil {
			userID = &parsed
		}
	}

	var expireAfter *time.Duration
	if req.ExpireAfter != "" {
		d, err := time.ParseDuration(req.ExpireAfter)
		if err != nil {
			c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Некорректный формат expire_after"})
			return
		}
		expireAfter = &d
	}

	shortLink, err := h.service.CreateShortLink(req.OriginalURL, userID, expireAfter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: "Ошибка создания короткой ссылки"})
		return
	}

	shortURL := fmt.Sprintf("%s/%s", h.cfg.Domain, shortLink.ShortCode)

	c.JSON(http.StatusOK, response.SuccessShortLinkResponse{
		ShortURL:    shortURL,
		OriginalURL: shortLink.OriginalURL,
		ExpireAt:    shortLink.ExpireAt,
	})
}

// GetOriginalURL godoc
// @Summary Перенаправление
// @Description Перенаправление на оригинальный URL по shortCode
// @Tags links
// @Accept json
// @Produce json
// @Param shortCode path string true "shortCode"
// @Success 302 {object} string "Успешное получение оригинального URL"
// @Failure 400 {object} response.ErrorResponse "Ошибка валидации"
// @Failure 404 {object} response.ErrorResponse "Короткая ссылка не найдена или неактивна/истекла"
// @Router /{shortCode} [get]
func (h *ShortLinkHandler) GetOriginalURL(c *gin.Context) {
	shortCode := c.Param("shortCode")
	if shortCode == "" {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "shortCode is required"})
		return
	}

	shortLink, err := h.service.GetShortLinkByCode(shortCode)
	if err != nil {
		c.JSON(http.StatusNotFound, response.ErrorResponse{Error: "Short link not found or inactive/expired"})
		return
	}
	originalURL := shortLink.OriginalURL

	if shortLink.UserID != nil {
		ip := c.ClientIP()
		userAgent := c.Request.UserAgent()
		go func() {
			_ = h.clickService.CreateClick(shortLink.ID, ip, userAgent)
		}()
	}

	c.Redirect(http.StatusFound, originalURL)
}
