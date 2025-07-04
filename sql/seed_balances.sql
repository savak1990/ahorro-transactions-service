-- seed_balances.sql: Sample balances/accounts data for Aurora PostgreSQL Ahorro Transactions Service
-- Note: Higher rank numbers indicate more frequently used accounts
-- Rank uses steps of 10 to allow easy insertion of new balances between existing ones

INSERT INTO balance (id, group_id, user_id, currency, title, description, rank, created_at, updated_at) VALUES
    -- Primary user balances (99bb2200-0011-2233-4455-667788990011)
    ('ba001111-1111-1111-1111-111111111111', '88aa1100-0011-2233-4455-667788990011', '99bb2200-0011-2233-4455-667788990011', 'EUR', 'BBVA Main Account', 'Primary checking account at BBVA', 60, NOW(), NOW()),
    ('ba004444-4444-4444-4444-444444444444', '88aa1100-0011-2233-4455-667788990011', '99bb2200-0011-2233-4455-667788990011', 'EUR', 'Cash Wallet', 'Physical cash wallet', 50, NOW(), NOW()),
    ('ba002222-2222-2222-2222-222222222222', '88aa1100-0011-2233-4455-667788990011', '99bb2200-0011-2233-4455-667788990011', 'EUR', 'Santander Savings', 'Savings account at Santander', 40, NOW(), NOW()),
    ('ba003333-3333-3333-3333-333333333333', '88aa1100-0011-2233-4455-667788990011', '99bb2200-0011-2233-4455-667788990011', 'EUR', 'ING Orange Account', 'Orange account at ING Bank', 30, NOW(), NOW()),
    ('ba005555-5555-5555-5555-555555555555', '88aa1100-0011-2233-4455-667788990011', '99bb2200-0011-2233-4455-667788990011', 'USD', 'USD Investment Account', 'Investment account in US Dollars', 20, NOW(), NOW()),
    ('ba006666-6666-6666-6666-666666666666', '88aa1100-0011-2233-4455-667788990011', '99bb2200-0011-2233-4455-667788990011', 'GBP', 'UK Travel Account', 'Account for UK travel expenses', 10, NOW(), NOW()),
    
    -- Secondary user balance (99bb3300-0011-2233-4455-667788990022) - same group
    ('ba007777-7777-7777-7777-777777777777', '88aa1100-0011-2233-4455-667788990011', '99bb3300-0011-2233-4455-667788990022', 'EUR', 'Default Account', 'Default checking account for secondary user', 60, NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
