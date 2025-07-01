package models

// ListCategoriesInput defines the input options for listing categories.
type ListCategoriesInput struct {
	UserID string
	Limit  int
}

// ListTransactionsInput defines the filter and pagination options for listing transactions.
type ListTransactionsInput struct {
	GroupID    string
	UserID     string
	BalanceID  string
	Type       string
	CategoryId string
	MerchantId string
	SortBy     string
	Order      string
	Count      int
}

// ListBalancesInput defines the filter and pagination options for list of balances
type ListBalancesInput struct {
	GroupID   string
	UserID    string
	BalanceID string
	SortBy    string
	Order     string
	Count     int
}

// ListMerchantsInput defines the filter and pagination options for list of merchants
type ListMerchantsInput struct {
	Name   string
	SortBy string
	Order  string
	Count  int
}
