#!/bin/bash
# filepath: /Users/savak/Projects/Ahorro/ahorro-transactions-service/scripts/seed_database.sh
# seed_database.sh: Execute all seed scripts in proper order for Aurora PostgreSQL Ahorro Transactions Service

set -e  # Exit on any error

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Get the directory of this script
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SQL_DIR="$SCRIPT_DIR/../sql"

# Database connection parameters (can be overridden by environment variables)
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5432}"
DB_NAME="${DB_NAME:-postgres}"
DB_USER="${DB_USER:-postgres}"
DB_PASSWORD="${DB_PASSWORD:-password}"

echo -e "${BLUE}Starting database seeding process...${NC}"

# Function to execute SQL file
execute_sql_file() {
    local file_path="$1"
    local file_name=$(basename "$file_path")
    
    echo -e "${BLUE}Executing: $file_name${NC}"
    
    if PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -f "$file_path"; then
        echo -e "${GREEN}✓ Successfully executed: $file_name${NC}"
    else
        echo -e "${RED}✗ Failed to execute: $file_name${NC}"
        exit 1
    fi
}

# Execute seed scripts in proper order
echo -e "${BLUE}Executing seed scripts in dependency order...${NC}"

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

echo -e "${GREEN}✓ Database seeding completed successfully!${NC}"

# Verify the data
echo -e "${BLUE}Verifying seeded data...${NC}"
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

echo -e "${GREEN}✓ Database seeding verification completed!${NC}"
