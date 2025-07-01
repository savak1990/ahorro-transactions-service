-- seed_balances.sql: Sample balances/accounts data for Aurora PostgreSQL Ahorro Transactions Service

INSERT INTO balance (id, group_id, user_id, currency, title, description, rank, created_at, updated_at) VALUES
    ('aa111111-1111-1111-1111-111111111111', '88aa1100-0011-2233-4455-667788990011', '99bb2200-0011-2233-4455-667788990011', 'EUR', 'BBVA Main Account', 'Primary checking account at BBVA', 1, NOW(), NOW()),
    ('bb222222-2222-2222-2222-222222222222', '88aa1100-0011-2233-4455-667788990011', '99bb2200-0011-2233-4455-667788990011', 'EUR', 'Santander Savings', 'Savings account at Santander', 2, NOW(), NOW()),
    ('cc333333-3333-3333-3333-333333333333', '88aa1100-0011-2233-4455-667788990011', '99bb2200-0011-2233-4455-667788990011', 'EUR', 'ING Orange Account', 'Orange account at ING Bank', 3, NOW(), NOW()),
    ('dd444444-4444-4444-4444-444444444444', '88aa1100-0011-2233-4455-667788990011', '99bb2200-0011-2233-4455-667788990011', 'EUR', 'Cash Wallet', 'Physical cash wallet', 4, NOW(), NOW()),
    ('ee555555-5555-5555-5555-555555555555', '88aa1100-0011-2233-4455-667788990011', '99bb2200-0011-2233-4455-667788990011', 'USD', 'USD Investment Account', 'Investment account in US Dollars', 5, NOW(), NOW()),
    ('ff666666-6666-6666-6666-666666666666', '88aa1100-0011-2233-4455-667788990011', '99bb2200-0011-2233-4455-667788990011', 'GBP', 'UK Travel Account', 'Account for UK travel expenses', 6, NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
