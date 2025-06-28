package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/savak1990/transactions-service/app/aws"
	"github.com/savak1990/transactions-service/app/config"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type CommonHandlerImpl struct {
	db     *gorm.DB
	config *config.AppConfig
}

func NewCommonHandlerImpl(db *gorm.DB) *CommonHandlerImpl {
	return &CommonHandlerImpl{
		db: db,
	}
}

// NewCommonHandlerImplWithConfig creates a new CommonHandlerImpl with lazy DB initialization
func NewCommonHandlerImplWithConfig(cfg config.AppConfig) *CommonHandlerImpl {
	return &CommonHandlerImpl{
		config: &cfg,
	}
}

// getDB returns the database connection, initializing it if necessary
func (h *CommonHandlerImpl) getDB() *gorm.DB {
	if h.db != nil {
		return h.db
	}

	if h.config != nil {
		// Lazy initialization - this will trigger connection and panic if failed
		h.db = aws.GetGormDB(*h.config)
		return h.db
	}

	panic(fmt.Errorf("CommonHandlerImpl not properly initialized"))
}

func (h *CommonHandlerImpl) HandleHealth(w http.ResponseWriter, r *http.Request) {
	// Test database connection with timeout context
	dbStatus := "healthy"

	// Use lazy initialization - try to get DB but catch panics
	var db *gorm.DB
	func() {
		defer func() {
			if r := recover(); r != nil {
				log.WithField("error", r).Warn("Database connection failed during health check")
				if dbErr, ok := r.(*aws.DatabaseConnectionError); ok {
					if isDatabaseMaintenanceError(dbErr) {
						dbStatus = "maintenance"
					} else if isDatabaseTimeoutError(dbErr) {
						dbStatus = "timeout"
					} else {
						dbStatus = "unhealthy"
					}
				} else {
					dbStatus = "unhealthy"
				}
			}
		}()
		db = h.getDB()
	}()

	if db != nil {
		// First check if the GORM connection is healthy using our helper
		if !aws.IsHealthy() {
			log.Warn("Database connection appears unhealthy via GORM health check")
			dbStatus = "unhealthy"
		} else {
			// Additional ping test with timeout for more detailed status
			sqlDB, err := db.DB()
			if err != nil {
				log.WithError(err).Error("Failed to get SQL DB from GORM")
				dbStatus = "unhealthy"
			} else {
				// Use a shorter timeout for health checks to detect cold starts quickly
				ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
				defer cancel()

				if err := sqlDB.PingContext(ctx); err != nil {
					log.WithError(err).Warn("Database health check failed")
					if isDatabaseMaintenanceError(err) {
						dbStatus = "maintenance"
					} else if isDatabaseTimeoutError(err) {
						dbStatus = "warming_up"
					} else {
						dbStatus = "unhealthy"
					}
				}
			}
		}

		switch dbStatus {
		case "maintenance":
			w.WriteHeader(http.StatusServiceUnavailable)
			response := map[string]interface{}{
				"status":   "maintenance",
				"database": "upgrading",
				"message":  "Database is undergoing maintenance, please retry in a few minutes",
			}
			json.NewEncoder(w).Encode(response)
			return
		case "unhealthy":
			w.WriteHeader(http.StatusServiceUnavailable)
			response := map[string]interface{}{
				"status":   "unhealthy",
				"database": "disconnected",
				"error":    "Database connection failed",
			}
			json.NewEncoder(w).Encode(response)
			return
		case "warming_up":
			w.WriteHeader(http.StatusServiceUnavailable)
			response := map[string]interface{}{
				"status":   "warming_up",
				"database": "unavailable",
				"message":  "Database is temporarily unavailable, please retry in a few moments",
			}
			json.NewEncoder(w).Encode(response)
			return
		}
	} else {
		// Try to get connection status without triggering a connection attempt
		hasConnection, lastAttempt, nextRetryIn := aws.GetConnectionStatus()

		if !hasConnection && !lastAttempt.IsZero() {
			if nextRetryIn > 0 {
				dbStatus = "cooling_down"
			} else {
				dbStatus = "ready_to_retry"
			}
		} else if !hasConnection {
			dbStatus = "not_configured"
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{
		"status":   "healthy",
		"database": dbStatus,
	}
	json.NewEncoder(w).Encode(response)
}

func (h *CommonHandlerImpl) HandleInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	info := map[string]string{
		"version":  "1.0.0",
		"status":   "running",
		"database": "postgresql",
	}
	if err := json.NewEncoder(w).Encode(info); err != nil {
		http.Error(w, "Failed to encode info response", http.StatusInternalServerError)
	}
}

func (h *CommonHandlerImpl) HandleDbReset(w http.ResponseWriter, r *http.Request) {
	aws.ResetConnection()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{
		"message": "Database connection reset - next request will attempt fresh connection",
		"status":  "reset",
	}
	json.NewEncoder(w).Encode(response)
}
