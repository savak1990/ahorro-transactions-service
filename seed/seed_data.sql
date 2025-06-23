-- seed_data.sql: Sample data for Aurora PostgreSQL Ahorro Transactions Service

-- Insert sample merchants
INSERT INTO merchant (id, name, created_at, updated_at) VALUES
    ('f47ac10b-58cc-4372-a567-0e02b2c3d479', 'Amazon', NOW(), NOW()),
    ('6ba7b810-9dad-11d1-80b4-00c04fd430c8', 'Starbucks', NOW(), NOW()),
    ('6ba7b811-9dad-11d1-80b4-00c04fd430c8', 'Uber', NOW(), NOW()),
    ('6ba7b812-9dad-11d1-80b4-00c04fd430c8', 'Netflix', NOW(), NOW()),
    ('6ba7b813-9dad-11d1-80b4-00c04fd430c8', 'Grocery Store', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- Insert sample categories
INSERT INTO category (id, category_name, category_group, priority, created_at, updated_at) VALUES
    ('c47ac10b-58cc-4372-a567-0e02b2c3d479', 'Food & Dining', 'Personal', 10, NOW(), NOW()),
    ('c47ac10b-58cc-4372-a567-0e02b2c3d480', 'Shopping', 'Personal', 8, NOW(), NOW()),
    ('c47ac10b-58cc-4372-a567-0e02b2c3d481', 'Transportation', 'Travel', 7, NOW(), NOW()),
    ('c47ac10b-58cc-4372-a567-0e02b2c3d482', 'Entertainment', 'Personal', 6, NOW(), NOW()),
    ('c47ac10b-58cc-4372-a567-0e02b2c3d483', 'Groceries', 'Personal', 9, NOW(), NOW()),
    ('c47ac10b-58cc-4372-a567-0e02b2c3d484', 'Salary', 'Income', 15, NOW(), NOW()),
    ('c47ac10b-58cc-4372-a567-0e02b2c3d485', 'Freelance', 'Income', 12, NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- Insert sample transactions
INSERT INTO transaction (id, group_id, user_id, balance_id, merchant_id, type, approved_at, transacted_at, created_at, updated_at) VALUES
    ('a47ac10b-58cc-4372-a567-0e02b2c3d479', 'b47ac10b-58cc-4372-a567-0e02b2c3d479', '847ac10b-58cc-4372-a567-0e02b2c3d479', '847ac10b-58cc-4372-a567-0e02b2c3d479', 'f47ac10b-58cc-4372-a567-0e02b2c3d479', 'expense', NOW(), NOW() - INTERVAL '1 day', NOW(), NOW()),
    ('a47ac10b-58cc-4372-a567-0e02b2c3d480', 'b47ac10b-58cc-4372-a567-0e02b2c3d479', '847ac10b-58cc-4372-a567-0e02b2c3d479', '847ac10b-58cc-4372-a567-0e02b2c3d479', '6ba7b810-9dad-11d1-80b4-00c04fd430c8', 'expense', NOW(), NOW() - INTERVAL '2 days', NOW(), NOW()),
    ('a47ac10b-58cc-4372-a567-0e02b2c3d481', 'b47ac10b-58cc-4372-a567-0e02b2c3d479', '847ac10b-58cc-4372-a567-0e02b2c3d479', '847ac10b-58cc-4372-a567-0e02b2c3d479', '6ba7b811-9dad-11d1-80b4-00c04fd430c8', 'expense', NOW(), NOW() - INTERVAL '3 days', NOW(), NOW()),
    ('a47ac10b-58cc-4372-a567-0e02b2c3d482', 'b47ac10b-58cc-4372-a567-0e02b2c3d479', '847ac10b-58cc-4372-a567-0e02b2c3d479', '847ac10b-58cc-4372-a567-0e02b2c3d479', NULL, 'income', NOW(), NOW() - INTERVAL '7 days', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- Insert sample transaction entries
INSERT INTO transaction_entry (id, transaction_id, description, amount, category_id, created_at, updated_at) VALUES
    ('e47ac10b-58cc-4372-a567-0e02b2c3d479', 'a47ac10b-58cc-4372-a567-0e02b2c3d479', 'Office supplies from Amazon', -29.99, 'c47ac10b-58cc-4372-a567-0e02b2c3d480', NOW(), NOW()),
    ('e47ac10b-58cc-4372-a567-0e02b2c3d480', 'a47ac10b-58cc-4372-a567-0e02b2c3d480', 'Morning coffee', -5.50, 'c47ac10b-58cc-4372-a567-0e02b2c3d479', NOW(), NOW()),
    ('e47ac10b-58cc-4372-a567-0e02b2c3d481', 'a47ac10b-58cc-4372-a567-0e02b2c3d481', 'Ride to airport', -25.00, 'c47ac10b-58cc-4372-a567-0e02b2c3d481', NOW(), NOW()),
    ('e47ac10b-58cc-4372-a567-0e02b2c3d482', 'a47ac10b-58cc-4372-a567-0e02b2c3d482', 'Monthly salary', 3500.00, 'c47ac10b-58cc-4372-a567-0e02b2c3d484', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- Verify the data
SELECT 'Merchants' as table_name, COUNT(*) as record_count FROM merchant
UNION ALL
SELECT 'Categories' as table_name, COUNT(*) as record_count FROM category
UNION ALL
SELECT 'Transactions' as table_name, COUNT(*) as record_count FROM transaction
UNION ALL
SELECT 'Transaction Entries' as table_name, COUNT(*) as record_count FROM transaction_entry;
