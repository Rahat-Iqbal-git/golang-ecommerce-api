package main

import (
	"log"

	"github.com/rahat-iqbal/ecommerce-api/internal/config"
	"github.com/rahat-iqbal/ecommerce-api/internal/routes"
)

func main() {
	cfg := config.Load()

	router := routes.Setup(cfg)

	log.Printf("server starting on :%s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
