package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// buildCurrencyAmountsMap converts TransactionEntryAmounts to a map of currency -> amount
func buildCurrencyAmountsMap(amounts []TransactionEntryAmount) map[string]int {
	if len(amounts) == 0 {
		return nil
	}

	currencyAmounts := make(map[string]int)
	for _, amount := range amounts {
		currencyAmounts[amount.Currency] = int(amount.Amount)
	}

	return currencyAmounts
}

// Conversion methods between DAO models (database) and DTO models (API)

// ToAPICategory converts Category (DAO) to CategoryDto (API model)
func ToAPICategory(c *Category) CategoryDto {
	if c == nil {
		return CategoryDto{}
	}

	dto := CategoryDto{
		CategoryID:  c.ID.String(),
		Name:        c.Name,
		Description: c.Description,
		ImageUrl:    c.ImageUrl,
		Rank:        c.Rank,
		IsDeleted:   c.DeletedAt != nil,
	}

	// Add category group information if available
	if c.CategoryGroup != nil {
		dto.CategoryGroupID = c.CategoryGroup.ID.String()
		dto.CategoryGroupName = c.CategoryGroup.Name
		dto.CategoryGroupImageUrl = c.CategoryGroup.ImageUrl
		dto.CategoryGroupRank = c.CategoryGroup.Rank

		// Set CategoryGroupDeleted flag if the category group is soft deleted
		dto.CategoryGroupDeleted = c.CategoryGroup.DeletedAt != nil
	} else if c.CategoryGroupId != "" {
		// If category has a CategoryGroupId but CategoryGroup is nil,
		// it means the category group is soft-deleted (due to preload filtering)
		dto.CategoryGroupID = c.CategoryGroupId
		dto.CategoryGroupDeleted = true
	}

	return dto
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
		IsDeleted:       cg.DeletedAt != nil,
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
		id = NewCategoryGroupID() // Generate new ID if not provided
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

	// Get merchant ID if available
	var merchantID string
	if t.MerchantID != nil {
		merchantID = t.MerchantID.String()
	}

	// Convert operation ID to string if available
	var operationID string
	if t.OperationID != nil {
		operationID = t.OperationID.String()
	}

	// Convert approvedAt to string
	var approvedAt string
	if !t.ApprovedAt.IsZero() {
		approvedAt = t.ApprovedAt.Format(time.RFC3339)
	}

	return CreateTransactionDto{
		TransactionID:      t.ID.String(),
		UserID:             t.UserID.String(),
		GroupID:            t.GroupID.String(),
		BalanceID:          t.BalanceID.String(),
		Type:               t.Type,
		MerchantID:         merchantID,
		OperationID:        operationID,
		ApprovedAt:         approvedAt,
		TransactedAt:       t.TransactedAt.Format(time.RFC3339),
		CreatedAt:          t.CreatedAt.Format(time.RFC3339),
		UpdatedAt:          t.UpdatedAt.Format(time.RFC3339),
		DeletedAt:          formatTimePtr(t.DeletedAt),
		TransactionEntries: entryDtos,
	}
}

// ToAPIUpdateTransaction converts Transaction (DAO) to UpdateTransactionDto (API model)
func ToAPIUpdateTransaction(t *Transaction) UpdateTransactionDto {
	if t == nil {
		return UpdateTransactionDto{}
	}

	// Convert transaction entries to DTOs
	var entryDtos []CreateTransactionEntryDto
	for _, entry := range t.TransactionEntries {
		entryDtos = append(entryDtos, ToAPICreateTransactionEntry(&entry))
	}

	// Get merchant ID if available
	var merchantID string
	if t.MerchantID != nil {
		merchantID = t.MerchantID.String()
	}

	// Convert operation ID to string if available
	var operationID string
	if t.OperationID != nil {
		operationID = t.OperationID.String()
	}

	// Convert approvedAt to string
	var approvedAt string
	if !t.ApprovedAt.IsZero() {
		approvedAt = t.ApprovedAt.Format(time.RFC3339)
	}

	return UpdateTransactionDto{
		TransactionID:      t.ID.String(),
		UserID:             t.UserID.String(),
		GroupID:            t.GroupID.String(),
		BalanceID:          t.BalanceID.String(),
		Type:               t.Type,
		MerchantID:         merchantID,
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
		id = NewTransactionID() // Generate new ID with prefix if not provided
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
	if t.OperationID != "" {
		opID, err := uuid.Parse(t.OperationID)
		if err != nil {
			return nil, fmt.Errorf("invalid operation ID format: %w", err)
		}
		operationID = &opID
	}

	// Parse merchant ID if provided
	var merchantID *uuid.UUID
	if t.MerchantID != "" {
		// Validate merchant ID format
		mID, err := uuid.Parse(t.MerchantID)
		if err != nil {
			return nil, fmt.Errorf("invalid merchant ID format: %w", err)
		}
		merchantID = &mID
	}

	// Parse timestamps
	var transactedAt time.Time
	if t.TransactedAt != "" {
		var err error
		transactedAt, err = time.Parse(time.RFC3339, t.TransactedAt)
		if err != nil {
			return nil, fmt.Errorf("invalid transacted_at format: %w", err)
		}
	} else {
		// Use current time if transactedAt is not provided
		transactedAt = time.Now().UTC()
	}

	approvedAt := transactedAt // Default to transacted_at
	if t.ApprovedAt != "" {
		approvedAt, err = time.Parse(time.RFC3339, t.ApprovedAt)
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
		MerchantID:   merchantID,
		Type:         t.Type,
		OperationID:  operationID,
		ApprovedAt:   approvedAt,
		TransactedAt: transactedAt,
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
		id = NewBalanceID() // Generate new ID if not provided
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
		id = NewCategoryID() // Generate new ID if not provided
	}

	// Parse CategoryGroupId if provided
	var categoryGroupId string
	if c.CategoryGroupID != "" {
		// Validate that it's a proper UUID format
		_, err = uuid.Parse(c.CategoryGroupID)
		if err != nil {
			return nil, fmt.Errorf("invalid category group ID format: %w", err)
		}
		categoryGroupId = c.CategoryGroupID
	}

	return &Category{
		ID:              id,
		CategoryGroupId: categoryGroupId,
		Name:            c.Name,
		Description:     c.Description,
		Rank:            c.Rank,
		ImageUrl:        c.ImageUrl,
		// Note: UserId and GroupId should be set by the handler from request context
	}, nil
}

// FromAPICreateCategory converts CreateCategoryDto (API model) to Category (DAO)
func FromAPICreateCategory(c CreateCategoryDto) (*Category, error) {
	// Generate new ID for creation
	id := NewCategoryID()

	// Parse UserID
	userID, err := uuid.Parse(c.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID format: %w", err)
	}

	// Parse GroupID
	groupID, err := uuid.Parse(c.GroupID)
	if err != nil {
		return nil, fmt.Errorf("invalid group ID format: %w", err)
	}

	// Parse CategoryGroupId if provided
	var categoryGroupId string
	if c.CategoryGroupID != "" {
		// Validate that it's a proper UUID format
		_, err = uuid.Parse(c.CategoryGroupID)
		if err != nil {
			return nil, fmt.Errorf("invalid category group ID format: %w", err)
		}
		categoryGroupId = c.CategoryGroupID
	}

	return &Category{
		ID:              id,
		UserId:          userID,
		GroupId:         groupID,
		CategoryGroupId: categoryGroupId,
		Name:            c.Name,
		Description:     c.Description,
		Rank:            c.Rank,
		ImageUrl:        c.ImageUrl,
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

	// Amount is already in cents (int64), use directly
	amountCents := int(te.Amount)

	return CreateTransactionEntryDto{
		ID:              te.ID.String(),
		Description:     desc,
		Amount:          amountCents,
		CategoryID:      categoryID,
		CurrencyAmounts: buildCurrencyAmountsMap(te.TransactionEntryAmounts),
		CreatedAt:       te.CreatedAt.Format(time.RFC3339),
		UpdatedAt:       te.UpdatedAt.Format(time.RFC3339),
		DeletedAt:       formatTimePtr(te.DeletedAt),
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
		id = NewTransactionEntryID()
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
		Amount:        int64(te.Amount),
		CategoryID:    categoryID,
	}, nil
}

// FromAPIUpdateTransactionEntry converts UpdateTransactionEntryDto (API model) to TransactionEntry (DAO)
func FromAPIUpdateTransactionEntry(te UpdateTransactionEntryDto, transactionID uuid.UUID) (*TransactionEntry, error) {
	// Parse entry ID (required for updates)
	id, err := uuid.Parse(te.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid transaction entry ID format: %w", err)
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
		Amount:        int64(te.Amount),
		CategoryID:    categoryID,
	}, nil
}

// ToAPIUpdateTransactionEntry converts TransactionEntry (DAO) to UpdateTransactionEntryDto (API model)
func ToAPIUpdateTransactionEntry(te *TransactionEntry) UpdateTransactionEntryDto {
	if te == nil {
		return UpdateTransactionEntryDto{}
	}

	desc := ""
	if te.Description != nil {
		desc = *te.Description
	}

	categoryID := ""
	if te.CategoryID != nil {
		categoryID = te.CategoryID.String()
	}

	// Amount is already in cents (int64), use directly
	amountCents := int(te.Amount)

	return UpdateTransactionEntryDto{
		ID:          te.ID.String(),
		Description: desc,
		Amount:      amountCents,
		CategoryID:  categoryID,
		CreatedAt:   te.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   te.UpdatedAt.Format(time.RFC3339),
		DeletedAt:   formatTimePtr(te.DeletedAt),
	}
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
	balanceDeleted := false
	if te.Transaction != nil && te.Transaction.Balance != nil {
		balanceTitle = te.Transaction.Balance.Title
		balanceCurrency = te.Transaction.Balance.Currency
		balanceDeleted = te.Transaction.Balance.DeletedAt != nil
	}

	// Get category info
	categoryID := ""
	categoryName := ""
	categoryImageUrl := ""
	categoryGroupID := ""
	categoryGroupName := ""
	categoryGroupImageUrl := ""
	categoryIsDeleted := false
	categoryGroupDeleted := false
	if te.Category != nil {
		categoryID = te.Category.ID.String()
		categoryName = te.Category.Name
		if te.Category.ImageUrl != nil {
			categoryImageUrl = *te.Category.ImageUrl
		}
		// Get category group info
		categoryGroupID = te.Category.CategoryGroupId
		categoryGroupName = te.Category.Group
		// Check if category is soft deleted
		categoryIsDeleted = te.Category.DeletedAt != nil
		// Check if category group is soft deleted
		if te.Category.CategoryGroup != nil {
			categoryGroupDeleted = te.Category.CategoryGroup.DeletedAt != nil
		} else if te.Category.CategoryGroupId != "" {
			// If category has a CategoryGroupId but CategoryGroup is nil,
			// it means the category group is soft-deleted (due to preload filtering)
			categoryGroupDeleted = true
		}
		// Note: To get categoryGroupImageUrl, we would need to look up the CategoryGroup by name
		// This would require an additional database query or preloaded relationship
		// For now, we'll leave it empty and could enhance this later
	}

	// Amount is already in cents (int64), use directly
	amountCents := te.Amount

	return TransactionEntryDto{
		GroupID:               groupID,
		UserID:                userID,
		BalanceID:             balanceID,
		TransactionID:         transactionID,
		TransactionEntryID:    te.ID.String(),
		Type:                  transactionType,
		Amount:                int(amountCents),
		CurrencyAmounts:       buildCurrencyAmountsMap(te.TransactionEntryAmounts),
		BalanceTitle:          balanceTitle,
		BalanceCurrency:       balanceCurrency,
		BalanceDeleted:        balanceDeleted,
		CategoryID:            categoryID,
		CategoryName:          categoryName,
		CategoryImageUrl:      categoryImageUrl,
		CategoryGroupName:     categoryGroupName,
		CategoryGroupImageUrl: &categoryGroupImageUrl,
		CategoryGroupID:       categoryGroupID,
		CategoryIsDeleted:     categoryIsDeleted,
		CategoryGroupDeleted:  categoryGroupDeleted,
		MerchantName:          merchantName,
		MerchantImageUrl:      merchantImageUrl,
		OperationID:           operationID,
		ApprovedAt:            approvedAt,
		TransactedAt:          transactedAt,
		CreatedAt:             te.CreatedAt.Format(time.RFC3339),
		UpdatedAt:             te.UpdatedAt.Format(time.RFC3339),
		DeletedAt: func() string {
			if te.DeletedAt != nil {
				return te.DeletedAt.Format(time.RFC3339)
			}
			return ""
		}(),
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
		id = NewMerchantID() // Generate new ID if not provided
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

// FromAPIUpdateCategory converts UpdateCategoryDto (API model) to Category (DAO)
func FromAPIUpdateCategory(c UpdateCategoryDto, categoryID string) (*Category, error) {
	// Parse CategoryID
	id, err := uuid.Parse(categoryID)
	if err != nil {
		return nil, fmt.Errorf("invalid category ID format: %w", err)
	}

	// Parse UserID
	userID, err := uuid.Parse(c.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID format: %w", err)
	}

	// Parse GroupID
	groupID, err := uuid.Parse(c.GroupID)
	if err != nil {
		return nil, fmt.Errorf("invalid group ID format: %w", err)
	}

	// Parse CategoryGroupId if provided
	var categoryGroupId string
	if c.CategoryGroupID != "" {
		// Validate that it's a proper UUID format
		_, err = uuid.Parse(c.CategoryGroupID)
		if err != nil {
			return nil, fmt.Errorf("invalid category group ID format: %w", err)
		}
		categoryGroupId = c.CategoryGroupID
	}

	return &Category{
		ID:              id,
		UserId:          userID,
		GroupId:         groupID,
		CategoryGroupId: categoryGroupId,
		Name:            c.Name,
		Description:     c.Description,
		Rank:            c.Rank,
		ImageUrl:        c.ImageUrl,
	}, nil
}

// FromAPICreateTransactionsRequest converts CreateTransactionsRequestDto to slice of Transaction DAOs
func FromAPICreateTransactionsRequest(req CreateTransactionsRequestDto) ([]Transaction, error) {
	var transactions []Transaction
	for _, txDto := range req.Transactions {
		tx, err := FromAPICreateTransaction(txDto)
		if err != nil {
			return nil, fmt.Errorf("error converting transaction: %w", err)
		}
		transactions = append(transactions, *tx)
	}
	return transactions, nil
}

// ToAPICreateTransactionsResponse converts slice of Transaction DAOs to CreateTransactionsResponseDto
func ToAPICreateTransactionsResponse(transactions []Transaction, operationID *string) CreateTransactionsResponseDto {
	var txDtos []CreateTransactionDto
	for _, tx := range transactions {
		txDtos = append(txDtos, ToAPICreateTransaction(&tx))
	}
	return CreateTransactionsResponseDto{
		Transactions: txDtos,
		OperationID:  operationID,
	}
}

// FromAPIUpdateTransaction converts UpdateTransactionDto (API model) to Transaction (DAO)
func FromAPIUpdateTransaction(t UpdateTransactionDto) (*Transaction, error) {
	// Parse UUIDs
	id, err := parseUUID(t.TransactionID)
	if err != nil {
		return nil, fmt.Errorf("transaction ID is required for update")
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
	if t.OperationID != "" {
		opID, err := uuid.Parse(t.OperationID)
		if err != nil {
			return nil, fmt.Errorf("invalid operation ID format: %w", err)
		}
		operationID = &opID
	}

	// Parse merchant ID if provided
	var merchantID *uuid.UUID
	if t.MerchantID != "" {
		// Validate merchant ID format
		mID, err := uuid.Parse(t.MerchantID)
		if err != nil {
			return nil, fmt.Errorf("invalid merchant ID format: %w", err)
		}
		merchantID = &mID
	}

	// Parse timestamps
	var transactedAt time.Time
	if t.TransactedAt != "" {
		var err error
		transactedAt, err = time.Parse(time.RFC3339, t.TransactedAt)
		if err != nil {
			return nil, fmt.Errorf("invalid transacted_at format: %w", err)
		}
	} else {
		// Use current time if transactedAt is not provided
		transactedAt = time.Now().UTC()
	}

	approvedAt := transactedAt // Default to transacted_at
	if t.ApprovedAt != "" {
		approvedAt, err = time.Parse(time.RFC3339, t.ApprovedAt)
		if err != nil {
			return nil, fmt.Errorf("invalid approved_at format: %w", err)
		}
	}

	// Create transaction
	transaction := &Transaction{
		ID:           id,
		GroupID:      groupID,
		UserID:       userID,
		BalanceID:    balanceID,
		MerchantID:   merchantID,
		Type:         t.Type,
		OperationID:  operationID,
		ApprovedAt:   approvedAt,
		TransactedAt: transactedAt,
	}

	// Create transaction entries
	var entries []TransactionEntry
	for _, entryDto := range t.TransactionEntries {
		entry, err := FromAPICreateTransactionEntry(entryDto, id)
		if err != nil {
			return nil, fmt.Errorf("error converting transaction entry: %w", err)
		}
		entries = append(entries, *entry)
	}

	transaction.TransactionEntries = entries
	return transaction, nil
}

// ToAPISingleTransactionEntry converts TransactionEntry (DAO) to SingleTransactionEntryDto with category details
func ToAPISingleTransactionEntry(te *TransactionEntry, category *Category) SingleTransactionEntryDto {
	if te == nil {
		return SingleTransactionEntryDto{}
	}

	dto := SingleTransactionEntryDto{
		TransactionEntryID: te.ID.String(),
		Amount:             int(te.Amount),
		CurrencyAmounts:    buildCurrencyAmountsMap(te.TransactionEntryAmounts),
		CreatedAt:          te.CreatedAt.Format(time.RFC3339),
		UpdatedAt:          te.UpdatedAt.Format(time.RFC3339),
		DeletedAt:          formatTimePtr(te.DeletedAt),
	}

	// Set description
	if te.Description != nil {
		dto.Description = *te.Description
	}

	// Add category information if available
	if category != nil {
		dto.CategoryID = category.ID.String()
		dto.CategoryName = category.Name
		if category.ImageUrl != nil && *category.ImageUrl != "" {
			dto.CategoryIcon = *category.ImageUrl
		}

		// Add category group information if available
		if category.CategoryGroup != nil {
			dto.CategoryGroupID = category.CategoryGroup.ID.String()
			dto.CategoryGroupName = category.CategoryGroup.Name
			if category.CategoryGroup.ImageUrl != nil && *category.CategoryGroup.ImageUrl != "" {
				dto.CategoryGroupIcon = *category.CategoryGroup.ImageUrl
			}
		}
	}

	return dto
}
