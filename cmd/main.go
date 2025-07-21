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

	userHandler := handler.NewUserHandler(service.NewUserService(repository.NewUserRepository(db), log))
	handlers := &router.Handlers{
		User: userHandler,
	}

	r := router.Router(db, log, handlers)
	if err := r.Run(); err != nil {
		log.Fatal("Не удалось запустить сервер", zap.Error(err))
	}
}
