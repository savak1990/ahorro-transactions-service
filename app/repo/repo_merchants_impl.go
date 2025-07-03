package repo

import (
	"context"
	"fmt"

	"github.com/savak1990/transactions-service/app/models"
	"gorm.io/gorm"
)

// CreateMerchant creates a new merchant in the database
func (r *PostgreSQLRepository) CreateMerchant(ctx context.Context, merchant models.Merchant) (*models.Merchant, error) {
	db := r.getDB()
	if err := db.WithContext(ctx).Create(&merchant).Error; err != nil {
		return nil, fmt.Errorf("failed to create merchant: %w", err)
	}
	return &merchant, nil
}

// ListMerchants retrieves merchants based on the filter
func (r *PostgreSQLRepository) ListMerchants(ctx context.Context, filter models.ListMerchantsInput) ([]models.Merchant, error) {
	var merchants []models.Merchant
	db := r.getDB()
	query := db.WithContext(ctx)

	// Apply filters
	if filter.Name != "" {
		query = query.Where("name ILIKE ?", "%"+filter.Name+"%")
	}

	// Apply ordering
	orderBy := "created_at"
	if filter.SortBy != "" {
		switch filter.SortBy {
		case "rank", "name", "created_at", "updated_at":
			orderBy = filter.SortBy
		default:
			orderBy = "created_at" // fallback to default if invalid sort field
		}
	}
	order := "ASC"
	if filter.Order != "" && (filter.Order == "DESC" || filter.Order == "desc") {
		order = "DESC"
	}
	query = query.Order(fmt.Sprintf("%s %s", orderBy, order))

	// Apply limit
	limit := 50 // default limit
	if filter.Limit > 0 && filter.Limit <= 100 {
		limit = filter.Limit
	}
	query = query.Limit(limit)

	if err := query.Find(&merchants).Error; err != nil {
		return nil, fmt.Errorf("failed to list merchants: %w", err)
	}

	return merchants, nil
}

// GetMerchant retrieves a merchant by ID
func (r *PostgreSQLRepository) GetMerchant(ctx context.Context, merchantId string) (*models.Merchant, error) {
	var merchant models.Merchant
	db := r.getDB()
	if err := db.WithContext(ctx).Where("id = ?", merchantId).First(&merchant).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("merchant not found: %s", merchantId)
		}
		return nil, fmt.Errorf("failed to get merchant: %w", err)
	}
	return &merchant, nil
}

// UpdateMerchant updates an existing merchant
func (r *PostgreSQLRepository) UpdateMerchant(ctx context.Context, merchant models.Merchant) (*models.Merchant, error) {
	db := r.getDB()
	if err := db.WithContext(ctx).Save(&merchant).Error; err != nil {
		return nil, fmt.Errorf("failed to update merchant: %w", err)
	}
	return &merchant, nil
}

// DeleteMerchant soft deletes a merchant by ID
func (r *PostgreSQLRepository) DeleteMerchant(ctx context.Context, merchantId string) error {
	db := r.getDB()
	if err := db.WithContext(ctx).Where("id = ?", merchantId).Delete(&models.Merchant{}).Error; err != nil {
		return fmt.Errorf("failed to delete merchant: %w", err)
	}
	return nil
}

// DeleteMerchantsByUserId deletes all merchants for a user ID
func (r *PostgreSQLRepository) DeleteMerchantsByUserId(ctx context.Context, userId string) error {
	db := r.getDB()
	if err := db.WithContext(ctx).Where("user_id = ?", userId).Delete(&models.Merchant{}).Error; err != nil {
		return fmt.Errorf("failed to delete merchants for user %s: %w", userId, err)
	}
	return nil
}

// GetMerchantByNameAndUserId checks if a merchant with the given name already exists for a user
func (r *PostgreSQLRepository) GetMerchantByNameAndUserId(ctx context.Context, name string, userId string) (*models.Merchant, error) {
	var merchant models.Merchant
	db := r.getDB()
	err := db.WithContext(ctx).Where("name = ? AND user_id = ? AND deleted_at IS NULL", name, userId).First(&merchant).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // No existing merchant found - this is expected when name is available
		}
		return nil, fmt.Errorf("failed to check for existing merchant: %w", err)
	}
	return &merchant, nil
}
