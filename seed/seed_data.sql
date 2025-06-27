-- seed_data.sql: Sample data for Aurora PostgreSQL Ahorro Transactions Service

-- Insert sample merchants
INSERT INTO merchant (id, name, description, created_at, updated_at) VALUES
    -- Spanish Supermarkets & Retail
    ('f47ac10b-58cc-4372-a567-0e02b2c3d479', 'Mercadona', 'Spanish supermarket chain', NOW(), NOW()),
    ('f47ac10b-58cc-4372-a567-0e02b2c3d480', 'El Corte Inglés', 'Spanish department store chain', NOW(), NOW()),
    ('f47ac10b-58cc-4372-a567-0e02b2c3d481', 'Carrefour', 'French multinational retail corporation', NOW(), NOW()),
    ('f47ac10b-58cc-4372-a567-0e02b2c3d482', 'Lidl', 'German discount supermarket chain', NOW(), NOW()),
    ('f47ac10b-58cc-4372-a567-0e02b2c3d483', 'Zara', 'Spanish fashion retailer', NOW(), NOW()),
    ('f47ac10b-58cc-4372-a567-0e02b2c3d484', 'MediaMarkt', 'German electronics retailer', NOW(), NOW()),
    
    -- Restaurants & Food
    ('6ba7b810-9dad-11d1-80b4-00c04fd430c8', 'Starbucks', 'American coffeehouse chain', NOW(), NOW()),
    ('6ba7b811-9dad-11d1-80b4-00c04fd430c8', 'McDonald''s', 'American fast food restaurant chain', NOW(), NOW()),
    ('6ba7b812-9dad-11d1-80b4-00c04fd430c8', 'Telepizza', 'Spanish pizza delivery chain', NOW(), NOW()),
    ('6ba7b813-9dad-11d1-80b4-00c04fd430c8', 'Domino''s Pizza', 'American pizza restaurant chain', NOW(), NOW()),
    ('6ba7b814-9dad-11d1-80b4-00c04fd430c8', 'Burger King', 'American fast food restaurant chain', NOW(), NOW()),
    
    -- Transportation
    ('6ba7b815-9dad-11d1-80b4-00c04fd430c8', 'Uber', 'American mobility company', NOW(), NOW()),
    ('6ba7b816-9dad-11d1-80b4-00c04fd430c8', 'Cabify', 'Spanish ride-hailing company', NOW(), NOW()),
    ('6ba7b817-9dad-11d1-80b4-00c04fd430c8', 'Renfe', 'Spanish railway company', NOW(), NOW()),
    ('6ba7b818-9dad-11d1-80b4-00c04fd430c8', 'Repsol', 'Spanish oil and gas company', NOW(), NOW()),
    ('6ba7b819-9dad-11d1-80b4-00c04fd430c8', 'Metro Madrid', 'Madrid metro system', NOW(), NOW()),
    
    -- Entertainment & Streaming
    ('6ba7b81a-9dad-11d1-80b4-00c04fd430c8', 'Netflix', 'American streaming service', NOW(), NOW()),
    ('6ba7b81b-9dad-11d1-80b4-00c04fd430c8', 'Amazon Prime', 'Amazon streaming service', NOW(), NOW()),
    ('6ba7b81c-9dad-11d1-80b4-00c04fd430c8', 'Spotify', 'Swedish music streaming service', NOW(), NOW()),
    ('6ba7b81d-9dad-11d1-80b4-00c04fd430c8', 'Disney+', 'American streaming service', NOW(), NOW()),
    ('6ba7b81e-9dad-11d1-80b4-00c04fd430c8', 'HBO Max', 'American streaming service', NOW(), NOW()),
    
    -- Technology & Online Services
    ('6ba7b81f-9dad-11d1-80b4-00c04fd430c8', 'Amazon', 'American e-commerce company', NOW(), NOW()),
    ('6ba7b820-9dad-11d1-80b4-00c04fd430c8', 'Apple', 'American technology company', NOW(), NOW()),
    ('6ba7b821-9dad-11d1-80b4-00c04fd430c8', 'Google', 'American technology company', NOW(), NOW()),
    ('6ba7b822-9dad-11d1-80b4-00c04fd430c8', 'Microsoft', 'American technology company', NOW(), NOW()),
    
    -- Utilities & Services
    ('6ba7b823-9dad-11d1-80b4-00c04fd430c8', 'Iberdrola', 'Spanish electric utility company', NOW(), NOW()),
    ('6ba7b824-9dad-11d1-80b4-00c04fd430c8', 'Telefónica', 'Spanish telecommunications company', NOW(), NOW()),
    ('6ba7b825-9dad-11d1-80b4-00c04fd430c8', 'Orange', 'French telecommunications company', NOW(), NOW()),
    ('6ba7b826-9dad-11d1-80b4-00c04fd430c8', 'Vodafone', 'British telecommunications company', NOW(), NOW()),
    
    -- Banks & Financial
    ('6ba7b827-9dad-11d1-80b4-00c04fd430c8', 'BBVA', 'Spanish multinational bank', NOW(), NOW()),
    ('6ba7b828-9dad-11d1-80b4-00c04fd430c8', 'Santander', 'Spanish multinational bank', NOW(), NOW()),
    ('6ba7b829-9dad-11d1-80b4-00c04fd430c8', 'CaixaBank', 'Spanish bank', NOW(), NOW()),
    ('6ba7b82a-9dad-11d1-80b4-00c04fd430c8', 'ING', 'Dutch multinational bank', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- Insert sample categories
INSERT INTO category (id, category_name, "group", rank, created_at, updated_at) VALUES
    -- Food & Dining Categories
    ('c47ac10b-58cc-4372-a567-0e02b2c3d479', 'Groceries', 'Food & Dining', 10, NOW(), NOW()),
    ('c47ac10b-58cc-4372-a567-0e02b2c3d480', 'Restaurants', 'Food & Dining', 9, NOW(), NOW()),
    ('c47ac10b-58cc-4372-a567-0e02b2c3d481', 'Fast Food', 'Food & Dining', 8, NOW(), NOW()),
    ('c47ac10b-58cc-4372-a567-0e02b2c3d482', 'Coffee & Tea', 'Food & Dining', 7, NOW(), NOW()),
    ('c47ac10b-58cc-4372-a567-0e02b2c3d483', 'Bars & Pubs', 'Food & Dining', 6, NOW(), NOW()),
    ('c47ac10b-58cc-4372-a567-0e02b2c3d484', 'Takeaway', 'Food & Dining', 5, NOW(), NOW()),
    
    -- Shopping Categories
    ('c47ac10b-58cc-4372-a567-0e02b2c3d485', 'Clothing', 'Shopping', 15, NOW(), NOW()),
    ('c47ac10b-58cc-4372-a567-0e02b2c3d486', 'Electronics', 'Shopping', 14, NOW(), NOW()),
    ('c47ac10b-58cc-4372-a567-0e02b2c3d487', 'Home & Garden', 'Shopping', 13, NOW(), NOW()),
    ('c47ac10b-58cc-4372-a567-0e02b2c3d488', 'Books & Magazines', 'Shopping', 12, NOW(), NOW()),
    ('c47ac10b-58cc-4372-a567-0e02b2c3d489', 'Pharmacy', 'Shopping', 11, NOW(), NOW()),
    ('c47ac10b-58cc-4372-a567-0e02b2c3d48a', 'Gifts', 'Shopping', 10, NOW(), NOW()),
    
    -- Transportation Categories
    ('c47ac10b-58cc-4372-a567-0e02b2c3d48b', 'Public Transport', 'Transportation', 20, NOW(), NOW()),
    ('c47ac10b-58cc-4372-a567-0e02b2c3d48c', 'Taxi & Rideshare', 'Transportation', 19, NOW(), NOW()),
    ('c47ac10b-58cc-4372-a567-0e02b2c3d48d', 'Fuel', 'Transportation', 18, NOW(), NOW()),
    ('c47ac10b-58cc-4372-a567-0e02b2c3d48e', 'Parking', 'Transportation', 17, NOW(), NOW()),
    ('c47ac10b-58cc-4372-a567-0e02b2c3d48f', 'Car Maintenance', 'Transportation', 16, NOW(), NOW()),
    ('c47ac10b-58cc-4372-a567-0e02b2c3d490', 'Flight', 'Transportation', 15, NOW(), NOW()),
    ('c47ac10b-58cc-4372-a567-0e02b2c3d491', 'Train', 'Transportation', 14, NOW(), NOW()),
    
    -- Entertainment Categories
    ('c47ac10b-58cc-4372-a567-0e02b2c3d492', 'Movies & Cinema', 'Entertainment', 25, NOW(), NOW()),
    ('c47ac10b-58cc-4372-a567-0e02b2c3d493', 'Streaming Services', 'Entertainment', 24, NOW(), NOW()),
    ('c47ac10b-58cc-4372-a567-0e02b2c3d494', 'Music', 'Entertainment', 23, NOW(), NOW()),
    ('c47ac10b-58cc-4372-a567-0e02b2c3d495', 'Games', 'Entertainment', 22, NOW(), NOW()),
    ('c47ac10b-58cc-4372-a567-0e02b2c3d496', 'Sports & Fitness', 'Entertainment', 21, NOW(), NOW()),
    ('c47ac10b-58cc-4372-a567-0e02b2c3d497', 'Events & Concerts', 'Entertainment', 20, NOW(), NOW()),
    
    -- Bills & Utilities Categories
    ('c47ac10b-58cc-4372-a567-0e02b2c3d498', 'Electricity', 'Bills & Utilities', 30, NOW(), NOW()),
    ('c47ac10b-58cc-4372-a567-0e02b2c3d499', 'Water', 'Bills & Utilities', 29, NOW(), NOW()),
    ('c47ac10b-58cc-4372-a567-0e02b2c3d49a', 'Gas', 'Bills & Utilities', 28, NOW(), NOW()),
    ('c47ac10b-58cc-4372-a567-0e02b2c3d49b', 'Internet', 'Bills & Utilities', 27, NOW(), NOW()),
    ('c47ac10b-58cc-4372-a567-0e02b2c3d49c', 'Mobile Phone', 'Bills & Utilities', 26, NOW(), NOW()),
    ('c47ac10b-58cc-4372-a567-0e02b2c3d49d', 'Insurance', 'Bills & Utilities', 25, NOW(), NOW()),
    ('c47ac10b-58cc-4372-a567-0e02b2c3d49e', 'Rent', 'Bills & Utilities', 24, NOW(), NOW()),
    
    -- Health & Wellness Categories
    ('c47ac10b-58cc-4372-a567-0e02b2c3d49f', 'Doctor Visits', 'Health & Wellness', 35, NOW(), NOW()),
    ('c47ac10b-58cc-4372-a567-0e02b2c3d4a0', 'Medications', 'Health & Wellness', 34, NOW(), NOW()),
    ('c47ac10b-58cc-4372-a567-0e02b2c3d4a1', 'Gym & Fitness', 'Health & Wellness', 33, NOW(), NOW()),
    ('c47ac10b-58cc-4372-a567-0e02b2c3d4a2', 'Beauty & Personal Care', 'Health & Wellness', 32, NOW(), NOW()),
    
    -- Income Categories
    ('c47ac10b-58cc-4372-a567-0e02b2c3d4a3', 'Salary', 'Income', 50, NOW(), NOW()),
    ('c47ac10b-58cc-4372-a567-0e02b2c3d4a4', 'Freelance', 'Income', 49, NOW(), NOW()),
    ('c47ac10b-58cc-4372-a567-0e02b2c3d4a5', 'Investment Returns', 'Income', 48, NOW(), NOW()),
    ('c47ac10b-58cc-4372-a567-0e02b2c3d4a6', 'Bonus', 'Income', 47, NOW(), NOW()),
    ('c47ac10b-58cc-4372-a567-0e02b2c3d4a7', 'Gift Money', 'Income', 46, NOW(), NOW()),
    ('c47ac10b-58cc-4372-a567-0e02b2c3d4a8', 'Refunds', 'Income', 45, NOW(), NOW()),
    
    -- Education Categories
    ('c47ac10b-58cc-4372-a567-0e02b2c3d4a9', 'Courses & Training', 'Education', 40, NOW(), NOW()),
    ('c47ac10b-58cc-4372-a567-0e02b2c3d4aa', 'Books & Materials', 'Education', 39, NOW(), NOW()),
    ('c47ac10b-58cc-4372-a567-0e02b2c3d4ab', 'School Fees', 'Education', 38, NOW(), NOW()),
    
    -- Financial Categories
    ('c47ac10b-58cc-4372-a567-0e02b2c3d4ac', 'Bank Fees', 'Financial', 45, NOW(), NOW()),
    ('c47ac10b-58cc-4372-a567-0e02b2c3d4ad', 'Loan Payments', 'Financial', 44, NOW(), NOW()),
    ('c47ac10b-58cc-4372-a567-0e02b2c3d4ae', 'Credit Card Payments', 'Financial', 43, NOW(), NOW()),
    ('c47ac10b-58cc-4372-a567-0e02b2c3d4af', 'Investments', 'Financial', 42, NOW(), NOW()),
    ('c47ac10b-58cc-4372-a567-0e02b2c3d4b0', 'Savings', 'Financial', 41, NOW(), NOW()),
    ('c47ac10b-58cc-4372-a567-0e02b2c3d4b1', 'Transfers', 'Financial', 40, NOW(), NOW()),
    
    -- Travel Categories
    ('c47ac10b-58cc-4372-a567-0e02b2c3d4b2', 'Hotels', 'Travel', 55, NOW(), NOW()),
    ('c47ac10b-58cc-4372-a567-0e02b2c3d4b3', 'Vacation', 'Travel', 54, NOW(), NOW()),
    ('c47ac10b-58cc-4372-a567-0e02b2c3d4b4', 'Business Travel', 'Travel', 53, NOW(), NOW()),
    
    -- Other Categories
    ('c47ac10b-58cc-4372-a567-0e02b2c3d4b5', 'Miscellaneous', 'Other', 1, NOW(), NOW()),
    ('c47ac10b-58cc-4372-a567-0e02b2c3d4b6', 'Cash Withdrawal', 'Other', 2, NOW(), NOW()),
    ('c47ac10b-58cc-4372-a567-0e02b2c3d4b7', 'Donations', 'Other', 3, NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- Insert sample balances/accounts
INSERT INTO balance (id, group_id, user_id, currency, title, description, created_at, updated_at) VALUES
    ('847ac10b-58cc-4372-a567-0e02b2c3d479', 'b47ac10b-58cc-4372-a567-0e02b2c3d479', '847ac10b-58cc-4372-a567-0e02b2c3d479', 'EUR', 'BBVA Main Account', 'Primary checking account at BBVA', NOW(), NOW()),
    ('847ac10b-58cc-4372-a567-0e02b2c3d480', 'b47ac10b-58cc-4372-a567-0e02b2c3d479', '847ac10b-58cc-4372-a567-0e02b2c3d479', 'EUR', 'Santander Savings', 'Savings account at Santander', NOW(), NOW()),
    ('847ac10b-58cc-4372-a567-0e02b2c3d481', 'b47ac10b-58cc-4372-a567-0e02b2c3d479', '847ac10b-58cc-4372-a567-0e02b2c3d479', 'EUR', 'ING Orange Account', 'Orange account at ING Bank', NOW(), NOW()),
    ('847ac10b-58cc-4372-a567-0e02b2c3d482', 'b47ac10b-58cc-4372-a567-0e02b2c3d479', '847ac10b-58cc-4372-a567-0e02b2c3d479', 'EUR', 'Cash Wallet', 'Physical cash wallet', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- Insert sample transactions
INSERT INTO transaction (id, group_id, user_id, balance_id, merchant_id, type, approved_at, transacted_at, created_at, updated_at) VALUES
    -- Groceries at Mercadona
    ('a47ac10b-58cc-4372-a567-0e02b2c3d479', 'b47ac10b-58cc-4372-a567-0e02b2c3d479', '847ac10b-58cc-4372-a567-0e02b2c3d479', '847ac10b-58cc-4372-a567-0e02b2c3d479', 'f47ac10b-58cc-4372-a567-0e02b2c3d479', 'expense', NOW(), NOW() - INTERVAL '1 day', NOW(), NOW()),
    
    -- Coffee at Starbucks
    ('a47ac10b-58cc-4372-a567-0e02b2c3d480', 'b47ac10b-58cc-4372-a567-0e02b2c3d479', '847ac10b-58cc-4372-a567-0e02b2c3d479', '847ac10b-58cc-4372-a567-0e02b2c3d479', '6ba7b810-9dad-11d1-80b4-00c04fd430c8', 'expense', NOW(), NOW() - INTERVAL '2 days', NOW(), NOW()),
    
    -- Taxi ride with Uber
    ('a47ac10b-58cc-4372-a567-0e02b2c3d481', 'b47ac10b-58cc-4372-a567-0e02b2c3d479', '847ac10b-58cc-4372-a567-0e02b2c3d479', '847ac10b-58cc-4372-a567-0e02b2c3d479', '6ba7b815-9dad-11d1-80b4-00c04fd430c8', 'expense', NOW(), NOW() - INTERVAL '3 days', NOW(), NOW()),
    
    -- Monthly salary
    ('a47ac10b-58cc-4372-a567-0e02b2c3d482', 'b47ac10b-58cc-4372-a567-0e02b2c3d479', '847ac10b-58cc-4372-a567-0e02b2c3d479', '847ac10b-58cc-4372-a567-0e02b2c3d479', NULL, 'income', NOW(), NOW() - INTERVAL '7 days', NOW(), NOW()),
    
    -- Shopping at El Corte Inglés
    ('a47ac10b-58cc-4372-a567-0e02b2c3d483', 'b47ac10b-58cc-4372-a567-0e02b2c3d479', '847ac10b-58cc-4372-a567-0e02b2c3d479', '847ac10b-58cc-4372-a567-0e02b2c3d479', 'f47ac10b-58cc-4372-a567-0e02b2c3d480', 'expense', NOW(), NOW() - INTERVAL '4 days', NOW(), NOW()),
    
    -- Netflix subscription
    ('a47ac10b-58cc-4372-a567-0e02b2c3d484', 'b47ac10b-58cc-4372-a567-0e02b2c3d479', '847ac10b-58cc-4372-a567-0e02b2c3d479', '847ac10b-58cc-4372-a567-0e02b2c3d479', '6ba7b81a-9dad-11d1-80b4-00c04fd430c8', 'expense', NOW(), NOW() - INTERVAL '5 days', NOW(), NOW()),
    
    -- Electronics at MediaMarkt
    ('a47ac10b-58cc-4372-a567-0e02b2c3d485', 'b47ac10b-58cc-4372-a567-0e02b2c3d479', '847ac10b-58cc-4372-a567-0e02b2c3d479', '847ac10b-58cc-4372-a567-0e02b2c3d479', 'f47ac10b-58cc-4372-a567-0e02b2c3d484', 'expense', NOW(), NOW() - INTERVAL '10 days', NOW(), NOW()),
    
    -- Clothing at Zara
    ('a47ac10b-58cc-4372-a567-0e02b2c3d486', 'b47ac10b-58cc-4372-a567-0e02b2c3d479', '847ac10b-58cc-4372-a567-0e02b2c3d479', '847ac10b-58cc-4372-a567-0e02b2c3d479', 'f47ac10b-58cc-4372-a567-0e02b2c3d483', 'expense', NOW(), NOW() - INTERVAL '6 days', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- Insert sample transaction entries
INSERT INTO transaction_entry (id, transaction_id, description, amount, category_id, created_at, updated_at) VALUES
    -- Mercadona groceries transaction
    ('e47ac10b-58cc-4372-a567-0e02b2c3d479', 'a47ac10b-58cc-4372-a567-0e02b2c3d479', 'Weekly groceries at Mercadona', -45.67, 'c47ac10b-58cc-4372-a567-0e02b2c3d479', NOW(), NOW()),
    
    -- Starbucks coffee
    ('e47ac10b-58cc-4372-a567-0e02b2c3d480', 'a47ac10b-58cc-4372-a567-0e02b2c3d480', 'Morning latte', -5.50, 'c47ac10b-58cc-4372-a567-0e02b2c3d482', NOW(), NOW()),
    
    -- Uber ride
    ('e47ac10b-58cc-4372-a567-0e02b2c3d481', 'a47ac10b-58cc-4372-a567-0e02b2c3d481', 'Ride to airport', -25.00, 'c47ac10b-58cc-4372-a567-0e02b2c3d48c', NOW(), NOW()),
    
    -- Monthly salary
    ('e47ac10b-58cc-4372-a567-0e02b2c3d482', 'a47ac10b-58cc-4372-a567-0e02b2c3d482', 'Monthly salary payment', 3500.00, 'c47ac10b-58cc-4372-a567-0e02b2c3d4a3', NOW(), NOW()),
    
    -- El Corte Inglés shopping - multiple items
    ('e47ac10b-58cc-4372-a567-0e02b2c3d483', 'a47ac10b-58cc-4372-a567-0e02b2c3d483', 'Home decoration items', -89.99, 'c47ac10b-58cc-4372-a567-0e02b2c3d487', NOW(), NOW()),
    ('e47ac10b-58cc-4372-a567-0e02b2c3d484', 'a47ac10b-58cc-4372-a567-0e02b2c3d483', 'Books and magazines', -25.50, 'c47ac10b-58cc-4372-a567-0e02b2c3d488', NOW(), NOW()),
    
    -- Netflix subscription
    ('e47ac10b-58cc-4372-a567-0e02b2c3d485', 'a47ac10b-58cc-4372-a567-0e02b2c3d484', 'Netflix monthly subscription', -12.99, 'c47ac10b-58cc-4372-a567-0e02b2c3d493', NOW(), NOW()),
    
    -- MediaMarkt electronics
    ('e47ac10b-58cc-4372-a567-0e02b2c3d486', 'a47ac10b-58cc-4372-a567-0e02b2c3d485', 'Wireless headphones', -149.99, 'c47ac10b-58cc-4372-a567-0e02b2c3d486', NOW(), NOW()),
    ('e47ac10b-58cc-4372-a567-0e02b2c3d487', 'a47ac10b-58cc-4372-a567-0e02b2c3d485', 'Phone charging cable', -19.99, 'c47ac10b-58cc-4372-a567-0e02b2c3d486', NOW(), NOW()),
    
    -- Zara clothing
    ('e47ac10b-58cc-4372-a567-0e02b2c3d488', 'a47ac10b-58cc-4372-a567-0e02b2c3d486', 'Summer t-shirt', -29.95, 'c47ac10b-58cc-4372-a567-0e02b2c3d485', NOW(), NOW()),
    ('e47ac10b-58cc-4372-a567-0e02b2c3d489', 'a47ac10b-58cc-4372-a567-0e02b2c3d486', 'Jeans', -59.95, 'c47ac10b-58cc-4372-a567-0e02b2c3d485', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- Verify the data
SELECT 'Merchants' as table_name, COUNT(*) as record_count FROM merchant
UNION ALL
SELECT 'Categories' as table_name, COUNT(*) as record_count FROM category
UNION ALL
SELECT 'Transactions' as table_name, COUNT(*) as record_count FROM transaction
UNION ALL
SELECT 'Transaction Entries' as table_name, COUNT(*) as record_count FROM transaction_entry;
