package main

import (
	"log"
	"theraclosure/users-service/internal/adapters/config"
	"theraclosure/users-service/internal/adapters/http"
	"theraclosure/users-service/internal/adapters/persistence"
	"theraclosure/users-service/internal/core/services"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database
	db, err := persistence.NewDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize repositories
	userRepo := persistence.NewUserProfileRepository(db.DB)
	enrollmentRepo := persistence.NewEnrollmentRepository(db.DB)

	// Initialize services
	userService := services.NewUserService(userRepo, enrollmentRepo)
	enrollmentService := services.NewEnrollmentService(enrollmentRepo)

	// Initialize HTTP server
	server := http.NewServer(cfg, userService, enrollmentService)

	// Start server
	if err := server.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
