#!/bin/bash
# filepath: /Users/savak/Projects/Ahorro/ahorro-transactions-service/scripts/verify_seed_data.sh
# verify_seed_data.sh: Verify that seed data was inserted correctly

set -e

# Database connection parameters (can be overridden by environment variables)
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5432}"
DB_NAME="${DB_NAME:-postgres}"
DB_USER="${DB_USER:-postgres}"
DB_PASSWORD="${DB_PASSWORD:-password}"

echo "Verifying Seed Data"
echo "Host: $DB_HOST:$DB_PORT"
echo "Database: $DB_NAME"
echo ""

# Function to run query and check results
verify_table() {
    local table_name="$1"
    local expected_min="$2"
    local display_name="$3"
    
    echo "Checking $display_name..."
    
    local count=$(PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -t -c "SELECT COUNT(*) FROM $table_name;" | tr -d ' ')
    
    if [ "$count" -ge "$expected_min" ]; then
        echo "OK $display_name: $count records (expected >= $expected_min)"
        return 0
    else
        echo "FAIL $display_name: $count records (expected >= $expected_min)"
        return 1
    fi
}

# Function to check data relationships
verify_relationships() {
    echo "Verifying data relationships..."
    
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
        echo "OK All categories have valid group references"
    else
        echo "WARNING Found $orphaned_categories categories with invalid group references"
    fi
    
    # Check transactions have valid balance references
    local orphaned_transactions=$(PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -t -c "
        SELECT COUNT(*) 
        FROM transaction t 
        WHERE NOT EXISTS (
            SELECT 1 FROM balance b WHERE b.id = t.balance_id
        );" | tr -d ' ')
    
    if [ "$orphaned_transactions" -eq 0 ]; then
        echo "OK All transactions have valid balance references"
    else
        echo "FAIL Found $orphaned_transactions transactions with invalid balance references"
    fi
    
    # Check transaction entries have valid transaction references
    local orphaned_entries=$(PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -t -c "
        SELECT COUNT(*) 
        FROM transaction_entry te 
        WHERE NOT EXISTS (
            SELECT 1 FROM transaction t WHERE t.id = te.transaction_id
        );" | tr -d ' ')
    
    if [ "$orphaned_entries" -eq 0 ]; then
        echo "OK All transaction entries have valid transaction references"
    else
        echo "FAIL Found $orphaned_entries transaction entries with invalid transaction references"
    fi
    
    # Check transaction entry amounts have valid transaction entry references
    local orphaned_amounts=$(PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -t -c "
        SELECT COUNT(*) 
        FROM transaction_entry_amount tea 
        WHERE NOT EXISTS (
            SELECT 1 FROM transaction_entry te WHERE te.id = tea.transaction_entry_id
        );" | tr -d ' ')
    
    if [ "$orphaned_amounts" -eq 0 ]; then
        echo "OK All transaction entry amounts have valid transaction entry references"
    else
        echo "FAIL Found $orphaned_amounts transaction entry amounts with invalid transaction entry references"
    fi
}

# Function to check data quality
verify_data_quality() {
    echo "Verifying data quality..."
    
    # Check for transactions with different currencies in same transaction
    local currency_mismatches=$(PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -t -c "
        SELECT COUNT(DISTINCT t.id)
        FROM transaction t
        JOIN balance b ON t.balance_id = b.id
        JOIN transaction_entry te ON t.id = te.transaction_id
        GROUP BY t.id, b.currency
        HAVING COUNT(DISTINCT b.currency) > 1;" | tr -d ' ')
    
    if [ -z "$currency_mismatches" ] || [ "$currency_mismatches" -eq 0 ]; then
        echo "OK No currency mismatches found"
    else
        echo "WARNING Found $currency_mismatches transactions with currency mismatches"
    fi
    
    # Check for transactions without entries
    local empty_transactions=$(PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -t -c "
        SELECT COUNT(*) 
        FROM transaction t 
        WHERE NOT EXISTS (
            SELECT 1 FROM transaction_entry te WHERE te.transaction_id = t.id
        );" | tr -d ' ')
    
    if [ "$empty_transactions" -eq 0 ]; then
        echo "OK All transactions have entries"
    else
        echo "WARNING Found $empty_transactions transactions without entries"
    fi
    
    # Check transaction entries without currency amounts
    local entries_without_amounts=$(PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -t -c "
        SELECT COUNT(*) 
        FROM transaction_entry te 
        WHERE NOT EXISTS (
            SELECT 1 FROM transaction_entry_amount tea WHERE tea.transaction_entry_id = te.id
        );" | tr -d ' ')
    
    if [ "$entries_without_amounts" -eq 0 ]; then
        echo "OK All transaction entries have currency amounts"
    else
        echo "WARNING Found $entries_without_amounts transaction entries without currency amounts"
    fi
    
    # Check if all entries have amounts in supported currencies (EUR, USD, GBP)
    local entries_with_full_currencies=$(PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -t -c "
        SELECT COUNT(DISTINCT te.id)
        FROM transaction_entry te
        WHERE (
            SELECT COUNT(DISTINCT currency) 
            FROM transaction_entry_amount tea 
            WHERE tea.transaction_entry_id = te.id 
            AND currency IN ('EUR', 'USD', 'GBP')
        ) = 3;" | tr -d ' ')
    
    local total_entries=$(PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -t -c "SELECT COUNT(*) FROM transaction_entry;" | tr -d ' ')
    
    if [ "$entries_with_full_currencies" -eq "$total_entries" ]; then
        echo "OK All transaction entries have amounts in all supported currencies (EUR, USD, GBP)"
    else
        echo "WARNING Found $((total_entries - entries_with_full_currencies)) transaction entries missing some currency amounts"
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
verify_table "transaction_entry_amount" 60 "Transaction Entry Amounts" || ((error_count++))

echo ""

# Verify relationships and data quality
verify_relationships
verify_data_quality

echo ""

# Final summary
if [ $error_count -eq 0 ]; then
    echo "SUCCESS All seed data verification checks passed!"
    exit 0
else
    echo "FAIL Found $error_count issues with seed data"
    echo "TIP Consider re-running the seed process: make seed"
    exit 1
fi
