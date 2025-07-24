package handler

import (
	"linkvault/internal/response"
	"linkvault/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ClickHandler struct {
	serviceClick *service.ClickService
	serviceLink  *service.ShortLinkService
}

func NewClickHandler(serviceClick *service.ClickService, serviceLink *service.ShortLinkService) *ClickHandler {
	return &ClickHandler{serviceClick: serviceClick, serviceLink: serviceLink}
}

// GetLinkStats godoc
// @Summary Статистика по короткой ссылке
// @Description Получить статистику: количество переходов, уникальные IP, география, график по дням
// @Tags links
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "ID короткой ссылки"
// @Success 200 {object} response.LinkStatsResponse "Статистика по ссылке"
// @Failure 400 {object} response.ErrorResponse "Ошибка валидации"
// @Failure 404 {object} response.ErrorResponse "Ссылка не найдена"
// @Failure 500 {object} response.ErrorResponse "Ошибка сервера"
// @Router /links/{id}/stats [get]
func (h *ClickHandler) GetLinkStats(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "id is required"})
		return
	}

	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.ErrorResponse{Error: "Missing or invalid token"})
		return
	}
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusUnauthorized, response.ErrorResponse{Error: "Invalid user id"})
		return
	}

	shortLink, err := h.serviceLink.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, response.ErrorResponse{Error: "Short link not found"})
		return
	}
	if shortLink.UserID == nil || *shortLink.UserID != userID {
		c.JSON(http.StatusForbidden, response.ErrorResponse{Error: "Access denied"})
		return
	}

	stats, err := h.serviceClick.GetStats(id)
	if err != nil {
		c.JSON(http.StatusNotFound, response.ErrorResponse{Error: "Short link not found or no stats"})
		return
	}

	c.JSON(http.StatusOK, response.LinkStatsResponse{Stats: stats})
}
