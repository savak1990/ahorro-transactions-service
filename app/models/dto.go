package models

// CreateTransactionDto represents a financial transaction for creation via POST requests.
type CreateTransactionDto struct {
	TransactionID      string                      `json:"transactionId,omitempty"`
	GroupID            string                      `json:"groupId"`
	UserID             string                      `json:"userId"`
	BalanceID          string                      `json:"balanceId"`
	Type               string                      `json:"type"` // Supported: init, income, expense, movement, move_in, move_out
	MerchantID         string                      `json:"merchantId,omitempty"`
	OperationID        string                      `json:"operationId,omitempty"`
	ApprovedAt         string                      `json:"approvedAt,omitempty"`
	TransactedAt       string                      `json:"transactedAt"`
	CreatedAt          string                      `json:"createdAt,omitempty"`
	UpdatedAt          string                      `json:"updatedAt,omitempty"`
	DeletedAt          string                      `json:"deletedAt,omitempty"`
	TransactionEntries []CreateTransactionEntryDto `json:"transactionEntries"`
}

// CreateTransactionsRequestDto represents a batch transaction creation request
type CreateTransactionsRequestDto struct {
	Transactions []CreateTransactionDto `json:"transactions"`
}

// CreateTransactionsResponseDto represents a batch transaction creation response
type CreateTransactionsResponseDto struct {
	Transactions []CreateTransactionDto `json:"transactions"`
	OperationID  *string                `json:"operationId,omitempty"`
}

// CreateTransactionEntryDto represents a single entry within a transaction for creation
type CreateTransactionEntryDto struct {
	ID          string `json:"id,omitempty"`
	Description string `json:"description"`
	Amount      int    `json:"amount"`
	CategoryID  string `json:"categoryId"`
	CreatedAt   string `json:"createdAt,omitempty"`
	UpdatedAt   string `json:"updatedAt,omitempty"`
	DeletedAt   string `json:"deletedAt,omitempty"`
}

// TransactionDto represents a financial transaction for API responses.
type TransactionEntryDto struct {
	GroupID               string  `json:"groupId"`
	UserID                string  `json:"userId"`
	BalanceID             string  `json:"balanceId"`
	TransactionID         string  `json:"transactionId"`
	TransactionEntryID    string  `json:"transactionEntryId"`
	Type                  string  `json:"type"` // Supported: init, income, expense, movement, move_in, move_out
	Amount                int     `json:"amount"`
	BalanceTitle          string  `json:"balanceTitle"`
	BalanceCurrency       string  `json:"balanceCurrency"`
	BalanceDeleted        bool    `json:"balanceDeleted,omitempty"`
	CategoryID            string  `json:"categoryId,omitempty"`
	CategoryName          string  `json:"categoryName"`
	CategoryImageUrl      string  `json:"categoryImageUrl,omitempty"`
	CategoryGroupName     string  `json:"categoryGroupName,omitempty"`
	CategoryGroupImageUrl *string `json:"cattegoryGroupImageUrl,omitempty"`
	CategoryGroupID       string  `json:"categoryGroupId,omitempty"`
	CategoryIsDeleted     bool    `json:"categoryIsDeleted,omitempty"`
	CategoryGroupDeleted  bool    `json:"categoryGroupDeleted,omitempty"`
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
	IsDeleted       bool    `json:"isDeleted,omitempty"`
}

// CreateCategoryDto represents a category for creation via POST requests.
type CreateCategoryDto struct {
	GroupID         string  `json:"groupId"`
	UserID          string  `json:"userId"`
	CategoryGroupID string  `json:"categoryGroupId"`
	Name            string  `json:"name"`
	Description     string  `json:"description,omitempty"`
	ImageUrl        *string `json:"imageUrl,omitempty"`
	Rank            *int    `json:"rank,omitempty"`
}

// UpdateCategoryDto represents a category for update via PUT requests.
type UpdateCategoryDto struct {
	GroupID         string  `json:"groupId"`
	UserID          string  `json:"userId"`
	CategoryGroupID string  `json:"categoryGroupId,omitempty"`
	Name            string  `json:"name,omitempty"`
	Description     string  `json:"description,omitempty"`
	ImageUrl        *string `json:"imageUrl,omitempty"`
	Rank            *int    `json:"rank,omitempty"`
}

// Category represents a user's category with a score for prioritization.
type CategoryDto struct {
	CategoryID            string  `json:"categoryId"`
	Name                  string  `json:"name"`
	Description           string  `json:"description,omitempty"`
	ImageUrl              *string `json:"imageUrl,omitempty"`
	Rank                  *int    `json:"rank,omitempty"`
	CategoryGroupID       string  `json:"categoryGroupId,omitempty"`
	CategoryGroupName     string  `json:"categoryGroupName,omitempty"`
	CategoryGroupImageUrl *string `json:"categoryGroupImageUrl,omitempty"`
	CategoryGroupRank     *int    `json:"categoryGroupRank,omitempty"`
	IsDeleted             bool    `json:"isDeleted,omitempty"`
	CategoryGroupDeleted  bool    `json:"categoryGroupDeleted,omitempty"`
}

// CategoryGroupWithCategoriesDto represents a category group with its associated categories
type CategoryGroupWithCategoriesDto struct {
	CategoryGroupID       string        `json:"categoryGroupId"`
	CategoryGroupName     string        `json:"categoryGroupName"`
	CategoryGroupImageUrl *string       `json:"categoryGroupImageUrl,omitempty"`
	CategoryGroupRank     *int          `json:"categoryGroupRank,omitempty"`
	Categories            []CategoryDto `json:"categories"`
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

// TransactionStatsDto represents aggregated transaction statistics by currency
type CurrencyStatsDto struct {
	Amount                  int `json:"amount"`                  // Total amount in cents
	TransactionsCount       int `json:"transactionsCount"`       // Number of transactions
	TransactionEntriesCount int `json:"transactionEntriesCount"` // Number of transaction entries
}

// TransactionStatsResponseDto represents the response for transaction statistics
type TransactionStatsResponseDto struct {
	Totals map[string]map[string]CurrencyStatsDto `json:"totals"` // Stats by transaction type, then by currency
}

// UpdateTransactionDto represents a financial transaction for update via PUT requests.
type UpdateTransactionDto struct {
	TransactionID      string                      `json:"transactionId,omitempty"`
	GroupID            string                      `json:"groupId"`
	UserID             string                      `json:"userId"`
	BalanceID          string                      `json:"balanceId"`
	Type               string                      `json:"type"` // Supported: init, income, expense, movement, move_in, move_out
	MerchantID         string                      `json:"merchantId,omitempty"`
	OperationID        string                      `json:"operationId,omitempty"`
	ApprovedAt         string                      `json:"approvedAt,omitempty"`
	TransactedAt       string                      `json:"transactedAt"`
	CreatedAt          string                      `json:"createdAt,omitempty"`
	UpdatedAt          string                      `json:"updatedAt,omitempty"`
	DeletedAt          string                      `json:"deletedAt,omitempty"`
	TransactionEntries []CreateTransactionEntryDto `json:"transactionEntries"`
}
