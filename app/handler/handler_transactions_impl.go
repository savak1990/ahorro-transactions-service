package handler

import (
	"encoding/json"
	"fmt"
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
	var rawBody json.RawMessage
	if err := json.NewDecoder(r.Body).Decode(&rawBody); err != nil {
		WriteJSONError(w, http.StatusBadRequest, models.ErrorCodeBadRequest, "Invalid request body: "+err.Error())
		return
	}

	// Try to parse as batch request first
	var batchReq models.CreateTransactionsRequestDto
	if err := json.Unmarshal(rawBody, &batchReq); err == nil && len(batchReq.Transactions) > 0 {
		// Handle batch request
		h.handleBatchTransactions(w, r, batchReq)
		return
	}

	// Fall back to single transaction
	var req models.CreateTransactionDto
	if err := json.Unmarshal(rawBody, &req); err != nil {
		WriteJSONError(w, http.StatusBadRequest, models.ErrorCodeBadRequest, "Invalid request body: "+err.Error())
		return
	}

	// Validate single transaction cannot be move_in or move_out
	if req.Type == "move_in" || req.Type == "move_out" {
		WriteJSONError(w, http.StatusBadRequest, models.ErrorCodeBadRequest, "Single transactions cannot be of type 'move_in' or 'move_out'. Use batch transactions for movement operations.")
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

// handleBatchTransactions handles creation of multiple transactions
func (h *HandlerImpl) handleBatchTransactions(w http.ResponseWriter, r *http.Request, batchReq models.CreateTransactionsRequestDto) {
	// Validate transaction count (max 5)
	if len(batchReq.Transactions) > 5 {
		WriteJSONError(w, http.StatusConflict, models.ErrorCodeConflict, "Maximum 5 transactions allowed per batch")
		return
	}

	// Validate movement operations: must have exactly 2 transactions with one move_in and one move_out
	if err := h.validateMovementTransactions(batchReq.Transactions); err != nil {
		WriteJSONError(w, http.StatusBadRequest, models.ErrorCodeBadRequest, err.Error())
		return
	}

	// Convert request DTOs to transactions
	transactions, err := models.FromAPICreateTransactionsRequest(batchReq)
	if err != nil {
		logrus.WithError(err).Error("Failed to convert batch request DTOs")
		WriteJSONError(w, http.StatusBadRequest, models.ErrorCodeBadRequest, "Invalid request: "+err.Error())
		return
	}

	// Create transactions using service
	created, operationID, err := h.Service.CreateTransactions(r.Context(), transactions)
	if err != nil {
		h.handleServiceError(w, err, "CreateTransactions")
		return
	}

	// Convert response to DTO
	responseDto := models.ToAPICreateTransactionsResponse(created, operationID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responseDto)
}

// GET /transactions
func (h *HandlerImpl) ListTransactions(w http.ResponseWriter, r *http.Request) {
	filter := models.ListTransactionsInput{
		UserID:  r.URL.Query().Get("userId"),
		GroupID: r.URL.Query().Get("groupId"),
		SortBy:  r.URL.Query().Get("sortedBy"),
		Order:   r.URL.Query().Get("order"),
	}

	// Parse types array from query parameter - support both formats:
	// 1. type=expense,income (comma-separated)
	// 2. type=expense&type=income (multiple parameters)
	filter.BalanceIds = ParseQueryStringArray(r.URL.Query(), "balanceId")
	filter.Types = ParseQueryStringArray(r.URL.Query(), "type")
	filter.CategoryIds = ParseQueryStringArray(r.URL.Query(), "categoryId")
	filter.CategoryGroupIds = ParseQueryStringArray(r.URL.Query(), "categoryGroupId")
	filter.TransactionIds = ParseQueryStringArray(r.URL.Query(), "transactionId")
	filter.MerchantIds = ParseQueryStringArray(r.URL.Query(), "merchantId")
	filter.OperationIds = ParseQueryStringArray(r.URL.Query(), "operationId")

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

// validateMovementTransactions validates that movement operations follow the correct rules
func (h *HandlerImpl) validateMovementTransactions(transactions []models.CreateTransactionDto) error {
	moveInCount := 0
	moveOutCount := 0
	hasMovementTypes := false

	// Count movement transaction types
	for _, tx := range transactions {
		switch tx.Type {
		case "move_in":
			moveInCount++
			hasMovementTypes = true
		case "move_out":
			moveOutCount++
			hasMovementTypes = true
		case "expense", "income", "movement":
			// Regular transaction types are allowed in batch, but not mixed with move_in/move_out
			if hasMovementTypes {
				return fmt.Errorf("cannot mix movement types (move_in/move_out) with regular transaction types (expense/income/movement) in the same batch")
			}
		default:
			return fmt.Errorf("invalid transaction type: %s", tx.Type)
		}
	}

	// If there are movement types, validate the rules
	if hasMovementTypes {
		if len(transactions) != 2 {
			return fmt.Errorf("movement operations must contain exactly 2 transactions, got %d", len(transactions))
		}
		if moveInCount != 1 {
			return fmt.Errorf("movement operations must contain exactly 1 'move_in' transaction, got %d", moveInCount)
		}
		if moveOutCount != 1 {
			return fmt.Errorf("movement operations must contain exactly 1 'move_out' transaction, got %d", moveOutCount)
		}
	}

	return nil
}
