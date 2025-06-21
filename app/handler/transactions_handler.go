package handler

import "net/http"

// TransactionsHandler interface
// Handles all transaction-related endpoints
type TransactionsHandler interface {
	CreateTransaction(http.ResponseWriter, *http.Request)
	ListTransactions(http.ResponseWriter, *http.Request)
	GetTransaction(http.ResponseWriter, *http.Request)
	UpdateTransaction(http.ResponseWriter, *http.Request)
	DeleteTransaction(http.ResponseWriter, *http.Request)
}
