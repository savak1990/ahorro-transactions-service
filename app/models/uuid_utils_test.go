package models

import (
	"strings"
	"testing"
)

func TestUUIDHexPrefixes(t *testing.T) {
	// Test all entity type UUID generators
	tests := []struct {
		name           string
		generator      func() string
		expectedPrefix string
		entityType     string
	}{
		{"Balance", func() string { return NewBalanceID().String() }, "ba", "Balance"},
		{"Category", func() string { return NewCategoryID().String() }, "ca", "Category"},
		{"CategoryGroup", func() string { return NewCategoryGroupID().String() }, "c9", "CategoryGroup"},
		{"Merchant", func() string { return NewMerchantID().String() }, "4e", "Merchant"},
		{"Transaction", func() string { return NewTransactionID().String() }, "7a", "Transaction"},
		{"TransactionEntry", func() string { return NewTransactionEntryID().String() }, "7e", "TransactionEntry"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Generate UUID
			uuidStr := tt.generator()

			// Remove hyphens to check prefix
			uuidHex := strings.ReplaceAll(uuidStr, "-", "")

			// Check if it starts with expected 2-character hex prefix
			if !strings.HasPrefix(uuidHex, tt.expectedPrefix) {
				t.Errorf("Expected UUID to start with '%s', but got: %s (first 2 chars: %s)",
					tt.expectedPrefix, uuidStr, uuidHex[:2])
			}

			// Verify it's a valid UUID format (36 characters with hyphens)
			if len(uuidStr) != 36 {
				t.Errorf("Expected UUID to be 36 characters long, got %d: %s", len(uuidStr), uuidStr)
			}

			// Test entity type detection
			uuid, err := parseUUID(uuidStr)
			if err != nil {
				t.Fatalf("Failed to parse generated UUID: %v", err)
			}

			entityType := GetEntityTypeFromUUID(uuid)
			if entityType != tt.entityType {
				t.Errorf("Expected entity type '%s', got '%s' for UUID: %s",
					tt.entityType, entityType, uuidStr)
			}
		})
	}
}

func TestUUIDExample(t *testing.T) {
	// Generate some examples to show what the UUIDs look like
	examples := map[string]func() string{
		"Balance":          func() string { return NewBalanceID().String() },
		"Category":         func() string { return NewCategoryID().String() },
		"CategoryGroup":    func() string { return NewCategoryGroupID().String() },
		"Merchant":         func() string { return NewMerchantID().String() },
		"Transaction":      func() string { return NewTransactionID().String() },
		"TransactionEntry": func() string { return NewTransactionEntryID().String() },
	}

	t.Log("Example UUIDs with 2-character hex prefixes:")
	for name, generator := range examples {
		uuid := generator()
		t.Logf("  %s: %s", name, uuid)
	}
}
