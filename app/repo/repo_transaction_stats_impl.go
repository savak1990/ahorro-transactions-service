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
		Where("b.deleted_at IS NULL").
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
		query = query.Joins("JOIN merchant m ON t.merchant_id = m.id AND m.deleted_at IS NULL").
			Where("t.merchant_id = ?", *filter.MerchantId)
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
