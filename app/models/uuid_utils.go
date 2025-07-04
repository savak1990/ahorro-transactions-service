package models

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

// Entity type prefixes using 2-character hex prefixes (0-9, a-f)
const (
	PrefixBalance          = "ba" // balance
	PrefixCategory         = "ca" // category
	PrefixCategoryGroup    = "c9" // category group
	PrefixMerchant         = "4e" // merchant
	PrefixTransaction      = "7a" // transaction
	PrefixTransactionEntry = "7e" // transaction entry
)

// GenerateUUIDWithPrefix creates a UUID with the specified 2-character hex prefix
func GenerateUUIDWithPrefix(prefix string) uuid.UUID {
	// Generate a new UUID
	id := uuid.New()

	// Get the original UUID as hex string (without hyphens)
	idStr := id.String()
	idHex := strings.ReplaceAll(idStr, "-", "")

	// Ensure prefix is exactly 2 hex characters
	if len(prefix) != 2 {
		// Fallback to original UUID if prefix is invalid
		return id
	}

	// Replace the first 2 hex characters with our prefix
	modifiedHex := prefix + idHex[2:]

	// Rebuild UUID string with hyphens
	uuidStr := fmt.Sprintf("%s-%s-%s-%s-%s",
		modifiedHex[0:8],
		modifiedHex[8:12],
		modifiedHex[12:16],
		modifiedHex[16:20],
		modifiedHex[20:32],
	)

	// Parse back to UUID
	result, err := uuid.Parse(uuidStr)
	if err != nil {
		// Fallback to original UUID if parsing fails
		return id
	}

	return result
}

// Helper functions for each entity type
func NewBalanceID() uuid.UUID {
	return GenerateUUIDWithPrefix(PrefixBalance)
}

func NewCategoryID() uuid.UUID {
	return GenerateUUIDWithPrefix(PrefixCategory)
}

func NewCategoryGroupID() uuid.UUID {
	return GenerateUUIDWithPrefix(PrefixCategoryGroup)
}

func NewMerchantID() uuid.UUID {
	return GenerateUUIDWithPrefix(PrefixMerchant)
}

func NewTransactionID() uuid.UUID {
	return GenerateUUIDWithPrefix(PrefixTransaction)
}

func NewTransactionEntryID() uuid.UUID {
	return GenerateUUIDWithPrefix(PrefixTransactionEntry)
}

// GetEntityTypeFromUUID extracts the entity type from a UUID by examining its 2-character hex prefix
func GetEntityTypeFromUUID(id uuid.UUID) string {
	idStr := strings.ReplaceAll(id.String(), "-", "")

	// Check the first 2 hex characters against known prefixes
	if len(idStr) < 2 {
		return "Unknown"
	}

	prefix := idStr[:2]

	switch prefix {
	case PrefixBalance:
		return "Balance"
	case PrefixCategory:
		return "Category"
	case PrefixCategoryGroup:
		return "CategoryGroup"
	case PrefixMerchant:
		return "Merchant"
	case PrefixTransaction:
		return "Transaction"
	case PrefixTransactionEntry:
		return "TransactionEntry"
	default:
		return "Unknown"
	}
}
