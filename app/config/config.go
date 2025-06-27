package config

import (
	"os"
	"strconv"
)

type AppConfig struct {
	AWSRegion  string
	AWSProfile string

	// Aurora PostgreSQL Database Configuration
	DBHost     string
	DBPort     int
	DBName     string
	DBUser     string
	DBPassword string
}

func LoadConfig() AppConfig {
	dbPort, err := strconv.Atoi(getEnv("DB_PORT", "5432"))
	if err != nil {
		dbPort = 5432
	}

	return AppConfig{
		AWSRegion:  getEnv("AWS_REGION", "eu-west-1"),
		AWSProfile: os.Getenv("AWS_PROFILE"),

		// Aurora PostgreSQL Database Configuration
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     dbPort,
		DBName:     os.Getenv("DB_NAME"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
