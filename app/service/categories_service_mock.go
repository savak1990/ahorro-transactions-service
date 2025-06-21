package service

import (
	"context"

	"github.com/savak1990/transactions-service/app/models"
	"github.com/stretchr/testify/mock"
)

// MockCategoriesService provides a mock implementation of CategoriesService for testing.
type MockCategoriesService struct {
	mock.Mock
}

func (m *MockCategoriesService) ListCategoriesForUser(ctx context.Context, input models.ListCategoriesInput) ([]models.Category, string, error) {
	args := m.Called(ctx, input)
	var categories []models.Category
	if v := args.Get(0); v != nil {
		categories = v.([]models.Category)
	}
	var lastEvaluatedKey string
	if v := args.Get(1); v != nil {
		lastEvaluatedKey = v.(string)
	}
	return categories, lastEvaluatedKey, args.Error(2)
}

// Ensure MockCategoriesService implements CategoriesService interface
var _ CategoriesService = (*MockCategoriesService)(nil)
