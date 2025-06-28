-- Comprehensive query to extract transaction entries with merchant, category, and balance names
-- for a specific group_id with all relevant transaction details

SELECT 
    -- Transaction Entry Details
    te.id AS entry_id,
    te.description AS entry_description,
    te.amount,
    te.created_at AS entry_created_at,
    
    -- Transaction Details
    t.id AS transaction_id,
    t.type AS transaction_type,
    t.transacted_at,
    t.approved_at,
    
    -- Merchant Information (LEFT JOIN because salary transactions don't have merchants)
    COALESCE(m.name, 'No Merchant') AS merchant_name,
    m.description AS merchant_description,
    
    -- Category Information
    c.category_name,
    c."group" AS category_group,
    c.rank AS category_rank,
    
    -- Balance Information
    b.title AS balance_name,
    b.description AS balance_description,
    b.currency,
    
    -- User and Group Information
    t.user_id,
    t.group_id
    
FROM transaction_entry te
    -- Join with transaction (required)
    INNER JOIN transaction t ON te.transaction_id = t.id
    
    -- Join with merchant (optional - income transactions may not have merchants)
    LEFT JOIN merchant m ON t.merchant_id = m.id
    
    -- Join with category (optional - some entries might not have categories)
    LEFT JOIN category c ON te.category_id = c.id
    
    -- Join with balance (required)
    INNER JOIN balance b ON t.balance_id = b.id

-- Filter by group_id from your seed data
WHERE t.group_id = 'b47ac10b-58cc-4372-a567-0e02b2c3d479'

-- Order by transaction date (newest first) and then by entry id for consistency
ORDER BY t.transacted_at DESC, te.id;
