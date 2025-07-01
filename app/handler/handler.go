package handler

import "net/http"

// Handler interface
// Handles all transaction-related endpoints
type Handler interface {
	CreateTransaction(http.ResponseWriter, *http.Request)
	ListTransactions(http.ResponseWriter, *http.Request)
	GetTransaction(http.ResponseWriter, *http.Request)
	UpdateTransaction(http.ResponseWriter, *http.Request)
	DeleteTransaction(http.ResponseWriter, *http.Request)

	CreateBalance(http.ResponseWriter, *http.Request)
	ListBalances(http.ResponseWriter, *http.Request)
	GetBalance(http.ResponseWriter, *http.Request)
	UpdateBalance(http.ResponseWriter, *http.Request)
	DeleteBalance(http.ResponseWriter, *http.Request)
	DeleteBalancesByUserId(http.ResponseWriter, *http.Request)

	CreateCategory(http.ResponseWriter, *http.Request)
	ListCategories(http.ResponseWriter, *http.Request)
	GetCategory(http.ResponseWriter, *http.Request)
	UpdateCategory(http.ResponseWriter, *http.Request)
	DeleteCategory(http.ResponseWriter, *http.Request)
	DeleteCategoriesByUserId(http.ResponseWriter, *http.Request)

	CreateCategoryGroup(http.ResponseWriter, *http.Request)
	ListCategoryGroups(http.ResponseWriter, *http.Request)
	GetCategoryGroup(http.ResponseWriter, *http.Request)
	UpdateCategoryGroup(http.ResponseWriter, *http.Request)
	DeleteCategoryGroup(http.ResponseWriter, *http.Request)

	CreateMerchant(http.ResponseWriter, *http.Request)
	ListMerchants(http.ResponseWriter, *http.Request)
	GetMerchant(http.ResponseWriter, *http.Request)
	UpdateMerchant(http.ResponseWriter, *http.Request)
	DeleteMerchant(http.ResponseWriter, *http.Request)
	DeleteMerchantsByUserId(http.ResponseWriter, *http.Request)
}
