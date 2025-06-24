// filepath: /Users/savak/Projects/Ahorro/ahorro-transactions-service/app/handler/transactions_handler_mock.go
package handler

import (
	"net/http"

	"github.com/stretchr/testify/mock"
)

// HandlerMock provides mock implementations for Handler interface
// (for use in tests or as a stub)
type HandlerMock struct {
	mock.Mock
}

func NewHandlerMock() *HandlerMock { return &HandlerMock{} }

func (h *HandlerMock) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	h.Called(w, r)
}
func (h *HandlerMock) ListTransactions(w http.ResponseWriter, r *http.Request) {
	h.Called(w, r)
}
func (h *HandlerMock) GetTransaction(w http.ResponseWriter, r *http.Request) {
	h.Called(w, r)
}
func (h *HandlerMock) UpdateTransaction(w http.ResponseWriter, r *http.Request) {
	h.Called(w, r)
}
func (h *HandlerMock) DeleteTransaction(w http.ResponseWriter, r *http.Request) {
	h.Called(w, r)
}
func (h *HandlerMock) CreateBalance(w http.ResponseWriter, r *http.Request) {
	h.Called(w, r)
}
func (h *HandlerMock) ListBalances(w http.ResponseWriter, r *http.Request) {
	h.Called(w, r)
}
func (h *HandlerMock) GetBalance(w http.ResponseWriter, r *http.Request) {
	h.Called(w, r)
}
func (h *HandlerMock) UpdateBalance(w http.ResponseWriter, r *http.Request) {
	h.Called(w, r)
}
func (h *HandlerMock) DeleteBalance(w http.ResponseWriter, r *http.Request) {
	h.Called(w, r)
}
func (h *HandlerMock) CreateCategory(w http.ResponseWriter, r *http.Request) {
	h.Called(w, r)
}
func (h *HandlerMock) ListCategories(w http.ResponseWriter, r *http.Request) {
	h.Called(w, r)
}
func (h *HandlerMock) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	h.Called(w, r)
}

var _ Handler = (*HandlerMock)(nil)
