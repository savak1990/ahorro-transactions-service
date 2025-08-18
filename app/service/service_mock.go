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

func (svc *MockService) CreateTransactions(ctx context.Context, transactions []models.Transaction) ([]models.Transaction, *string, error) {
	args := svc.Called(ctx, transactions)
	if args.Get(0) == nil {
		return nil, nil, args.Error(2)
	}
	var operationID *string
	if args.Get(1) != nil {
		operationID = args.Get(1).(*string)
	}
	return args.Get(0).([]models.Transaction), operationID, args.Error(2)
}

func (svc *MockService) GetTransaction(ctx context.Context, transactionID string) (*models.SingleTransactionDto, error) {
	args := svc.Called(ctx, transactionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.SingleTransactionDto), args.Error(1)
}

func (svc *MockService) UpdateTransaction(ctx context.Context, transactionID string, updateDto models.UpdateTransactionDto) (*models.Transaction, error) {
	args := svc.Called(ctx, transactionID, updateDto)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Transaction), args.Error(1)
}

func (svc *MockService) DeleteTransaction(ctx context.Context, transactionID string) error {
	args := svc.Called(ctx, transactionID)
	return args.Error(0)
}

func (svc *MockService) ListTransactions(ctx context.Context, filter models.ListTransactionsInput) ([]models.Transaction, error) {
	args := svc.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Transaction), args.Error(1)
}

func (svc *MockService) ListTransactionEntries(ctx context.Context, filter models.ListTransactionsInput) ([]models.TransactionEntry, error) {
	args := svc.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.TransactionEntry), args.Error(1)
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

func (svc *MockService) DeleteBalancesByUserId(ctx context.Context, userId string) error {
	args := svc.Called(ctx, userId)
	return args.Error(0)
}

func (svc *MockService) ListBalances(ctx context.Context, filter models.ListBalancesInput) ([]models.Balance, error) {
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

func (svc *MockService) GetCategory(ctx context.Context, categoryID string) (*models.Category, error) {
	args := svc.Called(ctx, categoryID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Category), args.Error(1)
}

func (svc *MockService) UpdateCategory(ctx context.Context, category models.Category) (*models.Category, error) {
	args := svc.Called(ctx, category)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Category), args.Error(1)
}

func (svc *MockService) DeleteCategory(ctx context.Context, categoryID string) error {
	args := svc.Called(ctx, categoryID)
	return args.Error(0)
}

func (svc *MockService) DeleteCategoriesByUserId(ctx context.Context, userId string) error {
	args := svc.Called(ctx, userId)
	return args.Error(0)
}

func (svc *MockService) CreateMerchant(ctx context.Context, merchant models.Merchant) (*models.Merchant, error) {
	args := svc.Called(ctx, merchant)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Merchant), args.Error(1)
}

func (svc *MockService) GetMerchant(ctx context.Context, merchantID string) (*models.Merchant, error) {
	args := svc.Called(ctx, merchantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Merchant), args.Error(1)
}

func (svc *MockService) UpdateMerchant(ctx context.Context, merchant models.Merchant) (*models.Merchant, error) {
	args := svc.Called(ctx, merchant)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Merchant), args.Error(1)
}

func (svc *MockService) DeleteMerchant(ctx context.Context, merchantID string) error {
	args := svc.Called(ctx, merchantID)
	return args.Error(0)
}

func (svc *MockService) DeleteMerchantsByUserId(ctx context.Context, userId string) error {
	args := svc.Called(ctx, userId)
	return args.Error(0)
}

func (svc *MockService) ListMerchants(ctx context.Context, filter models.ListMerchantsInput) ([]models.Merchant, error) {
	args := svc.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Merchant), args.Error(1)
}

// CategoryGroup service mock methods

func (svc *MockService) CreateCategoryGroup(ctx context.Context, categoryGroup models.CategoryGroup) (*models.CategoryGroup, error) {
	args := svc.Called(ctx, categoryGroup)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.CategoryGroup), args.Error(1)
}

func (svc *MockService) ListCategoryGroups(ctx context.Context, filter models.ListCategoryGroupsInput) ([]models.CategoryGroup, error) {
	args := svc.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.CategoryGroup), args.Error(1)
}

func (svc *MockService) GetCategoryGroup(ctx context.Context, categoryGroupID string) (*models.CategoryGroup, error) {
	args := svc.Called(ctx, categoryGroupID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.CategoryGroup), args.Error(1)
}

func (svc *MockService) UpdateCategoryGroup(ctx context.Context, categoryGroup models.CategoryGroup) (*models.CategoryGroup, error) {
	args := svc.Called(ctx, categoryGroup)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.CategoryGroup), args.Error(1)
}

func (svc *MockService) DeleteCategoryGroup(ctx context.Context, categoryGroupID string) error {
	args := svc.Called(ctx, categoryGroupID)
	return args.Error(0)
}

func (svc *MockService) GetTransactionStats(ctx context.Context, input models.TransactionStatsInput) ([]models.TransactionStatsItemDto, error) {
	args := svc.Called(ctx, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.TransactionStatsItemDto), args.Error(1)
}

// Ensure MockService implements Service
var _ Service = (*MockService)(nil)
