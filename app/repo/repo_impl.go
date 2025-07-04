package repo

import (
	"fmt"

	"github.com/savak1990/transactions-service/app/aws"
	"github.com/savak1990/transactions-service/app/config"
	"gorm.io/gorm"
)

// PostgreSQLRepository implements Repository interface using PostgreSQL
type PostgreSQLRepository struct {
	db     *gorm.DB
	config *config.AppConfig // Store config for lazy initialization
}

// NewPostgreSQLRepository creates a new PostgreSQL repository
func NewPostgreSQLRepository(db *gorm.DB) *PostgreSQLRepository {
	return &PostgreSQLRepository{
		db: db,
	}
}

// NewPostgreSQLRepositoryWithConfig creates a new PostgreSQL repository with lazy DB initialization
func NewPostgreSQLRepositoryWithConfig(cfg config.AppConfig) *PostgreSQLRepository {
	return &PostgreSQLRepository{
		config: &cfg,
	}
}

// getDB returns the database connection, initializing it if necessary
func (r *PostgreSQLRepository) getDB() *gorm.DB {
	if r.db != nil {
		return r.db
	}

	if r.config != nil {
		// Lazy initialization - this will trigger connection and panic if failed
		// We want this panic to bubble up to the maintenance middleware
		r.db = aws.GetGormDB(*r.config)
		return r.db
	}

	panic(fmt.Errorf("PostgreSQL repository not properly initialized"))
}

// Ensure PostgreSQLRepository implements Repository interface
var _ Repository = (*PostgreSQLRepository)(nil)
