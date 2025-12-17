package main

import (
	"fmt"
	"log"
	"os"

	"github.com/DXR3IN/auth-service-v2/internal/config"
	"github.com/DXR3IN/auth-service-v2/internal/http"
	"github.com/DXR3IN/auth-service-v2/internal/migration"
	repo "github.com/DXR3IN/auth-service-v2/internal/repository"
	"github.com/joho/godotenv"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	// Load environment variables
	if _, exists := os.LookupEnv("RUNNING_IN_DOCKER"); !exists {
		if err := godotenv.Load(); err != nil {
			log.Printf("Warning: .env file not found")
		}
	}

	cfg := config.NewConfigFromEnv()

	// Build database connection string
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		cfg.DBHost, cfg.DBUser, cfg.DBPass, cfg.DBName, cfg.DBPort)

	// Connect to database with logging
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Database connection established")

	// Run migrations
	log.Println("Running database migrations...")
	if err := migration.RunMigrations(db); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
	log.Println("All migrations completed successfully")

	// Optional: Show migration status
	if err := migration.ShowMigrationStatus(db); err != nil {
		log.Printf("Warning: Could not show migration status: %v", err)
	}

	// Initialize repository
	userRepo := repo.NewUserRepository(db)

	// Setup HTTP router
	r := http.NewRouter(cfg, userRepo)

	// Start server
	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("Starting server at %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
