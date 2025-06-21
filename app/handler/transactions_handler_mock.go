// filepath: /Users/savak/Projects/Ahorro/ahorro-transactions-service/app/handler/transactions_handler_mock.go
package handler

import (
	"net/http"

	"github.com/stretchr/testify/mock"
)

// TransactionsHandlerMock provides mock implementations for TransactionsHandler interface
// (for use in tests or as a stub)
type TransactionsHandlerMock struct {
	mock.Mock
}

func NewTransactionsHandlerMock() *TransactionsHandlerMock { return &TransactionsHandlerMock{} }

func (h *TransactionsHandlerMock) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	h.Called(w, r)
}
func (h *TransactionsHandlerMock) ListTransactions(w http.ResponseWriter, r *http.Request) {
	h.Called(w, r)
}
func (h *TransactionsHandlerMock) GetTransaction(w http.ResponseWriter, r *http.Request) {
	h.Called(w, r)
}
func (h *TransactionsHandlerMock) UpdateTransaction(w http.ResponseWriter, r *http.Request) {
	h.Called(w, r)
}
func (h *TransactionsHandlerMock) DeleteTransaction(w http.ResponseWriter, r *http.Request) {
	h.Called(w, r)
}

var _ TransactionsHandler = (*TransactionsHandlerMock)(nil)
