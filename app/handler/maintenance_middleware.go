package handler

import (
	"net/http"
	"strings"

	"github.com/savak1990/transactions-service/app/aws"
	"github.com/savak1990/transactions-service/app/models"
	log "github.com/sirupsen/logrus"
)

// DatabaseMaintenanceMiddleware catches database maintenance errors and converts them to 503
func DatabaseMaintenanceMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Use a response wrapper to capture errors before they are written
		wrappedWriter := &responseWriter{ResponseWriter: w, statusCode: 200}

		defer func() {
			if err := recover(); err != nil {
				log.WithField("panic", err).Error("Panic caught by maintenance middleware")

				// Check if the panic is due to database maintenance/connection issues
				if errStr, ok := err.(string); ok {
					if isDatabaseMaintenanceErrorString(errStr) || isDatabaseConnectionErrorString(errStr) {
						WriteJSONError(w, http.StatusServiceUnavailable, models.ErrorCodeDbTimeout, "Database is temporarily unavailable, please retry in a few minutes")
						return
					}
				}

				// Check if panic is due to an error type
				if errObj, ok := err.(error); ok {
					errMessage := errObj.Error()
					if isDatabaseMaintenanceErrorString(errMessage) || isDatabaseConnectionErrorString(errMessage) {
						WriteJSONError(w, http.StatusServiceUnavailable, models.ErrorCodeDbTimeout, "Database is temporarily unavailable, please retry in a few minutes")
						return
					}
				}

				// Check for specific database connection error type
				if dbErr, ok := err.(*aws.DatabaseConnectionError); ok {
					log.WithField("database_error", dbErr.Error()).Error("Database connection error caught by middleware")
					WriteJSONError(w, http.StatusServiceUnavailable, models.ErrorCodeDbTimeout, "Database is temporarily unavailable, please retry in a few minutes")
					return
				}

				// Re-panic if it's not a database maintenance error
				panic(err)
			}
		}()

		next.ServeHTTP(wrappedWriter, r)
	})
}

// responseWriter wraps http.ResponseWriter to capture response codes
type responseWriter struct {
	http.ResponseWriter
	statusCode int
	written    bool
}

func (rw *responseWriter) WriteHeader(code int) {
	if rw.written {
		return
	}
	rw.statusCode = code
	rw.written = true
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if !rw.written {
		rw.WriteHeader(http.StatusOK)
	}
	return rw.ResponseWriter.Write(b)
}

// isDatabaseMaintenanceErrorString checks if a string indicates database maintenance
func isDatabaseMaintenanceErrorString(errStr string) bool {
	errStr = strings.ToLower(errStr)
	return strings.Contains(errStr, "maintenance") ||
		strings.Contains(errStr, "upgrade") ||
		strings.Contains(errStr, "unavailable") ||
		strings.Contains(errStr, "database system is starting up") ||
		strings.Contains(errStr, "connection to server") ||
		strings.Contains(errStr, "could not connect to server") ||
		strings.Contains(errStr, "aurora cluster") ||
		strings.Contains(errStr, "cluster endpoint") ||
		strings.Contains(errStr, "aurora serverless") ||
		strings.Contains(errStr, "serverless cluster") ||
		strings.Contains(errStr, "scaling") ||
		strings.Contains(errStr, "cold start") ||
		strings.Contains(errStr, "warming up") ||
		strings.Contains(errStr, "paused") ||
		strings.Contains(errStr, "resuming")
}

// isDatabaseConnectionErrorString checks if a string indicates database connection failure
func isDatabaseConnectionErrorString(errStr string) bool {
	errStr = strings.ToLower(errStr)
	return strings.Contains(errStr, "connection refused") ||
		strings.Contains(errStr, "dial tcp") ||
		strings.Contains(errStr, "failed to connect after") ||
		strings.Contains(errStr, "database connection failed") ||
		strings.Contains(errStr, "no such host") ||
		strings.Contains(errStr, "i/o timeout") ||
		strings.Contains(errStr, "context deadline exceeded") ||
		strings.Contains(errStr, "server closed the connection") ||
		strings.Contains(errStr, "timeout") ||
		strings.Contains(errStr, "network is unreachable") ||
		strings.Contains(errStr, "connection reset by peer")
}
