package models

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestToAPITransactionEntry(t *testing.T) {
	// Create test data
	transactionID := uuid.New()
	entryID := uuid.New()
	categoryID := uuid.New()
	groupID := uuid.New()
	userID := uuid.New()
	balanceID := uuid.New()
	merchantID := uuid.New()

	// Create a test amount of $45.50 (which should be 4550 cents)
	amount := int64(4550) // Amount in cents

	// Create test category with group
	category := &Category{
		ID:          categoryID,
		Name:        "Groceries",
		Group:       "Food & Dining",
		ImageUrl:    stringPtr("https://example.com/groceries.png"),
		Description: "Food and grocery purchases",
	}

	// Create test merchant
	merchant := &Merchant{
		ID:       merchantID,
		Name:     "Mercadona",
		ImageUrl: stringPtr("https://example.com/mercadona.png"),
	}

	// Create test balance
	balance := &Balance{
		ID:       balanceID,
		Title:    "Main Checking",
		Currency: "EUR",
	}

	// Create test transaction
	transaction := &Transaction{
		ID:           transactionID,
		GroupID:      groupID,
		UserID:       userID,
		BalanceID:    balanceID,
		MerchantID:   &merchantID,
		Type:         "expense",
		ApprovedAt:   time.Now(),
		TransactedAt: time.Now(),
		Merchant:     merchant,
		Balance:      balance,
	}

	// Create test transaction entry
	entry := &TransactionEntry{
		ID:            entryID,
		TransactionID: transactionID,
		Amount:        amount,
		CategoryID:    &categoryID,
		Transaction:   transaction,
		Category:      category,
	}

	// Convert to API model
	result := ToAPITransactionEntry(entry)

	// Verify the conversion
	if result.Amount != 4550 {
		t.Errorf("Expected amount to be 4550 cents, got %d", result.Amount)
	}

	if result.CategoryName != "Groceries" {
		t.Errorf("Expected categoryName to be 'Groceries', got '%s'", result.CategoryName)
	}

	if result.CategoryGroupName != "Food & Dining" {
		t.Errorf("Expected categoryGroupName to be 'Food & Dining', got '%s'", result.CategoryGroupName)
	}

	if result.CategoryImageUrl != "https://example.com/groceries.png" {
		t.Errorf("Expected categoryImageUrl to be 'https://example.com/groceries.png', got '%s'", result.CategoryImageUrl)
	}

	if result.MerchantName != "Mercadona" {
		t.Errorf("Expected merchantName to be 'Mercadona', got '%s'", result.MerchantName)
	}

	if result.BalanceTitle != "Main Checking" {
		t.Errorf("Expected balanceTitle to be 'Main Checking', got '%s'", result.BalanceTitle)
	}

	if result.BalanceCurrency != "EUR" {
		t.Errorf("Expected balanceCurrency to be 'EUR', got '%s'", result.BalanceCurrency)
	}

	if result.Type != "expense" {
		t.Errorf("Expected type to be 'expense', got '%s'", result.Type)
	}

	// Test CategoryGroupImageUrl is properly set as pointer
	if result.CategoryGroupImageUrl == nil {
		t.Error("Expected categoryGroupImageUrl to not be nil")
	} else if *result.CategoryGroupImageUrl != "" {
		t.Errorf("Expected categoryGroupImageUrl to be empty string, got '%s'", *result.CategoryGroupImageUrl)
	}

	// Test categoryIsDeleted flag for non-deleted category
	if result.CategoryIsDeleted {
		t.Error("Expected categoryIsDeleted to be false for non-deleted category")
	}
}

func TestToAPITransactionEntry_ZeroAmount(t *testing.T) {
	// Test with zero amount
	entry := &TransactionEntry{
		ID:     uuid.New(),
		Amount: int64(0), // Zero amount in cents
	}

	result := ToAPITransactionEntry(entry)

	if result.Amount != 0 {
		t.Errorf("Expected amount to be 0 cents, got %d", result.Amount)
	}
}

func TestToAPITransactionEntry_NilCategory(t *testing.T) {
	// Test with nil category
	entry := &TransactionEntry{
		ID:       uuid.New(),
		Amount:   int64(1099), // $10.99 in cents
		Category: nil,
	}

	result := ToAPITransactionEntry(entry)

	if result.CategoryName != "" {
		t.Errorf("Expected categoryName to be empty, got '%s'", result.CategoryName)
	}

	if result.CategoryGroupName != "" {
		t.Errorf("Expected categoryGroupName to be empty, got '%s'", result.CategoryGroupName)
	}

	// Amount should still be converted correctly
	if result.Amount != 1099 {
		t.Errorf("Expected amount to be 1099 cents, got %d", result.Amount)
	}

	// Test categoryIsDeleted flag for nil category
	if result.CategoryIsDeleted {
		t.Error("Expected categoryIsDeleted to be false when category is nil")
	}
}

func TestToAPITransactionEntry_SoftDeletedCategory(t *testing.T) {
	// Test with soft-deleted category (deleted_at is not nil)
	now := time.Now()
	category := &Category{
		ID:        uuid.New(),
		Name:      "Deleted Category",
		DeletedAt: &now, // Category is soft deleted
	}

	entry := &TransactionEntry{
		ID:       uuid.New(),
		Amount:   int64(2500), // $25.00 in cents
		Category: category,
	}

	result := ToAPITransactionEntry(entry)

	if !result.CategoryIsDeleted {
		t.Error("Expected categoryIsDeleted to be true for soft-deleted category")
	}

	if result.CategoryName != "Deleted Category" {
		t.Errorf("Expected categoryName to be 'Deleted Category', got '%s'", result.CategoryName)
	}

	// Amount should still be converted correctly
	if result.Amount != 2500 {
		t.Errorf("Expected amount to be 2500 cents, got %d", result.Amount)
	}
}

func TestToAPITransactionEntry_SoftDeletedCategoryGroup(t *testing.T) {
	// Test with soft-deleted category group (deleted_at is not nil)
	now := time.Now()
	categoryGroup := &CategoryGroup{
		ID:        uuid.New(),
		Name:      "Deleted Group",
		DeletedAt: &now, // Category group is soft deleted
	}

	category := &Category{
		ID:              uuid.New(),
		Name:            "Active Category",
		CategoryGroup:   categoryGroup,
		CategoryGroupId: categoryGroup.ID.String(),
	}

	entry := &TransactionEntry{
		ID:       uuid.New(),
		Amount:   int64(3000), // $30.00 in cents
		Category: category,
	}

	result := ToAPITransactionEntry(entry)

	if result.CategoryIsDeleted {
		t.Error("Expected categoryIsDeleted to be false for active category")
	}

	if !result.CategoryGroupDeleted {
		t.Error("Expected categoryGroupDeleted to be true for soft-deleted category group")
	}

	if result.CategoryName != "Active Category" {
		t.Errorf("Expected categoryName to be 'Active Category', got '%s'", result.CategoryName)
	}

	// Amount should still be converted correctly
	if result.Amount != 3000 {
		t.Errorf("Expected amount to be 3000 cents, got %d", result.Amount)
	}
}

// Helper function to create string pointers
func stringPtr(s string) *string {
	return &s
}
