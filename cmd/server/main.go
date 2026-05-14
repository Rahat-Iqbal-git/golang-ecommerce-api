package main

import (
	"log"

	// "github.com/joho/godotenv"
	"github.com/rahat-iqbal/ecommerce-api/internal/config"
	"github.com/rahat-iqbal/ecommerce-api/internal/routes"
)

func main() {
	// _ = godotenv.Load() // loads .env if present; ignored in production where env vars are set directly

	cfg := config.Load()

	router := routes.Setup(cfg)

	log.Printf("server starting on :%s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
