package repo

import (
	"context"
	"fmt"

	"github.com/savak1990/transactions-service/app/models"
	"gorm.io/gorm"
)

// CreateCategory creates a new category in the database
func (r *PostgreSQLRepository) CreateCategory(ctx context.Context, category models.Category) (*models.Category, error) {
	db := r.getDB()
	// Create the category in the database
	if err := db.WithContext(ctx).Create(&category).Error; err != nil {
		return nil, fmt.Errorf("failed to create category: %w", err)
	}
	return &category, nil
}

// ListCategories retrieves all categories
func (r *PostgreSQLRepository) ListCategories(ctx context.Context, input models.ListCategoriesInput) ([]models.Category, error) {
	var categories []models.Category
	db := r.getDB()
	query := db.WithContext(ctx)

	// Apply ordering
	orderBy := "name"
	if input.SortBy != "" {
		switch input.SortBy {
		case "name", "rank", "created_at", "updated_at":
			orderBy = input.SortBy
		case "createdAt":
			orderBy = "created_at"
		case "updatedAt":
			orderBy = "updated_at"
		default:
			orderBy = "name" // fallback to default if invalid sort field
		}
	}
	order := "ASC"
	if input.Order != "" && (input.Order == "DESC" || input.Order == "desc") {
		order = "DESC"
	}
	query = query.Order(fmt.Sprintf("%s %s", orderBy, order))

	// Apply limit
	if input.Limit > 0 {
		query = query.Limit(input.Limit)
	}

	// For now, we'll ignore pagination (StartKey) since it's more complex
	// You can implement pagination later if needed

	if err := query.Preload("CategoryGroup").Find(&categories).Error; err != nil {
		return nil, fmt.Errorf("failed to list categories: %w", err)
	}

	return categories, nil
}

// GetCategory retrieves a category by ID
func (r *PostgreSQLRepository) GetCategory(ctx context.Context, categoryID string) (*models.Category, error) {
	var category models.Category
	db := r.getDB()
	if err := db.WithContext(ctx).Preload("CategoryGroup").Where("id = ?", categoryID).First(&category).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("category not found: %s", categoryID)
		}
		return nil, fmt.Errorf("failed to get category: %w", err)
	}
	return &category, nil
}

// UpdateCategory updates an existing category
func (r *PostgreSQLRepository) UpdateCategory(ctx context.Context, category models.Category) (*models.Category, error) {
	db := r.getDB()
	if err := db.WithContext(ctx).Save(&category).Error; err != nil {
		return nil, fmt.Errorf("failed to update category: %w", err)
	}

	// Fetch the updated category with preloaded relationships for response
	var updatedCategory models.Category
	if err := db.WithContext(ctx).Preload("CategoryGroup").Where("id = ?", category.ID).First(&updatedCategory).Error; err != nil {
		return nil, fmt.Errorf("failed to get updated category: %w", err)
	}

	return &updatedCategory, nil
}

// DeleteCategory deletes a category by ID
func (r *PostgreSQLRepository) DeleteCategory(ctx context.Context, categoryID string) error {
	db := r.getDB()
	if err := db.WithContext(ctx).Where("id = ?", categoryID).Delete(&models.Category{}).Error; err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}
	return nil
}

// DeleteCategoriesByUserId deletes all categories for a user ID
func (r *PostgreSQLRepository) DeleteCategoriesByUserId(ctx context.Context, userId string) error {
	db := r.getDB()
	if err := db.WithContext(ctx).Where("user_id = ?", userId).Delete(&models.Category{}).Error; err != nil {
		return fmt.Errorf("failed to delete categories for user %s: %w", userId, err)
	}
	return nil
}
