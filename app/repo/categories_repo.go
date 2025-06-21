package repo

import (
	"context"

	"github.com/savak1990/transactions-service/app/models"
)

// CategoriesRepo defines the repository interface for accessing categories in the data store.
type CategoriesRepo interface {
	ListCategoriesForUser(ctx context.Context, input models.ListCategoriesInput) ([]models.Category, string, error)
}
