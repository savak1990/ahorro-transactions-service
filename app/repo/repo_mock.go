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

func (m *MockRepository) ListTransactions(ctx context.Context, filter models.ListTransactionsFilter) ([]models.Transaction, string, error) {
	args := m.Called(ctx, filter)
	var transactions []models.Transaction
	if v := args.Get(0); v != nil {
		transactions = v.([]models.Transaction)
	}
	var nextToken string
	if v := args.Get(1); v != nil {
		nextToken = v.(string)
	}
	return transactions, nextToken, args.Error(2)
}

func (m *MockRepository) ListTransactionEntries(ctx context.Context, filter models.ListTransactionsFilter) ([]models.TransactionEntry, string, error) {
	args := m.Called(ctx, filter)
	var entries []models.TransactionEntry
	if v := args.Get(0); v != nil {
		entries = v.([]models.TransactionEntry)
	}
	var nextToken string
	if v := args.Get(1); v != nil {
		nextToken = v.(string)
	}
	return entries, nextToken, args.Error(2)
}

// Category methods

func (m *MockRepository) CreateCategory(ctx context.Context, category models.Category) (*models.Category, error) {
	args := m.Called(ctx, category)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Category), args.Error(1)
}

func (m *MockRepository) ListCategories(ctx context.Context, input models.ListCategoriesInput) ([]models.Category, string, error) {
	args := m.Called(ctx, input)
	var categories []models.Category
	if v := args.Get(0); v != nil {
		categories = v.([]models.Category)
	}
	var nextToken string
	if v := args.Get(1); v != nil {
		nextToken = v.(string)
	}
	return categories, nextToken, args.Error(2)
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

func (m *MockRepository) ListBalances(ctx context.Context, filter models.ListBalancesFilter) ([]models.Balance, error) {
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

// Helper methods for testing

// ExpectCreateTransaction sets up an expectation for CreateTransaction method
func (m *MockRepository) ExpectCreateTransaction(ctx context.Context, tx models.Transaction, result *models.Transaction, err error) *mock.Call {
	return m.On("CreateTransaction", ctx, tx).Return(result, err)
}

// ExpectListTransactions sets up an expectation for ListTransactions method
func (m *MockRepository) ExpectListTransactions(ctx context.Context, filter models.ListTransactionsFilter, result []models.Transaction, nextToken string, err error) *mock.Call {
	return m.On("ListTransactions", ctx, filter).Return(result, nextToken, err)
}

// ExpectGetTransaction sets up an expectation for GetTransaction method
func (m *MockRepository) ExpectGetTransaction(ctx context.Context, userID, transactionID string, result *models.Transaction, err error) *mock.Call {
	return m.On("GetTransaction", ctx, userID, transactionID).Return(result, err)
}

// ExpectUpdateTransaction sets up an expectation for UpdateTransaction method
func (m *MockRepository) ExpectUpdateTransaction(ctx context.Context, tx models.Transaction, result *models.Transaction, err error) *mock.Call {
	return m.On("UpdateTransaction", ctx, tx).Return(result, err)
}

// ExpectDeleteTransaction sets up an expectation for DeleteTransaction method
func (m *MockRepository) ExpectDeleteTransaction(ctx context.Context, userID, transactionID string, err error) *mock.Call {
	return m.On("DeleteTransaction", ctx, userID, transactionID).Return(err)
}

// ExpectCreateCategory sets up an expectation for CreateCategory method
func (m *MockRepository) ExpectCreateCategory(ctx context.Context, category models.Category, result *models.Category, err error) *mock.Call {
	return m.On("CreateCategory", ctx, category).Return(result, err)
}

// ExpectListCategories sets up an expectation for ListCategories method
func (m *MockRepository) ExpectListCategories(ctx context.Context, input models.ListCategoriesInput, result []models.Category, nextToken string, err error) *mock.Call {
	return m.On("ListCategories", ctx, input).Return(result, nextToken, err)
}

// ExpectDeleteCategory sets up an expectation for DeleteCategory method
func (m *MockRepository) ExpectDeleteCategory(ctx context.Context, userID, categoryID string, err error) *mock.Call {
	return m.On("DeleteCategory", ctx, userID, categoryID).Return(err)
}

// ExpectCreateBalance sets up an expectation for CreateBalance method
func (m *MockRepository) ExpectCreateBalance(ctx context.Context, balance models.Balance, result *models.Balance, err error) *mock.Call {
	return m.On("CreateBalance", ctx, balance).Return(result, err)
}

// ExpectListBalances sets up an expectation for ListBalances method
func (m *MockRepository) ExpectListBalances(ctx context.Context, filter models.ListBalancesFilter, result []models.Balance, err error) *mock.Call {
	return m.On("ListBalances", ctx, filter).Return(result, err)
}

// ExpectGetBalance sets up an expectation for GetBalance method
func (m *MockRepository) ExpectGetBalance(ctx context.Context, userID string, result *models.Balance, err error) *mock.Call {
	return m.On("GetBalance", ctx, userID).Return(result, err)
}

// ExpectUpdateBalance sets up an expectation for UpdateBalance method
func (m *MockRepository) ExpectUpdateBalance(ctx context.Context, balance models.Balance, result *models.Balance, err error) *mock.Call {
	return m.On("UpdateBalance", ctx, balance).Return(result, err)
}

// ExpectDeleteBalance sets up an expectation for DeleteBalance method
func (m *MockRepository) ExpectDeleteBalance(ctx context.Context, userID string, err error) *mock.Call {
	return m.On("DeleteBalance", ctx, userID).Return(err)
}

// Ensure MockRepository implements Repository interface
var _ Repository = (*MockRepository)(nil)
