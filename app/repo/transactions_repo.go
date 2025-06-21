package repo

import (
	"context"

	"github.com/savak1990/transactions-service/app/models"
)

type TransactionsRepo interface {
	CreateTransaction(ctx context.Context, tx models.Transaction) (*models.Transaction, error)
	ListTransactions(ctx context.Context, filter models.ListTransactionsFilter) ([]models.Transaction, string, error)
	GetTransaction(ctx context.Context, userID, transactionID string) (*models.Transaction, error)
	UpdateTransaction(ctx context.Context, tx models.Transaction) (*models.Transaction, error)
	DeleteTransaction(ctx context.Context, userID, transactionID string) error
}
