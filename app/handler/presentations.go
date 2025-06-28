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
		strings.Contains(errStr, "context deadline exceeded")
}
