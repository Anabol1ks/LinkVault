package router

import (
	"linkvault/internal/config"
	"linkvault/internal/handler"
	"linkvault/internal/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Handlers struct {
	User *handler.UserHandler
	Link *handler.ShortLinkHandler
}

func Router(db *gorm.DB, log *zap.Logger, handlers *Handlers, cfg *config.Config) *gin.Engine {
	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	auth := r.Group("/auth")
	{
		auth.POST("/register", handlers.User.Register)
		auth.POST("/login", handlers.User.Login)
		auth.POST("/refresh", handlers.User.Refresh)
	}

	links := r.Group("/links")
	links.POST("", middleware.OptionalJWTAuth(&cfg.JWT), handlers.Link.CreateShortLink)

	r.GET("/:shortCode", handlers.Link.GetOriginalURL)

	return r
}
