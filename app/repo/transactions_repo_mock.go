package repo

import (
	"context"

	"github.com/savak1990/transactions-service/app/models"
	"github.com/stretchr/testify/mock"
)

type MockTransactionsRepo struct {
	mock.Mock
}

func (m *MockTransactionsRepo) CreateTransaction(ctx context.Context, tx models.Transaction) (*models.Transaction, error) {
	args := m.Called(ctx, tx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Transaction), args.Error(1)
}

func (m *MockTransactionsRepo) ListTransactions(ctx context.Context, filter models.ListTransactionsFilter) ([]models.Transaction, string, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, "", args.Error(2)
	}
	return args.Get(0).([]models.Transaction), args.String(1), args.Error(2)
}

func (m *MockTransactionsRepo) GetTransaction(ctx context.Context, userID, transactionID string) (*models.Transaction, error) {
	args := m.Called(ctx, userID, transactionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Transaction), args.Error(1)
}

func (m *MockTransactionsRepo) UpdateTransaction(ctx context.Context, tx models.Transaction) (*models.Transaction, error) {
	args := m.Called(ctx, tx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Transaction), args.Error(1)
}

func (m *MockTransactionsRepo) DeleteTransaction(ctx context.Context, userID, transactionID string) error {
	args := m.Called(ctx, userID, transactionID)
	return args.Error(0)
}

var _ TransactionsRepo = (*MockTransactionsRepo)(nil)
