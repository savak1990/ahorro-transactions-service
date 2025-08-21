package config

import (
	"os"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm/logger"
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

	// Exchange Rate Configuration
	ExchangeRateApiKey string
	ExchangeRateDbName string

	// Application Configuration
	LogLevel string
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

		// Currency Exchange Rate Db
		ExchangeRateDbName: os.Getenv("EXCHANGE_RATE_DB_NAME"),

		// Application Configuration
		LogLevel: getEnv("LOG_LEVEL", "info"),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// GetLogrusLevel converts string log level to logrus.Level
func GetLogrusLevel(level string) logrus.Level {
	switch strings.ToLower(level) {
	case "debug":
		return logrus.DebugLevel
	case "info":
		return logrus.InfoLevel
	case "warn", "warning":
		return logrus.WarnLevel
	case "error":
		return logrus.ErrorLevel
	default:
		return logrus.InfoLevel
	}
}

// GetGormLogLevel converts string log level to gorm logger.LogLevel
func GetGormLogLevel(level string) logger.LogLevel {
	switch strings.ToLower(level) {
	case "debug":
		return logger.Info // GORM debug level shows all SQL queries
	case "info":
		return logger.Warn // GORM warn level shows errors and warnings but not slow queries
	case "warn", "warning":
		return logger.Error // GORM error level shows only errors
	case "error":
		return logger.Error
	default:
		return logger.Warn
	}
}
