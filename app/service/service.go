package service

import (
	"context"

	m "github.com/savak1990/transactions-service/app/models"
)

// Service defines the business logic for transactions.
type Service interface {
	CreateTransaction(ctx context.Context, tx m.Transaction) (*m.Transaction, error)
	CreateTransactions(ctx context.Context, transactions []m.Transaction) ([]m.Transaction, *string, error) // Returns transactions, operationID, error
	GetTransaction(ctx context.Context, transactionID string) (*m.SingleTransactionDto, error)
	UpdateTransaction(ctx context.Context, transactionID string, updateDto m.UpdateTransactionDto) (*m.Transaction, error)
	DeleteTransaction(ctx context.Context, transactionID string) error
	ListTransactions(ctx context.Context, filter m.ListTransactionsInput) ([]m.Transaction, error)
	ListTransactionEntries(ctx context.Context, filter m.ListTransactionsInput) ([]m.TransactionEntry, error)

	CreateBalance(ctx context.Context, balance m.Balance) (*m.Balance, error)
	GetBalance(ctx context.Context, balanceID string) (*m.Balance, error)
	UpdateBalance(ctx context.Context, balance m.Balance) (*m.Balance, error)
	DeleteBalance(ctx context.Context, balanceID string) error
	DeleteBalancesByUserId(ctx context.Context, userId string) error
	ListBalances(ctx context.Context, filter m.ListBalancesInput) ([]m.Balance, error)

	CreateCategory(ctx context.Context, category m.Category) (*m.Category, error)
	ListCategories(ctx context.Context, filter m.ListCategoriesInput) ([]m.Category, error)
	GetCategory(ctx context.Context, categoryID string) (*m.Category, error)
	UpdateCategory(ctx context.Context, category m.Category) (*m.Category, error)
	DeleteCategory(ctx context.Context, categoryID string) error
	DeleteCategoriesByUserId(ctx context.Context, userId string) error

	CreateCategoryGroup(ctx context.Context, categoryGroup m.CategoryGroup) (*m.CategoryGroup, error)
	ListCategoryGroups(ctx context.Context, filter m.ListCategoryGroupsInput) ([]m.CategoryGroup, error)
	GetCategoryGroup(ctx context.Context, categoryGroupID string) (*m.CategoryGroup, error)
	UpdateCategoryGroup(ctx context.Context, categoryGroup m.CategoryGroup) (*m.CategoryGroup, error)
	DeleteCategoryGroup(ctx context.Context, categoryGroupID string) error

	CreateMerchant(ctx context.Context, merchant m.Merchant) (*m.Merchant, error)
	GetMerchant(ctx context.Context, merchantID string) (*m.Merchant, error)
	UpdateMerchant(ctx context.Context, merchant m.Merchant) (*m.Merchant, error)
	DeleteMerchant(ctx context.Context, merchantID string) error
	DeleteMerchantsByUserId(ctx context.Context, userId string) error
	ListMerchants(ctx context.Context, filter m.ListMerchantsInput) ([]m.Merchant, error)

	// Transaction statistics
	GetTransactionStats(ctx context.Context, filter m.TransactionStatsInput) (*m.TransactionStatsResponseDto, error)
}
