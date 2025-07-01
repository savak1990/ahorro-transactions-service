#!/bin/bash
# filepath: /Users/savak/Projects/Ahorro/ahorro-transactions-service/scripts/verify_seed_data.sh
# verify_seed_data.sh: Verify that seed data was inserted correctly

set -e

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Database connection parameters (can be overridden by environment variables)
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5432}"
DB_NAME="${DB_NAME:-postgres}"
DB_USER="${DB_USER:-postgres}"
DB_PASSWORD="${DB_PASSWORD:-password}"

echo -e "${BLUE}🔍 Verifying Seed Data${NC}"
echo -e "${BLUE}Host: $DB_HOST:$DB_PORT${NC}"
echo -e "${BLUE}Database: $DB_NAME${NC}"
echo ""

# Function to run query and check results
verify_table() {
    local table_name="$1"
    local expected_min="$2"
    local display_name="$3"
    
    echo -e "${BLUE}Checking $display_name...${NC}"
    
    local count=$(PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -t -c "SELECT COUNT(*) FROM $table_name;" | tr -d ' ')
    
    if [ "$count" -ge "$expected_min" ]; then
        echo -e "${GREEN}✓ $display_name: $count records (expected ≥ $expected_min)${NC}"
        return 0
    else
        echo -e "${RED}✗ $display_name: $count records (expected ≥ $expected_min)${NC}"
        return 1
    fi
}

# Function to check data relationships
verify_relationships() {
    echo -e "${BLUE}Verifying data relationships...${NC}"
    
    # Check categories have valid category_group_id references
    local orphaned_categories=$(PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -t -c "
        SELECT COUNT(*) 
        FROM category c 
        WHERE NOT EXISTS (
            SELECT 1 FROM category_group cg 
            WHERE cg.id::text = c.category_group_id 
            OR cg.name = c.\"group\"
        );" | tr -d ' ')
    
    if [ "$orphaned_categories" -eq 0 ]; then
        echo -e "${GREEN}✓ All categories have valid group references${NC}"
    else
        echo -e "${YELLOW}⚠ Found $orphaned_categories categories with invalid group references${NC}"
    fi
    
    # Check transactions have valid balance references
    local orphaned_transactions=$(PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -t -c "
        SELECT COUNT(*) 
        FROM transaction t 
        WHERE NOT EXISTS (
            SELECT 1 FROM balance b WHERE b.id = t.balance_id
        );" | tr -d ' ')
    
    if [ "$orphaned_transactions" -eq 0 ]; then
        echo -e "${GREEN}✓ All transactions have valid balance references${NC}"
    else
        echo -e "${RED}✗ Found $orphaned_transactions transactions with invalid balance references${NC}"
    fi
    
    # Check transaction entries have valid transaction references
    local orphaned_entries=$(PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -t -c "
        SELECT COUNT(*) 
        FROM transaction_entry te 
        WHERE NOT EXISTS (
            SELECT 1 FROM transaction t WHERE t.id = te.transaction_id
        );" | tr -d ' ')
    
    if [ "$orphaned_entries" -eq 0 ]; then
        echo -e "${GREEN}✓ All transaction entries have valid transaction references${NC}"
    else
        echo -e "${RED}✗ Found $orphaned_entries transaction entries with invalid transaction references${NC}"
    fi
}

# Function to check data quality
verify_data_quality() {
    echo -e "${BLUE}Verifying data quality...${NC}"
    
    # Check for transactions with different currencies in same transaction
    local currency_mismatches=$(PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -t -c "
        SELECT COUNT(DISTINCT t.id)
        FROM transaction t
        JOIN balance b ON t.balance_id = b.id
        JOIN transaction_entry te ON t.id = te.transaction_id
        GROUP BY t.id, b.currency
        HAVING COUNT(DISTINCT b.currency) > 1;" | tr -d ' ')
    
    if [ -z "$currency_mismatches" ] || [ "$currency_mismatches" -eq 0 ]; then
        echo -e "${GREEN}✓ No currency mismatches found${NC}"
    else
        echo -e "${YELLOW}⚠ Found $currency_mismatches transactions with currency mismatches${NC}"
    fi
    
    # Check for transactions without entries
    local empty_transactions=$(PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -t -c "
        SELECT COUNT(*) 
        FROM transaction t 
        WHERE NOT EXISTS (
            SELECT 1 FROM transaction_entry te WHERE te.transaction_id = t.id
        );" | tr -d ' ')
    
    if [ "$empty_transactions" -eq 0 ]; then
        echo -e "${GREEN}✓ All transactions have entries${NC}"
    else
        echo -e "${YELLOW}⚠ Found $empty_transactions transactions without entries${NC}"
    fi
}

# Main verification
error_count=0

# Verify table counts
verify_table "category_group" 5 "Category Groups" || ((error_count++))
verify_table "category" 10 "Categories" || ((error_count++))
verify_table "merchant" 20 "Merchants" || ((error_count++))
verify_table "balance" 4 "Balances" || ((error_count++))
verify_table "transaction" 8 "Transactions" || ((error_count++))
verify_table "transaction_entry" 10 "Transaction Entries" || ((error_count++))

echo ""

# Verify relationships and data quality
verify_relationships
verify_data_quality

echo ""

# Final summary
if [ $error_count -eq 0 ]; then
    echo -e "${GREEN}🎉 All seed data verification checks passed!${NC}"
    exit 0
else
    echo -e "${RED}❌ Found $error_count issues with seed data${NC}"
    echo -e "${YELLOW}💡 Consider re-running the seed process: make seed${NC}"
    exit 1
fi
