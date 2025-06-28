package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
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
		logrus.WithError(err).Error("CreateTransaction failed")
		if isDatabaseTimeoutError(err) {
			WriteJSONError(w, http.StatusServiceUnavailable, models.ErrorCodeDbTimeout, "Database is warming up, please retry in a few seconds")
		} else if isDatabaseError(err) {
			WriteJSONError(w, http.StatusInternalServerError, models.ErrorCodeDbError, err.Error())
		} else {
			WriteJSONError(w, http.StatusInternalServerError, models.ErrorCodeInternalServer, err.Error())
		}
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
		logrus.WithError(err).Error("ListTransactionEntries failed")
		WriteJSONError(w, http.StatusInternalServerError, models.ErrorCodeInternalServer, err.Error())
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
		http.Error(w, "Missing transaction_id", http.StatusBadRequest)
		return
	}

	tx, err := h.Service.GetTransaction(r.Context(), transactionID)
	if err != nil {
		logrus.WithError(err).Error("GetTransaction failed")
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	// Parse UUID from string
	id, err := uuid.Parse(transactionID)
	if err != nil {
		http.Error(w, "Invalid transaction ID format", http.StatusBadRequest)
		return
	}
	tx.ID = id

	updated, err := h.Service.UpdateTransaction(r.Context(), tx)
	if err != nil {
		logrus.WithError(err).Error("UpdateTransaction failed")
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
		http.Error(w, "Missing transaction_id", http.StatusBadRequest)
		return
	}

	err := h.Service.DeleteTransaction(r.Context(), transactionID)
	if err != nil {
		logrus.WithError(err).Error("DeleteTransaction failed")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Balance handlers
func (h *HandlerImpl) CreateBalance(w http.ResponseWriter, r *http.Request) {
	var balanceDto models.BalanceDto
	if err := json.NewDecoder(r.Body).Decode(&balanceDto); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Convert DTO to DAO model
	balance, err := models.FromAPIBalance(balanceDto)
	if err != nil {
		http.Error(w, "Invalid balance data: "+err.Error(), http.StatusBadRequest)
		return
	}

	created, err := h.Service.CreateBalance(r.Context(), *balance)
	if err != nil {
		logrus.WithError(err).Error("CreateBalance failed")
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
		logrus.WithError(err).Error("ListBalances failed")
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
		http.Error(w, "Missing balance_id", http.StatusBadRequest)
		return
	}
	balance, err := h.Service.GetBalance(r.Context(), balanceID)
	if err != nil {
		logrus.WithError(err).Error("GetBalance failed")
		if err.Error() == fmt.Sprintf("balance not found: %s", balanceID) {
			http.Error(w, "Balance not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if balance == nil {
		http.Error(w, "Balance not found", http.StatusNotFound)
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
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Parse balance ID and set it in DTO
	id, err := uuid.Parse(balanceID)
	if err != nil {
		http.Error(w, "Invalid balance ID format", http.StatusBadRequest)
		return
	}
	balanceDto.BalanceID = id.String()

	// Convert DTO to DAO model
	balance, err := models.FromAPIBalance(balanceDto)
	if err != nil {
		http.Error(w, "Invalid balance data: "+err.Error(), http.StatusBadRequest)
		return
	}

	updated, err := h.Service.UpdateBalance(r.Context(), *balance)
	if err != nil {
		logrus.WithError(err).Error("UpdateBalance failed")
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
		http.Error(w, "Missing balance_id", http.StatusBadRequest)
		return
	}
	err := h.Service.DeleteBalance(r.Context(), balanceID)
	if err != nil {
		logrus.WithError(err).Error("DeleteBalance failed")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Category handlers
func (h *HandlerImpl) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var categoryDto models.CategoryDto
	if err := json.NewDecoder(r.Body).Decode(&categoryDto); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Convert DTO to DAO model
	category, err := models.FromAPICategory(categoryDto)
	if err != nil {
		http.Error(w, "Invalid category data: "+err.Error(), http.StatusBadRequest)
		return
	}

	created, err := h.Service.CreateCategory(r.Context(), *category)
	if err != nil {
		logrus.WithError(err).Error("CreateCategory failed")
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
		logrus.WithError(err).Error("ListCategories failed")
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
		http.Error(w, "Missing category_id", http.StatusBadRequest)
		return
	}

	category, err := h.Service.GetCategory(r.Context(), categoryID)
	if err != nil {
		logrus.WithError(err).Error("GetCategory failed")
		if err.Error() == fmt.Sprintf("category not found: %s", categoryID) {
			http.Error(w, "Category not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if category == nil {
		http.Error(w, "Category not found", http.StatusNotFound)
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
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Parse category ID and set it in DTO
	id, err := uuid.Parse(categoryID)
	if err != nil {
		http.Error(w, "Invalid category ID format", http.StatusBadRequest)
		return
	}
	categoryDto.CategoryID = id.String()

	// Convert DTO to DAO model
	category, err := models.FromAPICategory(categoryDto)
	if err != nil {
		http.Error(w, "Invalid category data: "+err.Error(), http.StatusBadRequest)
		return
	}

	updated, err := h.Service.UpdateCategory(r.Context(), *category)
	if err != nil {
		logrus.WithError(err).Error("UpdateCategory failed")
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
		http.Error(w, "Missing category_id", http.StatusBadRequest)
		return
	}
	err := h.Service.DeleteCategory(r.Context(), categoryID)
	if err != nil {
		logrus.WithError(err).Error("DeleteCategory failed")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Merchant handlers
func (h *HandlerImpl) CreateMerchant(w http.ResponseWriter, r *http.Request) {
	var merchantDto models.MerchantDto
	if err := json.NewDecoder(r.Body).Decode(&merchantDto); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Convert DTO to DAO model
	merchant, err := models.FromAPIMerchant(merchantDto)
	if err != nil {
		http.Error(w, "Invalid merchant data: "+err.Error(), http.StatusBadRequest)
		return
	}

	created, err := h.Service.CreateMerchant(r.Context(), *merchant)
	if err != nil {
		logrus.WithError(err).Error("CreateMerchant failed")
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
		logrus.WithError(err).Error("ListMerchants failed")
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
		http.Error(w, "Missing merchant_id", http.StatusBadRequest)
		return
	}

	merchant, err := h.Service.GetMerchant(r.Context(), merchantID)
	if err != nil {
		logrus.WithError(err).Error("GetMerchant failed")
		if err.Error() == fmt.Sprintf("merchant not found: %s", merchantID) {
			http.Error(w, "Merchant not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if merchant == nil {
		http.Error(w, "Merchant not found", http.StatusNotFound)
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
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Parse merchant ID and set it in DTO
	id, err := uuid.Parse(merchantID)
	if err != nil {
		http.Error(w, "Invalid merchant ID format", http.StatusBadRequest)
		return
	}
	merchantDto.MerchantID = id.String()

	// Convert DTO to DAO model
	merchant, err := models.FromAPIMerchant(merchantDto)
	if err != nil {
		http.Error(w, "Invalid merchant data: "+err.Error(), http.StatusBadRequest)
		return
	}

	updated, err := h.Service.UpdateMerchant(r.Context(), *merchant)
	if err != nil {
		logrus.WithError(err).Error("UpdateMerchant failed")
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
		http.Error(w, "Missing merchant_id", http.StatusBadRequest)
		return
	}
	err := h.Service.DeleteMerchant(r.Context(), merchantID)
	if err != nil {
		logrus.WithError(err).Error("DeleteMerchant failed")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// parseInt is a helper for parsing integers from query params
func parseInt(s string) (int, error) {
	return strconv.Atoi(strings.TrimSpace(s))
}

var _ Handler = (*HandlerImpl)(nil)
