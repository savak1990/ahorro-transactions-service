package service

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"

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
	displayCurrency = strings.ToUpper(displayCurrency)

	var wg sync.WaitGroup
	var statsList []models.TransactionStatsItemDto
	var exchangeRates map[string]float64
	var statsErr, exchangeRatesErr error

	// Launch two goroutines in parallel
	wg.Add(2)

	// Goroutine 1: Get transaction stats from repository
	go func() {
		defer wg.Done()
		statsList, statsErr = s.repo.GetTransactionStats(ctx, filter)
	}()

	// Goroutine 2: Get exchange rates for display currency
	go func() {
		defer wg.Done()
		exchangeRates, exchangeRatesErr = s.exchangeRatesClient.GetExchangeRates(ctx, displayCurrency)
	}()

	// Wait for both goroutines to complete
	wg.Wait()

	// Check for errors from either operation
	if statsErr != nil {
		logrus.Errorf("Error getting transaction stats: %v", statsErr)
		return nil, fmt.Errorf("failed to get transaction stats: %w", statsErr)
	}
	if exchangeRatesErr != nil {
		logrus.Errorf("Error getting exchange rates: %v", exchangeRatesErr)
		return nil, fmt.Errorf("failed to get exchange rates: %v", exchangeRatesErr)
	}

	// Convert currencies and merge items with same labels
	mergedStats, _ := s.convertAndMergeStats(statsList, displayCurrency, exchangeRates)

	// Sort the results after merging
	s.sortTransactionStatsItems(mergedStats, filter.Sort, filter.Order)

	// Apply limit with "Other" category aggregation
	finalStats := s.applyLimitWithOther(mergedStats, filter.Limit, displayCurrency)

	return finalStats, nil
}

// convertAndMergeStats converts all amounts to the display currency and merges items with the same label using pre-fetched exchange rates
func (s *ServiceImpl) convertAndMergeStats(stats []models.TransactionStatsItemDto, displayCurrency string, exchangeRates map[string]float64) ([]models.TransactionStatsItemDto, []string) {

	logrus.WithFields(logrus.Fields{
		"stats":           stats,
		"displayCurrency": displayCurrency,
		"exchangeRates":   exchangeRates,
	}).Debug("convertAndMergeStats input variables")

	// Group by label and merge amounts after currency conversion
	labelGroups := make(map[string]*models.TransactionStatsItemDto)

	// Currencies used
	currenciesUsed := make([]string, 0, 5)
	currenciesUsed = append(currenciesUsed, displayCurrency)

	for _, stat := range stats {
		convertedAmount := stat.Amount
		if displayCurrency == "" || stat.Currency != displayCurrency {
			currenciesUsed = append(currenciesUsed, stat.Currency)
			convertedAmount = s.convertCurrencyWithRates(stat.Amount, stat.Currency, exchangeRates)
		}

		if existing, exists := labelGroups[stat.Label]; exists {
			// Merge with existing item
			existing.Amount += convertedAmount
			existing.Count += stat.Count
		} else {
			// Create new item with converted currency
			labelGroups[stat.Label] = &models.TransactionStatsItemDto{
				Label:    stat.Label,
				Amount:   convertedAmount,
				Currency: displayCurrency, // Always use display currency
				Count:    stat.Count,
				Icon:     stat.Icon,
			}
		}
	}

	// Convert map back to slice
	result := make([]models.TransactionStatsItemDto, 0, len(labelGroups))
	for _, item := range labelGroups {
		result = append(result, *item)
	}

	return result, currenciesUsed
}

// convertCurrencyWithRates converts an amount from source currency to target currency using pre-fetched rates
func (s *ServiceImpl) convertCurrencyWithRates(amount int, fromCurrency string, exchangeRates map[string]float64) int {
	// Get the rate for target currency from pre-fetched rates
	rate, exists := exchangeRates[fromCurrency]
	if !exists {
		logrus.Warn("Exchange rate not found")
		// If target currency not found in rates, return original amount
		return amount
	}

	convertedAmount := float64(amount) / rate
	logrus.
		WithField("originalAmount", amount).
		WithField("convertedAmount", convertedAmount).
		WithField("currency", fromCurrency).
		Debug("Currency conversion")
	return int(convertedAmount + 0.5) // Round to nearest cent
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
