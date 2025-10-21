package main


import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"theraclosure/geolocation-service/internal/adapters/config"
	"theraclosure/geolocation-service/internal/adapters/http"
	"theraclosure/geolocation-service/internal/adapters/persistence"
	"theraclosure/geolocation-service/internal/core/services"
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
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Initialize repository
	repo := persistence.NewGeolocationRepository(db)

	// Run world data seeding if enabled
	if shouldSeedData() {
		log.Println("Seeding world data...")
		if err := seedWorldData(cfg); err != nil {
			log.Printf("Warning: Failed to seed world data: %v", err)
		} else {
			log.Println("World data seeding completed successfully")
		}
	}

	// Initialize service
	service := services.NewGeolocationService(repo)

	// Initialize HTTP server
	server := http.NewServer(service, cfg)

	// Start server
	log.Printf("Geolocation service starting on %s", cfg.GetServerAddress())
	if err := server.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// shouldSeedData checks if world data seeding is enabled
func shouldSeedData() bool {
	seedData := os.Getenv("SEED_DATA")
	return strings.ToLower(seedData) == "true"
}

// seedWorldData runs the Python world data seeder
func seedWorldData(cfg *config.Config) error {
	// Get the script path relative to the service root
	scriptPath := filepath.Join("data", "seeds", "world_data_seeder.py")
	
	// Check if script exists
	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		return err
	}

	// Prepare environment variables for the Python script
	env := os.Environ()
	env = append(env, "DB_HOST="+cfg.Database.Host)
	env = append(env, "DB_PORT="+cfg.Database.Port)
	env = append(env, "DB_USER="+cfg.Database.User)
	env = append(env, "DB_PASSWORD="+cfg.Database.Password)
	env = append(env, "DB_NAME="+cfg.Database.DBName)

	// Run the Python seeder script
	cmd := exec.Command("python3", scriptPath)
	cmd.Env = env
	cmd.Dir = "."

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Seeder output: %s", output)
		return err
	}

	log.Printf("Seeder completed: %s", output)
	return nil
}