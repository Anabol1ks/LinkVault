package handler

import (
	"linkvault/internal/response"
	"linkvault/internal/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ShortLinkHandler struct {
	service *service.ShortLinkService
}

func NewShortLinkHandler(service *service.ShortLinkService) *ShortLinkHandler {
	return &ShortLinkHandler{
		service: service,
	}
}

type CreateShortLinkRequest struct {
	OriginalURL string    `json:"original_url" binding:"required,url"`
	ExpireAt    time.Time `json:"expire_at"`
}

func (h *ShortLinkHandler) CreateShortLink(c *gin.Context) {
	var req CreateShortLinkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Ошибка валидации"})
		return
	}

	userIDStr, _ := c.GetQuery("user_id")
	userID, _ := uuid.Parse(userIDStr)

	shortLink, err := h.service.CreateShortLink(req.OriginalURL, userID, &req.ExpireAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: "Ошибка создания короткой ссылки"})
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse{Data: shortLink})
}
