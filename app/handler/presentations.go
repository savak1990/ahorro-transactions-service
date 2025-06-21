package handler

import (
	"encoding/json"
	"net/http"
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
