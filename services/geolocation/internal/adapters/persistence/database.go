package persistence


import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"theraclosure/geolocation-service/internal/adapters/config"
)

type Database struct {
	db *gorm.DB
}

func NewDatabase(cfg *config.Config) (*Database, error) {
	// Configure GORM logger
	var gormLogger logger.Interface
	switch cfg.App.LogLevel {
	case "debug":
		gormLogger = logger.Default.LogMode(logger.Info)
	case "info":
		gormLogger = logger.Default.LogMode(logger.Warn)
	default:
		gormLogger = logger.Default.LogMode(logger.Silent)
	}

	// Connect to database
	db, err := gorm.Open(postgres.Open(cfg.GetDatabaseDSN()), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return nil, err
	}

	// Skip auto-migration since we use manual SQL schema
	// Schema is managed via init-db.sql for better control
	log.Println("Skipping auto-migration - using manual schema from init-db.sql")

	// Database initialization successful
	log.Println("Database initialization completed successfully")

	return &Database{db: db}, nil
}

func createIndexes(db *gorm.DB) error {
	// Country indexes
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_countries_code ON countries(code)").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_countries_code2 ON countries(code2)").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_countries_name ON countries(name)").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_countries_active ON countries(active)").Error; err != nil {
		return err
	}

	// State indexes
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_states_country_id ON states(country_id)").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_states_name ON states(name)").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_states_code ON states(code)").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_states_active ON states(active)").Error; err != nil {
		return err
	}

	// City indexes
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_cities_state_id ON cities(state_id)").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_cities_name ON cities(name)").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_cities_zip_code ON cities(zip_code)").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_cities_active ON cities(active)").Error; err != nil {
		return err
	}

	return nil
}

func (d *Database) GetDB() *gorm.DB {
	return d.db
}