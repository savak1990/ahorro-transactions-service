package service

import (
	"context"

	"github.com/savak1990/transactions-service/app/models"
	"github.com/stretchr/testify/mock"
)

type MockService struct {
	mock.Mock
}

func (svc *MockService) CreateTransaction(ctx context.Context, tx models.Transaction) (*models.Transaction, error) {
	args := svc.Called(ctx, tx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Transaction), args.Error(1)
}

func (svc *MockService) GetTransaction(ctx context.Context, transactionID string) (*models.Transaction, error) {
	args := svc.Called(ctx, transactionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Transaction), args.Error(1)
}

func (svc *MockService) UpdateTransaction(ctx context.Context, tx models.Transaction) (*models.Transaction, error) {
	args := svc.Called(ctx, tx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Transaction), args.Error(1)
}

func (svc *MockService) DeleteTransaction(ctx context.Context, transactionID string) error {
	args := svc.Called(ctx, transactionID)
	return args.Error(0)
}

func (svc *MockService) ListTransactions(ctx context.Context, filter models.ListTransactionsFilter) ([]models.Transaction, string, error) {
	args := svc.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, "", args.Error(2)
	}
	return args.Get(0).([]models.Transaction), args.String(1), args.Error(2)
}

func (svc *MockService) ListTransactionEntries(ctx context.Context, filter models.ListTransactionsFilter) ([]models.TransactionEntry, string, error) {
	args := svc.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, "", args.Error(2)
	}
	return args.Get(0).([]models.TransactionEntry), args.String(1), args.Error(2)
}

func (svc *MockService) CreateBalance(ctx context.Context, balance models.Balance) (*models.Balance, error) {
	args := svc.Called(ctx, balance)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Balance), args.Error(1)
}

func (svc *MockService) GetBalance(ctx context.Context, balanceID string) (*models.Balance, error) {
	args := svc.Called(ctx, balanceID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Balance), args.Error(1)
}

func (svc *MockService) UpdateBalance(ctx context.Context, balance models.Balance) (*models.Balance, error) {
	args := svc.Called(ctx, balance)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Balance), args.Error(1)
}

func (svc *MockService) DeleteBalance(ctx context.Context, balanceID string) error {
	args := svc.Called(ctx, balanceID)
	return args.Error(0)
}

func (svc *MockService) ListBalances(ctx context.Context, filter models.ListBalancesFilter) ([]models.Balance, error) {
	args := svc.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Balance), args.Error(1)
}

func (svc *MockService) CreateCategory(ctx context.Context, category models.Category) (*models.Category, error) {
	args := svc.Called(ctx, category)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Category), args.Error(1)
}

func (svc *MockService) ListCategories(ctx context.Context, filter models.ListCategoriesInput) ([]models.Category, error) {
	args := svc.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Category), args.Error(1)
}

func (svc *MockService) DeleteCategory(ctx context.Context, categoryID string) error {
	args := svc.Called(ctx, categoryID)
	return args.Error(0)
}

// Ensure MockService implements Service
var _ Service = (*MockService)(nil)
