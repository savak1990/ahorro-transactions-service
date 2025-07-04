package repo

import (
	"context"
	"fmt"

	"github.com/savak1990/transactions-service/app/models"
	"gorm.io/gorm"
)

// CreateTransaction creates a new transaction in the database
func (r *PostgreSQLRepository) CreateTransaction(ctx context.Context, tx models.Transaction) (*models.Transaction, error) {

	db := r.getDB()

	// Create the transaction in the database with its transaction entries
	if err := db.WithContext(ctx).Create(&tx).Error; err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	// Reload the transaction with all relationships
	var createdTx models.Transaction
	if err := db.WithContext(ctx).
		Preload("Merchant").
		Preload("TransactionEntries").
		Preload("TransactionEntries.Category").
		Preload("TransactionEntries.Category.CategoryGroup"). // Load CategoryGroup without filtering to detect soft-deleted groups
		Where("id = ?", tx.ID).
		First(&createdTx).Error; err != nil {
		return nil, fmt.Errorf("failed to reload created transaction: %w", err)
	}

	return &createdTx, nil
}

// CreateTransactions creates multiple transactions atomically in a single database transaction
func (r *PostgreSQLRepository) CreateTransactions(ctx context.Context, transactions []models.Transaction) ([]models.Transaction, error) {
	db := r.getDB()

	var createdTransactions []models.Transaction

	// Use a database transaction to ensure atomicity
	err := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Create all transactions in the same database transaction
		for _, transaction := range transactions {
			if err := tx.Create(&transaction).Error; err != nil {
				return fmt.Errorf("failed to create transaction %s: %w", transaction.ID.String(), err)
			}

			// Reload the transaction with all relationships
			var createdTx models.Transaction
			if err := tx.
				Preload("Merchant").
				Preload("TransactionEntries").
				Preload("TransactionEntries.Category").
				Preload("TransactionEntries.Category.CategoryGroup").
				Where("id = ?", transaction.ID).
				First(&createdTx).Error; err != nil {
				return fmt.Errorf("failed to reload created transaction %s: %w", transaction.ID.String(), err)
			}

			createdTransactions = append(createdTransactions, createdTx)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return createdTransactions, nil
}

// GetTransaction retrieves a single transaction by ID
func (r *PostgreSQLRepository) GetTransaction(ctx context.Context, transactionID string) (*models.Transaction, error) {

	db := r.getDB()

	var tx models.Transaction
	if err := db.WithContext(ctx).
		Preload("Merchant", "deleted_at IS NULL"). // Only load non-deleted merchants
		Preload("TransactionEntries").
		Preload("TransactionEntries.Category").
		Preload("TransactionEntries.Category.CategoryGroup"). // Load CategoryGroup without filtering to detect soft-deleted groups
		Joins("JOIN balance ON transaction.balance_id = balance.id").
		Where("transaction.id = ? AND balance.deleted_at IS NULL", transactionID).
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

	db := r.getDB()

	if err := db.WithContext(ctx).Save(tx).Error; err != nil {
		return nil, fmt.Errorf("failed to update transaction: %w", err)
	}

	return &tx, nil
}

// DeleteTransaction soft deletes a transaction
func (r *PostgreSQLRepository) DeleteTransaction(ctx context.Context, transactionID string) error {
	db := r.getDB()
	if err := db.WithContext(ctx).Where("id = ?", transactionID).Delete(&models.Transaction{}).Error; err != nil {
		return fmt.Errorf("failed to delete transaction: %w", err)
	}
	return nil
}

// ListTransactions retrieves transaction entries based on the filter
func (r *PostgreSQLRepository) ListTransactions(ctx context.Context, filter models.ListTransactionsInput) ([]models.Transaction, error) {
	// For now, return empty result since we need to change the return type to []models.TransactionEntry
	// This will be updated when we change the interface
	return []models.Transaction{}, nil
}

// ListTransactionEntries retrieves transaction entries with all related data based on the filter
func (r *PostgreSQLRepository) ListTransactionEntries(ctx context.Context, filter models.ListTransactionsInput) ([]models.TransactionEntry, error) {
	var entries []models.TransactionEntry

	db := r.getDB()

	query := db.WithContext(ctx).
		Preload("Transaction").
		Preload("Transaction.Merchant", "deleted_at IS NULL"). // Only load non-deleted merchants
		Preload("Transaction.Balance").
		Preload("Category").
		Preload("Category.CategoryGroup") // Load CategoryGroup without filtering to detect soft-deleted groups

	// Always join with balance table to exclude soft-deleted balances
	query = query.Joins("JOIN transaction ON transaction_entry.transaction_id = transaction.id").
		Joins("JOIN balance ON transaction.balance_id = balance.id").
		Where("balance.deleted_at IS NULL")

	// Join category table if we need to filter by category or category group
	if len(filter.CategoryIds) > 0 || len(filter.CategoryGroupIds) > 0 {
		query = query.Joins("JOIN category ON transaction_entry.category_id = category.id")
	}

	// Join merchant table if we need to filter by merchant (exclude soft-deleted merchants)
	if filter.MerchantId != "" {
		query = query.Joins("JOIN merchant ON transaction.merchant_id = merchant.id AND merchant.deleted_at IS NULL")
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

	if filter.TransactionID != "" {
		query = query.Where("transaction.id = ?", filter.TransactionID)
	}

	if len(filter.Types) > 0 {
		query = query.Where("transaction.type IN ?", filter.Types)
	}

	// Apply date range filters
	if !filter.StartTime.IsZero() {
		query = query.Where("transaction.transacted_at >= ?", filter.StartTime)
	}

	if !filter.EndTime.IsZero() {
		query = query.Where("transaction.transacted_at <= ?", filter.EndTime)
	}

	// Apply category filter
	if len(filter.CategoryIds) > 0 {
		query = query.Where("transaction_entry.category_id IN ?", filter.CategoryIds)
	}

	// Apply category group filter
	if len(filter.CategoryGroupIds) > 0 {
		query = query.Where("category.category_group_id IN ?", filter.CategoryGroupIds)
	}

	// Apply merchant filter
	if filter.MerchantId != "" {
		query = query.Where("transaction.merchant_id = ?", filter.MerchantId)
	}

	// Apply sorting
	orderBy := "transaction_entry.created_at"
	if filter.SortBy != "" {
		switch filter.SortBy {
		case "transactedAt": // API uses camelCase, map to database field
			orderBy = "transaction.transacted_at"
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
	if filter.Limit > 0 && filter.Limit <= 100 {
		limit = filter.Limit
	}
	query = query.Limit(limit)

	if err := query.Find(&entries).Error; err != nil {
		return nil, fmt.Errorf("failed to list transaction entries: %w", err)
	}

	return entries, nil
}
