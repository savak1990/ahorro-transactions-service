package models

// CreateTransactionDto represents a financial transaction for creation via POST requests.
type CreateTransactionDto struct {
	TransactionID      string                      `json:"transactionId,omitempty"`
	GroupID            string                      `json:"groupId"`
	UserID             string                      `json:"userId"`
	BalanceID          string                      `json:"balanceId"`
	Type               string                      `json:"type"`
	Merchant           *string                     `json:"merchant,omitempty"`
	OperationID        *string                     `json:"operationId,omitempty"`
	ApprovedAt         *string                     `json:"approvedAt,omitempty"`
	TransactedAt       string                      `json:"transactedAt"`
	CreatedAt          string                      `json:"createdAt,omitempty"`
	UpdatedAt          string                      `json:"updatedAt,omitempty"`
	DeletedAt          string                      `json:"deletedAt,omitempty"`
	TransactionEntries []CreateTransactionEntryDto `json:"transactionEntries"`
}

// CreateTransactionEntryDto represents a single entry within a transaction for creation
type CreateTransactionEntryDto struct {
	ID          string  `json:"id,omitempty"`
	Description string  `json:"description"`
	Amount      float64 `json:"amount"`
	CategoryID  string  `json:"categoryId"`
	CreatedAt   string  `json:"createdAt,omitempty"`
	UpdatedAt   string  `json:"updatedAt,omitempty"`
	DeletedAt   string  `json:"deletedAt,omitempty"`
}

// TransactionDto represents a financial transaction for API responses.
type TransactionEntryDto struct {
	GroupID               string  `json:"groupId"`
	UserID                string  `json:"userId"`
	BalanceID             string  `json:"balanceId"`
	TransactionID         string  `json:"transactionId"`
	TransactionEntryID    string  `json:"transactionEntryId"`
	Type                  string  `json:"type"`
	Amount                int     `json:"amount"`
	BalanceTitle          string  `json:"balanceTitle"`
	BalanceCurrency       string  `json:"balanceCurrency"`
	CategoryName          string  `json:"categoryName"`
	CategoryImageUrl      string  `json:"categoryImageUrl,omitempty"`
	CategoryGroupName     string  `json:"categoryGroupName,omitempty"`
	CategoryGroupImageUrl *string `json:"categoryGroupImageUrl,omitempty"`
	CategoryGroupID       string  `json:"categoryGroupId,omitempty"`
	MerchantName          string  `json:"merchantName,omitempty"`
	MerchantImageUrl      string  `json:"merchantImageUrl,omitempty"`
	OperationID           string  `json:"operationId,omitempty"`
	ApprovedAt            string  `json:"approvedAt,omitempty"`
	TransactedAt          string  `json:"transactedAt"`
}

// Balance represents a user's balance/account for API responses.
type BalanceDto struct {
	BalanceID   string `json:"balanceId"`
	GroupID     string `json:"groupId"`
	UserID      string `json:"userId"`
	Currency    string `json:"currency"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	Rank        *int   `json:"rank,omitempty"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
	DeletedAt   string `json:"deletedAt,omitempty"`
}

// CategoryGroup represents a group of categories for API responses.
type CategoryGroupDto struct {
	CategoryGroupId string  `json:"categoryGroupId"`
	Name            string  `json:"name"`
	Description     string  `json:"description,omitempty"`
	ImageUrl        *string `json:"imageUrl,omitempty"`
	Rank            *int    `json:"rank,omitempty"`
}

// Category represents a user's category with a score for prioritization.
type CategoryDto struct {
	CategoryID  string  `json:"categoryId"`
	Name        string  `json:"name"`
	Description string  `json:"description,omitempty"`
	ImageUrl    *string `json:"imageUrl,omitempty"`
	Rank        *int    `json:"rank,omitempty"`
}

// Merchant represents a merchant for API responses.
type MerchantDto struct {
	MerchantID  string `json:"merchantId"`
	GroupID     string `json:"groupId,omitempty"`
	UserID      string `json:"userId,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	ImageUrl    string `json:"imageUrl,omitempty"`
	Rank        *int   `json:"rank,omitempty"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
	DeletedAt   string `json:"deletedAt,omitempty"`
}
