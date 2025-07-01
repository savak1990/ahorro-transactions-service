package handler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/savak1990/transactions-service/app/aws"
	"github.com/savak1990/transactions-service/app/models"
	"github.com/sirupsen/logrus"
)

// handleServiceError is a centralized error handler for service layer errors
// It checks for database connection errors and converts them to proper panics for middleware
// Returns true if the error was handled (response was written), false if caller should continue
func (h *HandlerImpl) handleServiceError(w http.ResponseWriter, err error, operation string) bool {
	if err == nil {
		return false
	}

	logrus.WithError(err).Errorf("%s failed", operation)

	// Check if this is a database connection error - if so, convert to proper panic
	if isDatabaseConnectionError(err) {
		// Convert the regular error back to a DatabaseConnectionError panic
		// so the maintenance middleware can catch it
		panic(&aws.DatabaseConnectionError{
			Message: fmt.Sprintf("Database connection failed during %s", operation),
			Cause:   err,
		})
	}

	// Handle other common error types
	if isDatabaseMaintenanceError(err) {
		WriteJSONError(w, http.StatusServiceUnavailable, models.ErrorCodeDbTimeout, "Database is undergoing maintenance, please retry in a few minutes")
		return true
	}

	if isDatabaseTimeoutError(err) {
		WriteJSONError(w, http.StatusServiceUnavailable, models.ErrorCodeDbTimeout, "Database is temporarily unavailable, please retry in a few moments")
		return true
	}

	if isDatabaseError(err) {
		WriteJSONError(w, http.StatusInternalServerError, models.ErrorCodeDbError, err.Error())
		return true
	}

	// Default to internal server error
	WriteJSONError(w, http.StatusInternalServerError, models.ErrorCodeInternalServer, err.Error())
	return true
}

// handleNotFoundError handles "not found" errors with a specific pattern
func (h *HandlerImpl) handleNotFoundError(w http.ResponseWriter, err error, resourceType, resourceID string) bool {
	expectedMessage := fmt.Sprintf("%s not found: %s", resourceType, resourceID)
	if err.Error() == expectedMessage {
		WriteJSONError(w, http.StatusNotFound, models.ErrorCodeNotFound, fmt.Sprintf("%s not found", strings.Title(resourceType)))
		return true
	}
	return false
}

// isDatabaseConnectionError checks if an error is a database connection error
func isDatabaseConnectionError(err error) bool {
	if err == nil {
		return false
	}

	errStr := strings.ToLower(err.Error())
	return strings.Contains(errStr, "connection refused") ||
		strings.Contains(errStr, "dial tcp") ||
		strings.Contains(errStr, "failed to connect") ||
		strings.Contains(errStr, "database connection failed") ||
		strings.Contains(errStr, "no such host") ||
		strings.Contains(errStr, "i/o timeout") ||
		strings.Contains(errStr, "context deadline exceeded") ||
		strings.Contains(errStr, "server closed the connection") ||
		strings.Contains(errStr, "timeout") ||
		strings.Contains(errStr, "network is unreachable") ||
		strings.Contains(errStr, "connection reset by peer")
}
