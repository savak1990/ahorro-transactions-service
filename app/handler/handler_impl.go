package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/savak1990/transactions-service/app/aws"
	"github.com/savak1990/transactions-service/app/models"
	"github.com/savak1990/transactions-service/app/service"
	"github.com/sirupsen/logrus"
)

// HandlerImpl implements Handler and wires to TransactionsService
// Add a field for the service dependency
type HandlerImpl struct {
	Service service.Service
}

func NewHandlerImpl(svc service.Service) *HandlerImpl {
	return &HandlerImpl{Service: svc}
}

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
	filter := models.ListTransactionsFilter{
		UserID:    r.URL.Query().Get("userId"),
		GroupID:   r.URL.Query().Get("groupId"),
		BalanceID: r.URL.Query().Get("balanceId"), // Now from query parameter
		Type:      r.URL.Query().Get("type"),
		Category:  r.URL.Query().Get("category"),
		SortBy:    r.URL.Query().Get("sortedBy"),
		Order:     r.URL.Query().Get("order"),
	}
	if count := r.URL.Query().Get("count"); count != "" {
		// parse count as int
		if n, err := parseInt(count); err == nil {
			filter.Count = n
		}
	}
	filter.StartKey = r.URL.Query().Get("startKey")

	entries, nextToken, err := h.Service.ListTransactionEntries(r.Context(), filter)
	if err != nil {
		h.handleServiceError(w, err, "ListTransactionEntries")
		return
	}

	// Convert to DTOs for response
	entryDtos := make([]models.TransactionEntryDto, len(entries))
	for i, entry := range entries {
		entryDtos[i] = models.ToAPITransactionEntry(&entry)
	}

	WriteJSONListResponse(w, entryDtos, nextToken)
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
	filter := models.ListBalancesFilter{
		UserID:  r.URL.Query().Get("userId"),
		GroupID: r.URL.Query().Get("groupId"),
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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"balances": balanceDtos,
	})
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
		WriteJSONError(w, http.StatusBadRequest, models.ErrorCodeBadRequest, "Missing balance_id")
		return
	}

	err := h.Service.DeleteBalance(r.Context(), balanceID)
	if err != nil {
		h.handleServiceError(w, err, "DeleteBalance")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Category handlers
func (h *HandlerImpl) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var categoryDto models.CategoryDto
	if err := json.NewDecoder(r.Body).Decode(&categoryDto); err != nil {
		WriteJSONError(w, http.StatusBadRequest, models.ErrorCodeBadRequest, "Invalid request body: "+err.Error())
		return
	}

	// Convert DTO to DAO model
	category, err := models.FromAPICategory(categoryDto)
	if err != nil {
		WriteJSONError(w, http.StatusBadRequest, models.ErrorCodeBadRequest, "Invalid category data: "+err.Error())
		return
	}

	created, err := h.Service.CreateCategory(r.Context(), *category)
	if err != nil {
		h.handleServiceError(w, err, "CreateCategory")
		return
	}

	// Convert back to DTO for response
	responseDto := models.ToAPICategory(created)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responseDto)
}

func (h *HandlerImpl) ListCategories(w http.ResponseWriter, r *http.Request) {
	filter := models.ListCategoriesInput{
		UserID: r.URL.Query().Get("userId"),
	}
	if limit := r.URL.Query().Get("limit"); limit != "" {
		if n, err := parseInt(limit); err == nil {
			filter.Limit = n
		}
	}
	filter.StartKey = r.URL.Query().Get("startKey")
	results, err := h.Service.ListCategories(r.Context(), filter)
	if err != nil {
		h.handleServiceError(w, err, "ListCategories")
		return
	}

	// Convert to DTOs for response
	categoryDtos := make([]models.CategoryDto, len(results))
	for i, category := range results {
		categoryDtos[i] = models.ToAPICategory(&category)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"categories": categoryDtos,
	})
}

func (h *HandlerImpl) GetCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	categoryID := vars["category_id"]
	if categoryID == "" {
		WriteJSONError(w, http.StatusBadRequest, models.ErrorCodeBadRequest, "Missing category_id")
		return
	}

	category, err := h.Service.GetCategory(r.Context(), categoryID)
	if err != nil {
		// Try to handle as "not found" error first
		if h.handleNotFoundError(w, err, "category", categoryID) {
			return
		}
		// Handle all other errors (including database connection errors)
		h.handleServiceError(w, err, "GetCategory")
		return
	}

	if category == nil {
		WriteJSONError(w, http.StatusNotFound, models.ErrorCodeNotFound, "Category not found")
		return
	}

	// Convert to DTO for response
	responseDto := models.ToAPICategory(category)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responseDto)
}

func (h *HandlerImpl) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	categoryID := vars["category_id"]
	var categoryDto models.CategoryDto
	if err := json.NewDecoder(r.Body).Decode(&categoryDto); err != nil {
		WriteJSONError(w, http.StatusBadRequest, models.ErrorCodeBadRequest, "Invalid request body: "+err.Error())
		return
	}

	// Parse category ID and set it in DTO
	id, err := uuid.Parse(categoryID)
	if err != nil {
		WriteJSONError(w, http.StatusBadRequest, models.ErrorCodeBadRequest, "Invalid category ID format")
		return
	}
	categoryDto.CategoryID = id.String()

	// Convert DTO to DAO model
	category, err := models.FromAPICategory(categoryDto)
	if err != nil {
		WriteJSONError(w, http.StatusBadRequest, models.ErrorCodeBadRequest, "Invalid category data: "+err.Error())
		return
	}

	updated, err := h.Service.UpdateCategory(r.Context(), *category)
	if err != nil {
		h.handleServiceError(w, err, "UpdateCategory")
		return
	}

	// Convert back to DTO for response
	responseDto := models.ToAPICategory(updated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responseDto)
}

func (h *HandlerImpl) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	categoryID := vars["category_id"]
	if categoryID == "" {
		WriteJSONError(w, http.StatusBadRequest, models.ErrorCodeBadRequest, "Missing category_id")
		return
	}
	err := h.Service.DeleteCategory(r.Context(), categoryID)
	if err != nil {
		h.handleServiceError(w, err, "DeleteCategory")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

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
	filter := models.ListMerchantsFilter{
		Name:   r.URL.Query().Get("name"),
		SortBy: r.URL.Query().Get("sortBy"),
		Order:  r.URL.Query().Get("order"),
	}
	if count := r.URL.Query().Get("limit"); count != "" {
		if n, err := parseInt(count); err == nil {
			filter.Count = n
		}
	}
	filter.StartKey = r.URL.Query().Get("startKey")

	results, nextToken, err := h.Service.ListMerchants(r.Context(), filter)
	if err != nil {
		h.handleServiceError(w, err, "ListMerchants")
		return
	}

	// Convert to DTOs for response
	merchantDtos := make([]models.MerchantDto, len(results))
	for i, merchant := range results {
		merchantDtos[i] = models.ToAPIMerchant(&merchant)
	}

	response := map[string]interface{}{
		"merchants": merchantDtos,
	}
	if nextToken != "" {
		response["nextToken"] = nextToken
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
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

// parseInt is a helper for parsing integers from query params
func parseInt(s string) (int, error) {
	return strconv.Atoi(strings.TrimSpace(s))
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

// handleDatabaseOperation wraps database operations to ensure connection errors are properly panicked
func handleDatabaseOperation(operation func() error) {
	err := operation()
	if err != nil && isDatabaseConnectionError(err) {
		// Convert database connection errors to proper panics for middleware
		panic(&aws.DatabaseConnectionError{
			Message: "Database connection failed during operation",
			Cause:   err,
		})
	}
	if err != nil {
		// Re-throw other errors normally
		panic(err)
	}
}

var _ Handler = (*HandlerImpl)(nil)

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
