package handler

import (
	"encoding/json"
	"net/http"
)

type CommonHandlerImpl struct{}

func NewCommonHandlerImpl() *CommonHandlerImpl {
	return &CommonHandlerImpl{}
}

func (h *CommonHandlerImpl) HandleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (h *CommonHandlerImpl) HandleInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	info := map[string]string{
		"version": "1.0.0",
		"status":  "running",
	}
	if err := json.NewEncoder(w).Encode(info); err != nil {
		http.Error(w, "Failed to encode info response", http.StatusInternalServerError)
	}
}
