package repo

import (
	"context"
	"encoding/json"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/mock"
)

type ExchangeRatesDb interface {
	GetExchangeRates(ctx context.Context, baseCurrency string, date ...time.Time) (map[string]float64, error)
	GetSupportedCurrencies(ctx context.Context) ([]string, error)
}

type ExchangeRatesDbImpl struct {
	dbName         string
	dbClient       *dynamodb.Client
	defaultTimeout time.Duration
}

var _ ExchangeRatesDb = (*ExchangeRatesDbImpl)(nil)

func NewExchangeRatesDb(dbName string, awsConfig aws.Config) ExchangeRatesDb {
	return &ExchangeRatesDbImpl{
		dbName:         dbName,
		dbClient:       dynamodb.NewFromConfig(awsConfig),
		defaultTimeout: 2 * time.Second,
	}
}

// GetExchangeRates retrieves exchange rates for the given base currency from DynamoDB
func (db *ExchangeRatesDbImpl) GetExchangeRates(ctx context.Context, baseCurrency string, date ...time.Time) (map[string]float64, error) {
	// Create a timeout context if none provided
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), db.defaultTimeout)
		defer cancel()
	}

	// Determine the date to use
	var targetDate string
	if len(date) > 0 {
		targetDate = date[0].Format("2006-01-02")
	} else {
		targetDate = time.Now().Format("2006-01-02")
	}

	input := &dynamodb.GetItemInput{
		TableName: aws.String(db.dbName),
		Key: map[string]types.AttributeValue{
			"Key": &types.AttributeValueMemberS{
				Value: baseCurrency,
			},
			"SortKey": &types.AttributeValueMemberS{
				Value: targetDate,
			},
		},
	}

	result, err := db.dbClient.GetItem(ctx, input)
	if err != nil {
		return nil, err
	}

	// Check if item was found
	if result.Item == nil {
		return nil, nil // Item not found, return nil
	}

	// Extract and decode the exchange rates from the ExchangeRates attribute
	if ratesAttr, exists := result.Item["ExchangeRates"]; exists {
		if mapVal, ok := ratesAttr.(*types.AttributeValueMemberM); ok {
			rates := make(map[string]float64)
			for currency, value := range mapVal.Value {
				if numVal, ok := value.(*types.AttributeValueMemberN); ok {
					if rate, err := parseFloat64(numVal.Value); err == nil {
						rates[currency] = rate
					}
				}
			}
			return rates, nil
		}
	}

	return nil, nil
}

// GetSupportedCurrencies retrieves the list of supported currencies from DynamoDB
func (db *ExchangeRatesDbImpl) GetSupportedCurrencies(ctx context.Context) ([]string, error) {
	// Create a timeout context if none provided
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), db.defaultTimeout)
		defer cancel()
	}

	input := &dynamodb.GetItemInput{
		TableName: aws.String(db.dbName),
		Key: map[string]types.AttributeValue{
			"Key": &types.AttributeValueMemberS{
				Value: "SupportedCurrencies",
			},
			"SortKey": &types.AttributeValueMemberS{
				Value: "-",
			},
		},
	}

	result, err := db.dbClient.GetItem(ctx, input)
	if err != nil {
		return nil, err
	}

	// Check if item was found
	if result.Item == nil {
		return nil, nil // Item not found, return nil
	}

	// Extract and decode the supported currencies from the SupportedCurrencies attribute
	if currenciesAttr, exists := result.Item["SupportedCurrencies"]; exists {
		if listVal, ok := currenciesAttr.(*types.AttributeValueMemberL); ok {
			currencies := make([]string, 0, len(listVal.Value))
			for _, item := range listVal.Value {
				if strVal, ok := item.(*types.AttributeValueMemberS); ok {
					currencies = append(currencies, strVal.Value)
				}
			}
			return currencies, nil
		}
	}

	return nil, nil
}

// parseFloat64 safely parses a string to float64
func parseFloat64(s string) (float64, error) {
	var f float64
	err := json.Unmarshal([]byte(s), &f)
	return f, err
}

// ExchangeRatesStaticImpl provides hardcoded exchange rates for testing
type ExchangeRatesStaticImpl struct{}

var _ ExchangeRatesDb = (*ExchangeRatesStaticImpl)(nil)

func NewExchangeRatesStaticDb() ExchangeRatesDb {
	return &ExchangeRatesStaticImpl{}
}

// GetExchangeRates returns hardcoded exchange rates for all supported currencies
func (db *ExchangeRatesStaticImpl) GetExchangeRates(ctx context.Context, baseCurrency string, date ...time.Time) (map[string]float64, error) {
	// Hardcoded exchange rates (approximate values as of 2025)
	// Note: Static implementation ignores the date parameter since rates are hardcoded
	rates := getStaticRates(baseCurrency)
	if rates == nil {
		return nil, nil // Currency not supported
	}
	return rates, nil
}

// GetSupportedCurrencies returns hardcoded list of supported currencies
func (db *ExchangeRatesStaticImpl) GetSupportedCurrencies(ctx context.Context) ([]string, error) {
	// Return all currencies that have static rates defined
	return []string{"USD", "EUR", "GBP", "CHF", "JPY"}, nil
}

// getStaticRates returns predefined exchange rates for supported currencies
func getStaticRates(baseCurrency string) map[string]float64 {
	allRates := map[string]map[string]float64{
		"USD": {
			"EUR": 0.85, "GBP": 0.73, "CHF": 0.88, "SEK": 10.20, "NOK": 10.50,
			"DKK": 6.35, "PLN": 3.95, "CZK": 22.50, "HUF": 350.00, "RON": 4.20,
			"UAH": 37.00, "BYN": 2.60, "RUB": 75.00, "JPY": 148.00, "CAD": 1.35,
			"AUD": 1.48, "CNY": 7.20,
		},
		"EUR": {
			"USD": 1.18, "GBP": 0.86, "CHF": 1.04, "SEK": 12.00, "NOK": 12.35,
			"DKK": 7.46, "PLN": 4.65, "CZK": 26.50, "HUF": 412.00, "RON": 4.95,
			"UAH": 43.50, "BYN": 3.06, "RUB": 88.50, "JPY": 174.00, "CAD": 1.59,
			"AUD": 1.74, "CNY": 8.48,
		},
		"GBP": {
			"USD": 1.37, "EUR": 1.16, "CHF": 1.21, "SEK": 13.95, "NOK": 14.35,
			"DKK": 8.67, "PLN": 5.41, "CZK": 30.85, "HUF": 479.00, "RON": 5.76,
			"UAH": 50.65, "BYN": 3.56, "RUB": 103.00, "JPY": 202.50, "CAD": 1.85,
			"AUD": 2.03, "CNY": 9.87,
		},
		"CHF": {
			"USD": 1.14, "EUR": 0.96, "GBP": 0.83, "SEK": 11.59, "NOK": 11.93,
			"DKK": 7.20, "PLN": 4.49, "CZK": 25.61, "HUF": 398.00, "RON": 4.78,
			"UAH": 42.05, "BYN": 2.95, "RUB": 85.50, "JPY": 168.50, "CAD": 1.54,
			"AUD": 1.68, "CNY": 8.21,
		},
		"JPY": {
			"USD": 0.0068, "EUR": 0.0057, "GBP": 0.0049, "CHF": 0.0059, "SEK": 0.069,
			"NOK": 0.071, "DKK": 0.043, "PLN": 0.027, "CZK": 0.152, "HUF": 2.36,
			"RON": 0.028, "UAH": 0.25, "BYN": 0.018, "RUB": 0.51, "CAD": 0.0091,
			"AUD": 0.010, "CNY": 0.049,
		},
	}

	if rates, exists := allRates[baseCurrency]; exists {
		return rates
	}
	return nil
}

// ExchangeRatesDbMock provides a mock implementation for testing
type ExchangeRatesDbMock struct {
	mock.Mock
}

var _ ExchangeRatesDb = (*ExchangeRatesDbMock)(nil)

func NewExchangeRatesDbMock() ExchangeRatesDb {
	return &ExchangeRatesDbMock{}
}

// GetExchangeRates mocks the GetExchangeRates method
func (db *ExchangeRatesDbMock) GetExchangeRates(ctx context.Context, baseCurrency string, date ...time.Time) (map[string]float64, error) {
	args := db.Called(ctx, baseCurrency, date)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]float64), args.Error(1)
}

// GetSupportedCurrencies mocks the GetSupportedCurrencies method
func (db *ExchangeRatesDbMock) GetSupportedCurrencies(ctx context.Context) ([]string, error) {
	args := db.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

// Ensure ExchangeRatesDbMock implements ExchangeRatesDb interface
var _ ExchangeRatesDb = (*ExchangeRatesDbMock)(nil)
