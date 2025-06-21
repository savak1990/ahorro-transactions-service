package handler

import (
	"net/http"

	"github.com/stretchr/testify/mock"
)

// CategoriesHandlerMock provides mock implementations for CategoriesHandler interface
// (for use in tests or as a stub)
type CategoriesHandlerMock struct {
	mock.Mock
}

func NewCategoriesHandlerMock() *CategoriesHandlerMock { return &CategoriesHandlerMock{} }

func (h *CategoriesHandlerMock) ListCategoriesForUser(w http.ResponseWriter, r *http.Request) {
	h.Called(w, r)
}

var _ CategoriesHandler = (*CategoriesHandlerMock)(nil)
