package aws

import (
	"fmt"
	"sync"
	"time"

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
		// Enhanced DSN with connection timeouts for Aurora Serverless v2
		dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=require TimeZone=UTC connect_timeout=30 statement_timeout=30000",
			cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

		var err error
		gormDB, err = retryConnect(dsn, 5) // Retry up to 5 times for cold starts
		if err != nil {
			log.WithError(err).Fatal("Failed to connect to PostgreSQL database with GORM")
		}

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

// retryConnect attempts to establish database connection with retries for cold starts
func retryConnect(dsn string, maxRetries int) (*gorm.DB, error) {
	var lastErr error

	for i := 0; i < maxRetries; i++ {
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})

		if err == nil {
			// Test the connection
			sqlDB, sqlErr := db.DB()
			if sqlErr == nil {
				pingErr := sqlDB.Ping()
				if pingErr == nil {
					// Configure connection pool for Aurora Serverless v2
					sqlDB.SetMaxIdleConns(2)                  // Reduced for serverless
					sqlDB.SetMaxOpenConns(10)                 // Reduced for serverless
					sqlDB.SetConnMaxLifetime(5 * time.Minute) // Shorter for serverless
					sqlDB.SetConnMaxIdleTime(1 * time.Minute) // Shorter for serverless
					return db, nil
				}
				lastErr = pingErr
			} else {
				lastErr = sqlErr
			}
		} else {
			lastErr = err
		}

		// If this isn't the last attempt, wait before retrying
		if i < maxRetries-1 {
			waitTime := time.Duration(i+1) * 10 * time.Second // 10s, 20s, 30s...
			log.WithFields(log.Fields{
				"attempt":      i + 1,
				"max_attempts": maxRetries,
				"wait_time":    waitTime,
				"error":        lastErr,
			}).Warn("Database connection failed, retrying...")
			time.Sleep(waitTime)
		}
	}

	return nil, fmt.Errorf("failed to connect after %d attempts, last error: %w", maxRetries, lastErr)
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
