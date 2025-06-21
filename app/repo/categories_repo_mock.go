package repo

import (
	"context"

	"github.com/savak1990/transactions-service/app/models"
	"github.com/stretchr/testify/mock"
)

// MockCategoriesRepo provides a mock implementation of CategoriesRepo for testing.
type MockCategoriesRepo struct {
	mock.Mock
}

func NewMockCategoriesRepo() *MockCategoriesRepo {
	return &MockCategoriesRepo{}
}

func (m *MockCategoriesRepo) ListCategoriesForUser(ctx context.Context, input models.ListCategoriesInput) ([]models.Category, string, error) {
	args := m.Called(ctx, input)
	var categories []models.Category
	if v := args.Get(0); v != nil {
		categories = v.([]models.Category)
	}
	var startKey string
	if v := args.Get(1); v != nil {
		startKey = v.(string)
	}
	return categories, startKey, args.Error(2)
}

// Ensure MockCategoriesRepo implements CategoriesRepo interface
var _ CategoriesRepo = (*MockCategoriesRepo)(nil)
