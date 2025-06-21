package service

import (
	"context"

	"github.com/savak1990/transactions-service/app/models"
)

// CategoriesService defines business logic for categories.
type CategoriesService interface {
	// ListCategoriesForUser returns a list of categories for the given user, with pagination support.
	ListCategoriesForUser(ctx context.Context, input models.ListCategoriesInput) ([]models.Category, string, error)
}
