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
	ListTransactions(ctx context.Context, filter models.ListTransactionsInput) ([]models.Transaction, error)
	ListTransactionEntries(ctx context.Context, filter models.ListTransactionsInput) ([]models.TransactionEntry, error)
	UpdateTransaction(ctx context.Context, tx models.Transaction) (*models.Transaction, error)
	DeleteTransaction(ctx context.Context, transactionID string) error

	// Category methods
	CreateCategory(ctx context.Context, category models.Category) (*models.Category, error)
	ListCategories(ctx context.Context, input models.ListCategoriesInput) ([]models.Category, error)
	GetCategory(ctx context.Context, categoryID string) (*models.Category, error)
	UpdateCategory(ctx context.Context, category models.Category) (*models.Category, error)
	DeleteCategory(ctx context.Context, categoryID string) error

	// Balance methods
	CreateBalance(ctx context.Context, balance models.Balance) (*models.Balance, error)
	ListBalances(ctx context.Context, filter models.ListBalancesInput) ([]models.Balance, error)
	GetBalance(ctx context.Context, balanceId string) (*models.Balance, error)
	UpdateBalance(ctx context.Context, balance models.Balance) (*models.Balance, error)
	DeleteBalance(ctx context.Context, balanceId string) error

	// Merchant methods
	CreateMerchant(ctx context.Context, merchant models.Merchant) (*models.Merchant, error)
	ListMerchants(ctx context.Context, filter models.ListMerchantsInput) ([]models.Merchant, error)
	GetMerchant(ctx context.Context, merchantId string) (*models.Merchant, error)
	UpdateMerchant(ctx context.Context, merchant models.Merchant) (*models.Merchant, error)
	DeleteMerchant(ctx context.Context, merchantId string) error
}
