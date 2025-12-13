package main

import (
	"fmt"
	"log"
	"os"

	"github.com/DXR3IN/device-service-v2/internal/config"
	"github.com/DXR3IN/device-service-v2/internal/http"
	repo "github.com/DXR3IN/device-service-v2/internal/repository"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	if _, exists := os.LookupEnv("RUNNING_IN_DOCKER"); !exists {
		godotenv.Load()
	}

	cfg := config.NewConfigFromEnv()

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		cfg.DBHost, cfg.DBUser, cfg.DBPass, cfg.DBName, cfg.DBPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	if err := db.AutoMigrate(&repo.Device{}); err != nil {
		log.Fatalf("failed to migrate: %v", err)
	}

	deviceRepo := repo.NewDeviceRepository(db)

	r := http.NewRouter(cfg, deviceRepo)

	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("Starting server at %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
