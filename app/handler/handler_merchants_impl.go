package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/savak1990/transactions-service/app/helpers"
	"github.com/savak1990/transactions-service/app/models"
)

// Merchant handlers
func (h *HandlerImpl) CreateMerchant(w http.ResponseWriter, r *http.Request) {
	var merchantDto models.MerchantDto
	if err := json.NewDecoder(r.Body).Decode(&merchantDto); err != nil {
		WriteJSONError(w, http.StatusBadRequest, models.ErrorCodeBadRequest, "Invalid request body: "+err.Error())
		return
	}

	// Convert DTO to DAO model
	merchant, err := models.FromAPIMerchant(merchantDto)
	if err != nil {
		WriteJSONError(w, http.StatusBadRequest, models.ErrorCodeBadRequest, "Invalid merchant data: "+err.Error())
		return
	}

	created, err := h.Service.CreateMerchant(r.Context(), *merchant)
	if err != nil {
		// Check if it's a duplicate merchant error
		errorMsg := err.Error()
		if strings.Contains(errorMsg, "already exists for this user") {
			WriteJSONError(w, http.StatusConflict, models.ErrorCodeConflict, err.Error())
			return
		}
		h.handleServiceError(w, err, "CreateMerchant")
		return
	}

	// Convert back to DTO for response
	responseDto := models.ToAPIMerchant(created)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(responseDto)
}

func (h *HandlerImpl) ListMerchants(w http.ResponseWriter, r *http.Request) {
	filter := models.ListMerchantsInput{
		Name:   r.URL.Query().Get("name"),
		SortBy: r.URL.Query().Get("sortBy"),
		Order:  r.URL.Query().Get("order"),
	}
	if count := r.URL.Query().Get("limit"); count != "" {
		if n, err := helpers.ParseInt(count); err == nil {
			filter.Limit = n
		}
	}

	results, err := h.Service.ListMerchants(r.Context(), filter)
	if err != nil {
		h.handleServiceError(w, err, "ListMerchants")
		return
	}

	// Convert to DTOs for response
	merchantDtos := make([]models.MerchantDto, len(results))
	for i, merchant := range results {
		merchantDtos[i] = models.ToAPIMerchant(&merchant)
	}

	// Use paginated response structure (without actual pagination for now)
	WriteJSONListResponse(w, merchantDtos, "")
}

func (h *HandlerImpl) GetMerchant(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	merchantID := vars["merchant_id"]
	if merchantID == "" {
		WriteJSONError(w, http.StatusBadRequest, models.ErrorCodeBadRequest, "Missing merchant_id")
		return
	}

	merchant, err := h.Service.GetMerchant(r.Context(), merchantID)
	if err != nil {
		// Try to handle as "not found" error first
		if h.handleNotFoundError(w, err, "merchant", merchantID) {
			return
		}
		// Handle all other errors (including database connection errors)
		h.handleServiceError(w, err, "GetMerchant")
		return
	}

	if merchant == nil {
		WriteJSONError(w, http.StatusNotFound, models.ErrorCodeNotFound, "Merchant not found")
		return
	}

	// Convert to DTO for response
	responseDto := models.ToAPIMerchant(merchant)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responseDto)
}

func (h *HandlerImpl) UpdateMerchant(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	merchantID := vars["merchant_id"]
	var merchantDto models.MerchantDto
	if err := json.NewDecoder(r.Body).Decode(&merchantDto); err != nil {
		WriteJSONError(w, http.StatusBadRequest, models.ErrorCodeBadRequest, "Invalid request body: "+err.Error())
		return
	}

	// Parse merchant ID and set it in DTO
	id, err := uuid.Parse(merchantID)
	if err != nil {
		WriteJSONError(w, http.StatusBadRequest, models.ErrorCodeBadRequest, "Invalid merchant ID format")
		return
	}
	merchantDto.MerchantID = id.String()

	// Convert DTO to DAO model
	merchant, err := models.FromAPIMerchant(merchantDto)
	if err != nil {
		WriteJSONError(w, http.StatusBadRequest, models.ErrorCodeBadRequest, "Invalid merchant data: "+err.Error())
		return
	}

	updated, err := h.Service.UpdateMerchant(r.Context(), *merchant)
	if err != nil {
		h.handleServiceError(w, err, "UpdateMerchant")
		return
	}

	// Convert back to DTO for response
	responseDto := models.ToAPIMerchant(updated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responseDto)
}

func (h *HandlerImpl) DeleteMerchant(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	merchantID := vars["merchant_id"]
	if merchantID == "" {
		WriteJSONError(w, http.StatusBadRequest, models.ErrorCodeBadRequest, "Missing merchant_id")
		return
	}
	err := h.Service.DeleteMerchant(r.Context(), merchantID)
	if err != nil {
		h.handleServiceError(w, err, "DeleteMerchant")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *HandlerImpl) DeleteMerchantsByUserId(w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query().Get("userId")

	if userId == "" {
		WriteJSONError(w, http.StatusBadRequest, models.ErrorCodeBadRequest, "Missing userId parameter")
		return
	}

	err := h.Service.DeleteMerchantsByUserId(r.Context(), userId)
	if err != nil {
		h.handleServiceError(w, err, "DeleteMerchantsByUserId")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
