package service

import (
	"errors"
	"linkvault/internal/models"
	"linkvault/internal/repository"

	"golang.org/x/crypto/bcrypt"

	"go.uber.org/zap"
)

var ErrUserExists = errors.New("user already exists")

type UserService struct {
	repo *repository.UserRepository
	log  *zap.Logger
}

func NewUserService(repo *repository.UserRepository, log *zap.Logger) *UserService {
	return &UserService{
		repo: repo,
		log:  log,
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
