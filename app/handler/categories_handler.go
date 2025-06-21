package handler

import "net/http"

// CategoriesHandler interface
// Handles all category-related endpoints
type CategoriesHandler interface {
	ListCategoriesForUser(http.ResponseWriter, *http.Request)
}
