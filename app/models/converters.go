package models

import (
	"time"

	"github.com/google/uuid"
)

// Conversion methods between DAO models (database) and DTO models (API)

// ToAPICategory converts Category (DAO) to CategoryDto (API model)
func ToAPICategory(c *Category) CategoryDto {
	score := 0
	if c.Rank != nil {
		score = *c.Rank
	}

	return CategoryDto{
		Name:     c.CategoryName,
		ImageUrl: "", // No image URL in current schema
		Score:    score,
	}
}

// ToAPIBalance converts Balance (DAO) to BalanceDto (API model)
func ToAPIBalance(b *Balance) BalanceDto {
	desc := ""
	if b.Desc != nil {
		desc = *b.Desc
	}

	return BalanceDto{
		BalanceID:   b.ID.String(),
		GroupID:     b.GroupID.String(),
		UserID:      b.UserID.String(),
		Currency:    b.Currency,
		Title:       b.Title,
		Description: desc,
		CreatedAt:   b.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   b.UpdatedAt.Format(time.RFC3339),
		DeletedAt:   formatTimePtr(b.DeletedAt),
	}
}

// ToAPITransaction converts Transaction (DAO) to TransactionDto (API model)
func ToAPITransaction(t *Transaction) TransactionDto {
	// Calculate total amount from transaction entries
	var totalAmount float64
	for _, entry := range t.TransactionEntries {
		totalAmount += entry.Amount
	}

	// Get primary category and description from first entry
	var category, description string
	if len(t.TransactionEntries) > 0 {
		if t.TransactionEntries[0].Category != nil {
			category = t.TransactionEntries[0].Category.CategoryName
		}
		if t.TransactionEntries[0].Description != nil {
			description = *t.TransactionEntries[0].Description
		}
	}

	return TransactionDto{
		TransactionID: t.ID.String(),
		UserID:        t.UserID.String(),
		GroupID:       t.GroupID.String(),
		Type:          t.Type,
		Amount:        totalAmount,
		BalanceID:     t.BalanceID.String(),
		FromBalanceID: t.BalanceID.String(), // Assuming from and to are same for now
		ToBalanceID:   "",                   // Can be populated based on business logic
		Category:      category,
		Description:   description,
		ApprovedAt:    t.ApprovedAt.Format(time.RFC3339),
		TransactedAt:  t.TransactedAt.Format(time.RFC3339),
		CreatedAt:     t.CreatedAt.Format(time.RFC3339),
		UpdatedAt:     t.UpdatedAt.Format(time.RFC3339),
		DeletedAt:     formatTimePtr(t.DeletedAt),
	}
}

// FromAPITransaction converts TransactionDto (API model) to Transaction (DAO)
func FromAPITransaction(t TransactionDto) (*Transaction, error) {
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

	// Parse timestamps
	approvedAt, err := time.Parse(time.RFC3339, t.ApprovedAt)
	if err != nil {
		approvedAt = time.Now()
	}

	transactedAt, err := time.Parse(time.RFC3339, t.TransactedAt)
	if err != nil {
		transactedAt = time.Now()
	}

	return &Transaction{
		ID:           id,
		GroupID:      groupID,
		UserID:       userID,
		BalanceID:    balanceID,
		Type:         t.Type,
		ApprovedAt:   approvedAt,
		TransactedAt: transactedAt,
		TransactionEntries: []TransactionEntry{
			{
				Description: &t.Description,
				Amount:      t.Amount,
			},
		},
	}, nil
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

	return &Balance{
		ID:       id,
		GroupID:  groupID,
		UserID:   userID,
		Currency: b.Currency,
		Title:    b.Title,
		Desc:     desc,
	}, nil
}
