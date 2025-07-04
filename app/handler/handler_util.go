package handler

import (
	"net/url"
	"strings"
)

// ParseQueryStringArray parses query parameters that can be provided in multiple formats:
// 1. Comma-separated values: ?param=value1,value2,value3
// 2. Multiple parameters: ?param=value1&param=value2&param=value3
// 3. Single value: ?param=value1
//
// Returns a slice of trimmed, non-empty strings.
//
// Example usage:
//
//	types := ParseQueryStringArray(query, "type")
//	categories := ParseQueryStringArray(query, "categoryId")
func ParseQueryStringArray(query url.Values, paramName string) []string {
	paramValues := query[paramName]
	if len(paramValues) == 0 {
		return nil
	}

	var result []string
	for _, paramValue := range paramValues {
		// Split by comma to handle comma-separated values
		if strings.Contains(paramValue, ",") {
			commaSeparated := strings.Split(paramValue, ",")
			for _, item := range commaSeparated {
				trimmed := strings.TrimSpace(item)
				if trimmed != "" {
					result = append(result, trimmed)
				}
			}
		} else {
			// Handle single values
			trimmed := strings.TrimSpace(paramValue)
			if trimmed != "" {
				result = append(result, trimmed)
			}
		}
	}

	return result
}
