package helpers

import (
	"strconv"
	"strings"
)

// parseInt is a helper for parsing integers from query params
func ParseInt(s string) (int, error) {
	return strconv.Atoi(strings.TrimSpace(s))
}
