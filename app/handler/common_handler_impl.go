package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type CommonHandlerImpl struct {
	db *gorm.DB
}

func NewCommonHandlerImpl(db *gorm.DB) *CommonHandlerImpl {
	return &CommonHandlerImpl{
		db: db,
	}
}

func (h *CommonHandlerImpl) HandleHealth(w http.ResponseWriter, r *http.Request) {
	// Test database connection with timeout context
	dbStatus := "healthy"

	if h.db != nil {
		sqlDB, err := h.db.DB()
		if err != nil {
			log.WithError(err).Error("Failed to get SQL DB from GORM")
			dbStatus = "unhealthy"
		} else {
			// Use a shorter timeout for health checks to detect cold starts
			ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
			defer cancel()

			if err := sqlDB.PingContext(ctx); err != nil {
				log.WithError(err).Warn("Database health check failed (possibly cold start)")
				if isDatabaseTimeoutError(err) {
					dbStatus = "warming_up"
				} else {
					dbStatus = "unhealthy"
				}
			}
		}

		if dbStatus == "unhealthy" {
			w.WriteHeader(http.StatusServiceUnavailable)
			response := map[string]interface{}{
				"status":   "unhealthy",
				"database": "disconnected",
				"error":    "Database connection failed",
			}
			json.NewEncoder(w).Encode(response)
			return
		} else if dbStatus == "warming_up" {
			w.WriteHeader(http.StatusServiceUnavailable)
			response := map[string]interface{}{
				"status":   "warming_up",
				"database": "cold_start",
				"message":  "Database is starting up, please retry in 10-30 seconds",
			}
			json.NewEncoder(w).Encode(response)
			return
		}
	} else {
		dbStatus = "not_configured"
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
