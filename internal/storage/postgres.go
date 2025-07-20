package storage

import (
	"fmt"
	"linkvault/internal/config"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB(cfg *config.DBConfig, log *zap.Logger) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name, cfg.SSL)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{PrepareStmt: false})
	if err != nil {
		log.Fatal("Не удалось подключиться к базе данных", zap.Error(err))
		return nil, err
	}

	log.Info("Подключение к базе данных успешно установлено")
	return db, nil
}
