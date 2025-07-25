package service

import (
	"errors"
	"linkvault/internal/config"
	"linkvault/internal/jwt"
	"linkvault/internal/models"
	"linkvault/internal/repository"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"go.uber.org/zap"
)

var ErrUserExists = errors.New("user already exists")

type UserService struct {
	repo *repository.UserRepository
	log  *zap.Logger
	cfg  *config.Config
}

func NewUserService(repo *repository.UserRepository, log *zap.Logger, cfg *config.Config) *UserService {
	return &UserService{
		repo: repo,
		log:  log,
		cfg:  cfg,
	}
}

func (s *UserService) Register(name, email, password string) (*models.User, error) {
	if _, err := s.repo.FindByEmail(email); err == nil {
		s.log.Warn("User already exists", zap.String("email", email))
		return nil, ErrUserExists
	}

	hashedPassword, err := hashedPassword(password)
	if err != nil {
		s.log.Error("Failed to hash password", zap.Error(err))
		return nil, err
	}

	user := &models.User{
		Name:         name,
		Email:        email,
		PasswordHash: hashedPassword,
	}

	if err := s.repo.Create(user); err != nil {
		s.log.Error("Failed to register user", zap.Error(err))
		return nil, err
	}

	return user, nil
}

func hashedPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

var ErrUserNotFound = errors.New("user not found")
var ErrInvalidPassword = errors.New("invalid password")

func (s *UserService) Login(email, password string) (access, refresh string, err error) {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		s.log.Warn("User not found", zap.String("email", email))
		return "", "", ErrUserNotFound
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		s.log.Warn("Invalid password", zap.String("email", email))
		return "", "", ErrInvalidPassword
	}

	access, err = jwt.GenerateAccessToken(user.ID.String(), &s.cfg.JWT)
	if err != nil {
		s.log.Error("Failed to generate access token", zap.Error(err))
		return "", "", err
	}

	refresh, err = jwt.GenerateRefreshToken(user.ID.String(), &s.cfg.JWT)
	if err != nil {
		s.log.Error("Failed to generate refresh token", zap.Error(err))
		return "", "", err
	}

	return access, refresh, nil
}

var ErrInvalidToken = errors.New("invalid token")

func (s *UserService) Refresh(refreshToken string) (access, refresh string, err error) {
	claims, err := jwt.ParseRefreshToken(refreshToken, s.cfg.JWT.Refresh)
	if err != nil {
		return "", "", ErrInvalidToken
	}

	access, err = jwt.GenerateAccessToken(claims.UserID, &s.cfg.JWT)
	if err != nil {
		return "", "", err
	}

	refresh, err = jwt.GenerateRefreshToken(claims.UserID, &s.cfg.JWT)
	if err != nil {
		return "", "", err
	}

	return access, refresh, nil
}

func (s *UserService) Profile(userID uuid.UUID) (*models.User, error) {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}
