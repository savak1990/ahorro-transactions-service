package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// PostgreSQL GORM Models

// Merchant represents a merchant in the system
type Merchant struct {
	ID        uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name      string     `gorm:"not null" json:"name"`
	CreatedAt time.Time  `gorm:"default:now()" json:"created_at"`
	UpdatedAt time.Time  `gorm:"default:now()" json:"updated_at"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at,omitempty"`

	// Relationships
	Transactions []TransactionDB `gorm:"foreignKey:MerchantID" json:"-"`
}

// CategoryDB represents a category in the PostgreSQL database
type CategoryDB struct {
	ID            uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	CategoryName  string     `gorm:"not null" json:"category_name"`
	CategoryGroup *string    `json:"category_group"`
	Priority      *int       `json:"priority"`
	CreatedAt     time.Time  `gorm:"default:now()" json:"created_at"`
	UpdatedAt     time.Time  `gorm:"default:now()" json:"updated_at"`
	DeletedAt     *time.Time `gorm:"index" json:"deleted_at,omitempty"`

	// Relationships
	TransactionEntries []TransactionEntry `gorm:"foreignKey:CategoryID" json:"-"`
}

// TransactionDB represents a transaction in the PostgreSQL database
type TransactionDB struct {
	ID           uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	GroupID      uuid.UUID  `gorm:"type:uuid;not null" json:"group_id"`
	UserID       uuid.UUID  `gorm:"type:uuid;not null;index" json:"user_id"`
	BalanceID    uuid.UUID  `gorm:"type:uuid;not null" json:"balance_id"`
	MerchantID   *uuid.UUID `gorm:"type:uuid" json:"merchant_id"`
	Type         string     `gorm:"check:type IN ('income', 'expense', 'movement');not null" json:"type"`
	OperationID  *uuid.UUID `gorm:"type:uuid" json:"operation_id"`
	ApprovedAt   time.Time  `gorm:"not null" json:"approved_at"`
	TransactedAt time.Time  `gorm:"not null" json:"transacted_at"`
	CreatedAt    time.Time  `gorm:"default:now()" json:"created_at"`
	UpdatedAt    time.Time  `gorm:"default:now()" json:"updated_at"`
	DeletedAt    *time.Time `gorm:"index" json:"deleted_at,omitempty"`

	// Relationships
	Merchant           *Merchant          `gorm:"foreignKey:MerchantID" json:"merchant,omitempty"`
	TransactionEntries []TransactionEntry `gorm:"foreignKey:TransactionID" json:"transaction_entries"`
}

// TransactionEntry represents a transaction entry (line item) in the PostgreSQL database
type TransactionEntry struct {
	ID            uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TransactionID uuid.UUID  `gorm:"type:uuid;not null;index" json:"transaction_id"`
	Description   *string    `json:"description"`
	Amount        float64    `gorm:"type:numeric(18,2);not null" json:"amount"`
	CategoryID    *uuid.UUID `gorm:"type:uuid" json:"category_id"`
	BudgetID      *uuid.UUID `gorm:"type:uuid" json:"budget_id"`
	CreatedAt     time.Time  `gorm:"default:now()" json:"created_at"`
	UpdatedAt     time.Time  `gorm:"default:now()" json:"updated_at"`
	DeletedAt     *time.Time `gorm:"index" json:"deleted_at,omitempty"`

	// Relationships
	Transaction *TransactionDB `gorm:"foreignKey:TransactionID" json:"-"`
	Category    *CategoryDB    `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
}

// GORM Hooks for automatic timestamp updates
func (m *Merchant) BeforeUpdate(tx *gorm.DB) error {
	m.UpdatedAt = time.Now()
	return nil
}

func (c *CategoryDB) BeforeUpdate(tx *gorm.DB) error {
	c.UpdatedAt = time.Now()
	return nil
}

func (t *TransactionDB) BeforeUpdate(tx *gorm.DB) error {
	t.UpdatedAt = time.Now()
	return nil
}

func (te *TransactionEntry) BeforeUpdate(tx *gorm.DB) error {
	te.UpdatedAt = time.Now()
	return nil
}

// Conversion methods between API models and DB models

// ToAPICategory converts CategoryDB to Category (API model)
func (c *CategoryDB) ToAPICategory() Category {
	score := 0
	if c.Priority != nil {
		score = *c.Priority
	}

	return Category{
		Name:     c.CategoryName,
		ImageUrl: "", // No image URL in current schema
		Score:    score,
	}
}

// ToAPITransaction converts TransactionDB to Transaction (API model)
func (t *TransactionDB) ToAPITransaction() Transaction {
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

	return Transaction{
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

// FromAPITransaction converts Transaction (API model) to TransactionDB
func FromAPITransaction(t Transaction) (*TransactionDB, error) {
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

	return &TransactionDB{
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

// Helper functions
func parseUUID(s string) (uuid.UUID, error) {
	if s == "" {
		return uuid.UUID{}, nil
	}
	return uuid.Parse(s)
}

func formatTimePtr(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format(time.RFC3339)
}
