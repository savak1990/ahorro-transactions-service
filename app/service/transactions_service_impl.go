package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/savak1990/transactions-service/app/models"
	repo "github.com/savak1990/transactions-service/app/repo"
)

type TransactionsServiceImpl struct {
	repo repo.TransactionsRepo
}

func NewTransactionsServiceImpl(repo repo.TransactionsRepo) *TransactionsServiceImpl {
	return &TransactionsServiceImpl{
		repo: repo,
	}
}

func (s *TransactionsServiceImpl) CreateTransaction(ctx context.Context, tx models.Transaction) (*models.Transaction, error) {
	tx.TransactionID = uuid.NewString()
	now := time.Now().UTC().Format(time.RFC3339)
	tx.CreatedAt = now
	tx.UpdatedAt = now
	return s.repo.CreateTransaction(ctx, tx)
}

func (s *TransactionsServiceImpl) ListTransactions(ctx context.Context, filter models.ListTransactionsFilter) ([]models.Transaction, string, error) {
	return s.repo.ListTransactions(ctx, filter)
}

func (s *TransactionsServiceImpl) GetTransaction(ctx context.Context, userID, transactionID string) (*models.Transaction, error) {
	if transactionID == "" {
		return nil, errors.New("transactionID required")
	}
	return s.repo.GetTransaction(ctx, userID, transactionID)
}

func (s *TransactionsServiceImpl) UpdateTransaction(ctx context.Context, tx models.Transaction) (*models.Transaction, error) {
	tx.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	return s.repo.UpdateTransaction(ctx, tx)
}

func (s *TransactionsServiceImpl) DeleteTransaction(ctx context.Context, userID, transactionID string) error {
	if transactionID == "" {
		return errors.New("transactionID required")
	}
	return s.repo.DeleteTransaction(ctx, userID, transactionID)
}

// Ensure TransactionsServiceImpl implements TransactionsService
var _ TransactionsService = (*TransactionsServiceImpl)(nil)
