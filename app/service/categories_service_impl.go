package service

import (
	"context"

	"github.com/savak1990/transactions-service/app/models"
	"github.com/savak1990/transactions-service/app/repo"
)

// CategoriesServiceImpl provides the implementation for CategoriesService.
type CategoriesServiceImpl struct {
	Repo repo.CategoriesRepo
}

// NewCategoriesServiceImpl creates a new CategoriesServiceImpl.
func NewCategoriesServiceImpl(repo repo.CategoriesRepo) *CategoriesServiceImpl {
	return &CategoriesServiceImpl{Repo: repo}
}

// ListCategoriesForUser returns a list of categories for the given user, with pagination.
func (s *CategoriesServiceImpl) ListCategoriesForUser(ctx context.Context, input models.ListCategoriesInput) ([]models.Category, string, error) {
	return s.Repo.ListCategoriesForUser(ctx, input)
}

// Ensure CategoriesServiceImpl implements CategoriesService interface
var _ CategoriesService = (*CategoriesServiceImpl)(nil)
