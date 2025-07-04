package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/savak1990/transactions-service/app/models"
	"gorm.io/gorm"
)

// CategoryGroup repository methods

// CreateCategoryGroup creates a new category group in the database
func (r *PostgreSQLRepository) CreateCategoryGroup(ctx context.Context, categoryGroup models.CategoryGroup) (*models.CategoryGroup, error) {
	db := r.getDB()
	if err := db.WithContext(ctx).Create(&categoryGroup).Error; err != nil {
		return nil, fmt.Errorf("failed to create category group: %w", err)
	}
	return &categoryGroup, nil
}

// ListCategoryGroups retrieves category groups based on the filter
func (r *PostgreSQLRepository) ListCategoryGroups(ctx context.Context, filter models.ListCategoryGroupsInput) ([]models.CategoryGroup, error) {
	var categoryGroups []models.CategoryGroup
	db := r.getDB()
	query := db.WithContext(ctx).Where("deleted_at IS NULL")

	// Apply ordering
	orderBy := "rank"
	if filter.SortBy != "" {
		switch filter.SortBy {
		case "rank", "name", "created_at", "updated_at":
			orderBy = filter.SortBy
		case "createdAt":
			orderBy = "created_at"
		case "updatedAt":
			orderBy = "updated_at"
		default:
			orderBy = "rank" // fallback to default if invalid sort field
		}
	}
	order := "DESC"
	if filter.Order != "" && (filter.Order == "ASC" || filter.Order == "asc") {
		order = "ASC"
	}
	query = query.Order(fmt.Sprintf("%s %s", orderBy, order))

	// Apply limit
	if filter.Limit > 0 && filter.Limit <= 100 {
		query = query.Limit(filter.Limit)
	}

	if err := query.Find(&categoryGroups).Error; err != nil {
		return nil, fmt.Errorf("failed to list category groups: %w", err)
	}

	return categoryGroups, nil
}

// GetCategoryGroup retrieves a category group by ID
func (r *PostgreSQLRepository) GetCategoryGroup(ctx context.Context, categoryGroupID string) (*models.CategoryGroup, error) {
	var categoryGroup models.CategoryGroup
	db := r.getDB()

	if err := db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", categoryGroupID).First(&categoryGroup).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get category group: %w", err)
	}

	return &categoryGroup, nil
}

// UpdateCategoryGroup updates an existing category group
func (r *PostgreSQLRepository) UpdateCategoryGroup(ctx context.Context, categoryGroup models.CategoryGroup) (*models.CategoryGroup, error) {
	db := r.getDB()
	if err := db.WithContext(ctx).Save(&categoryGroup).Error; err != nil {
		return nil, fmt.Errorf("failed to update category group: %w", err)
	}
	return &categoryGroup, nil
}

// DeleteCategoryGroup deletes a category group by ID using soft deletion
func (r *PostgreSQLRepository) DeleteCategoryGroup(ctx context.Context, categoryGroupID string) error {
	db := r.getDB()

	// Simply set deleted_at to current time - no foreign key constraints involved
	now := time.Now()
	result := db.WithContext(ctx).Model(&models.CategoryGroup{}).
		Where("id = ? AND deleted_at IS NULL", categoryGroupID).
		Update("deleted_at", now)

	if result.Error != nil {
		return fmt.Errorf("failed to soft delete category group: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("category group not found or already deleted: %s", categoryGroupID)
	}

	return nil
}
