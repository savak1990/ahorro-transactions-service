package repo

import (
	"context"

	"github.com/savak1990/transactions-service/app/models"
	"github.com/stretchr/testify/mock"
)

// MockRepository provides a mock implementation of Repository interface for testing.
type MockRepository struct {
	mock.Mock
}

// NewMockRepository creates a new MockRepository instance.
func NewMockRepository() *MockRepository {
	return &MockRepository{}
}

// Transaction methods

func (m *MockRepository) CreateTransaction(ctx context.Context, tx models.Transaction) (*models.Transaction, error) {
	args := m.Called(ctx, tx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Transaction), args.Error(1)
}

func (m *MockRepository) GetTransaction(ctx context.Context, transactionID string) (*models.Transaction, error) {
	args := m.Called(ctx, transactionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Transaction), args.Error(1)
}

func (m *MockRepository) UpdateTransaction(ctx context.Context, tx models.Transaction) (*models.Transaction, error) {
	args := m.Called(ctx, tx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Transaction), args.Error(1)
}

func (m *MockRepository) DeleteTransaction(ctx context.Context, transactionID string) error {
	args := m.Called(ctx, transactionID)
	return args.Error(0)
}

func (m *MockRepository) ListTransactions(ctx context.Context, filter models.ListTransactionsInput) ([]models.Transaction, error) {
	args := m.Called(ctx, filter)
	var transactions []models.Transaction
	if v := args.Get(0); v != nil {
		transactions = v.([]models.Transaction)
	}
	return transactions, args.Error(1)
}

func (m *MockRepository) ListTransactionEntries(ctx context.Context, filter models.ListTransactionsInput) ([]models.TransactionEntry, error) {
	args := m.Called(ctx, filter)
	var entries []models.TransactionEntry
	if v := args.Get(0); v != nil {
		entries = v.([]models.TransactionEntry)
	}
	return entries, args.Error(1)
}

// Category methods

func (m *MockRepository) CreateCategory(ctx context.Context, category models.Category) (*models.Category, error) {
	args := m.Called(ctx, category)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Category), args.Error(1)
}

func (m *MockRepository) ListCategories(ctx context.Context, input models.ListCategoriesInput) ([]models.Category, error) {
	args := m.Called(ctx, input)
	var categories []models.Category
	if v := args.Get(0); v != nil {
		categories = v.([]models.Category)
	}
	return categories, args.Error(1)
}

func (m *MockRepository) GetCategory(ctx context.Context, categoryID string) (*models.Category, error) {
	args := m.Called(ctx, categoryID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Category), args.Error(1)
}

func (m *MockRepository) UpdateCategory(ctx context.Context, category models.Category) (*models.Category, error) {
	args := m.Called(ctx, category)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Category), args.Error(1)
}

func (m *MockRepository) DeleteCategory(ctx context.Context, categoryID string) error {
	args := m.Called(ctx, categoryID)
	return args.Error(0)
}

// Balance methods

func (m *MockRepository) CreateBalance(ctx context.Context, balance models.Balance) (*models.Balance, error) {
	args := m.Called(ctx, balance)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Balance), args.Error(1)
}

func (m *MockRepository) ListBalances(ctx context.Context, filter models.ListBalancesInput) ([]models.Balance, error) {
	args := m.Called(ctx, filter)
	var balances []models.Balance
	if v := args.Get(0); v != nil {
		balances = v.([]models.Balance)
	}
	return balances, args.Error(1)
}

func (m *MockRepository) GetBalance(ctx context.Context, userID string) (*models.Balance, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Balance), args.Error(1)
}

func (m *MockRepository) UpdateBalance(ctx context.Context, balance models.Balance) (*models.Balance, error) {
	args := m.Called(ctx, balance)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Balance), args.Error(1)
}

func (m *MockRepository) DeleteBalance(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

// Merchant methods

func (m *MockRepository) CreateMerchant(ctx context.Context, merchant models.Merchant) (*models.Merchant, error) {
	args := m.Called(ctx, merchant)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Merchant), args.Error(1)
}

func (m *MockRepository) ListMerchants(ctx context.Context, filter models.ListMerchantsInput) ([]models.Merchant, error) {
	args := m.Called(ctx, filter)
	var merchants []models.Merchant
	if v := args.Get(0); v != nil {
		merchants = v.([]models.Merchant)
	}
	return merchants, args.Error(1)
}

func (m *MockRepository) GetMerchant(ctx context.Context, merchantId string) (*models.Merchant, error) {
	args := m.Called(ctx, merchantId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Merchant), args.Error(1)
}

func (m *MockRepository) UpdateMerchant(ctx context.Context, merchant models.Merchant) (*models.Merchant, error) {
	args := m.Called(ctx, merchant)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Merchant), args.Error(1)
}

func (m *MockRepository) DeleteMerchant(ctx context.Context, merchantId string) error {
	args := m.Called(ctx, merchantId)
	return args.Error(0)
}

// Helper methods for testing

// ExpectCreateTransaction sets up an expectation for CreateTransaction method
func (m *MockRepository) ExpectCreateTransaction(ctx context.Context, tx models.Transaction, result *models.Transaction, err error) *mock.Call {
	return m.On("CreateTransaction", ctx, tx).Return(result, err)
}

// ExpectListTransactions sets up an expectation for ListTransactions method
func (m *MockRepository) ExpectListTransactions(ctx context.Context, filter models.ListTransactionsInput, result []models.Transaction, err error) *mock.Call {
	return m.On("ListTransactions", ctx, filter).Return(result, err)
}

// ExpectGetTransaction sets up an expectation for GetTransaction method
func (m *MockRepository) ExpectGetTransaction(ctx context.Context, transactionID string, result *models.Transaction, err error) *mock.Call {
	return m.On("GetTransaction", ctx, transactionID).Return(result, err)
}

// ExpectUpdateTransaction sets up an expectation for UpdateTransaction method
func (m *MockRepository) ExpectUpdateTransaction(ctx context.Context, tx models.Transaction, result *models.Transaction, err error) *mock.Call {
	return m.On("UpdateTransaction", ctx, tx).Return(result, err)
}

// ExpectDeleteTransaction sets up an expectation for DeleteTransaction method
func (m *MockRepository) ExpectDeleteTransaction(ctx context.Context, transactionID string, err error) *mock.Call {
	return m.On("DeleteTransaction", ctx, transactionID).Return(err)
}

// ExpectCreateCategory sets up an expectation for CreateCategory method
func (m *MockRepository) ExpectCreateCategory(ctx context.Context, category models.Category, result *models.Category, err error) *mock.Call {
	return m.On("CreateCategory", ctx, category).Return(result, err)
}

// ExpectListCategories sets up an expectation for ListCategories method
func (m *MockRepository) ExpectListCategories(ctx context.Context, input models.ListCategoriesInput, result []models.Category, err error) *mock.Call {
	return m.On("ListCategories", ctx, input).Return(result, err)
}

// ExpectDeleteCategory sets up an expectation for DeleteCategory method
func (m *MockRepository) ExpectDeleteCategory(ctx context.Context, categoryID string, err error) *mock.Call {
	return m.On("DeleteCategory", ctx, categoryID).Return(err)
}

// ExpectCreateBalance sets up an expectation for CreateBalance method
func (m *MockRepository) ExpectCreateBalance(ctx context.Context, balance models.Balance, result *models.Balance, err error) *mock.Call {
	return m.On("CreateBalance", ctx, balance).Return(result, err)
}

// ExpectListBalances sets up an expectation for ListBalances method
func (m *MockRepository) ExpectListBalances(ctx context.Context, filter models.ListBalancesInput, result []models.Balance, err error) *mock.Call {
	return m.On("ListBalances", ctx, filter).Return(result, err)
}

// ExpectGetBalance sets up an expectation for GetBalance method
func (m *MockRepository) ExpectGetBalance(ctx context.Context, balanceId string, result *models.Balance, err error) *mock.Call {
	return m.On("GetBalance", ctx, balanceId).Return(result, err)
}

// ExpectUpdateBalance sets up an expectation for UpdateBalance method
func (m *MockRepository) ExpectUpdateBalance(ctx context.Context, balance models.Balance, result *models.Balance, err error) *mock.Call {
	return m.On("UpdateBalance", ctx, balance).Return(result, err)
}

// ExpectDeleteBalance sets up an expectation for DeleteBalance method
func (m *MockRepository) ExpectDeleteBalance(ctx context.Context, balanceId string, err error) *mock.Call {
	return m.On("DeleteBalance", ctx, balanceId).Return(err)
}

// ExpectCreateMerchant sets up an expectation for CreateMerchant method
func (m *MockRepository) ExpectCreateMerchant(ctx context.Context, merchant models.Merchant, result *models.Merchant, err error) *mock.Call {
	return m.On("CreateMerchant", ctx, merchant).Return(result, err)
}

// ExpectListMerchants sets up an expectation for ListMerchants method
func (m *MockRepository) ExpectListMerchants(ctx context.Context, filter models.ListMerchantsInput, result []models.Merchant, err error) *mock.Call {
	return m.On("ListMerchants", ctx, filter).Return(result, err)
}

// ExpectGetMerchant sets up an expectation for GetMerchant method
func (m *MockRepository) ExpectGetMerchant(ctx context.Context, merchantId string, result *models.Merchant, err error) *mock.Call {
	return m.On("GetMerchant", ctx, merchantId).Return(result, err)
}

// ExpectUpdateMerchant sets up an expectation for UpdateMerchant method
func (m *MockRepository) ExpectUpdateMerchant(ctx context.Context, merchant models.Merchant, result *models.Merchant, err error) *mock.Call {
	return m.On("UpdateMerchant", ctx, merchant).Return(result, err)
}

// ExpectDeleteMerchant sets up an expectation for DeleteMerchant method
func (m *MockRepository) ExpectDeleteMerchant(ctx context.Context, merchantId string, err error) *mock.Call {
	return m.On("DeleteMerchant", ctx, merchantId).Return(err)
}

// Ensure MockRepository implements Repository interface
var _ Repository = (*MockRepository)(nil)
