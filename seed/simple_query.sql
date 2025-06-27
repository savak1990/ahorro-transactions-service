-- Simple query to get transaction entries with merchant, category, and balance names
-- for group_id 'b47ac10b-58cc-4372-a567-0e02b2c3d479'

SELECT 
    te.description AS entry_description,
    te.amount,
    COALESCE(m.name, 'No Merchant') AS merchant_name,
    c.category_name,
    b.title AS balance_name,
    t.transacted_at
    
FROM transaction_entry te
    INNER JOIN transaction t ON te.transaction_id = t.id
    LEFT JOIN merchant m ON t.merchant_id = m.id
    LEFT JOIN category c ON te.category_id = c.id
    INNER JOIN balance b ON t.balance_id = b.id

WHERE t.group_id = 'b47ac10b-58cc-4372-a567-0e02b2c3d479'

ORDER BY t.transacted_at DESC;
