package models

// Transaction represents a financial transaction in the system.
type Transaction struct {
	TransactionID string  `json:"transaction_id" dynamodbav:"transaction_id"`
	UserID        string  `json:"user_id" dynamodbav:"user_id"`
	GroupID       string  `json:"group_id" dynamodbav:"group_id"`
	Type          string  `json:"type" dynamodbav:"type"`
	Amount        float64 `json:"amount" dynamodbav:"amount"`
	BalanceID     string  `json:"balance_id" dynamodbav:"balance_id"`
	FromBalanceID string  `json:"from_balance_id" dynamodbav:"from_balance_id"`
	ToBalanceID   string  `json:"to_balance_id,omitempty" dynamodbav:"to_balance_id"`
	Category      string  `json:"category" dynamodbav:"category"`
	Description   string  `json:"description" dynamodbav:"description"`
	ApprovedAt    string  `json:"approved_at,omitempty" dynamodbav:"approved_at"`
	TransactedAt  string  `json:"transacted_at" dynamodbav:"transacted_at"`
	CreatedAt     string  `json:"created_at" dynamodbav:"created_at"`
	UpdatedAt     string  `json:"updated_at" dynamodbav:"updated_at"`
	DeletedAt     string  `json:"deleted_at,omitempty" dynamodbav:"deleted_at"`
}

// Category represents a user's category with a score for prioritization.
type Category struct {
	Name     string `json:"name" dynamodbav:"name"`
	ImageUrl string `json:"image_url,omitempty" dynamodbav:"image_url"`
	Score    int    `json:"-" dynamodbav:"score"`
}

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
	UserID   string
	Type     string
	Category string
	SortBy   string
	Order    string
	Count    int
	StartKey string
}
