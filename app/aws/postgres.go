package aws

import (
	"fmt"
	"sync"

	"github.com/savak1990/transactions-service/app/config"
	"github.com/savak1990/transactions-service/app/models"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	gormDB     *gorm.DB
	gormDBOnce sync.Once
)

func GetGormDB(cfg config.AppConfig) *gorm.DB {
	gormDBOnce.Do(func() {
		dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=require TimeZone=UTC",
			cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

		var err error
		gormDB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
		if err != nil {
			log.WithError(err).Fatal("Failed to connect to PostgreSQL database with GORM")
		}

		// Get underlying SQL DB for connection testing
		sqlDB, err := gormDB.DB()
		if err != nil {
			log.WithError(err).Fatal("Failed to get underlying SQL DB")
		}

		// Test the connection
		if err = sqlDB.Ping(); err != nil {
			log.WithError(err).Fatal("Failed to ping PostgreSQL database")
		}

		// Configure connection pool
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetMaxOpenConns(100)

		log.WithFields(log.Fields{
			"host":   cfg.DBHost,
			"port":   cfg.DBPort,
			"dbname": cfg.DBName,
		}).Info("Successfully connected to PostgreSQL database with GORM")

		// Auto-migrate the schema
		if err := autoMigrate(gormDB); err != nil {
			log.WithError(err).Fatal("Failed to auto-migrate database schema")
		}
	})

	if gormDB == nil {
		log.Fatal("GORM DB is not initialized")
	}
	return gormDB
}

// autoMigrate runs GORM auto-migration for all models
func autoMigrate(db *gorm.DB) error {
	log.Info("Running GORM auto-migration...")

	// Migrate all models (using DAO models for PostgreSQL)
	err := db.AutoMigrate(
		&models.Balance{},
		&models.Merchant{},
		&models.Category{},
		&models.Transaction{},
		&models.TransactionEntry{},
	)

	if err != nil {
		return fmt.Errorf("auto-migration failed: %w", err)
	}

	log.Info("GORM auto-migration completed successfully")
	return nil
}

// GetPostgreSQLClient returns the GORM DB for consistency with existing code
func GetPostgreSQLClient(cfg config.AppConfig) *gorm.DB {
	return GetGormDB(cfg)
}
