package persistence

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"theraclosure/payments-service/internal/adapters/config"
	"theraclosure/payments-service/internal/core/domain"
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

	// Auto-migrate schemas
	err = db.AutoMigrate(
		&domain.Customer{},
		&domain.Subscription{},
		&domain.Payment{},
	)
	if err != nil {
		return nil, err
	}

	log.Println("Database connected and migrated successfully")

	return &Database{db: db}, nil
}

func (d *Database) GetDB() *gorm.DB {
	return d.db
}
