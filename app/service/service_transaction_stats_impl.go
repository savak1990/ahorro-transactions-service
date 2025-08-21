package service

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/savak1990/transactions-service/app/models"
	"github.com/sirupsen/logrus"
)

// GetTransactionStats retrieves aggregated transaction statistics
func (s *ServiceImpl) GetTransactionStats(ctx context.Context, filter models.TransactionStatsInput) ([]models.TransactionStatsItemDto, error) {
	// Set default display currency if not provided
	displayCurrency := filter.DisplayCurrency
	if displayCurrency == "" {
		displayCurrency = "EUR" // Default to EUR
	}

	// Get transaction stats from repository (currency conversion handled in DB layer)
	statsList, err := s.repo.GetTransactionStats(ctx, filter)
	if err != nil {
		logrus.Errorf("Error getting transaction stats: %v", err)
		return nil, fmt.Errorf("failed to get transaction stats: %w", err)
	}

	// Sort the results
	s.sortTransactionStatsItems(statsList, filter.Sort, filter.Order)

	// Apply limit with "Other" category aggregation
	finalStats := s.applyLimitWithOther(statsList, filter.Limit, displayCurrency)

	return finalStats, nil
}

// sortTransactionStatsItems sorts the transaction stats items based on the provided sort field and order
func (s *ServiceImpl) sortTransactionStatsItems(items []models.TransactionStatsItemDto, sortBy, order string) {
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

// applyLimitWithOther applies limit and creates an "Other" category for remaining items
func (s *ServiceImpl) applyLimitWithOther(items []models.TransactionStatsItemDto, limit int, displayCurrency string) []models.TransactionStatsItemDto {
	// If no limit or limit is greater than items, return all items
	if limit <= 0 || len(items) <= limit {
		return items
	}

	// If limit is 1, return only "Other" with all items combined
	if limit == 1 {
		otherAmount := 0
		otherCount := 0
		for _, item := range items {
			otherAmount += item.Amount
			otherCount += item.Count
		}
		return []models.TransactionStatsItemDto{
			{
				Label:    "Other",
				Amount:   otherAmount,
				Currency: displayCurrency,
				Count:    otherCount,
				Icon:     nil,
			},
		}
	}

	// Take top (limit-1) items and aggregate the rest into "Other"
	topItems := items[:limit-1]
	remainingItems := items[limit-1:]

	// Calculate aggregated values for "Other"
	otherAmount := 0
	otherCount := 0
	for _, item := range remainingItems {
		otherAmount += item.Amount
		otherCount += item.Count
	}

	// Create "Other" item
	otherItem := models.TransactionStatsItemDto{
		Label:    "Other",
		Amount:   otherAmount,
		Currency: displayCurrency,
		Count:    otherCount,
		Icon:     nil,
	}

	// Combine top items with "Other"
	result := make([]models.TransactionStatsItemDto, 0, limit)
	result = append(result, topItems...)
	result = append(result, otherItem)

	return result
}
