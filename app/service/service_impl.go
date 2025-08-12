package service

import (
	"context"
	"fmt"
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
	// Get the base transaction
	tx, err := s.repo.GetTransaction(ctx, transactionID)
	if err != nil {
		return nil, fmt.Errorf("transaction with ID %s not found: %w", transactionID, err)
	}

	// Get balance details - return error if balance not found
	balance, err := s.repo.GetBalance(ctx, tx.BalanceID.String())
	if err != nil {
		return nil, fmt.Errorf("balance with ID %s not found: %w", tx.BalanceID.String(), err)
	}

	// Create the main DTO
	dto := models.SingleTransactionDto{
		TransactionID:   tx.ID.String(),
		GroupID:         tx.GroupID.String(),
		UserID:          tx.UserID.String(),
		BalanceID:       tx.BalanceID.String(),
		BalanceTitle:    balance.Title,
		BalanceCurrency: balance.Currency,
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
		// Convert DTO entries to DAO entries
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

			// Create new entry (preserve existing ID if updating)
			var entryID uuid.UUID
			if entryDto.ID != "" {
				parsedID, err := uuid.Parse(entryDto.ID)
				if err != nil {
					return nil, fmt.Errorf("invalid entry ID format for entry %d: %w", i, err)
				}
				entryID = parsedID
			} else {
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
func (s *ServiceImpl) GetTransactionStats(ctx context.Context, filter models.TransactionStatsInput) (*models.TransactionStatsResponseDto, error) {
	// Get raw stats from repository
	rawStats, err := s.repo.GetTransactionStats(ctx, filter)
	if err != nil {
		return nil, err
	}

	// Convert to response format
	response := &models.TransactionStatsResponseDto{
		Totals: rawStats,
	}

	return response, nil
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

		// Validate move operations have negative amounts for move_out and positive for move_in
		if transactions[i].Type == "move_out" {
			for j := range transactions[i].TransactionEntries {
				if transactions[i].TransactionEntries[j].Amount > 0 {
					transactions[i].TransactionEntries[j].Amount = -transactions[i].TransactionEntries[j].Amount
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
