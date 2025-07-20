package main

import (
	"linkvault/internal/config"
	"linkvault/internal/logger"
	"linkvault/internal/router"
	"linkvault/internal/storage"
	"os"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

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

	r := router.Router(db, log)
	if err := r.Run(); err != nil {
		log.Fatal("Не удалось запустить сервер", zap.Error(err))
	}
}
