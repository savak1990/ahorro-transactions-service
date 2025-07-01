-- seed_categories.sql: Sample categories data for Aurora PostgreSQL Ahorro Transactions Service

INSERT INTO category (id, user_id, group_id, category_group_id, name, "group", description, rank, created_at, updated_at) VALUES
    -- Food & Dining Categories (linked to Food & Dining group)
    ('ca001111-1111-1111-1111-111111111111', '99bb2200-0011-2233-4455-667788990011', '88aa1100-0011-2233-4455-667788990011', 'c9001234-7890-1234-5678-901234567890', 'Groceries', 'Food & Dining', 'Grocery shopping and food supplies', 10, NOW(), NOW()),
    ('ca002222-2222-2222-2222-222222222222', '99bb2200-0011-2233-4455-667788990011', '88aa1100-0011-2233-4455-667788990011', 'c9001234-7890-1234-5678-901234567890', 'Coffee & Tea', 'Food & Dining', 'Coffee shops and tea houses', 9, NOW(), NOW()),
    ('ca003333-3333-3333-3333-333333333333', '99bb2200-0011-2233-4455-667788990011', '88aa1100-0011-2233-4455-667788990011', 'c9001234-7890-1234-5678-901234567890', 'Restaurants', 'Food & Dining', 'Dining out at restaurants', 8, NOW(), NOW()),
    ('ca004444-4444-4444-4444-444444444444', '99bb2200-0011-2233-4455-667788990011', '88aa1100-0011-2233-4455-667788990011', 'c9001234-7890-1234-5678-901234567890', 'Fast Food', 'Food & Dining', 'Quick service restaurants', 7, NOW(), NOW()),
    
    -- Shopping Categories (linked to Shopping group)
    ('ca005555-5555-5555-5555-555555555555', '99bb2200-0011-2233-4455-667788990011', '88aa1100-0011-2233-4455-667788990011', 'c9012345-7890-1234-5678-901234567890', 'Home & Garden', 'Shopping', 'Home improvement and garden supplies', 15, NOW(), NOW()),
    ('ca006666-6666-6666-6666-666666666666', '99bb2200-0011-2233-4455-667788990011', '88aa1100-0011-2233-4455-667788990011', 'c9012345-7890-1234-5678-901234567890', 'Books & Magazines', 'Shopping', 'Books, magazines, and reading materials', 14, NOW(), NOW()),
    ('ca007777-7777-7777-7777-777777777777', '99bb2200-0011-2233-4455-667788990011', '88aa1100-0011-2233-4455-667788990011', 'c9034567-7890-1234-5678-901234567890', 'Streaming Services', 'Entertainment', 'Netflix, Amazon Prime, etc.', 13, NOW(), NOW()),
    ('ca008888-8888-8888-8888-888888888888', '99bb2200-0011-2233-4455-667788990011', '88aa1100-0011-2233-4455-667788990011', 'c9012345-7890-1234-5678-901234567890', 'Electronics', 'Shopping', 'Electronic devices and gadgets', 12, NOW(), NOW()),
    ('ca009999-9999-9999-9999-999999999999', '99bb2200-0011-2233-4455-667788990011', '88aa1100-0011-2233-4455-667788990011', 'c9012345-7890-1234-5678-901234567890', 'Clothing', 'Shopping', 'Clothing and fashion', 11, NOW(), NOW()),
    
    -- Income & Financial Categories
    ('ca00aaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', '99bb2200-0011-2233-4455-667788990011', '88aa1100-0011-2233-4455-667788990011', 'c9067890-7890-1234-5678-901234567890', 'Freelance', 'Income', 'Freelance work and consulting', 20, NOW(), NOW()),
    ('ca00bbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', '99bb2200-0011-2233-4455-667788990011', '88aa1100-0011-2233-4455-667788990011', 'c9067890-7890-1234-5678-901234567890', 'Investment Returns', 'Income', 'Dividends and investment income', 19, NOW(), NOW()),
    ('ca00cccc-cccc-cccc-cccc-cccccccccccc', '99bb2200-0011-2233-4455-667788990011', '88aa1100-0011-2233-4455-667788990011', 'c9023456-7890-1234-5678-901234567890', 'Hotels', 'Travel', 'Hotel accommodations', 18, NOW(), NOW()),
    ('ca00dddd-dddd-dddd-dddd-dddddddddddd', '99bb2200-0011-2233-4455-667788990011', '88aa1100-0011-2233-4455-667788990011', 'c9078901-7890-1234-5678-901234567890', 'Transfers', 'Financial', 'Money transfers between accounts', 17, NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
