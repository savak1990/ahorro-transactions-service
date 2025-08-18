package service

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/savak1990/transactions-service/app/models"
	repo "github.com/savak1990/transactions-service/app/repo"
)

type ServiceImpl struct {
	repo repo.Repository
}

func NewServiceImpl(repo repo.Repository) *ServiceImpl {
	return &ServiceImpl{
		repo: repo,
	}
}

func (s *ServiceImpl) CreateTransaction(ctx context.Context, tx models.Transaction) (*models.Transaction, error) {
	tx.ID = models.NewTransactionID()

	// Validate balance exists (required)
	_, err := s.repo.GetBalance(ctx, tx.BalanceID.String())
	if err != nil {
		return nil, fmt.Errorf("balance with ID %s not found: %w", tx.BalanceID.String(), err)
	}

	// Validate merchant exists if merchantID is provided
	if tx.MerchantID != nil {
		_, err := s.repo.GetMerchant(ctx, tx.MerchantID.String())
		if err != nil {
			return nil, fmt.Errorf("merchant with ID %s not found: %w", tx.MerchantID.String(), err)
		}
	}

	// Validate categories exist for all transaction entries
	for i, entry := range tx.TransactionEntries {
		if entry.CategoryID != nil {
			_, err := s.repo.GetCategory(ctx, entry.CategoryID.String())
			if err != nil {
				return nil, fmt.Errorf("category with ID %s not found for transaction entry %d: %w", entry.CategoryID.String(), i, err)
			}
		}
	}

	return s.repo.CreateTransaction(ctx, tx)
}

func (s *ServiceImpl) GetTransaction(ctx context.Context, transactionID string) (*models.SingleTransactionDto, error) {
	// Get the base transaction with preloaded balance
	tx, err := s.repo.GetTransaction(ctx, transactionID)
	if err != nil {
		return nil, fmt.Errorf("transaction with ID %s not found: %w", transactionID, err)
	}

	// Use the preloaded balance from the transaction
	balance := tx.Balance

	// Create the main DTO
	dto := models.SingleTransactionDto{
		TransactionID:   tx.ID.String(),
		GroupID:         tx.GroupID.String(),
		UserID:          tx.UserID.String(),
		BalanceID:       tx.BalanceID.String(),
		BalanceTitle:    balance.Title,
		BalanceCurrency: balance.Currency,
		BalanceDeleted:  balance.DeletedAt != nil,
		Type:            tx.Type,
		ApprovedAt:      tx.ApprovedAt.Format(time.RFC3339),
		TransactedAt:    tx.TransactedAt.Format(time.RFC3339),
		CreatedAt:       tx.CreatedAt.Format(time.RFC3339),
		UpdatedAt:       tx.UpdatedAt.Format(time.RFC3339),
	}

	// Add merchant details if available
	if tx.MerchantID != nil {
		merchant, err := s.repo.GetMerchant(ctx, tx.MerchantID.String())
		if err == nil && merchant != nil {
			dto.MerchantID = merchant.ID.String()
			dto.MerchantName = merchant.Name
			if merchant.ImageUrl != nil && *merchant.ImageUrl != "" {
				dto.MerchantLogo = *merchant.ImageUrl
			}
		}
	}

	// Add operation ID if available
	if tx.OperationID != nil {
		dto.OperationID = tx.OperationID.String()
	}

	// Convert transaction entries with category details - return error if any category not found
	var entryDtos []models.SingleTransactionEntryDto
	for _, entry := range tx.TransactionEntries {
		var category *models.Category
		if entry.CategoryID != nil {
			category, err = s.repo.GetCategory(ctx, entry.CategoryID.String())
			if err != nil {
				return nil, fmt.Errorf("category with ID %s not found for transaction entry: %w", entry.CategoryID.String(), err)
			}
		}
		entryDto := models.ToAPISingleTransactionEntry(&entry, category)
		entryDtos = append(entryDtos, entryDto)
	}
	dto.TransactionEntries = entryDtos

	return &dto, nil
}

func (s *ServiceImpl) ListTransactions(ctx context.Context, filter models.ListTransactionsInput) ([]models.Transaction, error) {
	return s.repo.ListTransactions(ctx, filter)
}

func (s *ServiceImpl) UpdateTransaction(ctx context.Context, transactionID string, updateDto models.UpdateTransactionDto) (*models.Transaction, error) {
	// First, fetch the existing transaction
	existingTx, err := s.repo.GetTransaction(ctx, transactionID)
	if err != nil {
		return nil, fmt.Errorf("transaction with ID %s not found: %w", transactionID, err)
	}

	// Validate and update BalanceID if provided (cannot be null)
	if updateDto.BalanceID != "" {
		balanceUUID, err := uuid.Parse(updateDto.BalanceID)
		if err != nil {
			return nil, fmt.Errorf("invalid balance ID format: %w", err)
		}

		// Validate balance exists
		_, err = s.repo.GetBalance(ctx, updateDto.BalanceID)
		if err != nil {
			return nil, fmt.Errorf("balance with ID %s not found: %w", updateDto.BalanceID, err)
		}

		existingTx.BalanceID = balanceUUID
	}

	// Validate and update MerchantID (can be null)
	if updateDto.MerchantID == "" {
		// Set merchant to null
		existingTx.MerchantID = nil
	} else {
		merchantUUID, err := uuid.Parse(updateDto.MerchantID)
		if err != nil {
			return nil, fmt.Errorf("invalid merchant ID format: %w", err)
		}

		// Validate merchant exists
		_, err = s.repo.GetMerchant(ctx, updateDto.MerchantID)
		if err != nil {
			return nil, fmt.Errorf("merchant with ID %s not found: %w", updateDto.MerchantID, err)
		}

		existingTx.MerchantID = &merchantUUID
	}

	// Update other fields if provided
	if updateDto.Type != "" {
		existingTx.Type = updateDto.Type
	}

	if updateDto.OperationID != "" {
		operationUUID, err := uuid.Parse(updateDto.OperationID)
		if err != nil {
			return nil, fmt.Errorf("invalid operation ID format: %w", err)
		}
		existingTx.OperationID = &operationUUID
	}

	// Update dates only if provided
	if updateDto.ApprovedAt != "" {
		approvedAt, err := time.Parse(time.RFC3339, updateDto.ApprovedAt)
		if err != nil {
			return nil, fmt.Errorf("invalid approved_at format: %w", err)
		}
		existingTx.ApprovedAt = approvedAt
	}

	if updateDto.TransactedAt != "" {
		transactedAt, err := time.Parse(time.RFC3339, updateDto.TransactedAt)
		if err != nil {
			return nil, fmt.Errorf("invalid transacted_at format: %w", err)
		}
		existingTx.TransactedAt = transactedAt
	}

	// Handle transaction entries
	if len(updateDto.TransactionEntries) > 0 {
		// Convert DTO entries to DAO entries with intelligent update logic
		var newEntries []models.TransactionEntry
		for i, entryDto := range updateDto.TransactionEntries {
			// Validate category if provided
			var categoryID *uuid.UUID
			if entryDto.CategoryID != "" {
				catUUID, err := uuid.Parse(entryDto.CategoryID)
				if err != nil {
					return nil, fmt.Errorf("invalid category ID format for entry %d: %w", i, err)
				}

				// Validate category exists
				_, err = s.repo.GetCategory(ctx, entryDto.CategoryID)
				if err != nil {
					return nil, fmt.Errorf("category with ID %s not found for entry %d: %w", entryDto.CategoryID, i, err)
				}

				categoryID = &catUUID
			}

			// Parse entry ID - if not provided, generate new ID for creation
			var entryID uuid.UUID
			if entryDto.ID != "" {
				// Update existing entry
				var err error
				entryID, err = uuid.Parse(entryDto.ID)
				if err != nil {
					return nil, fmt.Errorf("invalid entry ID format for entry %d: %w", i, err)
				}
			} else {
				// Create new entry
				entryID = models.NewTransactionEntryID()
			}

			var description *string
			if entryDto.Description != "" {
				description = &entryDto.Description
			}

			entry := models.TransactionEntry{
				ID:            entryID,
				TransactionID: existingTx.ID,
				Description:   description,
				Amount:        int64(entryDto.Amount),
				CategoryID:    categoryID,
			}

			newEntries = append(newEntries, entry)
		}

		existingTx.TransactionEntries = newEntries
	} else {
		// If no transaction entries provided, clear the entries slice to avoid
		// the repository thinking they need to be created/updated
		existingTx.TransactionEntries = nil
	}

	// Set updated timestamp
	existingTx.UpdatedAt = time.Now().UTC()

	// Update the transaction
	return s.repo.UpdateTransaction(ctx, *existingTx)
}

func (s *ServiceImpl) DeleteTransaction(ctx context.Context, transactionID string) error {
	return s.repo.DeleteTransaction(ctx, transactionID)
}

func (s *ServiceImpl) CreateBalance(ctx context.Context, balance models.Balance) (*models.Balance, error) {
	balance.ID = models.NewBalanceID()
	return s.repo.CreateBalance(ctx, balance)
}

func (s *ServiceImpl) GetBalance(ctx context.Context, balanceID string) (*models.Balance, error) {
	return s.repo.GetBalance(ctx, balanceID)
}

func (s *ServiceImpl) ListBalances(ctx context.Context, filter models.ListBalancesInput) ([]models.Balance, error) {
	return s.repo.ListBalances(ctx, filter)
}

func (s *ServiceImpl) UpdateBalance(ctx context.Context, balance models.Balance) (*models.Balance, error) {
	balance.UpdatedAt = time.Now().UTC()
	return s.repo.UpdateBalance(ctx, balance)
}

func (s *ServiceImpl) DeleteBalance(ctx context.Context, balanceID string) error {
	return s.repo.DeleteBalance(ctx, balanceID)
}

func (s *ServiceImpl) DeleteBalancesByUserId(ctx context.Context, userId string) error {
	return s.repo.DeleteBalancesByUserId(ctx, userId)
}

func (s *ServiceImpl) CreateCategory(ctx context.Context, category models.Category) (*models.Category, error) {
	category.ID = models.NewCategoryID()
	return s.repo.CreateCategory(ctx, category)
}

func (s *ServiceImpl) ListCategories(ctx context.Context, filter models.ListCategoriesInput) ([]models.Category, error) {
	return s.repo.ListCategories(ctx, filter)
}

func (s *ServiceImpl) GetCategory(ctx context.Context, categoryID string) (*models.Category, error) {
	return s.repo.GetCategory(ctx, categoryID)
}

func (s *ServiceImpl) UpdateCategory(ctx context.Context, category models.Category) (*models.Category, error) {
	category.UpdatedAt = time.Now().UTC()
	return s.repo.UpdateCategory(ctx, category)
}

func (s *ServiceImpl) DeleteCategory(ctx context.Context, categoryID string) error {
	return s.repo.DeleteCategory(ctx, categoryID)
}

func (s *ServiceImpl) DeleteCategoriesByUserId(ctx context.Context, userId string) error {
	return s.repo.DeleteCategoriesByUserId(ctx, userId)
}

func (s *ServiceImpl) CreateMerchant(ctx context.Context, merchant models.Merchant) (*models.Merchant, error) {
	// Check if a merchant with the same name already exists for this user
	existingMerchant, err := s.repo.GetMerchantByNameAndUserId(ctx, merchant.Name, merchant.UserID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to check for existing merchant: %w", err)
	}
	if existingMerchant != nil {
		return nil, fmt.Errorf("merchant with name '%s' already exists for this user", merchant.Name)
	}

	merchant.ID = models.NewMerchantID()
	return s.repo.CreateMerchant(ctx, merchant)
}

func (s *ServiceImpl) GetMerchant(ctx context.Context, merchantID string) (*models.Merchant, error) {
	return s.repo.GetMerchant(ctx, merchantID)
}

func (s *ServiceImpl) ListMerchants(ctx context.Context, filter models.ListMerchantsInput) ([]models.Merchant, error) {
	return s.repo.ListMerchants(ctx, filter)
}

func (s *ServiceImpl) UpdateMerchant(ctx context.Context, merchant models.Merchant) (*models.Merchant, error) {
	merchant.UpdatedAt = time.Now().UTC()
	return s.repo.UpdateMerchant(ctx, merchant)
}

func (s *ServiceImpl) DeleteMerchant(ctx context.Context, merchantID string) error {
	return s.repo.DeleteMerchant(ctx, merchantID)
}

func (s *ServiceImpl) DeleteMerchantsByUserId(ctx context.Context, userId string) error {
	return s.repo.DeleteMerchantsByUserId(ctx, userId)
}

func (s *ServiceImpl) ListTransactionEntries(ctx context.Context, filter models.ListTransactionsInput) ([]models.TransactionEntry, error) {
	return s.repo.ListTransactionEntries(ctx, filter)
}

// CategoryGroup service methods

func (s *ServiceImpl) CreateCategoryGroup(ctx context.Context, categoryGroup models.CategoryGroup) (*models.CategoryGroup, error) {
	categoryGroup.ID = models.NewCategoryGroupID()
	return s.repo.CreateCategoryGroup(ctx, categoryGroup)
}

func (s *ServiceImpl) ListCategoryGroups(ctx context.Context, filter models.ListCategoryGroupsInput) ([]models.CategoryGroup, error) {
	return s.repo.ListCategoryGroups(ctx, filter)
}

func (s *ServiceImpl) GetCategoryGroup(ctx context.Context, categoryGroupID string) (*models.CategoryGroup, error) {
	return s.repo.GetCategoryGroup(ctx, categoryGroupID)
}

func (s *ServiceImpl) UpdateCategoryGroup(ctx context.Context, categoryGroup models.CategoryGroup) (*models.CategoryGroup, error) {
	categoryGroup.UpdatedAt = time.Now().UTC()
	return s.repo.UpdateCategoryGroup(ctx, categoryGroup)
}

func (s *ServiceImpl) DeleteCategoryGroup(ctx context.Context, categoryGroupID string) error {
	return s.repo.DeleteCategoryGroup(ctx, categoryGroupID)
}

// GetTransactionStats retrieves aggregated transaction statistics
func (s *ServiceImpl) GetTransactionStats(ctx context.Context, filter models.TransactionStatsInput) ([]models.TransactionStatsItemDto, error) {
	// Get raw stats from repository
	statsList, err := s.repo.GetTransactionStats(ctx, filter)
	if err != nil {
		return nil, err
	}

	// Set default display currency if not provided
	displayCurrency := filter.DisplayCurrency
	if displayCurrency == "" {
		displayCurrency = "EUR" // Default to EUR
	}
	displayCurrency = strings.ToUpper(displayCurrency)

	// Convert currencies and merge items with same labels
	mergedStats := s.convertAndMergeStats(statsList, displayCurrency)

	// Sort the results after merging
	s.sortTransactionStatsItems(mergedStats, filter.Sort, filter.Order)

	// Apply limit with "Other" category aggregation
	finalStats := s.applyLimitWithOther(mergedStats, filter.Limit, displayCurrency)

	return finalStats, nil
}

// convertAndMergeStats converts all amounts to the display currency and merges items with the same label
func (s *ServiceImpl) convertAndMergeStats(stats []models.TransactionStatsItemDto, displayCurrency string) []models.TransactionStatsItemDto {
	// Group by label and merge amounts after currency conversion
	labelGroups := make(map[string]*models.TransactionStatsItemDto)

	for _, stat := range stats {
		// Convert amount to display currency (for now, using 1:1 rate as requested)
		convertedAmount := s.convertCurrency(stat.Amount, stat.Currency, displayCurrency)

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

	return result
}

// convertCurrency converts an amount from source currency to target currency
// For now, using 1:1 conversion rate as requested
func (s *ServiceImpl) convertCurrency(amount int, fromCurrency, toCurrency string) int {
	// TODO: In the future, integrate with external exchange rate service
	// For now, return the same amount (1:1 conversion)
	return amount
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

// CreateTransactions creates multiple transactions atomically and generates operation ID if needed
func (s *ServiceImpl) CreateTransactions(ctx context.Context, transactions []models.Transaction) ([]models.Transaction, *string, error) {
	// Validate maximum number of transactions
	if len(transactions) > 5 {
		return nil, nil, fmt.Errorf("too many transactions: maximum 5 allowed, got %d", len(transactions))
	}

	// Validate minimum number of transactions
	if len(transactions) == 0 {
		return nil, nil, fmt.Errorf("at least one transaction is required")
	}

	var operationID *string

	// Generate operation ID if there are multiple transactions or any transaction has move_in/move_out type
	if len(transactions) > 1 || hasMovementTransactions(transactions) {
		opID := generateOperationID()
		operationID = &opID

		// Set the operation ID for all transactions
		for i := range transactions {
			opUUID, _ := uuid.Parse(opID)
			transactions[i].OperationID = &opUUID
		}
	}

	// Generate transaction IDs and validate move operations
	for i := range transactions {
		transactions[i].ID = models.NewTransactionID()

		// Validate balance exists (required)
		_, err := s.repo.GetBalance(ctx, transactions[i].BalanceID.String())
		if err != nil {
			return nil, nil, fmt.Errorf("balance with ID %s not found for transaction %d: %w", transactions[i].BalanceID.String(), i, err)
		}

		// Validate merchant exists if merchantID is provided
		if transactions[i].MerchantID != nil {
			_, err := s.repo.GetMerchant(ctx, transactions[i].MerchantID.String())
			if err != nil {
				return nil, nil, fmt.Errorf("merchant with ID %s not found for transaction %d: %w", transactions[i].MerchantID.String(), i, err)
			}
		}

		// Validate categories exist for all transaction entries
		for j, entry := range transactions[i].TransactionEntries {
			if entry.CategoryID != nil {
				_, err := s.repo.GetCategory(ctx, entry.CategoryID.String())
				if err != nil {
					return nil, nil, fmt.Errorf("category with ID %s not found for transaction %d entry %d: %w", entry.CategoryID.String(), i, j, err)
				}
			}
		}
	}

	// Validate move operations come in pairs
	if err := validateMovementOperations(transactions); err != nil {
		return nil, nil, err
	}

	created, err := s.repo.CreateTransactions(ctx, transactions)
	if err != nil {
		return nil, nil, err
	}

	return created, operationID, nil
}

// hasMovementTransactions checks if any transaction is of move_in or move_out type
func hasMovementTransactions(transactions []models.Transaction) bool {
	for _, tx := range transactions {
		if tx.Type == "move_in" || tx.Type == "move_out" {
			return true
		}
	}
	return false
}

// generateOperationID generates a new operation ID with 'fa' prefix
func generateOperationID() string {
	id := uuid.New().String()
	return "fa" + id[2:] // Replace first 2 characters with 'fa'
}

// validateMovementOperations ensures move operations are properly paired
func validateMovementOperations(transactions []models.Transaction) error {
	moveInCount := 0
	moveOutCount := 0

	for _, tx := range transactions {
		if tx.Type == "move_in" {
			moveInCount++
		} else if tx.Type == "move_out" {
			moveOutCount++
		}
	}

	// If there are move operations, we need at least one move_in and one move_out
	if (moveInCount > 0 || moveOutCount > 0) && (moveInCount == 0 || moveOutCount == 0) {
		return fmt.Errorf("movement operations require both move_in and move_out transactions")
	}

	return nil
}

// Ensure ServiceImpl implements Service
var _ Service = (*ServiceImpl)(nil)
