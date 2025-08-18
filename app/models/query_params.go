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
	BalanceIds       []string
	CategoryIds      []string
	CategoryGroupIds []string
	MerchantIds      []string
	TransactionIds   []string
	OperationIds     []string
	Types            []string
	StartTime        time.Time
	EndTime          time.Time
	IncludeDeleted   bool
	SortBy           string
	Order            string
	Limit            int
}

// ListBalancesInput defines the filter and pagination options for list of balances
type ListBalancesInput struct {
	GroupID        string
	UserID         string
	IncludeDeleted bool
	SortBy         string
	Order          string
	Limit          int
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

const (
	GroupingCategoryGroup string = "categoryGroup"
	GroupingCategory      string = "category"
	GroupingMerchant      string = "merchant"
	GroupingBalance       string = "balance"
	GroupingCurrency      string = "currency"
	GroupingQuarter       string = "quarter"
	GroupingYear          string = "year"
	GroupingMonth         string = "month"
	GroupingWeek          string = "week"
	GroupingDay           string = "day"
)

// TransactionStatsInput defines the filter options for transaction statistics
type TransactionStatsInput struct {
	Type            string
	GroupID         string
	UserID          string
	BalanceID       []string
	CategoryId      []string
	CategoryGroupId []string
	MerchantId      []string
	TransactionId   []string
	Grouping        string
	Limit           int
	Sort            string
	Order           string
	DisplayCurrency string
	StartTime       time.Time
	EndTime         *time.Time
}
