package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rahat-iqbal/ecommerce-api/internal/config"
	"github.com/rahat-iqbal/ecommerce-api/internal/handlers"
	"github.com/rahat-iqbal/ecommerce-api/internal/middleware"
	"gorm.io/gorm"
)

func Setup(cfg *config.Config, db *gorm.DB) *gin.Engine {
	r := gin.Default()

	r.GET("/health", handlers.HealthCheck)

	api := r.Group("/api/v1")
	api.Use(middleware.AuthRequired(cfg.JWTSecret))
	{
		// protected routes go here
	}

	return r
}
