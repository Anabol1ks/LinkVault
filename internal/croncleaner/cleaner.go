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

	// 3. Для ссылок без владельца, статистику (клики) оставляем на неделю после истечения срока
	weekAgo := time.Now().Add(-7 * 24 * time.Hour)
	db.Joins("JOIN short_links ON short_links.id = clicks.short_link_id").
		Where("short_links.user_id IS NULL AND short_links.expire_at IS NOT NULL AND short_links.expire_at < ? AND clicks.clicked_at < ?", time.Now(), weekAgo).
		Delete(&models.Click{})
	log.Info("Удалены старые клики по анонимным ссылкам")

	// 4. Для ссылок с владельцем, удаляем только если истёк срок и деактивированы
	var userLinks []models.ShortLink
	db.Where("user_id IS NOT NULL AND is_active = false AND expire_at IS NOT NULL AND expire_at < ?", time.Now()).Find(&userLinks)
	for _, link := range userLinks {
		db.Where("short_link_id = ?", link.ID).Delete(&models.Click{})
		db.Delete(&link)
		log.Info("Удалена пользовательская деактивированная истёкшая ссылка", zap.String("short_code", link.ShortCode))
	}
}
