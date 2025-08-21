package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/savak1990/transactions-service/app/models"
)

// createTransactionEntryAmounts creates TransactionEntryAmount records for all supported currencies
func (s *ServiceImpl) createTransactionEntryAmounts(
	entryID uuid.UUID,
	baseAmount int64,
	baseCurrency string,
	supportedCurrencies []string,
	exchangeRates map[string]float64,
) []models.TransactionEntryAmount {
	var entryAmounts []models.TransactionEntryAmount

	for _, currency := range supportedCurrencies {
		var amount int64
		var exchangeRate float64

		if currency == baseCurrency {
			// Base currency: use original amount and exchange rate of 1.0
			amount = baseAmount
			exchangeRate = 1.0
		} else {
			// Get exchange rate for this currency
			if rate, exists := exchangeRates[currency]; exists {
				// Convert to target currency
				amount = int64(float64(baseAmount) * rate)
				exchangeRate = rate
			} else {
				// Skip currency if no exchange rate available
				continue
			}
		}

		entryAmount := models.TransactionEntryAmount{
			TransactionEntryID: entryID,
			Currency:           currency,
			Amount:             amount,
			ExchangeRate:       exchangeRate,
		}
		entryAmounts = append(entryAmounts, entryAmount)
	}

	return entryAmounts
}

func (s *ServiceImpl) CreateTransaction(ctx context.Context, tx models.Transaction) (*models.Transaction, error) {
	tx.ID = models.NewTransactionID()

	// Validate balance exists (required) and get it for currency determination
	balance, err := s.repo.GetBalance(ctx, tx.BalanceID.String())
	if err != nil {
		return nil, fmt.Errorf("balance with ID %s not found: %w", tx.BalanceID.String(), err)
	}
	baseCurrency := balance.Currency

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

	// Fetch supported currencies and exchange rates
	supportedCurrencies, err := s.exchangeRatesDb.GetSupportedCurrencies(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get supported currencies: %w", err)
	}

	// Fetch exchange rates for all supported currencies
	exchangeRates, err := s.exchangeRatesDb.GetSupportedCurrenciesRates(ctx, baseCurrency)
	if err != nil {
		return nil, fmt.Errorf("failed to get exchange rates for base currency %s: %w", baseCurrency, err)
	}

	// Create TransactionEntryAmount for each transaction entry and each supported currency
	for i := range tx.TransactionEntries {
		entry := &tx.TransactionEntries[i]
		entry.ID = models.NewTransactionEntryID()
		entry.TransactionID = tx.ID

		entry.TransactionEntryAmounts = s.createTransactionEntryAmounts(
			entry.ID,
			entry.Amount,
			baseCurrency,
			supportedCurrencies,
			exchangeRates)
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
		// Get the balance to determine the base currency (might have changed)
		var baseCurrency string
		if updateDto.BalanceID != "" {
			// Balance was updated in this request
			balance, err := s.repo.GetBalance(ctx, updateDto.BalanceID)
			if err != nil {
				return nil, fmt.Errorf("failed to get updated balance for currency determination: %w", err)
			}
			baseCurrency = balance.Currency
		} else {
			// Use existing balance
			balance, err := s.repo.GetBalance(ctx, existingTx.BalanceID.String())
			if err != nil {
				return nil, fmt.Errorf("failed to get balance for currency determination: %w", err)
			}
			baseCurrency = balance.Currency
		}

		// Fetch supported currencies and exchange rates
		supportedCurrencies, err := s.exchangeRatesDb.GetSupportedCurrencies(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get supported currencies: %w", err)
		}

		// Fetch exchange rates for all supported currencies
		exchangeRates, err := s.exchangeRatesDb.GetSupportedCurrenciesRates(ctx, baseCurrency)
		if err != nil {
			return nil, fmt.Errorf("failed to get exchange rates for base currency %s: %w", baseCurrency, err)
		}

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

			entry.TransactionEntryAmounts = s.createTransactionEntryAmounts(
				entryID,
				entry.Amount,
				baseCurrency,
				supportedCurrencies,
				exchangeRates,
			)

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
