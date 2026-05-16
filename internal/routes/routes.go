package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rahat-iqbal/ecommerce-api/internal/config"
	"github.com/rahat-iqbal/ecommerce-api/internal/handlers"
	"github.com/rahat-iqbal/ecommerce-api/internal/middleware"
	"github.com/rahat-iqbal/ecommerce-api/internal/repository"
	"github.com/rahat-iqbal/ecommerce-api/internal/service"
	"gorm.io/gorm"
)

func Setup(cfg *config.Config, db *gorm.DB) *gin.Engine {
	r := gin.Default()

	r.GET("/health", handlers.HealthCheck)

	auth := handlers.NewAuthHandler(db, cfg.JWTSecret)
	r.POST("/api/v1/register", auth.Register)
	r.POST("/api/v1/login", auth.Login)

	api := r.Group("/api/v1")
	api.Use(middleware.AuthRequired(cfg.JWTSecret))
	{
		products := handlers.NewProductHandler(db)
		api.GET("/products", products.List)
		api.GET("/products/:id", products.Get)
		api.POST("/products", products.Create)
		api.PUT("/products/:id", products.Update)
		api.DELETE("/products/:id", products.Delete)

		cart := handlers.NewCartHandler(
			service.NewCartService(
				repository.NewCartRepository(db),
				repository.NewProductRepository(db),
			),
		)
		api.GET("/cart", cart.ListItems)
		api.POST("/cart", cart.AddItem)
		api.PUT("/cart/:id", cart.UpdateItem)
		api.DELETE("/cart/:id", cart.RemoveItem)
	}

	return r
}
