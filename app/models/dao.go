package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// PostgreSQL GORM Models

// Balance represents a balance/account in the system
type Balance struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	GroupID     uuid.UUID  `gorm:"type:uuid;not null;index:idx_balance_group_id"`
	UserID      uuid.UUID  `gorm:"type:uuid;not null;index:idx_balance_user_id"`
	Currency    string     `gorm:"type:varchar(3);not null;default:'EUR'"` // ISO 4217 currency codes (3 chars)
	Title       string     `gorm:"type:varchar(100);not null"`             // Limited to 100 characters
	Description *string    `gorm:"type:varchar(500)"`                      // Optional description, limited to 500 characters
	Rank        int        `gorm:"column:rank"`                            // Optional rank for ordering
	CreatedAt   time.Time  `gorm:"default:now()"`
	UpdatedAt   time.Time  `gorm:"default:now()"`
	DeletedAt   *time.Time `gorm:"index"`

	// Relationships
	Transactions []Transaction `gorm:"foreignKey:BalanceID"`
}

// TableName specifies the table name for GORM
func (Balance) TableName() string {
	return "balance"
}

// Merchant represents a merchant in the system
type Merchant struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	GroupID     uuid.UUID  `gorm:"type:uuid;not null;index:idx_merchant_group_id"`
	UserID      uuid.UUID  `gorm:"type:uuid;not null;index:idx_merchant_user_id"`
	Name        string     `gorm:"not null"`
	Description *string    `gorm:"type:varchar(255)"`
	ImageUrl    *string    `gorm:"type:varchar(255)"`
	Rank        int        `gorm:"column:rank"`
	CreatedAt   time.Time  `gorm:"default:now()"`
	UpdatedAt   time.Time  `gorm:"default:now()"`
	DeletedAt   *time.Time `gorm:"index"`

	// Relationships
	Transactions []Transaction `gorm:"foreignKey:MerchantID"`
}

// TableName specifies the table name for GORM
func (Merchant) TableName() string {
	return "merchant"
}

// Category group represents a group of categories
type CategoryGroup struct {
	ID        uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name      string     `gorm:"not null"`
	Rank      *int       `gorm:"column:rank"`
	ImageUrl  *string    `gorm:"column:image_url;type:varchar(255)"`
	CreatedAt time.Time  `gorm:"default:now()"`
	UpdatedAt time.Time  `gorm:"default:now()"`
	DeletedAt *time.Time `gorm:"index"`
}

// TableName specifies the table name for GORM
func (CategoryGroup) TableName() string {
	return "category_group"
}

// Category represents a category in the PostgreSQL database
type Category struct {
	ID              uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserId          uuid.UUID  `gorm:"type:uuid;not null;index:idx_category_user_id"`
	GroupId         uuid.UUID  `gorm:"type:uuid;not null;index:idx_category_group_id"`
	CategoryGroupId string     `gorm:"not null;index:idx_category_group_id"`
	Name            string     `gorm:"not null"`
	Group           string     `gorm:"column:group;not null"` // Category group name from "group" column
	Description     string     `gorm:"type:varchar(255)"`
	Rank            *int       `gorm:"column:rank"`
	ImageUrl        *string    `gorm:"column:image_url;type:varchar(255)"`
	CreatedAt       time.Time  `gorm:"default:now()"`
	UpdatedAt       time.Time  `gorm:"default:now()"`
	DeletedAt       *time.Time `gorm:"index"`

	// Relationships
	CategoryGroup *CategoryGroup `gorm:"foreignKey:CategoryGroupId;references:ID"`
}

// TableName specifies the table name for GORM
func (Category) TableName() string {
	return "category"
}

// Transaction represents a transaction in the PostgreSQL database
type Transaction struct {
	ID           uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	GroupID      uuid.UUID  `gorm:"type:uuid;not null;index:idx_transaction_group_id"`
	UserID       uuid.UUID  `gorm:"type:uuid;not null;index:idx_transaction_user_id"`
	BalanceID    uuid.UUID  `gorm:"type:uuid;not null;index:idx_transaction_balance_id"`
	MerchantID   *uuid.UUID `gorm:"type:uuid;index:idx_transaction_merchant_id"`
	Type         string     `gorm:"type:varchar(20);not null;index:idx_transaction_type"`
	OperationID  *uuid.UUID `gorm:"type:uuid;index:idx_transaction_operation_id"`
	ApprovedAt   time.Time  `gorm:"not null"`
	TransactedAt time.Time  `gorm:"not null;index:idx_transaction_transacted_at"`
	CreatedAt    time.Time  `gorm:"default:now()"`
	UpdatedAt    time.Time  `gorm:"default:now()"`
	DeletedAt    *time.Time `gorm:"index"`

	// Relationships
	Merchant           *Merchant          `gorm:"foreignKey:MerchantID;constraint:OnDelete:SET NULL"`
	Balance            *Balance           `gorm:"foreignKey:BalanceID"`
	TransactionEntries []TransactionEntry `gorm:"foreignKey:TransactionID"`
}

// TableName specifies the table name for GORM
func (Transaction) TableName() string {
	return "transaction"
}

// TransactionEntry represents a transaction entry (line item) in the PostgreSQL database
type TransactionEntry struct {
	ID            uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	TransactionID uuid.UUID `gorm:"type:uuid;not null;index:idx_transaction_entry_transaction_id"`
	Description   *string
	Amount        int64      `gorm:"type:bigint;not null"` // Amount in cents (base currency)
	CategoryID    *uuid.UUID `gorm:"type:uuid;index:idx_transaction_entry_category_id"`
	CreatedAt     time.Time  `gorm:"default:now()"`
	UpdatedAt     time.Time  `gorm:"default:now()"`
	DeletedAt     *time.Time `gorm:"index"`

	// Relationships
	Transaction             *Transaction             `gorm:"foreignKey:TransactionID"`
	Category                *Category                `gorm:"foreignKey:CategoryID;references:ID"`
	TransactionEntryAmounts []TransactionEntryAmount `gorm:"foreignKey:TransactionEntryID"`
}

// TransactionEntryAmount represents amounts in different currencies for a transaction entry
type TransactionEntryAmount struct {
	TransactionEntryID uuid.UUID `gorm:"type:uuid;not null;primaryKey;index:idx_transaction_entry_amount_entry_id"`
	Currency           string    `gorm:"type:varchar(3);not null;primaryKey"` // ISO 4217 currency codes (3 chars)
	Amount             int64     `gorm:"type:bigint;not null"`                // Amount in cents for the specific currency
	ExchangeRate       float64   `gorm:"type:decimal(10,6);not null"`         // Exchange rate used for conversion
	CreatedAt          time.Time `gorm:"default:now()"`
	UpdatedAt          time.Time `gorm:"default:now()"`

	// Relationships
	TransactionEntry *TransactionEntry `gorm:"foreignKey:TransactionEntryID"`
}

// TableName specifies the table name for GORM
func (TransactionEntry) TableName() string {
	return "transaction_entry"
}

// TableName specifies the table name for GORM
func (TransactionEntryAmount) TableName() string {
	return "transaction_entry_amount"
}

// GORM Hooks for automatic timestamp updates
func (b *Balance) BeforeUpdate(tx *gorm.DB) error {
	b.UpdatedAt = time.Now()
	return nil
}

func (m *Merchant) BeforeUpdate(tx *gorm.DB) error {
	m.UpdatedAt = time.Now()
	return nil
}

func (c *Category) BeforeUpdate(tx *gorm.DB) error {
	c.UpdatedAt = time.Now()
	return nil
}

func (t *Transaction) BeforeUpdate(tx *gorm.DB) error {
	t.UpdatedAt = time.Now()
	return t.validateType()
}

func (t *Transaction) BeforeCreate(tx *gorm.DB) error {
	return t.validateType()
}

func (t *Transaction) validateType() error {
	// Skip validation if type is empty (happens during partial updates like setting merchant_id to NULL)
	if t.Type == "" {
		return nil
	}

	validTypes := []string{"init", "income", "expense", "move_in", "move_out"}
	for _, validType := range validTypes {
		if t.Type == validType {
			return nil
		}
	}
	return fmt.Errorf("invalid transaction type: %s. Must be one of: %v", t.Type, validTypes)
}

func (te *TransactionEntry) BeforeUpdate(tx *gorm.DB) error {
	te.UpdatedAt = time.Now()
	return nil
}

func (tea *TransactionEntryAmount) BeforeUpdate(tx *gorm.DB) error {
	tea.UpdatedAt = time.Now()
	return nil
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

// TransactionStatsRaw represents raw statistics data from database aggregation
type TransactionStatsRaw struct {
	TransactionType         string `gorm:"column:transaction_type"`
	Currency                string `gorm:"column:currency"`
	TotalAmount             int64  `gorm:"column:total_amount"` // Amount in cents
	TransactionsCount       int64  `gorm:"column:transactions_count"`
	TransactionEntriesCount int64  `gorm:"column:transaction_entries_count"`
}
