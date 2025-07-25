package croncleaner

import (
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func StartCleanerCron(db *gorm.DB, log *zap.Logger, schedule string) {
	c := cron.New()
	_, err := c.AddFunc(schedule, func() {
		CleanOldLinksAndClicks(db, log)
	})
	if err != nil {
		log.Error("Ошибка запуска CRON для очистки", zap.Error(err))
		return
	}
	c.Start()
	log.Info("CRON очистки старых ссылок и кликов запущен", zap.String("schedule", schedule))
}
