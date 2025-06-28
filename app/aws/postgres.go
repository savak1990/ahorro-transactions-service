package aws

import (
	"context"
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
	gormDB        *gorm.DB
	gormDBMutex   sync.RWMutex
	lastAttempt   time.Time
	retryInterval = 30 * time.Second // Wait 30 seconds before retrying after failure
	cfg           config.AppConfig
	isConfigured  bool
)

// DatabaseConnectionError represents a database connection failure
type DatabaseConnectionError struct {
	Message string
	Cause   error
}

func (e *DatabaseConnectionError) Error() string {
	return fmt.Sprintf("database connection failed: %s", e.Message)
}

func (e *DatabaseConnectionError) Unwrap() error {
	return e.Cause
}

// SetConfig stores the database configuration without connecting
func SetConfig(appCfg config.AppConfig) {
	cfg = appCfg
	isConfigured = true
}

func GetGormDB(appCfg config.AppConfig) *gorm.DB {
	// Store config if not already set
	if !isConfigured {
		SetConfig(appCfg)
	}

	gormDBMutex.RLock()

	// Check if we have a healthy connection
	if gormDB != nil {
		sqlDB, err := gormDB.DB()
		if err == nil {
			// Quick ping to verify connection is still alive
			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
			defer cancel()
			if err := sqlDB.PingContext(ctx); err == nil {
				gormDBMutex.RUnlock()
				return gormDB
			}
			log.WithError(err).Warn("Existing database connection is unhealthy, will attempt reconnection")
		}
	}

	// Check if we should retry (rate limiting)
	now := time.Now()
	if !lastAttempt.IsZero() && now.Sub(lastAttempt) < retryInterval {
		gormDBMutex.RUnlock()
		log.WithFields(log.Fields{
			"last_attempt": lastAttempt,
			"retry_after":  retryInterval - now.Sub(lastAttempt),
		}).Warn("Database connection failed recently, not retrying yet")
		panic(&DatabaseConnectionError{
			Message: fmt.Sprintf("Database connection failed recently, retry after %v", retryInterval-now.Sub(lastAttempt)),
			Cause:   fmt.Errorf("rate limited"),
		})
	}

	gormDBMutex.RUnlock()
	gormDBMutex.Lock()
	defer gormDBMutex.Unlock()

	// Double-check pattern - another goroutine might have established connection
	if gormDB != nil {
		sqlDB, err := gormDB.DB()
		if err == nil {
			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
			defer cancel()
			if err := sqlDB.PingContext(ctx); err == nil {
				return gormDB
			}
		}
		// Close broken connection
		if sqlDB, err := gormDB.DB(); err == nil {
			sqlDB.Close()
		}
		gormDB = nil
	}

	lastAttempt = now
	log.Info("Attempting to establish new database connection...")

	// Enhanced DSN with short timeout for fast-fail behavior
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=require TimeZone=UTC connect_timeout=2 statement_timeout=5000",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	var err error
	gormDB, err = retryConnect(dsn, 1) // Only 1 attempt for fast-fail
	if err != nil {
		// Connection failed
		gormDB = nil
		dbErr := &DatabaseConnectionError{
			Message: fmt.Sprintf("Failed to connect to PostgreSQL database at %s:%d", cfg.DBHost, cfg.DBPort),
			Cause:   err,
		}
		log.WithError(err).Error("Failed to connect to PostgreSQL database with GORM - will return error to caller")
		panic(dbErr)
	}

	log.WithFields(log.Fields{
		"host":   cfg.DBHost,
		"port":   cfg.DBPort,
		"dbname": cfg.DBName,
	}).Info("Successfully connected to PostgreSQL database with GORM")

	// Auto-migrate the schema
	if err := autoMigrate(gormDB); err != nil {
		// Close the connection and reset
		if sqlDB, sqlErr := gormDB.DB(); sqlErr == nil {
			sqlDB.Close()
		}
		gormDB = nil

		dbErr := &DatabaseConnectionError{
			Message: "Failed to auto-migrate database schema",
			Cause:   err,
		}
		log.WithError(err).Error("Failed to auto-migrate database schema - will return error to caller")
		panic(dbErr)
	}

	return gormDB
}

// retryConnect attempts to establish database connection with fast-fail behavior
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
					// Configure connection pool for RDS PostgreSQL
					sqlDB.SetMaxIdleConns(2)                  // Reduced for cost optimization
					sqlDB.SetMaxOpenConns(10)                 // Reduced for cost optimization
					sqlDB.SetConnMaxLifetime(5 * time.Minute) // Shorter for cost optimization
					sqlDB.SetConnMaxIdleTime(1 * time.Minute) // Shorter for cost optimization
					return db, nil
				}
				lastErr = pingErr
			} else {
				lastErr = sqlErr
			}
		} else {
			lastErr = err
		}

		// If this isn't the last attempt, wait before retrying (reduced wait time)
		if i < maxRetries-1 {
			waitTime := 1 * time.Second // Fast retry for quick failure
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

// WithDatabaseTimeout wraps a context with a database-specific timeout
func WithDatabaseTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, 8*time.Second) // 8 seconds max for DB operations
}

// TestDatabaseConnection tests if the database is reachable without panicking
func TestDatabaseConnection(appCfg config.AppConfig) error {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=require TimeZone=UTC connect_timeout=1 statement_timeout=2000",
		appCfg.DBHost, appCfg.DBPort, appCfg.DBUser, appCfg.DBPassword, appCfg.DBName)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // Silent mode for testing
	})
	if err != nil {
		return &DatabaseConnectionError{
			Message: fmt.Sprintf("Failed to test connection to database at %s:%d", appCfg.DBHost, appCfg.DBPort),
			Cause:   err,
		}
	}

	// Test the connection with a simple ping
	sqlDB, err := db.DB()
	if err != nil {
		return &DatabaseConnectionError{
			Message: "Failed to get underlying database connection for testing",
			Cause:   err,
		}
	}

	if err := sqlDB.Ping(); err != nil {
		return &DatabaseConnectionError{
			Message: fmt.Sprintf("Database ping failed for %s:%d", appCfg.DBHost, appCfg.DBPort),
			Cause:   err,
		}
	}

	// Close the test connection
	sqlDB.Close()
	return nil
}

// IsHealthy checks if the current database connection is healthy
func IsHealthy() bool {
	gormDBMutex.RLock()
	defer gormDBMutex.RUnlock()

	if gormDB == nil {
		return false
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		return false
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	return sqlDB.PingContext(ctx) == nil
}

// ResetConnection forces a fresh connection attempt on next request
// This can be called when you know the database has been restarted
func ResetConnection() {
	gormDBMutex.Lock()
	defer gormDBMutex.Unlock()

	if gormDB != nil {
		if sqlDB, err := gormDB.DB(); err == nil {
			sqlDB.Close()
		}
		gormDB = nil
	}

	lastAttempt = time.Time{} // Reset retry timer
	log.Info("Database connection reset - next request will attempt fresh connection")
}

// GetConnectionStatus returns the current connection status
func GetConnectionStatus() (bool, time.Time, time.Duration) {
	gormDBMutex.RLock()
	defer gormDBMutex.RUnlock()

	hasConnection := gormDB != nil
	nextRetryIn := time.Duration(0)

	if !lastAttempt.IsZero() {
		elapsed := time.Since(lastAttempt)
		if elapsed < retryInterval {
			nextRetryIn = retryInterval - elapsed
		}
	}

	return hasConnection, lastAttempt, nextRetryIn
}
