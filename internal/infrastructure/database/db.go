package database

import (
	"fmt"

	"github.com/Dubjay18/ecom-api/internal/config"
	"github.com/Dubjay18/ecom-api/internal/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Database struct {
	*gorm.DB
}

// NewPostgresDB creates a new database connection
func NewPostgresDB(cfg *config.DBConfig) (*Database, error) {
	db, err := gorm.Open(postgres.Open(cfg.GetDSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get underlying SQL database to set connection pool settings
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(cfg.MaxLifetime)

	// Run migrations
	if err := runMigrations(db); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return &Database{db}, nil
}

func runMigrations(db *gorm.DB) error {
	return db.AutoMigrate(
		&domain.User{},
		&domain.Product{},
		&domain.Category{},
		&domain.Order{},
		&domain.OrderItem{},
		&domain.Address{},
	)
}
