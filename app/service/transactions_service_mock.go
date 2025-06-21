package service

import (
	"context"

	"github.com/savak1990/transactions-service/app/models"
	"github.com/stretchr/testify/mock"
)

type MockTransactionsService struct {
	mock.Mock
}

func (svc *MockTransactionsService) CreateTransaction(ctx context.Context, tx models.Transaction) (*models.Transaction, error) {
	args := svc.Called(ctx, tx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Transaction), args.Error(1)
}

func (svc *MockTransactionsService) ListTransactions(ctx context.Context, filter models.ListTransactionsFilter) ([]models.Transaction, string, error) {
	args := svc.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, "", args.Error(2)
	}
	return args.Get(0).([]models.Transaction), args.String(1), args.Error(2)
}

func (svc *MockTransactionsService) GetTransaction(ctx context.Context, userID, transactionID string) (*models.Transaction, error) {
	args := svc.Called(ctx, userID, transactionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Transaction), args.Error(1)
}

func (svc *MockTransactionsService) UpdateTransaction(ctx context.Context, tx models.Transaction) (*models.Transaction, error) {
	args := svc.Called(ctx, tx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Transaction), args.Error(1)
}

func (svc *MockTransactionsService) DeleteTransaction(ctx context.Context, userID, transactionID string) error {
	args := svc.Called(ctx, userID, transactionID)
	return args.Error(0)
}

// Ensure MockTransactionsService implements TransactionsService
var _ TransactionsService = (*MockTransactionsService)(nil)
