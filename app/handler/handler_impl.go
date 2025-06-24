// filepath: /Users/savak/Projects/Ahorro/ahorro-transactions-service/app/handler/transactions_handler_impl.go
package handler

import (
	"encoding/json"
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
	var tx models.Transaction
	if err := json.NewDecoder(r.Body).Decode(&tx); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	created, err := h.Service.CreateTransaction(r.Context(), tx)
	if err != nil {
		logrus.WithError(err).Error("CreateTransaction failed")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(created)
}

// GET /transactions
func (h *HandlerImpl) ListTransactions(w http.ResponseWriter, r *http.Request) {
	filter := models.ListTransactionsFilter{
		UserID:   r.URL.Query().Get("user_id"),
		Type:     r.URL.Query().Get("type"),
		Category: r.URL.Query().Get("category"),
		SortBy:   r.URL.Query().Get("sorted_by"),
		Order:    r.URL.Query().Get("order"),
	}
	if count := r.URL.Query().Get("count"); count != "" {
		// parse count as int
		if n, err := parseInt(count); err == nil {
			filter.Count = n
		}
	}
	filter.StartKey = r.URL.Query().Get("startKey")
	results, nextToken, err := h.Service.ListTransactions(r.Context(), filter)
	if err != nil {
		logrus.WithError(err).Error("ListTransactions failed")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"transactions": results,
		"nextToken":    nextToken,
	})
}

// GET /transactions/{transaction_id}
func (h *HandlerImpl) GetTransaction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := r.URL.Query().Get("user_id")
	transactionID := vars["transaction_id"]
	if userID == "" || transactionID == "" {
		http.Error(w, "Missing user_id or transaction_id", http.StatusBadRequest)
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
	var balance models.Balance
	if err := json.NewDecoder(r.Body).Decode(&balance); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	created, err := h.Service.CreateBalance(r.Context(), balance)
	if err != nil {
		logrus.WithError(err).Error("CreateBalance failed")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(created)
}

func (h *HandlerImpl) ListBalances(w http.ResponseWriter, r *http.Request) {
	filter := models.ListBalancesFilter{
		UserID:  r.URL.Query().Get("user_id"),
		GroupID: r.URL.Query().Get("group_id"),
	}
	results, err := h.Service.ListBalances(r.Context(), filter)
	if err != nil {
		logrus.WithError(err).Error("ListBalances failed")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"balances": results,
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(balance)
}

func (h *HandlerImpl) UpdateBalance(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	balanceID := vars["balance_id"]
	var balance models.Balance
	if err := json.NewDecoder(r.Body).Decode(&balance); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	id, err := uuid.Parse(balanceID)
	if err != nil {
		http.Error(w, "Invalid balance ID format", http.StatusBadRequest)
		return
	}
	balance.ID = id
	updated, err := h.Service.UpdateBalance(r.Context(), balance)
	if err != nil {
		logrus.WithError(err).Error("UpdateBalance failed")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updated)
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
	var category models.Category
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	created, err := h.Service.CreateCategory(r.Context(), category)
	if err != nil {
		logrus.WithError(err).Error("CreateCategory failed")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(created)
}

func (h *HandlerImpl) ListCategories(w http.ResponseWriter, r *http.Request) {
	filter := models.ListCategoriesInput{
		UserID: r.URL.Query().Get("user_id"),
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
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"categories": results,
	})
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

// parseInt is a helper for parsing integers from query params
func parseInt(s string) (int, error) {
	return strconv.Atoi(strings.TrimSpace(s))
}

var _ Handler = (*HandlerImpl)(nil)
