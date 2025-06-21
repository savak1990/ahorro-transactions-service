package service

import (
	"context"

	"github.com/savak1990/transactions-service/app/models"
	m "github.com/savak1990/transactions-service/app/models"
)

// TransactionsService defines the business logic for transactions.
type TransactionsService interface {
	CreateTransaction(ctx context.Context, tx m.Transaction) (*m.Transaction, error)
	ListTransactions(ctx context.Context, filter models.ListTransactionsFilter) ([]m.Transaction, string, error)
	GetTransaction(ctx context.Context, userID, transactionID string) (*m.Transaction, error)
	UpdateTransaction(ctx context.Context, tx m.Transaction) (*m.Transaction, error)
	DeleteTransaction(ctx context.Context, userID, transactionID string) error
}
