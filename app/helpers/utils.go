package helpers

import (
	"strconv"
	"strings"
)

// parseInt is a helper for parsing integers from query params
func ParseInt(s string) (int, error) {
	return strconv.Atoi(strings.TrimSpace(s))
}

var supportedCurrencies = map[string]struct{}{
	"usd": {}, "eur": {}, "jpy": {}, "gbp": {}, "aud": {},
	"cad": {}, "chf": {}, "cny": {}, "hkd": {}, "nzd": {},
	"sek": {}, "krw": {}, "sgd": {}, "nok": {}, "mxn": {},
	"inr": {}, "rub": {}, "zar": {}, "try": {}, "brl": {},
	"uah": {}, // Ukrainian Hryvnia
	"byn": {}, // Belarusian Ruble
}

func IsValidCurrencyCode(code string) bool {
	code = strings.ToLower(strings.TrimSpace(code))
	_, ok := supportedCurrencies[code]
	return ok
}
