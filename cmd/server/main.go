package main

import (
	"log"

	"github.com/rahat-iqbal/ecommerce-api/internal/config"
	"github.com/rahat-iqbal/ecommerce-api/internal/database"
	"github.com/rahat-iqbal/ecommerce-api/internal/models"
	"github.com/rahat-iqbal/ecommerce-api/internal/routes"
)

func main() {
	// _ = godotenv.Load() // loads .env if present; ignored in production

	cfg := config.Load()

	db, err := database.Connect(cfg.DBUrl)
	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}
	log.Println("database connected")

	if err := db.AutoMigrate(&models.User{}, &models.Product{}); err != nil {
		log.Fatal("migration failed:", err)
	}

	router := routes.Setup(cfg, db)

	log.Printf("server starting on :%s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
