package repo

import (
	"context"
	"fmt"

	"github.com/savak1990/transactions-service/app/models"
	"gorm.io/gorm"
)

// PostgreSQLRepository implements Repository interface using PostgreSQL
type PostgreSQLRepository struct {
	db *gorm.DB
}

// NewPostgreSQLRepository creates a new PostgreSQL repository
func NewPostgreSQLRepository(db *gorm.DB) *PostgreSQLRepository {
	return &PostgreSQLRepository{
		db: db,
	}
}

// CreateTransaction creates a new transaction in the database
func (r *PostgreSQLRepository) CreateTransaction(ctx context.Context, tx models.Transaction) (*models.Transaction, error) {

	// Create the transaction in the database
	if err := r.db.WithContext(ctx).Create(tx).Error; err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}
	return &tx, nil
}

// GetTransaction retrieves a single transaction by ID
func (r *PostgreSQLRepository) GetTransaction(ctx context.Context, transactionID string) (*models.Transaction, error) {

	var tx models.Transaction
	if err := r.db.WithContext(ctx).
		Preload("Merchant").
		Preload("TransactionEntries").
		Preload("TransactionEntries.Category").
		Where("id = ?", transactionID).
		First(&tx).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("transaction not found: %s", transactionID)
		}
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}

	return &tx, nil
}

// UpdateTransaction updates an existing transaction
func (r *PostgreSQLRepository) UpdateTransaction(ctx context.Context, tx models.Transaction) (*models.Transaction, error) {

	if err := r.db.WithContext(ctx).Save(tx).Error; err != nil {
		return nil, fmt.Errorf("failed to update transaction: %w", err)
	}

	return &tx, nil
}

// DeleteTransaction soft deletes a transaction
func (r *PostgreSQLRepository) DeleteTransaction(ctx context.Context, transactionID string) error {
	if err := r.db.WithContext(ctx).Where("id = ?", transactionID).Delete(&models.Transaction{}).Error; err != nil {
		return fmt.Errorf("failed to delete transaction: %w", err)
	}
	return nil
}

// ListTransactions retrieves transactions based on the filter
func (r *PostgreSQLRepository) ListTransactions(ctx context.Context, filter models.ListTransactionsFilter) ([]models.Transaction, string, error) {
	return nil, "", nil
}

// CreateCategory creates a new category in the database
func (r *PostgreSQLRepository) CreateCategory(ctx context.Context, category models.Category) (*models.Category, error) {
	return nil, nil
}

// ListCategories retrieves all categories
func (r *PostgreSQLRepository) ListCategories(ctx context.Context, input models.ListCategoriesInput) ([]models.Category, string, error) {
	return nil, "", nil
}

// DeleteCategory deletes a category by ID
func (r *PostgreSQLRepository) DeleteCategory(ctx context.Context, categoryID string) error {
	return nil
}

// CreateBalance creates a new balance in the database
func (r *PostgreSQLRepository) CreateBalance(ctx context.Context, balance models.Balance) (*models.Balance, error) {
	return nil, nil
}

// ListBalances retrieves all balances
func (r *PostgreSQLRepository) ListBalances(ctx context.Context, filter models.ListBalancesFilter) ([]models.Balance, error) {
	return nil, nil
}

// GetBalance retrieves a balance by ID
func (r *PostgreSQLRepository) GetBalance(ctx context.Context, balanceID string) (*models.Balance, error) {
	return nil, nil
}

// UpdateBalance updates an existing balance
func (r *PostgreSQLRepository) UpdateBalance(ctx context.Context, balance models.Balance) (*models.Balance, error) {
	return nil, nil
}

// DeleteBalance deletes a balance by ID
func (r *PostgreSQLRepository) DeleteBalance(ctx context.Context, balanceID string) error {
	return nil
}

// Ensure MockRepository implements Repository interface
var _ Repository = (*PostgreSQLRepository)(nil)
