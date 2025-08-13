package repo

import (
	"context"
	"fmt"

	"github.com/savak1990/transactions-service/app/models"
)

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
		Where("te.deleted_at IS NULL").
		Group("t.type, b.currency")

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

	// Filter by multiple transaction types (OR operation)
	if len(filter.Type) > 0 {
		query = query.Where("t.type IN ?", filter.Type)
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
	// This handles filtering by category or category group independently
	var categoryConditions []string
	var categoryArgs []interface{}

	if len(filter.CategoryId) > 0 {
		categoryConditions = append(categoryConditions, "te.category_id IN ?")
		categoryArgs = append(categoryArgs, filter.CategoryId)
	}

	if len(filter.CategoryGroupId) > 0 {
		categoryConditions = append(categoryConditions, "c.category_group_id IN ?")
		categoryArgs = append(categoryArgs, filter.CategoryGroupId)
	}

	if len(categoryConditions) > 0 {
		// Join with category table if we have category-related filters
		query = query.Joins("JOIN category c ON te.category_id = c.id AND c.deleted_at IS NULL")

		// Apply OR condition between category filters
		if len(categoryConditions) == 1 {
			query = query.Where(categoryConditions[0], categoryArgs...)
		} else {
			// Multiple conditions with OR
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

	if err := query.Find(&results).Error; err != nil {
		return nil, fmt.Errorf("failed to get transaction stats: %w", err)
	}

	// Convert results to the expected format
	statsMap := make(map[string]map[string]models.CurrencyStatsDto)

	for _, result := range results {
		if statsMap[result.TransactionType] == nil {
			statsMap[result.TransactionType] = make(map[string]models.CurrencyStatsDto)
		}

		// Amount is already in cents (int64)
		amountInCents := int(result.TotalAmount)

		statsMap[result.TransactionType][result.Currency] = models.CurrencyStatsDto{
			Amount:                  amountInCents,
			TransactionsCount:       int(result.TransactionsCount),
			TransactionEntriesCount: int(result.TransactionEntriesCount),
		}
	}

	return statsMap, nil
}
