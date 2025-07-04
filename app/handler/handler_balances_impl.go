package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/savak1990/transactions-service/app/models"
)

// Balance handlers
func (h *HandlerImpl) CreateBalance(w http.ResponseWriter, r *http.Request) {
	var balanceDto models.BalanceDto
	if err := json.NewDecoder(r.Body).Decode(&balanceDto); err != nil {
		WriteJSONError(w, http.StatusBadRequest, models.ErrorCodeBadRequest, "Invalid request body: "+err.Error())
		return
	}

	// Convert DTO to DAO model
	balance, err := models.FromAPIBalance(balanceDto)
	if err != nil {
		WriteJSONError(w, http.StatusBadRequest, models.ErrorCodeBadRequest, "Invalid balance data: "+err.Error())
		return
	}

	created, err := h.Service.CreateBalance(r.Context(), *balance)
	if err != nil {
		h.handleServiceError(w, err, "CreateBalance")
		return
	}

	// Convert back to DTO for response
	responseDto := models.ToAPIBalance(created)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responseDto)
}

func (h *HandlerImpl) ListBalances(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	// Parse limit
	limit := 50 // default
	if limitStr := query.Get("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 && parsedLimit <= 100 {
			limit = parsedLimit
		}
	}

	filter := models.ListBalancesInput{
		UserID:    query.Get("userId"),
		GroupID:   query.Get("groupId"),
		BalanceID: query.Get("balanceId"),
		SortBy:    query.Get("sortBy"),
		Order:     query.Get("order"),
		Limit:     limit,
	}

	results, err := h.Service.ListBalances(r.Context(), filter)
	if err != nil {
		h.handleServiceError(w, err, "ListBalances")
		return
	}

	// Convert to DTOs for response
	balanceDtos := make([]models.BalanceDto, len(results))
	for i, balance := range results {
		balanceDtos[i] = models.ToAPIBalance(&balance)
	}

	// Use paginated response structure (without actual pagination for now)
	WriteJSONListResponse(w, balanceDtos, "")
}

func (h *HandlerImpl) GetBalance(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	balanceID := vars["balance_id"]
	if balanceID == "" {
		WriteJSONError(w, http.StatusBadRequest, models.ErrorCodeBadRequest, "Missing balance_id")
		return
	}

	balance, err := h.Service.GetBalance(r.Context(), balanceID)
	if err != nil {
		// Try to handle as "not found" error first
		if h.handleNotFoundError(w, err, "balance", balanceID) {
			return
		}
		// Handle all other errors (including database connection errors)
		h.handleServiceError(w, err, "GetBalance")
		return
	}

	if balance == nil {
		WriteJSONError(w, http.StatusNotFound, models.ErrorCodeNotFound, "Balance not found")
		return
	}

	// Convert to DTO for response
	responseDto := models.ToAPIBalance(balance)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responseDto)
}

func (h *HandlerImpl) UpdateBalance(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	balanceID := vars["balance_id"]
	var balanceDto models.BalanceDto
	if err := json.NewDecoder(r.Body).Decode(&balanceDto); err != nil {
		WriteJSONError(w, http.StatusBadRequest, models.ErrorCodeBadRequest, "Invalid request body: "+err.Error())
		return
	}

	// Parse balance ID and set it in DTO
	id, err := uuid.Parse(balanceID)
	if err != nil {
		WriteJSONError(w, http.StatusBadRequest, models.ErrorCodeBadRequest, "Invalid balance ID format")
		return
	}
	balanceDto.BalanceID = id.String()

	// Convert DTO to DAO model
	balance, err := models.FromAPIBalance(balanceDto)
	if err != nil {
		WriteJSONError(w, http.StatusBadRequest, models.ErrorCodeBadRequest, "Invalid balance data: "+err.Error())
		return
	}

	updated, err := h.Service.UpdateBalance(r.Context(), *balance)
	if err != nil {
		h.handleServiceError(w, err, "UpdateBalance")
		return
	}

	// Convert back to DTO for response
	responseDto := models.ToAPIBalance(updated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responseDto)
}

func (h *HandlerImpl) DeleteBalance(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	balanceID := vars["balance_id"]

	if balanceID == "" {
		WriteJSONError(w, http.StatusBadRequest, models.ErrorCodeBadRequest, "Missing balance_id parameter")
		return
	}

	err := h.Service.DeleteBalance(r.Context(), balanceID)
	if err != nil {
		h.handleServiceError(w, err, "DeleteBalance")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *HandlerImpl) DeleteBalancesByUserId(w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query().Get("userId")

	if userId == "" {
		WriteJSONError(w, http.StatusBadRequest, models.ErrorCodeBadRequest, "Missing userId parameter")
		return
	}

	err := h.Service.DeleteBalancesByUserId(r.Context(), userId)
	if err != nil {
		h.handleServiceError(w, err, "DeleteBalancesByUserId")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
