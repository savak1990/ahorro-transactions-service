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

	// Create the transaction in the database with its transaction entries
	if err := r.db.WithContext(ctx).Create(&tx).Error; err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	// Reload the transaction with all relationships
	var createdTx models.Transaction
	if err := r.db.WithContext(ctx).
		Preload("Merchant").
		Preload("TransactionEntries").
		Preload("TransactionEntries.Category").
		Where("id = ?", tx.ID).
		First(&createdTx).Error; err != nil {
		return nil, fmt.Errorf("failed to reload created transaction: %w", err)
	}

	return &createdTx, nil
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

// ListTransactions retrieves transaction entries based on the filter
func (r *PostgreSQLRepository) ListTransactions(ctx context.Context, filter models.ListTransactionsFilter) ([]models.Transaction, string, error) {
	// For now, return empty result since we need to change the return type to []models.TransactionEntry
	// This will be updated when we change the interface
	return []models.Transaction{}, "", nil
}

// ListTransactionEntries retrieves transaction entries with all related data based on the filter
func (r *PostgreSQLRepository) ListTransactionEntries(ctx context.Context, filter models.ListTransactionsFilter) ([]models.TransactionEntry, string, error) {
	var entries []models.TransactionEntry

	query := r.db.WithContext(ctx).
		Preload("Transaction").
		Preload("Transaction.Merchant").
		Preload("Transaction.Balance").
		Preload("Category")

	// Join transaction table once if we need to filter by transaction fields
	needsTransactionJoin := filter.GroupID != "" || filter.UserID != "" || filter.BalanceID != "" || filter.Type != ""
	if needsTransactionJoin {
		query = query.Joins("JOIN transaction ON transaction_entry.transaction_id = transaction.id")
	}

	// Join category table if we need to filter by category
	if filter.Category != "" {
		query = query.Joins("JOIN category ON transaction_entry.category_id = category.id")
	}

	// Apply transaction-related filters
	if filter.GroupID != "" {
		query = query.Where("transaction.group_id = ?", filter.GroupID)
	}

	if filter.UserID != "" {
		query = query.Where("transaction.user_id = ?", filter.UserID)
	}

	if filter.BalanceID != "" {
		query = query.Where("transaction.balance_id = ?", filter.BalanceID)
	}

	if filter.Type != "" {
		query = query.Where("transaction.type = ?", filter.Type)
	}

	// Apply category filter
	if filter.Category != "" {
		query = query.Where("category.category_name = ?", filter.Category)
	}

	// Handle cursor-based pagination
	if filter.StartKey != "" {
		query = query.Where("transaction_entry.id < ?", filter.StartKey)
	}

	// Apply sorting
	orderBy := "transaction_entry.created_at"
	if filter.SortBy != "" {
		switch filter.SortBy {
		case "transactedAt": // API uses camelCase, map to database field
			if needsTransactionJoin {
				orderBy = "transaction.transacted_at"
			} else {
				// If we didn't join transaction table yet, we need to join it for sorting
				query = query.Joins("JOIN transaction ON transaction_entry.transaction_id = transaction.id")
				orderBy = "transaction.transacted_at"
			}
		case "amount":
			orderBy = "transaction_entry.amount"
		case "createdAt": // API uses camelCase, map to database field
			orderBy = "transaction_entry.created_at"
		}
	}

	order := "DESC"
	if filter.Order != "" && (filter.Order == "ASC" || filter.Order == "asc") {
		order = "ASC"
	}

	query = query.Order(fmt.Sprintf("%s %s", orderBy, order))

	// Apply limit
	limit := 50 // default limit
	if filter.Count > 0 && filter.Count <= 100 {
		limit = filter.Count
	}
	query = query.Limit(limit + 1) // Get one extra to check if there are more records

	if err := query.Find(&entries).Error; err != nil {
		return nil, "", fmt.Errorf("failed to list transaction entries: %w", err)
	}

	// Handle pagination
	var nextToken string
	if len(entries) > limit {
		// Remove the extra record and set next token
		entries = entries[:limit]
		if len(entries) > 0 {
			// Use the last entry's ID as the next token (cursor-based pagination)
			nextToken = entries[len(entries)-1].ID.String()
		}
	}

	return entries, nextToken, nil
}

// CreateCategory creates a new category in the database
func (r *PostgreSQLRepository) CreateCategory(ctx context.Context, category models.Category) (*models.Category, error) {
	// Create the category in the database
	if err := r.db.WithContext(ctx).Create(&category).Error; err != nil {
		return nil, fmt.Errorf("failed to create category: %w", err)
	}
	return &category, nil
}

// ListCategories retrieves all categories
func (r *PostgreSQLRepository) ListCategories(ctx context.Context, input models.ListCategoriesInput) ([]models.Category, string, error) {
	var categories []models.Category
	query := r.db.WithContext(ctx)

	// Apply limit
	if input.Limit > 0 {
		query = query.Limit(input.Limit)
	}

	// For now, we'll ignore pagination (StartKey) since it's more complex
	// You can implement pagination later if needed

	if err := query.Order("category_name ASC").Find(&categories).Error; err != nil {
		return nil, "", fmt.Errorf("failed to list categories: %w", err)
	}

	// Return empty nextToken for now (no pagination implemented)
	return categories, "", nil
}

// DeleteCategory deletes a category by ID
func (r *PostgreSQLRepository) DeleteCategory(ctx context.Context, categoryID string) error {
	if err := r.db.WithContext(ctx).Where("id = ?", categoryID).Delete(&models.Category{}).Error; err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}
	return nil
}

// CreateBalance creates a new balance in the database
func (r *PostgreSQLRepository) CreateBalance(ctx context.Context, balance models.Balance) (*models.Balance, error) {
	// Create the balance in the database
	if err := r.db.WithContext(ctx).Create(&balance).Error; err != nil {
		return nil, fmt.Errorf("failed to create balance: %w", err)
	}
	return &balance, nil
}

// ListBalances retrieves all balances
func (r *PostgreSQLRepository) ListBalances(ctx context.Context, filter models.ListBalancesFilter) ([]models.Balance, error) {
	var balances []models.Balance
	query := r.db.WithContext(ctx)

	// Apply filters
	if filter.UserID != "" {
		query = query.Where("user_id = ?", filter.UserID)
	}
	if filter.GroupID != "" {
		query = query.Where("group_id = ?", filter.GroupID)
	}
	if filter.BalanceID != "" {
		query = query.Where("id = ?", filter.BalanceID)
	}

	// Apply ordering
	orderBy := "created_at ASC"
	if filter.SortBy != "" {
		direction := "ASC"
		if filter.Order == "desc" {
			direction = "DESC"
		}
		orderBy = fmt.Sprintf("%s %s", filter.SortBy, direction)
	}
	query = query.Order(orderBy)

	// Apply limit
	if filter.Count > 0 {
		query = query.Limit(filter.Count)
	}

	if err := query.Find(&balances).Error; err != nil {
		return nil, fmt.Errorf("failed to list balances: %w", err)
	}

	return balances, nil
}

// GetBalance retrieves a balance by ID
func (r *PostgreSQLRepository) GetBalance(ctx context.Context, balanceID string) (*models.Balance, error) {
	var balance models.Balance
	if err := r.db.WithContext(ctx).Where("id = ?", balanceID).First(&balance).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("balance not found: %s", balanceID)
		}
		return nil, fmt.Errorf("failed to get balance: %w", err)
	}
	return &balance, nil
}

// UpdateBalance updates an existing balance
func (r *PostgreSQLRepository) UpdateBalance(ctx context.Context, balance models.Balance) (*models.Balance, error) {
	if err := r.db.WithContext(ctx).Save(&balance).Error; err != nil {
		return nil, fmt.Errorf("failed to update balance: %w", err)
	}
	return &balance, nil
}

// DeleteBalance deletes a balance by ID
func (r *PostgreSQLRepository) DeleteBalance(ctx context.Context, balanceID string) error {
	if err := r.db.WithContext(ctx).Where("id = ?", balanceID).Delete(&models.Balance{}).Error; err != nil {
		return fmt.Errorf("failed to delete balance: %w", err)
	}
	return nil
}

// Ensure MockRepository implements Repository interface
var _ Repository = (*PostgreSQLRepository)(nil)
