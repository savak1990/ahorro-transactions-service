package repo

import (
	"context"

	"github.com/savak1990/transactions-service/app/models"
)

// Repo defines the repository interface for accessing transactions and categories
type Repository interface {
	// Transaction methods
	CreateTransaction(ctx context.Context, tx models.Transaction) (*models.Transaction, error)
	GetTransaction(ctx context.Context, transactionID string) (*models.Transaction, error)
	ListTransactions(ctx context.Context, filter models.ListTransactionsFilter) ([]models.Transaction, string, error)
	UpdateTransaction(ctx context.Context, tx models.Transaction) (*models.Transaction, error)
	DeleteTransaction(ctx context.Context, transactionID string) error

	// Category methods
	CreateCategory(ctx context.Context, category models.Category) (*models.Category, error)
	ListCategories(ctx context.Context, input models.ListCategoriesInput) ([]models.Category, string, error)
	DeleteCategory(ctx context.Context, categoryID string) error

	// Balance methods
	CreateBalance(ctx context.Context, balance models.Balance) (*models.Balance, error)
	ListBalances(ctx context.Context, filter models.ListBalancesFilter) ([]models.Balance, error)
	GetBalance(ctx context.Context, balanceId string) (*models.Balance, error)
	UpdateBalance(ctx context.Context, balance models.Balance) (*models.Balance, error)
	DeleteBalance(ctx context.Context, balanceId string) error
}
