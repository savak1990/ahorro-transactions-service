package models

import "time"

// ListCategoryGroupsInput defines the input for only category groups
type ListCategoryGroupsInput struct {
	Limit   int
	GroupBy string
	SortBy  string
	Order   string
}

// ListCategoriesInput defines the input options for listing categories.
type ListCategoriesInput struct {
	GroupID string
	UserID  string
	Limit   int
	GroupBy string
	SortBy  string
	Order   string
}

// ListTransactionsInput defines the filter and pagination options for listing transactions.
type ListTransactionsInput struct {
	GroupID          string
	UserID           string
	BalanceID        string
	CategoryIds      []string
	CategoryGroupIds []string
	MerchantId       string
	TransactionID    string
	Types            []string // Changed from Type string to Types []string to support multiple types
	StartTime        time.Time
	EndTime          time.Time
	SortBy           string
	Order            string
	Limit            int
}

// ListBalancesInput defines the filter and pagination options for list of balances
type ListBalancesInput struct {
	GroupID   string
	UserID    string
	BalanceID string
	SortBy    string
	Order     string
	Limit     int
}

// ListMerchantsInput defines the filter and pagination options for list of merchants
type ListMerchantsInput struct {
	GroupID    string
	UserID     string
	MerchantID string
	Name       string
	SortBy     string
	Order      string
	Limit      int
}

// TransactionStatsInput defines the filter options for transaction statistics
type TransactionStatsInput struct {
	GroupID         *string
	UserID          *string
	BalanceID       *string
	CategoryId      *string
	CategoryGroupId *string
	MerchantId      *string
	Type            *string // "expense", "income", or empty for both
	Currency        *string // Future feature, not implemented yet
	StartTime       *time.Time
	EndTime         *time.Time
}
