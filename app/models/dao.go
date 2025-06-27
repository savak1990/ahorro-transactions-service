package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
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
	Name        string     `gorm:"not null"`
	Description *string    `gorm:"type:varchar(500)"`
	ImageUrl    *string    `gorm:"type:varchar(255)"`
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

// Category represents a category in the PostgreSQL database
type Category struct {
	ID           uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	CategoryName string     `gorm:"not null"`
	Group        *string    `gorm:"column:group"`
	Rank         *int       `gorm:"column:rank"`
	ImageUrl     *string    `gorm:"column:image_url;type:varchar(255)"`
	CreatedAt    time.Time  `gorm:"default:now()"`
	UpdatedAt    time.Time  `gorm:"default:now()"`
	DeletedAt    *time.Time `gorm:"index"`
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
	Merchant           *Merchant          `gorm:"foreignKey:MerchantID"`
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
	Amount        decimal.Decimal `gorm:"type:numeric(18,2);not null"`
	CategoryID    *uuid.UUID      `gorm:"type:uuid;index:idx_transaction_entry_category_id"`
	CreatedAt     time.Time       `gorm:"default:now()"`
	UpdatedAt     time.Time       `gorm:"default:now()"`
	DeletedAt     *time.Time      `gorm:"index"`

	// Relationships
	Transaction *Transaction `gorm:"foreignKey:TransactionID"`
	Category    *Category    `gorm:"foreignKey:CategoryID"`
}

// TableName specifies the table name for GORM
func (TransactionEntry) TableName() string {
	return "transaction_entry"
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
	validTypes := []string{"income", "expense", "movement"}
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
