#!/bin/bash
# filepath: /Users/savak/Projects/Ahorro/ahorro-transactions-service/scripts/seed_database.sh
# seed_database.sh: Execute all seed scripts in proper order for Aurora PostgreSQL Ahorro Transactions Service

set -e  # Exit on any error

# Get the directory of this script
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SQL_DIR="$SCRIPT_DIR/../sql"

# Database connection parameters (can be overridden by environment variables)
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5432}"
DB_NAME="${DB_NAME:-postgres}"
DB_USER="${DB_USER:-postgres}"
DB_PASSWORD="${DB_PASSWORD:-password}"

echo "Starting database seeding process..."

# Function to execute SQL file
execute_sql_file() {
    local file_path="$1"
    local file_name=$(basename "$file_path")
    
    echo "Executing: $file_name"
    
    if PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -f "$file_path"; then
        echo "Successfully executed: $file_name"
    else
        echo "Failed to execute: $file_name"
        exit 1
    fi
}

# Execute seed scripts in proper order
echo "Executing seed scripts in dependency order..."

# 1. Category Groups (no dependencies)
execute_sql_file "$SQL_DIR/seed_category_groups.sql"

# 2. Categories (depends on category groups)
execute_sql_file "$SQL_DIR/seed_categories.sql"

# 3. Merchants (no dependencies on other entities)
execute_sql_file "$SQL_DIR/seed_merchants.sql"

# 4. Balances (no dependencies on other entities)
execute_sql_file "$SQL_DIR/seed_balances.sql"

# 5. Transactions (depends on balances and merchants)
execute_sql_file "$SQL_DIR/seed_transactions.sql"

# 6. Transaction Entries (depends on transactions and categories)
execute_sql_file "$SQL_DIR/seed_transaction_entries.sql"

echo "Database seeding completed successfully!"

# Verify the data
echo "Verifying seeded data..."
PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -c "
SELECT 'Category Groups' as table_name, COUNT(*) as record_count FROM category_group
UNION ALL
SELECT 'Categories' as table_name, COUNT(*) as record_count FROM category
UNION ALL
SELECT 'Merchants' as table_name, COUNT(*) as record_count FROM merchant
UNION ALL
SELECT 'Balances' as table_name, COUNT(*) as record_count FROM balance
UNION ALL
SELECT 'Transactions' as table_name, COUNT(*) as record_count FROM transaction
UNION ALL
SELECT 'Transaction Entries' as table_name, COUNT(*) as record_count FROM transaction_entry
ORDER BY table_name;
"

echo "Database seeding verification completed!"
