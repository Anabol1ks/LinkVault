package croncleaner

import (
	"linkvault/internal/models"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Очистка старых ссылок и кликов
func CleanOldLinksAndClicks(db *gorm.DB, log *zap.Logger) {
	// 1. Удаляем ссылки без владельца, если истёк срок
	var anonLinks []models.ShortLink
	db.Where("user_id IS NULL AND expire_at IS NOT NULL AND expire_at < ?", time.Now()).Find(&anonLinks)
	for _, link := range anonLinks {
		// Удаляем клики
		db.Where("short_link_id = ?", link.ID).Delete(&models.Click{})
		// Удаляем саму ссылку
		db.Delete(&link)
		log.Info("Удалена анонимная истёкшая ссылка", zap.String("short_code", link.ShortCode))
	}

	// 2. Удаляем деактивированные ссылки без владельца (срок истёк)
	var anonInactiveLinks []models.ShortLink
	db.Where("user_id IS NULL AND is_active = false AND expire_at IS NOT NULL AND expire_at < ?", time.Now()).Find(&anonInactiveLinks)
	for _, link := range anonInactiveLinks {
		db.Where("short_link_id = ?", link.ID).Delete(&models.Click{})
		db.Delete(&link)
		log.Info("Удалена анонимная деактивированная ссылка", zap.String("short_code", link.ShortCode))
	}

	// 4. Для ссылок с владельцем, удаляем только если истёк срок и деактивированы
	var userLinks []models.ShortLink
	db.Where("user_id IS NOT NULL AND is_active = false AND expire_at IS NOT NULL AND expire_at < ?", time.Now()).Find(&userLinks)
	for _, link := range userLinks {
		db.Where("short_link_id = ?", link.ID).Delete(&models.Click{})
		db.Delete(&link)
		log.Info("Удалена пользовательская деактивированная истёкшая ссылка", zap.String("short_code", link.ShortCode))
	}
}
