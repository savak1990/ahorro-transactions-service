package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// Conversion methods between DAO models (database) and DTO models (API)

// ToAPICategory converts Category (DAO) to CategoryDto (API model)
func ToAPICategory(c *Category) CategoryDto {
	if c == nil {
		return CategoryDto{}
	}

	return CategoryDto{
		CategoryID:  c.ID.String(),
		Name:        c.Name,
		Description: c.Description,
		ImageUrl:    c.ImageUrl,
		Rank:        c.Rank,
	}
}

// ToAPICategoryGroup converts CategoryGroup (DAO) to CategoryGroupDto (API model)
func ToAPICategoryGroup(cg *CategoryGroup) CategoryGroupDto {
	if cg == nil {
		return CategoryGroupDto{}
	}

	return CategoryGroupDto{
		CategoryGroupId: cg.ID.String(),
		Name:            cg.Name,
		ImageUrl:        cg.ImageUrl,
		Rank:            cg.Rank,
	}
}

// FromAPICategoryGroup converts CategoryGroupDto (API model) to CategoryGroup (DAO)
func FromAPICategoryGroup(dto CategoryGroupDto) (*CategoryGroup, error) {
	// Parse CategoryGroupId if provided, otherwise generate new ID
	var id uuid.UUID
	var err error
	if dto.CategoryGroupId != "" {
		id, err = uuid.Parse(dto.CategoryGroupId)
		if err != nil {
			return nil, fmt.Errorf("invalid category group ID format: %w", err)
		}
	} else {
		id = uuid.New() // Generate new ID if not provided
	}

	return &CategoryGroup{
		ID:       id,
		Name:     dto.Name,
		ImageUrl: dto.ImageUrl,
		Rank:     dto.Rank,
	}, nil
}

// ToAPIBalance converts Balance (DAO) to BalanceDto (API model)
func ToAPIBalance(b *Balance) BalanceDto {
	if b == nil {
		return BalanceDto{}
	}

	desc := ""
	if b.Description != nil {
		desc = *b.Description
	}

	return BalanceDto{
		BalanceID:   b.ID.String(),
		GroupID:     b.GroupID.String(),
		UserID:      b.UserID.String(),
		Currency:    b.Currency,
		Title:       b.Title,
		Description: desc,
		Rank:        &b.Rank,
		CreatedAt:   b.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   b.UpdatedAt.Format(time.RFC3339),
		DeletedAt:   formatTimePtr(b.DeletedAt),
	}
}

// ToAPICreateTransaction converts Transaction (DAO) to CreateTransactionDto (API model)
func ToAPICreateTransaction(t *Transaction) CreateTransactionDto {
	if t == nil {
		return CreateTransactionDto{}
	}

	// Convert transaction entries to DTOs
	var entryDtos []CreateTransactionEntryDto
	for _, entry := range t.TransactionEntries {
		entryDtos = append(entryDtos, ToAPICreateTransactionEntry(&entry))
	}

	// Get merchant name if available
	var merchant *string
	if t.Merchant != nil {
		merchant = &t.Merchant.Name
	}

	// Convert operation ID to string if available
	var operationID *string
	if t.OperationID != nil {
		opID := t.OperationID.String()
		operationID = &opID
	}

	// Convert approvedAt to pointer string for consistency
	var approvedAt *string
	if !t.ApprovedAt.IsZero() {
		approvedAtStr := t.ApprovedAt.Format(time.RFC3339)
		approvedAt = &approvedAtStr
	}

	return CreateTransactionDto{
		TransactionID:      t.ID.String(),
		UserID:             t.UserID.String(),
		GroupID:            t.GroupID.String(),
		BalanceID:          t.BalanceID.String(),
		Type:               t.Type,
		Merchant:           merchant,
		OperationID:        operationID,
		ApprovedAt:         approvedAt,
		TransactedAt:       t.TransactedAt.Format(time.RFC3339),
		CreatedAt:          t.CreatedAt.Format(time.RFC3339),
		UpdatedAt:          t.UpdatedAt.Format(time.RFC3339),
		DeletedAt:          formatTimePtr(t.DeletedAt),
		TransactionEntries: entryDtos,
	}
}

// FromAPICreateTransaction converts CreateTransactionDto (API model) to Transaction (DAO)
func FromAPICreateTransaction(t CreateTransactionDto) (*Transaction, error) {
	// Parse UUIDs
	id, err := parseUUID(t.TransactionID)
	if err != nil {
		id = uuid.New() // Generate new ID if not provided
	}

	userID, err := parseUUID(t.UserID)
	if err != nil {
		return nil, err
	}

	groupID, err := parseUUID(t.GroupID)
	if err != nil {
		return nil, err
	}

	balanceID, err := parseUUID(t.BalanceID)
	if err != nil {
		return nil, err
	}

	// Parse operation ID if provided
	var operationID *uuid.UUID
	if t.OperationID != nil && *t.OperationID != "" {
		opID, err := uuid.Parse(*t.OperationID)
		if err != nil {
			return nil, fmt.Errorf("invalid operation ID format: %w", err)
		}
		operationID = &opID
	}

	// Parse timestamps
	transactedAt, err := time.Parse(time.RFC3339, t.TransactedAt)
	if err != nil {
		return nil, fmt.Errorf("invalid transacted_at format: %w", err)
	}

	approvedAt := transactedAt // Default to transacted_at
	if t.ApprovedAt != nil && *t.ApprovedAt != "" {
		approvedAt, err = time.Parse(time.RFC3339, *t.ApprovedAt)
		if err != nil {
			return nil, fmt.Errorf("invalid approved_at format: %w", err)
		}
	}

	// Create transaction
	transactionID := id
	transaction := &Transaction{
		ID:           transactionID,
		GroupID:      groupID,
		UserID:       userID,
		BalanceID:    balanceID,
		Type:         t.Type,
		OperationID:  operationID,
		ApprovedAt:   approvedAt,
		TransactedAt: transactedAt,
		// Note: Merchant handling would require lookup/creation of Merchant entity
		// For now, merchant information is not stored directly on transaction
		// This could be implemented later with a merchant lookup/creation service
	}

	// Create transaction entries
	var entries []TransactionEntry
	for _, entryDto := range t.TransactionEntries {
		entry, err := FromAPICreateTransactionEntry(entryDto, transactionID)
		if err != nil {
			return nil, fmt.Errorf("error converting transaction entry: %w", err)
		}
		entries = append(entries, *entry)
	}

	transaction.TransactionEntries = entries
	return transaction, nil
}

// FromAPIBalance converts BalanceDto (API model) to Balance (DAO)
func FromAPIBalance(b BalanceDto) (*Balance, error) {
	// Parse UUIDs
	id, err := parseUUID(b.BalanceID)
	if err != nil {
		id = uuid.New() // Generate new ID if not provided
	}

	userID, err := parseUUID(b.UserID)
	if err != nil {
		return nil, err
	}

	groupID, err := parseUUID(b.GroupID)
	if err != nil {
		return nil, err
	}

	var desc *string
	if b.Description != "" {
		desc = &b.Description
	}

	rank := 0
	if b.Rank != nil {
		rank = *b.Rank
	}

	return &Balance{
		ID:          id,
		GroupID:     groupID,
		UserID:      userID,
		Currency:    b.Currency,
		Title:       b.Title,
		Description: desc,
		Rank:        rank,
	}, nil
}

// FromAPICategory converts CategoryDto (API model) to Category (DAO)
func FromAPICategory(c CategoryDto) (*Category, error) {

	// Parse CategoryID if provided, otherwise generate new ID
	var id uuid.UUID
	var err error
	if c.CategoryID != "" {
		id, err = uuid.Parse(c.CategoryID)
		if err != nil {
			return nil, fmt.Errorf("invalid category ID format: %w", err)
		}
	} else {
		id = uuid.New() // Generate new ID if not provided
	}

	return &Category{
		ID:       id,
		Name:     c.Name,
		ImageUrl: c.ImageUrl,
	}, nil
}

// ToAPICreateTransactionEntry converts TransactionEntry (DAO) to CreateTransactionEntryDto (API model)
func ToAPICreateTransactionEntry(te *TransactionEntry) CreateTransactionEntryDto {
	if te == nil {
		return CreateTransactionEntryDto{}
	}

	desc := ""
	if te.Description != nil {
		desc = *te.Description
	}

	categoryID := ""
	if te.CategoryID != nil {
		categoryID = te.CategoryID.String()
	}

	// Convert decimal to float64 for API response
	amountFloat, _ := te.Amount.Float64()

	return CreateTransactionEntryDto{
		ID:          te.ID.String(),
		Description: desc,
		Amount:      amountFloat,
		CategoryID:  categoryID,
		CreatedAt:   te.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   te.UpdatedAt.Format(time.RFC3339),
		DeletedAt:   formatTimePtr(te.DeletedAt),
	}
}

// FromAPICreateTransactionEntry converts CreateTransactionEntryDto (API model) to TransactionEntry (DAO)
func FromAPICreateTransactionEntry(te CreateTransactionEntryDto, transactionID uuid.UUID) (*TransactionEntry, error) {
	// Parse entry ID if provided, otherwise generate new ID
	var id uuid.UUID
	var err error
	if te.ID != "" {
		id, err = uuid.Parse(te.ID)
		if err != nil {
			return nil, fmt.Errorf("invalid transaction entry ID format: %w", err)
		}
	} else {
		id = uuid.New()
	}

	// Parse category ID if provided
	var categoryID *uuid.UUID
	if te.CategoryID != "" {
		catID, err := uuid.Parse(te.CategoryID)
		if err != nil {
			return nil, fmt.Errorf("invalid category ID format: %w", err)
		}
		categoryID = &catID
	}

	var desc *string
	if te.Description != "" {
		desc = &te.Description
	}

	return &TransactionEntry{
		ID:            id,
		TransactionID: transactionID,
		Description:   desc,
		Amount:        decimal.NewFromFloat(te.Amount),
		CategoryID:    categoryID,
	}, nil
}

// ToAPITransactionEntry converts TransactionEntry (DAO) to TransactionEntryDto (API model for GET responses)
func ToAPITransactionEntry(te *TransactionEntry) TransactionEntryDto {
	if te == nil {
		return TransactionEntryDto{}
	}

	// Get transaction info
	transactionID := ""
	groupID := ""
	userID := ""
	balanceID := ""
	transactionType := ""
	merchantName := ""
	merchantImageUrl := ""
	operationID := ""
	approvedAt := ""
	transactedAt := ""

	if te.Transaction != nil {
		transactionID = te.Transaction.ID.String()
		groupID = te.Transaction.GroupID.String()
		userID = te.Transaction.UserID.String()
		balanceID = te.Transaction.BalanceID.String()
		transactionType = te.Transaction.Type
		if te.Transaction.Merchant != nil {
			merchantName = te.Transaction.Merchant.Name
			if te.Transaction.Merchant.ImageUrl != nil {
				merchantImageUrl = *te.Transaction.Merchant.ImageUrl
			}
		}
		if te.Transaction.OperationID != nil {
			operationID = te.Transaction.OperationID.String()
		}
		if !te.Transaction.ApprovedAt.IsZero() {
			approvedAt = te.Transaction.ApprovedAt.Format(time.RFC3339)
		}
		transactedAt = te.Transaction.TransactedAt.Format(time.RFC3339)
	}

	// Get balance info
	balanceTitle := ""
	balanceCurrency := ""
	if te.Transaction != nil && te.Transaction.Balance != nil {
		balanceTitle = te.Transaction.Balance.Title
		balanceCurrency = te.Transaction.Balance.Currency
	}

	// Get category info
	categoryName := ""
	categoryImageUrl := ""
	categoryGroupName := ""
	categoryGroupImageUrl := ""
	if te.Category != nil {
		categoryName = te.Category.Name
		if te.Category.ImageUrl != nil {
			categoryImageUrl = *te.Category.ImageUrl
		}
		// Get category group name from the category's Group field
		categoryGroupName = te.Category.Group
		// Note: To get categoryGroupImageUrl, we would need to look up the CategoryGroup by name
		// This would require an additional database query or preloaded relationship
		// For now, we'll leave it empty and could enhance this later
	}

	// Convert decimal amount to integer cents for API response
	amountCents := int64(0)
	if !te.Amount.IsZero() {
		// Multiply by 100 to convert to cents and round to int64
		amountFloat, _ := te.Amount.Float64()
		amountCents = int64(amountFloat * 100)
	}

	return TransactionEntryDto{
		GroupID:               groupID,
		UserID:                userID,
		BalanceID:             balanceID,
		TransactionID:         transactionID,
		TransactionEntryID:    te.ID.String(),
		Type:                  transactionType,
		Amount:                int(amountCents),
		BalanceTitle:          balanceTitle,
		BalanceCurrency:       balanceCurrency,
		CategoryName:          categoryName,
		CategoryImageUrl:      categoryImageUrl,
		CategoryGroupName:     categoryGroupName,
		CategoryGroupImageUrl: &categoryGroupImageUrl,
		MerchantName:          merchantName,
		MerchantImageUrl:      merchantImageUrl,
		OperationID:           operationID,
		ApprovedAt:            approvedAt,
		TransactedAt:          transactedAt,
	}
}

// ToAPIMerchant converts Merchant (DAO) to MerchantDto (API model)
func ToAPIMerchant(m *Merchant) MerchantDto {
	if m == nil {
		return MerchantDto{}
	}

	desc := ""
	if m.Description != nil {
		desc = *m.Description
	}

	imageUrl := ""
	if m.ImageUrl != nil {
		imageUrl = *m.ImageUrl
	}

	return MerchantDto{
		MerchantID:  m.ID.String(),
		GroupID:     m.GroupID.String(),
		UserID:      m.UserID.String(),
		Name:        m.Name,
		Description: desc,
		ImageUrl:    imageUrl,
		Rank:        &m.Rank,
		CreatedAt:   m.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   m.UpdatedAt.Format(time.RFC3339),
		DeletedAt:   formatTimePtr(m.DeletedAt),
	}
}

// FromAPIMerchant converts MerchantDto (API model) to Merchant (DAO)
func FromAPIMerchant(m MerchantDto) (*Merchant, error) {
	// Parse UUID
	id, err := parseUUID(m.MerchantID)
	if err != nil {
		id = uuid.New() // Generate new ID if not provided
	}

	// Parse GroupID
	groupID, err := parseUUID(m.GroupID)
	if err != nil {
		return nil, fmt.Errorf("invalid group ID format: %w", err)
	}

	// Parse UserID
	userID, err := parseUUID(m.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID format: %w", err)
	}

	var desc *string
	if m.Description != "" {
		desc = &m.Description
	}

	var imageUrl *string
	if m.ImageUrl != "" {
		imageUrl = &m.ImageUrl
	}

	rank := 0
	if m.Rank != nil {
		rank = *m.Rank
	}

	return &Merchant{
		ID:          id,
		GroupID:     groupID,
		UserID:      userID,
		Name:        m.Name,
		Description: desc,
		ImageUrl:    imageUrl,
		Rank:        rank,
	}, nil
}
