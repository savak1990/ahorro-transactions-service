package handler

import (
	"encoding/json"
	"net/http"
	"strings"
)

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type PaginatedResponse[T any] struct {
	Items   []T    `json:"items"`
	NextKey string `json:"nextKey,omitempty"`
}

func WriteJSONError(w http.ResponseWriter, status int, code string, errMsg string) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	response := map[string]string{"code": code, "error": errMsg}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode error response", http.StatusInternalServerError)
		w.Write([]byte(`{"code":"InternalServerError","error":"Failed to encode error response"}`))
	}
}

func WriteJSONListResponse[T any](w http.ResponseWriter, items []T, nextKey string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := PaginatedResponse[T]{
		Items:   items,
		NextKey: nextKey,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// isDatabaseError checks if the error is related to database operations
func isDatabaseError(err error) bool {
	if err == nil {
		return false
	}
	errStr := strings.ToLower(err.Error())
	return strings.Contains(errStr, "sqlstate") ||
		strings.Contains(errStr, "constraint") ||
		strings.Contains(errStr, "foreign key") ||
		strings.Contains(errStr, "database") ||
		strings.Contains(errStr, "gorm") ||
		strings.Contains(errStr, "postgres")
}

// isDatabaseTimeoutError checks if the error is specifically a connection timeout
func isDatabaseTimeoutError(err error) bool {
	if err == nil {
		return false
	}
	errStr := strings.ToLower(err.Error())
	return strings.Contains(errStr, "timeout") ||
		strings.Contains(errStr, "connection refused") ||
		strings.Contains(errStr, "no such host") ||
		strings.Contains(errStr, "i/o timeout") ||
		strings.Contains(errStr, "context deadline exceeded") ||
		strings.Contains(errStr, "dial tcp") ||
		strings.Contains(errStr, "connect: connection refused") ||
		strings.Contains(errStr, "failed to connect after") ||
		strings.Contains(errStr, "database connection failed")
}

// isDatabaseMaintenanceError checks if the error is due to database maintenance/upgrade
func isDatabaseMaintenanceError(err error) bool {
	if err == nil {
		return false
	}
	errStr := strings.ToLower(err.Error())
	return strings.Contains(errStr, "maintenance") ||
		strings.Contains(errStr, "upgrade") ||
		strings.Contains(errStr, "unavailable") ||
		strings.Contains(errStr, "temporarily unavailable") ||
		strings.Contains(errStr, "server is not currently available") ||
		strings.Contains(errStr, "cluster is not available") ||
		strings.Contains(errStr, "database system is starting up") ||
		strings.Contains(errStr, "database system is shutting down") ||
		strings.Contains(errStr, "connection to server") ||
		strings.Contains(errStr, "could not connect to server") ||
		strings.Contains(errStr, "server closed the connection unexpectedly") ||
		strings.Contains(errStr, "aurora cluster") ||
		strings.Contains(errStr, "cluster endpoint") ||
		strings.Contains(errStr, "aurora serverless") ||
		strings.Contains(errStr, "serverless cluster") ||
		strings.Contains(errStr, "scaling") ||
		strings.Contains(errStr, "cold start") ||
		strings.Contains(errStr, "warming up") ||
		strings.Contains(errStr, "paused") ||
		strings.Contains(errStr, "resuming") ||
		strings.Contains(errStr, "dial tcp") ||
		strings.Contains(errStr, "connection refused") ||
		strings.Contains(errStr, "failed to connect after") ||
		strings.Contains(errStr, "database connection failed")
}
