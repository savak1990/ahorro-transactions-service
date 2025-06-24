package models

const (
	ErrorCodeBadRequest     = "BadRequest"
	ErrorCodeInternalServer = "InternalServerError"
	ErrorCodeNotFound       = "NotFound"
	ErrorCodeBadResponse    = "BadResponse"
)

// ListCategoriesInput defines the input options for listing categories.
type ListCategoriesInput struct {
	UserID   string
	Limit    int
	StartKey string
}

// ListTransactionsFilter defines the filter and pagination options for listing transactions.
type ListTransactionsFilter struct {
	GroupID   string
	UserID    string
	BalanceID string
	Type      string
	Category  string
	SortBy    string
	Order     string
	Count     int
	StartKey  string
}

// ListBalancesFilter defines the filter and pagination options for list of balances
type ListBalancesFilter struct {
	GroupID   string
	UserID    string
	BalanceID string
	SortBy    string
	Order     string
	Count     int
	StartKey  string
}
