# Database Seed System Documentation

## Overview

The database seed system has been completely restructured to use modular SQL scripts instead of a single monolithic `seed_data.sql` file. This provides better maintainability, dependency management, and ensures that seed data doesn't conflict with test data used in `requests/*.http` files.

## File Structure

```
sql/
├── seed_category_groups.sql      # Category groups (no dependencies)
├── seed_categories.sql           # Categories (depends on category_groups)
├── seed_merchants.sql           # Merchants (no dependencies)
├── seed_balances.sql            # Balance accounts (no dependencies)
├── seed_transactions.sql        # Transactions (depends on balances, merchants)
├── seed_transaction_entries.sql # Transaction entries (depends on transactions, categories)
└── seed_data.sql               # Legacy monolithic file (kept for reference)

scripts/
└── seed_database.sh            # Execution script with proper dependency order
```

## Execution Order

The seed script executes files in the correct dependency order:

1. **Category Groups** (`seed_category_groups.sql`) - No dependencies
2. **Categories** (`seed_categories.sql`) - Depends on category groups
3. **Merchants** (`seed_merchants.sql`) - No dependencies
4. **Balances** (`seed_balances.sql`) - No dependencies  
5. **Transactions** (`seed_transactions.sql`) - Depends on balances and merchants
6. **Transaction Entries** (`seed_transaction_entries.sql`) - Depends on transactions and categories

## Schema Alignment

All seed scripts are aligned with the latest DAO schema requirements:

### CategoryGroup Model
- Table: `category_group`
- Required fields: `id` (UUID), `name` (string)
- Optional fields: `rank` (int), `image_url` (string)

### Category Model  
- Table: `category`
- Required fields: `id` (UUID), `user_id` (UUID), `group_id` (UUID), `category_group_id` (string), `name` (string), `group` (string)
- Optional fields: `description` (string), `rank` (int), `image_url` (string)

### Merchant Model
- Table: `merchant`
- Required fields: `id` (UUID), `group_id` (UUID), `user_id` (UUID), `name` (string)
- Optional fields: `description` (string), `rank` (int)

### Balance Model
- Table: `balance`
- Required fields: `id` (UUID), `group_id` (UUID), `user_id` (UUID), `currency` (string), `title` (string)
- Optional fields: `description` (string), `rank` (int)

### Transaction Model
- Table: `transaction`
- Required fields: `id` (UUID), `group_id` (UUID), `user_id` (UUID), `balance_id` (UUID), `type` (string), `approved_at` (timestamp), `transacted_at` (timestamp)
- Optional fields: `merchant_id` (UUID), `operation_id` (UUID)

### TransactionEntry Model
- Table: `transaction_entry`
- Required fields: `id` (UUID), `transaction_id` (UUID), `amount` (decimal)
- Optional fields: `description` (string), `category_id` (UUID)

## Data Isolation

The seed data uses completely different UUIDs from those used in `requests/*.http` files to prevent conflicts during testing:

### Seed Data UUIDs
- Group ID: `88aa1100-0011-2233-4455-667788990011`
- User ID: `99bb2200-0011-2233-4455-667788990011`
- All entity IDs follow patterns like `11111111-1111-1111-1111-111111111111`

### Request File UUIDs  
- Group ID: `6a785a55-fced-4f13-af78-5c19a39c9abc`
- User IDs: `02c514a4-2021-708d-efff-ea6cd5e4eac9`, `12c514a4-2021-708d-efff-ea6cd5e4eac8`
- Balance ID: `28e2d53a-22e9-4c7e-9c06-0b91a9d091f4`

## Usage

### Via Makefile (Recommended)
```bash
make seed
```

### Direct Script Execution
```bash
# Set environment variables
export DB_HOST="your-db-host"
export DB_PORT="5432"
export DB_USER="postgres"
export DB_PASSWORD="your-password"
export DB_NAME="your-database"

# Run the script
./scripts/seed_database.sh
```

### Using Docker (Production-like)
The Makefile target automatically uses Docker with the PostgreSQL client, passing environment variables securely.

## Features

- ✅ **Modular Design**: Each entity type has its own SQL file
- ✅ **Dependency Management**: Automatic execution in correct order
- ✅ **Conflict Prevention**: Uses `ON CONFLICT (id) DO NOTHING` clauses
- ✅ **Data Isolation**: Separate UUIDs from test request files
- ✅ **Error Handling**: Script exits on any SQL error
- ✅ **Verification**: Automatically counts records after seeding
- ✅ **Schema Compliant**: Matches latest DAO model requirements
- ✅ **Environment Configurable**: All connection parameters via env vars
- ✅ **Docker Support**: Works with containerized PostgreSQL client

## Sample Data

The seed system creates realistic sample data including:

- **10 Category Groups**: Food & Dining, Shopping, Transportation, etc.
- **13 Categories**: Groceries, Coffee & Tea, Restaurants, Electronics, etc.
- **33 Merchants**: Mercadona, Starbucks, Uber, Netflix, etc.
- **6 Balance Accounts**: BBVA, Santander, ING accounts in EUR/USD/GBP
- **12 Transactions**: Mix of expenses, income, and transfers
- **16 Transaction Entries**: Detailed line items with proper categorization

## Migration from Legacy System

The old monolithic `seed_data.sql` file is preserved for reference but is no longer used. The new system provides:

1. **Better Maintainability**: Each entity type can be updated independently
2. **Clearer Dependencies**: Explicit execution order prevents foreign key errors
3. **Improved Testing**: Separate data sets for seed vs. request testing
4. **Enhanced Debugging**: Individual script failures are easier to diagnose
5. **Flexible Deployment**: Can seed partial data sets if needed

## Troubleshooting

### Common Issues
- **Foreign Key Errors**: Ensure scripts run in dependency order
- **UUID Conflicts**: Check that seed UUIDs don't overlap with test data
- **Permission Errors**: Verify database user has INSERT permissions
- **Connection Issues**: Check database endpoint and credentials

### Debug Mode
Add `set -x` to the top of `seed_database.sh` for verbose output:
```bash
#!/bin/bash
set -e
set -x  # Add this line for debug output
```

## Future Enhancements

Potential improvements to consider:

1. **Incremental Updates**: Support for updating existing seed data
2. **Environment-Specific Data**: Different data sets for dev/staging/prod
3. **Data Validation**: JSON schema validation for seed data consistency
4. **Performance Optimization**: Bulk insert operations for large data sets
5. **Rollback Support**: Ability to remove seed data cleanly
