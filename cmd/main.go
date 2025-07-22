package main

import (
	_ "linkvault/docs"
	"linkvault/internal/config"
	"linkvault/internal/handler"
	"linkvault/internal/logger"
	"linkvault/internal/repository"
	"linkvault/internal/router"
	"linkvault/internal/service"
	"linkvault/internal/storage"
	"os"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

// @Title TZ_OZON API
// @Version 1.0
// @securityDefinitions.apikey	BearerAuth
// @in							header
// @name						Authorization
func main() {
	_ = godotenv.Load()
	isDev := os.Getenv("ENV") == "development"
	if err := logger.Init(isDev); err != nil {
		panic(err)
	}
	defer logger.Sync()

	log := logger.L()

	cfg := config.Load(log)

	db, err := storage.ConnectDB(&cfg.DB, log)
	if err != nil {
		log.Fatal("Не удалось подключиться к базе данных", zap.Error(err))
		return
	}

	storage.Migrate(db, log)

	// 1. Репозитории
	userRepo := repository.NewUserRepository(db)
	shortLinkRepo := repository.NewShortLinkRepository(db)
	clickRepo := repository.NewClickRepository(db)

	// 2. Сервисы
	userService := service.NewUserService(userRepo, log, cfg)
	shortLinkService := service.NewShortLinkService(shortLinkRepo, log)
	clickService := service.NewClickService(clickRepo, log)

	// 3. Хендлеры
	userHandler := handler.NewUserHandler(userService)
	linkHandler := handler.NewShortLinkHandler(shortLinkService, clickService, cfg)

	// 4. Handlers для роутера
	handlers := &router.Handlers{
		User: userHandler,
		Link: linkHandler,
	}

	r := router.Router(db, log, handlers, cfg)
	if err := r.Run(); err != nil {
		log.Fatal("Не удалось запустить сервер", zap.Error(err))
	}
}
