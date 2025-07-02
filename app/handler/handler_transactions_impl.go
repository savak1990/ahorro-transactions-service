package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/savak1990/transactions-service/app/helpers"
	"github.com/savak1990/transactions-service/app/models"
	"github.com/sirupsen/logrus"
)

// POST /transactions
func (h *HandlerImpl) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var req models.CreateTransactionDto
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteJSONError(w, http.StatusBadRequest, models.ErrorCodeBadRequest, "Invalid request body: "+err.Error())
		return
	}

	// Convert request DTO to transaction
	transaction, err := models.FromAPICreateTransaction(req)
	if err != nil {
		logrus.WithError(err).Error("Failed to convert request DTO")
		WriteJSONError(w, http.StatusBadRequest, models.ErrorCodeBadRequest, "Invalid request: "+err.Error())
		return
	}

	created, err := h.Service.CreateTransaction(r.Context(), *transaction)
	if err != nil {
		h.handleServiceError(w, err, "CreateTransaction")
		return
	}

	// Convert response to DTO
	responseDto := models.ToAPICreateTransaction(created)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responseDto)
}

// GET /transactions
func (h *HandlerImpl) ListTransactions(w http.ResponseWriter, r *http.Request) {
	filter := models.ListTransactionsInput{
		UserID:          r.URL.Query().Get("userId"),
		GroupID:         r.URL.Query().Get("groupId"),
		BalanceID:       r.URL.Query().Get("balanceId"), // Now from query parameter
		Type:            r.URL.Query().Get("type"),
		CategoryId:      r.URL.Query().Get("categoryId"),
		CategoryGroupId: r.URL.Query().Get("categoryGroupId"),
		MerchantId:      r.URL.Query().Get("merchantId"),
		SortBy:          r.URL.Query().Get("sortedBy"),
		Order:           r.URL.Query().Get("order"),
	}

	// Parse startTime if provided
	if startTimeStr := r.URL.Query().Get("startTime"); startTimeStr != "" {
		if startTime, err := time.Parse(time.RFC3339, startTimeStr); err == nil {
			filter.StartTime = startTime
		}
	}

	// Parse endTime if provided
	if endTimeStr := r.URL.Query().Get("endTime"); endTimeStr != "" {
		if endTime, err := time.Parse(time.RFC3339, endTimeStr); err == nil {
			filter.EndTime = endTime
		}
	}

	if count := r.URL.Query().Get("count"); count != "" {
		// parse count as int
		if n, err := helpers.ParseInt(count); err == nil {
			filter.Limit = n
		}
	}

	entries, err := h.Service.ListTransactionEntries(r.Context(), filter)
	if err != nil {
		h.handleServiceError(w, err, "ListTransactionEntries")
		return
	}

	// Convert to DTOs for response
	entryDtos := make([]models.TransactionEntryDto, len(entries))
	for i, entry := range entries {
		entryDtos[i] = models.ToAPITransactionEntry(&entry)
	}

	// Use paginated response structure (without actual pagination for now)
	WriteJSONListResponse(w, entryDtos, "")
}

// GET /transactions/{transaction_id}
func (h *HandlerImpl) GetTransaction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	transactionID := vars["transaction_id"]
	if transactionID == "" {
		WriteJSONError(w, http.StatusBadRequest, models.ErrorCodeBadRequest, "Missing transaction_id")
		return
	}

	tx, err := h.Service.GetTransaction(r.Context(), transactionID)
	if err != nil {
		// Try to handle as "not found" error first
		if h.handleNotFoundError(w, err, "transaction", transactionID) {
			return
		}
		// Handle all other errors (including database connection errors)
		h.handleServiceError(w, err, "GetTransaction")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tx)
}

// PUT /transactions/{transaction_id}
func (h *HandlerImpl) UpdateTransaction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	transactionID := vars["transaction_id"]
	var tx models.Transaction
	if err := json.NewDecoder(r.Body).Decode(&tx); err != nil {
		WriteJSONError(w, http.StatusBadRequest, models.ErrorCodeBadRequest, "Invalid request body: "+err.Error())
		return
	}
	// Parse UUID from string
	id, err := uuid.Parse(transactionID)
	if err != nil {
		WriteJSONError(w, http.StatusBadRequest, models.ErrorCodeBadRequest, "Invalid transaction ID format")
		return
	}
	tx.ID = id

	updated, err := h.Service.UpdateTransaction(r.Context(), tx)
	if err != nil {
		h.handleServiceError(w, err, "UpdateTransaction")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updated)
}

// DELETE /transactions/{transaction_id}
func (h *HandlerImpl) DeleteTransaction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	transactionID := vars["transaction_id"]
	if transactionID == "" {
		WriteJSONError(w, http.StatusBadRequest, models.ErrorCodeBadRequest, "Missing transaction_id")
		return
	}

	err := h.Service.DeleteTransaction(r.Context(), transactionID)
	if err != nil {
		h.handleServiceError(w, err, "DeleteTransaction")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
