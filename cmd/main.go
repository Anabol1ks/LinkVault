package main

import (
	"linkvault/internal/config"
	"linkvault/internal/db"
	"linkvault/internal/logger"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load(".env")
	isDev := os.Getenv("ENV") == "development"
	if err := logger.Init(isDev); err != nil {
		panic(err)
	}
	defer logger.Sync()

	log := logger.L()

	cfg := config.Load(log)

	db.ConnectDB(&cfg.DB, log)

}
