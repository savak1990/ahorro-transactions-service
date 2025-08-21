package repo

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/savak1990/transactions-service/app/models"
)

// TransactionStatsRawGrouped represents raw statistics data from database aggregation for grouping
type TransactionStatsRawGrouped struct {
	GroupKey                string  `gorm:"column:group_key"`
	GroupLabel              string  `gorm:"column:group_label"`
	TotalAmount             int64   `gorm:"column:total_amount"` // Amount in cents
	TransactionsCount       int64   `gorm:"column:transactions_count"`
	TransactionEntriesCount int64   `gorm:"column:transaction_entries_count"`
	Icon                    *string `gorm:"column:icon"`
}

// GetTransactionStats retrieves aggregated transaction statistics based on grouping
func (r *PostgreSQLRepository) GetTransactionStats(ctx context.Context, filter models.TransactionStatsInput) ([]models.TransactionStatsItemDto, error) {
	var results []TransactionStatsRawGrouped
	db := r.getDB()

	// Set default display currency if not provided
	displayCurrency := filter.DisplayCurrency
	if displayCurrency == "" {
		displayCurrency = "EUR" // Default to EUR
	}

	// Build the base query with transaction_entry_amount join for display currency
	query := db.WithContext(ctx).Table("transaction_entry te").
		Joins("JOIN transaction t ON te.transaction_id = t.id").
		Joins("JOIN balance b ON t.balance_id = b.id").
		Joins("LEFT JOIN transaction_entry_amount tea ON te.id = tea.transaction_entry_id AND tea.currency = ?", displayCurrency).
		Where("te.deleted_at IS NULL AND t.deleted_at IS NULL AND b.deleted_at IS NULL")

	// Apply filters
	if filter.GroupID != "" {
		query = query.Where("t.group_id = ?", filter.GroupID)
	}

	if filter.UserID != "" {
		query = query.Where("t.user_id = ?", filter.UserID)
	}

	// Filter by multiple balance IDs (OR operation)
	if len(filter.BalanceID) > 0 {
		query = query.Where("t.balance_id IN ?", filter.BalanceID)
	}

	// Filter by single transaction type
	if filter.Type != "" {
		query = query.Where("t.type = ?", filter.Type)
	}

	// Filter by multiple transaction IDs (OR operation)
	if len(filter.TransactionId) > 0 {
		query = query.Where("t.id IN ?", filter.TransactionId)
	}

	// Filter by multiple merchant IDs (OR operation)
	if len(filter.MerchantId) > 0 {
		query = query.Joins("JOIN merchant m ON t.merchant_id = m.id AND m.deleted_at IS NULL").
			Where("t.merchant_id IN ?", filter.MerchantId)
	}

	// Filter by multiple category IDs OR category group IDs (OR operation)
	var categoryConditions []string
	var categoryArgs []interface{}

	if len(filter.CategoryId) > 0 {
		categoryConditions = append(categoryConditions, "te.category_id IN ?")
		categoryArgs = append(categoryArgs, filter.CategoryId)
	}

	if len(filter.CategoryGroupId) > 0 {
		categoryConditions = append(categoryConditions, "c_filter.category_group_id IN ?")
		categoryArgs = append(categoryArgs, filter.CategoryGroupId)
	}

	if len(categoryConditions) > 0 {
		query = query.Joins("JOIN category c_filter ON te.category_id = c_filter.id AND c_filter.deleted_at IS NULL")

		if len(categoryConditions) == 1 {
			query = query.Where(categoryConditions[0], categoryArgs...)
		} else {
			orCondition := fmt.Sprintf("(%s)", categoryConditions[0])
			for i := 1; i < len(categoryConditions); i++ {
				orCondition += fmt.Sprintf(" OR (%s)", categoryConditions[i])
			}
			query = query.Where(orCondition, categoryArgs...)
		}
	}

	// Apply time filters
	if !filter.StartTime.IsZero() {
		query = query.Where("t.transacted_at >= ?", filter.StartTime)
	}

	if filter.EndTime != nil {
		query = query.Where("t.transacted_at <= ?", *filter.EndTime)
	}

	// Add grouping-specific SELECT and GROUP BY clauses
	selectClause, groupByClause, joinClause := r.buildGroupingQuery(filter.Grouping, displayCurrency)

	// Add additional joins if needed
	if joinClause != "" {
		query = query.Joins(joinClause)
	}

	query = query.Select(selectClause).Group(groupByClause)

	if err := query.Find(&results).Error; err != nil {
		return nil, fmt.Errorf("failed to get transaction stats: %w", err)
	}

	// Convert results to TransactionStatsItemDto
	statsItems := make([]models.TransactionStatsItemDto, 0, len(results))

	for _, result := range results {
		statsItems = append(statsItems, models.TransactionStatsItemDto{
			Label:    result.GroupLabel,
			Amount:   int(result.TotalAmount),
			Currency: displayCurrency, // Use display currency instead of original currency from DB
			Count:    int(result.TransactionsCount),
			Icon:     result.Icon,
		})
	}

	// Note: Sorting and limiting will be done in service layer
	return statsItems, nil
}

// buildGroupingQuery returns the SELECT clause, GROUP BY clause, and additional JOIN clause for the specified grouping
func (r *PostgreSQLRepository) buildGroupingQuery(grouping string, displayCurrency string) (string, string, string) {
	baseSelect := `
		COALESCE(SUM(COALESCE(tea.amount, te.amount)), 0) as total_amount,
		COUNT(DISTINCT t.id) as transactions_count,
		COUNT(te.id) as transaction_entries_count`

	switch grouping {
	case models.GroupingCategory:
		return baseSelect + `,
			COALESCE(c.id::text, 'Unknown') as group_key,
			COALESCE(c.name, 'Unknown') as group_label,
			c.image_url as icon`,
			"c.id, c.name, c.image_url",
			"LEFT JOIN category c ON te.category_id = c.id AND c.deleted_at IS NULL"

	case models.GroupingCategoryGroup:
		return baseSelect + `,
			COALESCE(cg.id::text, 'Unknown') as group_key,
			COALESCE(cg.name, 'Unknown') as group_label,
			cg.image_url as icon`,
			"cg.id, cg.name, cg.image_url",
			`LEFT JOIN category c ON te.category_id = c.id AND c.deleted_at IS NULL
			 LEFT JOIN category_group cg ON c.category_group_id = cg.id AND cg.deleted_at IS NULL`

	case models.GroupingMerchant:
		return baseSelect + `,
			COALESCE(m.id::text, 'Unknown') as group_key,
			COALESCE(m.name, 'Unknown') as group_label,
			m.image_url as icon`,
			"m.id, m.name, m.image_url",
			"LEFT JOIN merchant m ON t.merchant_id = m.id AND m.deleted_at IS NULL"

	case models.GroupingBalance:
		return baseSelect + `,
			b.id::text as group_key,
			b.title as group_label,
			NULL as icon`,
			"b.id, b.title",
			""

	case models.GroupingCurrency:
		return baseSelect + `,
			'` + displayCurrency + `' as group_key,
			'` + displayCurrency + `' as group_label,
			NULL as icon`,
			"",
			""

	case models.GroupingMonth:
		return baseSelect + `,
			TO_CHAR(t.transacted_at, 'YYYY-MM') as group_key,
			TRIM(TO_CHAR(t.transacted_at, 'Month')) || ' ' || TO_CHAR(t.transacted_at, 'YYYY') as group_label,
			NULL as icon`,
			"TO_CHAR(t.transacted_at, 'YYYY-MM'), TRIM(TO_CHAR(t.transacted_at, 'Month')) || ' ' || TO_CHAR(t.transacted_at, 'YYYY')",
			""

	case models.GroupingQuarter:
		return baseSelect + `,
			EXTRACT(YEAR FROM t.transacted_at)::text || '-Q' || EXTRACT(QUARTER FROM t.transacted_at)::text as group_key,
			'Year ' || EXTRACT(YEAR FROM t.transacted_at)::text || ' Q' || EXTRACT(QUARTER FROM t.transacted_at)::text as group_label,
			NULL as icon`,
			"EXTRACT(YEAR FROM t.transacted_at), EXTRACT(QUARTER FROM t.transacted_at)",
			""

	case models.GroupingYear:
		return baseSelect + `,
			EXTRACT(YEAR FROM t.transacted_at)::text as group_key,
			'Year ' || EXTRACT(YEAR FROM t.transacted_at)::text as group_label,
			NULL as icon`,
			"EXTRACT(YEAR FROM t.transacted_at)",
			""

	case models.GroupingWeek:
		return baseSelect + `,
			EXTRACT(YEAR FROM t.transacted_at)::text || '-W' || EXTRACT(WEEK FROM t.transacted_at)::text as group_key,
			'Week ' || EXTRACT(WEEK FROM t.transacted_at)::text as group_label,
			NULL as icon`,
			"EXTRACT(YEAR FROM t.transacted_at), EXTRACT(WEEK FROM t.transacted_at)",
			""

	case models.GroupingDay:
		return baseSelect + `,
			TO_CHAR(t.transacted_at, 'YYYY-MM-DD') as group_key,
			TO_CHAR(t.transacted_at, 'DD Mon YYYY') as group_label,
			NULL as icon`,
			"TO_CHAR(t.transacted_at, 'YYYY-MM-DD'), TO_CHAR(t.transacted_at, 'DD Mon YYYY')",
			""

	default:
		// Default to currency grouping
		return baseSelect + `,
			'` + displayCurrency + `' as group_key,
			'` + displayCurrency + `' as group_label,
			NULL as icon`,
			"",
			""
	}
}

// sortTransactionStatsItems sorts the transaction stats items based on the provided sort field and order
func (r *PostgreSQLRepository) sortTransactionStatsItems(items []models.TransactionStatsItemDto, sortBy, order string) {
	sort.Slice(items, func(i, j int) bool {
		var less bool
		switch sortBy {
		case "count":
			less = items[i].Count < items[j].Count
		case "label":
			less = strings.ToLower(items[i].Label) < strings.ToLower(items[j].Label)
		default: // "amount"
			less = items[i].Amount < items[j].Amount
		}

		if order == "asc" {
			return less
		}
		return !less
	})
}
