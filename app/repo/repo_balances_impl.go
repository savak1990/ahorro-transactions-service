package repo

import (
	"context"
	"fmt"

	"github.com/savak1990/transactions-service/app/models"
	"gorm.io/gorm"
)

// CreateBalance creates a new balance in the database
func (r *PostgreSQLRepository) CreateBalance(ctx context.Context, balance models.Balance) (*models.Balance, error) {
	db := r.getDB()
	// Create the balance in the database
	if err := db.WithContext(ctx).Create(&balance).Error; err != nil {
		return nil, fmt.Errorf("failed to create balance: %w", err)
	}
	return &balance, nil
}

// ListBalances retrieves all balances
func (r *PostgreSQLRepository) ListBalances(ctx context.Context, filter models.ListBalancesInput) ([]models.Balance, error) {
	var balances []models.Balance
	db := r.getDB()
	query := db.WithContext(ctx)

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
		if filter.Order == "desc" || filter.Order == "DESC" {
			direction = "DESC"
		}

		// Validate sortBy field
		validSortFields := map[string]string{
			"rank":      "rank",
			"createdAt": "created_at",
			"updatedAt": "updated_at",
			"title":     "title",
			"name":      "title", // Map 'name' to the actual column name 'title'
		}

		if dbField, valid := validSortFields[filter.SortBy]; valid {
			orderBy = fmt.Sprintf("%s %s", dbField, direction)
		} else {
			orderBy = "created_at ASC"
		}
	}
	query = query.Order(orderBy)

	// Apply limit
	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	}

	if err := query.Find(&balances).Error; err != nil {
		return nil, fmt.Errorf("failed to list balances: %w", err)
	}

	return balances, nil
}

// GetBalance retrieves a balance by ID
func (r *PostgreSQLRepository) GetBalance(ctx context.Context, balanceID string) (*models.Balance, error) {
	var balance models.Balance
	db := r.getDB()
	if err := db.WithContext(ctx).Where("id = ?", balanceID).First(&balance).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("balance not found: %s", balanceID)
		}
		return nil, fmt.Errorf("failed to get balance: %w", err)
	}
	return &balance, nil
}

// UpdateBalance updates an existing balance
func (r *PostgreSQLRepository) UpdateBalance(ctx context.Context, balance models.Balance) (*models.Balance, error) {
	db := r.getDB()
	if err := db.WithContext(ctx).Save(&balance).Error; err != nil {
		return nil, fmt.Errorf("failed to update balance: %w", err)
	}
	return &balance, nil
}

// DeleteBalance deletes a balance by ID
func (r *PostgreSQLRepository) DeleteBalance(ctx context.Context, balanceID string) error {
	db := r.getDB()
	if err := db.WithContext(ctx).Where("id = ?", balanceID).Delete(&models.Balance{}).Error; err != nil {
		return fmt.Errorf("failed to delete balance: %w", err)
	}
	return nil
}

// DeleteBalancesByUserId deletes all balances for a user ID
func (r *PostgreSQLRepository) DeleteBalancesByUserId(ctx context.Context, userId string) error {
	db := r.getDB()
	if err := db.WithContext(ctx).Where("user_id = ?", userId).Delete(&models.Balance{}).Error; err != nil {
		return fmt.Errorf("failed to delete balances for user %s: %w", userId, err)
	}
	return nil
}
