package repo

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/savak1990/transactions-service/app/models"
	"gorm.io/gorm"
)

// PostgreSQLTransactionsRepository implements both TransactionsRepo and CategoriesRepo
type PostgreSQLTransactionsRepository struct {
	db *gorm.DB
}

// NewPostgreSQLTransactionsRepository creates a new PostgreSQL repository
func NewPostgreSQLTransactionsRepository(db *gorm.DB) *PostgreSQLTransactionsRepository {
	return &PostgreSQLTransactionsRepository{
		db: db,
	}
}

// Transaction Repository Methods

func (r *PostgreSQLTransactionsRepository) CreateTransaction(ctx context.Context, tx models.Transaction) (*models.Transaction, error) {
	// Convert API model to DB model
	dbTx, err := models.FromAPITransaction(tx)
	if err != nil {
		return nil, fmt.Errorf("failed to convert transaction: %w", err)
	}

	// Start transaction
	err = r.db.WithContext(ctx).Transaction(func(txDB *gorm.DB) error {
		// Create the transaction
		if err := txDB.Create(dbTx).Error; err != nil {
			return fmt.Errorf("failed to create transaction: %w", err)
		}

		// Update transaction entries with the transaction ID
		for i := range dbTx.TransactionEntries {
			dbTx.TransactionEntries[i].TransactionID = dbTx.ID
		}

		// Create transaction entries
		if len(dbTx.TransactionEntries) > 0 {
			if err := txDB.Create(&dbTx.TransactionEntries).Error; err != nil {
				return fmt.Errorf("failed to create transaction entries: %w", err)
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Load the created transaction with relationships
	var createdTx models.TransactionDB
	err = r.db.WithContext(ctx).
		Preload("TransactionEntries.Category").
		Preload("Merchant").
		First(&createdTx, "id = ?", dbTx.ID).Error
	if err != nil {
		return nil, fmt.Errorf("failed to load created transaction: %w", err)
	}

	// Convert back to API model
	result := createdTx.ToAPITransaction()
	return &result, nil
}

func (r *PostgreSQLTransactionsRepository) ListTransactions(ctx context.Context, filter models.ListTransactionsFilter) ([]models.Transaction, string, error) {
	var dbTransactions []models.TransactionDB
	query := r.db.WithContext(ctx).
		Preload("TransactionEntries.Category").
		Preload("Merchant")

	// Parse user ID
	userID, err := uuid.Parse(filter.UserID)
	if err != nil {
		return nil, "", fmt.Errorf("invalid user ID: %w", err)
	}

	// Apply filters
	query = query.Where("user_id = ?", userID)

	if filter.Type != "" {
		query = query.Where("type = ?", filter.Type)
	}

	if filter.Category != "" {
		// Join with transaction_entries and categories to filter by category
		query = query.Joins("JOIN transaction_entry ON transaction.id = transaction_entry.transaction_id").
			Joins("JOIN category ON transaction_entry.category_id = category.id").
			Where("category.category_name = ?", filter.Category)
	}

	// Apply sorting
	orderBy := "created_at DESC" // Default sorting
	if filter.SortBy != "" {
		direction := "ASC"
		if strings.ToUpper(filter.Order) == "DESC" {
			direction = "DESC"
		}

		switch filter.SortBy {
		case "amount":
			// Need to calculate amount from entries, for now use created_at
			orderBy = fmt.Sprintf("created_at %s", direction)
		case "date", "transacted_at":
			orderBy = fmt.Sprintf("transacted_at %s", direction)
		case "created_at":
			orderBy = fmt.Sprintf("created_at %s", direction)
		default:
			orderBy = fmt.Sprintf("created_at %s", direction)
		}
	}

	query = query.Order(orderBy)

	// Apply pagination
	limit := 20 // Default limit
	if filter.Count > 0 && filter.Count <= 100 {
		limit = filter.Count
	}

	// Handle start key (cursor-based pagination)
	if filter.StartKey != "" {
		// Parse start key as transaction ID
		startID, err := uuid.Parse(filter.StartKey)
		if err == nil {
			var startTx models.TransactionDB
			if err := r.db.WithContext(ctx).First(&startTx, "id = ?", startID).Error; err == nil {
				// Use the timestamp of the start transaction for pagination
				if strings.Contains(orderBy, "DESC") {
					query = query.Where("created_at < ?", startTx.CreatedAt)
				} else {
					query = query.Where("created_at > ?", startTx.CreatedAt)
				}
			}
		}
	}

	query = query.Limit(limit + 1) // Fetch one extra to determine if there's a next page

	// Execute query
	if err := query.Find(&dbTransactions).Error; err != nil {
		return nil, "", fmt.Errorf("failed to list transactions: %w", err)
	}

	// Determine next page token
	var nextToken string
	if len(dbTransactions) > limit {
		// Remove the extra record and use the last record's ID as next token
		lastTx := dbTransactions[limit-1]
		nextToken = lastTx.ID.String()
		dbTransactions = dbTransactions[:limit]
	}

	// Convert to API models
	result := make([]models.Transaction, len(dbTransactions))
	for i, dbTx := range dbTransactions {
		result[i] = dbTx.ToAPITransaction()
	}

	return result, nextToken, nil
}

func (r *PostgreSQLTransactionsRepository) GetTransaction(ctx context.Context, userID, transactionID string) (*models.Transaction, error) {
	// Parse IDs
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	txUUID, err := uuid.Parse(transactionID)
	if err != nil {
		return nil, fmt.Errorf("invalid transaction ID: %w", err)
	}

	var dbTx models.TransactionDB
	err = r.db.WithContext(ctx).
		Preload("TransactionEntries.Category").
		Preload("Merchant").
		Where("id = ? AND user_id = ?", txUUID, userUUID).
		First(&dbTx).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("transaction not found")
		}
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}

	result := dbTx.ToAPITransaction()
	return &result, nil
}

func (r *PostgreSQLTransactionsRepository) UpdateTransaction(ctx context.Context, tx models.Transaction) (*models.Transaction, error) {
	// Parse transaction ID
	txUUID, err := uuid.Parse(tx.TransactionID)
	if err != nil {
		return nil, fmt.Errorf("invalid transaction ID: %w", err)
	}

	userUUID, err := uuid.Parse(tx.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	var result *models.Transaction
	err = r.db.WithContext(ctx).Transaction(func(txDB *gorm.DB) error {
		// Get existing transaction
		var existingTx models.TransactionDB
		if err := txDB.Where("id = ? AND user_id = ?", txUUID, userUUID).First(&existingTx).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return fmt.Errorf("transaction not found")
			}
			return fmt.Errorf("failed to find transaction: %w", err)
		}

		// Convert API model to updates
		updates := map[string]interface{}{
			"type":          tx.Type,
			"approved_at":   parseTimeOrNow(tx.ApprovedAt),
			"transacted_at": parseTimeOrNow(tx.TransactedAt),
			"updated_at":    time.Now(),
		}

		// Update transaction
		if err := txDB.Model(&existingTx).Updates(updates).Error; err != nil {
			return fmt.Errorf("failed to update transaction: %w", err)
		}

		// Update transaction entries (simple approach: delete and recreate)
		if err := txDB.Where("transaction_id = ?", txUUID).Delete(&models.TransactionEntry{}).Error; err != nil {
			return fmt.Errorf("failed to delete old transaction entries: %w", err)
		}

		// Create new entry
		newEntry := models.TransactionEntry{
			TransactionID: txUUID,
			Description:   &tx.Description,
			Amount:        tx.Amount,
		}

		if err := txDB.Create(&newEntry).Error; err != nil {
			return fmt.Errorf("failed to create new transaction entry: %w", err)
		}

		// Load updated transaction
		var updatedTx models.TransactionDB
		if err := txDB.Preload("TransactionEntries.Category").Preload("Merchant").First(&updatedTx, "id = ?", txUUID).Error; err != nil {
			return fmt.Errorf("failed to load updated transaction: %w", err)
		}

		apiResult := updatedTx.ToAPITransaction()
		result = &apiResult
		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *PostgreSQLTransactionsRepository) DeleteTransaction(ctx context.Context, userID, transactionID string) error {
	// Parse IDs
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	txUUID, err := uuid.Parse(transactionID)
	if err != nil {
		return fmt.Errorf("invalid transaction ID: %w", err)
	}

	return r.db.WithContext(ctx).Transaction(func(txDB *gorm.DB) error {
		// Soft delete transaction entries
		if err := txDB.Where("transaction_id = ?", txUUID).Delete(&models.TransactionEntry{}).Error; err != nil {
			return fmt.Errorf("failed to delete transaction entries: %w", err)
		}

		// Soft delete transaction
		result := txDB.Where("id = ? AND user_id = ?", txUUID, userUUID).Delete(&models.TransactionDB{})
		if result.Error != nil {
			return fmt.Errorf("failed to delete transaction: %w", result.Error)
		}

		if result.RowsAffected == 0 {
			return fmt.Errorf("transaction not found")
		}

		return nil
	})
}

// Categories Repository Methods

func (r *PostgreSQLTransactionsRepository) ListCategoriesForUser(ctx context.Context, input models.ListCategoriesInput) ([]models.Category, string, error) {
	var dbCategories []models.CategoryDB
	query := r.db.WithContext(ctx)

	// Apply ordering by priority and name
	query = query.Order("priority DESC, category_name ASC")

	// Apply limit
	limit := 20 // Default limit
	if input.Limit > 0 && input.Limit <= 100 {
		limit = input.Limit
	}

	// Handle start key pagination
	if input.StartKey != "" {
		startID, err := uuid.Parse(input.StartKey)
		if err == nil {
			var startCategory models.CategoryDB
			if err := r.db.WithContext(ctx).First(&startCategory, "id = ?", startID).Error; err == nil {
				// Use priority and name for consistent pagination
				query = query.Where("(priority < ?) OR (priority = ? AND category_name > ?)",
					startCategory.Priority, startCategory.Priority, startCategory.CategoryName)
			}
		}
	}

	query = query.Limit(limit + 1) // Fetch one extra to determine next page

	if err := query.Find(&dbCategories).Error; err != nil {
		return nil, "", fmt.Errorf("failed to list categories: %w", err)
	}

	// Determine next page token
	var nextToken string
	if len(dbCategories) > limit {
		lastCategory := dbCategories[limit-1]
		nextToken = lastCategory.ID.String()
		dbCategories = dbCategories[:limit]
	}

	// Convert to API models
	result := make([]models.Category, len(dbCategories))
	for i, dbCategory := range dbCategories {
		result[i] = dbCategory.ToAPICategory()
	}

	return result, nextToken, nil
}

// Helper functions

func parseTimeOrNow(timeStr string) time.Time {
	if timeStr == "" {
		return time.Now()
	}

	t, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		return time.Now()
	}

	return t
}
