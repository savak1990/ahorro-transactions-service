// filepath: /Users/savak/Projects/Ahorro/ahorro-transactions-service/app/handler/transactions_handler_impl.go
package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/savak1990/transactions-service/app/models"
	"github.com/savak1990/transactions-service/app/service"
	"github.com/sirupsen/logrus"
)

// TransactionsHandlerImpl implements TransactionsHandler and wires to TransactionsService
// Add a field for the service dependency
type TransactionsHandlerImpl struct {
	Service service.TransactionsService
}

func NewTransactionsHandlerImpl(svc service.TransactionsService) *TransactionsHandlerImpl {
	return &TransactionsHandlerImpl{Service: svc}
}

// POST /transactions
func (h *TransactionsHandlerImpl) CreateTransaction(w http.ResponseWriter, r *http.Request) {
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
func (h *TransactionsHandlerImpl) ListTransactions(w http.ResponseWriter, r *http.Request) {
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
func (h *TransactionsHandlerImpl) GetTransaction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := r.URL.Query().Get("user_id")
	transactionID := vars["transaction_id"]
	if userID == "" || transactionID == "" {
		http.Error(w, "Missing user_id or transaction_id", http.StatusBadRequest)
		return
	}
	tx, err := h.Service.GetTransaction(r.Context(), userID, transactionID)
	if err != nil {
		logrus.WithError(err).Error("GetTransaction failed")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tx)
}

// PUT /transactions/{transaction_id}
func (h *TransactionsHandlerImpl) UpdateTransaction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	transactionID := vars["transaction_id"]
	var tx models.Transaction
	if err := json.NewDecoder(r.Body).Decode(&tx); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	tx.TransactionID = transactionID
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
func (h *TransactionsHandlerImpl) DeleteTransaction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := r.URL.Query().Get("user_id")
	transactionID := vars["transaction_id"]
	if userID == "" || transactionID == "" {
		http.Error(w, "Missing user_id or transaction_id", http.StatusBadRequest)
		return
	}
	err := h.Service.DeleteTransaction(r.Context(), userID, transactionID)
	if err != nil {
		logrus.WithError(err).Error("DeleteTransaction failed")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// parseInt is a helper for parsing integers from query params
func parseInt(s string) (int, error) {
	return strconv.Atoi(strings.TrimSpace(s))
}

var _ TransactionsHandler = (*TransactionsHandlerImpl)(nil)
