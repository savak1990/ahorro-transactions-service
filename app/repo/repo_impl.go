package repo

import (
	"context"
	"fmt"

	"github.com/savak1990/transactions-service/app/aws"
	"github.com/savak1990/transactions-service/app/config"
	"github.com/savak1990/transactions-service/app/models"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// PostgreSQLRepository implements Repository interface using PostgreSQL
type PostgreSQLRepository struct {
	db     *gorm.DB
	config *config.AppConfig // Store config for lazy initialization
}

// NewPostgreSQLRepository creates a new PostgreSQL repository
func NewPostgreSQLRepository(db *gorm.DB) *PostgreSQLRepository {
	return &PostgreSQLRepository{
		db: db,
	}
}

// NewPostgreSQLRepositoryWithConfig creates a new PostgreSQL repository with lazy DB initialization
func NewPostgreSQLRepositoryWithConfig(cfg config.AppConfig) *PostgreSQLRepository {
	return &PostgreSQLRepository{
		config: &cfg,
	}
}

// getDB returns the database connection, initializing it if necessary
func (r *PostgreSQLRepository) getDB() *gorm.DB {
	if r.db != nil {
		return r.db
	}

	if r.config != nil {
		// Lazy initialization - this will trigger connection and panic if failed
		// We want this panic to bubble up to the maintenance middleware
		r.db = aws.GetGormDB(*r.config)
		return r.db
	}

	panic(fmt.Errorf("PostgreSQL repository not properly initialized"))
}

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
		Where("id = ?", tx.ID).
		First(&createdTx).Error; err != nil {
		return nil, fmt.Errorf("failed to reload created transaction: %w", err)
	}

	return &createdTx, nil
}

// GetTransaction retrieves a single transaction by ID
func (r *PostgreSQLRepository) GetTransaction(ctx context.Context, transactionID string) (*models.Transaction, error) {

	db := r.getDB()

	var tx models.Transaction
	if err := db.WithContext(ctx).
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
		Preload("Transaction.Merchant").
		Preload("Transaction.Balance").
		Preload("Category")

	// Join transaction table once if we need to filter by transaction fields
	needsTransactionJoin := filter.GroupID != "" || filter.UserID != "" || filter.BalanceID != "" || filter.Type != "" || !filter.StartTime.IsZero() || !filter.EndTime.IsZero()
	if needsTransactionJoin {
		query = query.Joins("JOIN transaction ON transaction_entry.transaction_id = transaction.id")
	}

	// Join category table if we need to filter by category or category group
	if filter.CategoryId != "" || filter.CategoryGroupId != "" {
		query = query.Joins("JOIN category ON transaction_entry.category_id = category.id")
	}

	// Join merchant table if we need to filter by merchant
	if filter.MerchantId != "" {
		query = query.Joins("JOIN merchant ON transaction.merchant_id = merchant.id")
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

	// Apply date range filters
	if !filter.StartTime.IsZero() {
		query = query.Where("transaction.transacted_at >= ?", filter.StartTime)
	}

	if !filter.EndTime.IsZero() {
		query = query.Where("transaction.transacted_at <= ?", filter.EndTime)
	}

	// Apply category filter
	if filter.CategoryId != "" {
		query = query.Where("transaction_entry.category_id = ?", filter.CategoryId)
	}

	// Apply category group filter
	if filter.CategoryGroupId != "" {
		query = query.Where("category.\"group\" = ?", filter.CategoryGroupId)
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
	if filter.Limit > 0 && filter.Limit <= 100 {
		limit = filter.Limit
	}
	query = query.Limit(limit)

	if err := query.Find(&entries).Error; err != nil {
		return nil, fmt.Errorf("failed to list transaction entries: %w", err)
	}

	return entries, nil
}

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
	orderBy := "category_name"
	if input.SortBy != "" {
		switch input.SortBy {
		case "name", "category_name", "rank", "created_at", "updated_at":
			if input.SortBy == "name" {
				orderBy = "category_name" // Map 'name' to the actual column name
			} else {
				orderBy = input.SortBy
			}
		default:
			orderBy = "category_name" // fallback to default if invalid sort field
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

	if err := query.Find(&categories).Error; err != nil {
		return nil, fmt.Errorf("failed to list categories: %w", err)
	}

	return categories, nil
}

// GetCategory retrieves a category by ID
func (r *PostgreSQLRepository) GetCategory(ctx context.Context, categoryID string) (*models.Category, error) {
	var category models.Category
	db := r.getDB()
	if err := db.WithContext(ctx).Where("id = ?", categoryID).First(&category).Error; err != nil {
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
	return &category, nil
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
	query := db.WithContext(ctx)

	// Apply ordering
	orderBy := "rank"
	if filter.SortBy != "" {
		orderBy = filter.SortBy
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

	if err := db.WithContext(ctx).Where("id = ?", categoryGroupID).First(&categoryGroup).Error; err != nil {
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

// DeleteCategoryGroup deletes a category group by ID
func (r *PostgreSQLRepository) DeleteCategoryGroup(ctx context.Context, categoryGroupID string) error {
	db := r.getDB()
	if err := db.WithContext(ctx).Where("id = ?", categoryGroupID).Delete(&models.CategoryGroup{}).Error; err != nil {
		return fmt.Errorf("failed to delete category group: %w", err)
	}
	return nil
}

// GetTransactionStats retrieves aggregated transaction statistics
func (r *PostgreSQLRepository) GetTransactionStats(ctx context.Context, filter models.TransactionStatsInput) (map[string]map[string]models.CurrencyStatsDto, error) {
	var results []models.TransactionStatsRaw
	db := r.getDB()

	// Build the aggregation query
	query := db.WithContext(ctx).
		Table("transaction_entry te").
		Select(`
			t.type as transaction_type,
			b.currency as currency,
			COALESCE(SUM(te.amount), 0) as total_amount,
			COUNT(DISTINCT t.id) as transactions_count,
			COUNT(te.id) as transaction_entries_count
		`).
		Joins("JOIN transaction t ON te.transaction_id = t.id").
		Joins("JOIN balance b ON t.balance_id = b.id").
		Group("t.type, b.currency")

	// Apply filters
	if filter.GroupID != nil && *filter.GroupID != "" {
		query = query.Where("t.group_id = ?", *filter.GroupID)
	}

	if filter.UserID != nil && *filter.UserID != "" {
		query = query.Where("t.user_id = ?", *filter.UserID)
	}

	if filter.BalanceID != nil && *filter.BalanceID != "" {
		query = query.Where("t.balance_id = ?", *filter.BalanceID)
	}

	if filter.Type != nil && *filter.Type != "" {
		query = query.Where("t.type = ?", *filter.Type)
	}

	if filter.CategoryId != nil && *filter.CategoryId != "" {
		query = query.Where("te.category_id = ?", *filter.CategoryId)
	}

	if filter.CategoryGroupId != nil && *filter.CategoryGroupId != "" {
		query = query.Joins("JOIN category c ON te.category_id = c.id").
			Where("c.\"group\" = ?", *filter.CategoryGroupId)
	}

	if filter.MerchantId != nil && *filter.MerchantId != "" {
		query = query.Where("t.merchant_id = ?", *filter.MerchantId)
	}

	if filter.StartTime != nil {
		query = query.Where("t.transacted_at >= ?", *filter.StartTime)
	}

	if filter.EndTime != nil {
		query = query.Where("t.transacted_at <= ?", *filter.EndTime)
	}

	if err := query.Find(&results).Error; err != nil {
		return nil, fmt.Errorf("failed to get transaction stats: %w", err)
	}

	// Convert results to the expected format
	statsMap := make(map[string]map[string]models.CurrencyStatsDto)

	for _, result := range results {
		if statsMap[result.TransactionType] == nil {
			statsMap[result.TransactionType] = make(map[string]models.CurrencyStatsDto)
		}

		// Convert decimal amount to cents (multiply by 100)
		amountInCents := int(result.TotalAmount.Mul(decimal.NewFromInt(100)).IntPart())

		statsMap[result.TransactionType][result.Currency] = models.CurrencyStatsDto{
			Amount:                  amountInCents,
			TransactionsCount:       int(result.TransactionsCount),
			TransactionEntriesCount: int(result.TransactionEntriesCount),
		}
	}

	return statsMap, nil
}

// Ensure PostgreSQLRepository implements Repository interface
var _ Repository = (*PostgreSQLRepository)(nil)
