package handler

import (
	"errors"
	"linkvault/internal/response"
	"linkvault/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

type UserRegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// Register godoc
// @Summary Регистрация пользователя
// @Description Регистрация нового пользователя
// @Tags users
// @Accept json
// @Produce json
// @Param user body UserRegisterRequest true "Параметры регистрации пользователя"
// @Success 201 {object} response.UserRegisterResponse "Успешная регистрация пользователя"
// @Failure 400 {object} response.ErrorResponse "Ошибка валидации"
// @Failure 409 {object} response.ErrorResponse "Пользователь уже существует"
// @Failure 500 {object} response.ErrorResponse "Ошибка сервера"
// @Router /auth/register [post]
func (h *UserHandler) Register(c *gin.Context) {
	var req UserRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Ошибка валидации"})
		return
	}

	user, err := h.service.Register(req.Name, req.Email, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrUserExists):
			c.JSON(http.StatusConflict, response.ErrorResponse{Error: "Пользователь уже существует"})
		default:
			c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: "Ошибка сервера"})
		}
		return
	}

	resp := response.UserRegisterResponse{
		ID:    user.ID.String(),
		Name:  user.Name,
		Email: user.Email,
	}

	c.JSON(http.StatusCreated, resp)
}
