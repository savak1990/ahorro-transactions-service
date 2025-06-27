package service

import (
	"context"

	m "github.com/savak1990/transactions-service/app/models"
)

// Service defines the business logic for transactions.
type Service interface {
	CreateTransaction(ctx context.Context, tx m.Transaction) (*m.Transaction, error)
	GetTransaction(ctx context.Context, transactionID string) (*m.Transaction, error)
	UpdateTransaction(ctx context.Context, tx m.Transaction) (*m.Transaction, error)
	DeleteTransaction(ctx context.Context, transactionID string) error
	ListTransactions(ctx context.Context, filter m.ListTransactionsFilter) ([]m.Transaction, string, error)
	ListTransactionEntries(ctx context.Context, filter m.ListTransactionsFilter) ([]m.TransactionEntry, string, error)

	CreateBalance(ctx context.Context, balance m.Balance) (*m.Balance, error)
	GetBalance(ctx context.Context, balanceID string) (*m.Balance, error)
	UpdateBalance(ctx context.Context, balance m.Balance) (*m.Balance, error)
	DeleteBalance(ctx context.Context, balanceID string) error
	ListBalances(ctx context.Context, filter m.ListBalancesFilter) ([]m.Balance, error)

	CreateCategory(ctx context.Context, category m.Category) (*m.Category, error)
	ListCategories(ctx context.Context, filter m.ListCategoriesInput) ([]m.Category, error)
	DeleteCategory(ctx context.Context, categoryID string) error
}
