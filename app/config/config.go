package config

import "os"

type AppConfig struct {
	AWSRegion                 string
	AWSProfile                string
	CategoriesDynamoDbTable   string
	TransactionsDynamoDbTable string
}

func LoadConfig() AppConfig {
	return AppConfig{
		AWSRegion:                 getEnv("AWS_REGION", "eu-west-1"),
		AWSProfile:                os.Getenv("AWS_PROFILE"),
		CategoriesDynamoDbTable:   os.Getenv("CATEGORIES_DYNAMODB_TABLE"),
		TransactionsDynamoDbTable: os.Getenv("TRANSACTIONS_DYNAMODB_TABLE"),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
