package models

// Transaction represents a financial transaction in the system.
type TransactionDto struct {
	TransactionID string  `json:"transaction_id"`
	GroupID       string  `json:"group_id"`
	UserID        string  `json:"user_id"`
	BalanceID     string  `json:"balance_id"`
	Type          string  `json:"type"`
	Amount        float64 `json:"amount"`
	FromBalanceID string  `json:"from_balance_id"`
	ToBalanceID   string  `json:"to_balance_id,omitempty"`
	Category      string  `json:"category"`
	Description   string  `json:"description"`
	ApprovedAt    string  `json:"approved_at,omitempty"`
	TransactedAt  string  `json:"transacted_at"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`
	DeletedAt     string  `json:"deleted_at,omitempty"`
}

// Balance represents a user's balance/account for API responses.
type BalanceDto struct {
	BalanceID   string `json:"balance_id"`
	GroupID     string `json:"group_id"`
	UserID      string `json:"user_id"`
	Currency    string `json:"currency"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	DeletedAt   string `json:"deleted_at,omitempty"`
}

// Category represents a user's category with a score for prioritization.
type CategoryDto struct {
	Name     string `json:"name"`
	ImageUrl string `json:"image_url,omitempty"`
	Score    int    `json:"-"`
}
