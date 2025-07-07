-- seed_transaction_entries.sql: Sample transaction entries data for Aurora PostgreSQL Ahorro Transactions Service

INSERT INTO transaction_entry (id, transaction_id, description, amount, category_id, created_at, updated_at) VALUES
    -- Mercadona groceries transaction
    ('7e001111-1111-1111-1111-111111111111', '7a001111-1111-1111-1111-111111111111', 'Weekly groceries at Mercadona', 45.67, 'ca001111-1111-1111-1111-111111111111', NOW(), NOW()),
    
    -- Starbucks coffee
    ('7e002222-2222-2222-2222-222222222222', '7a002222-2222-2222-2222-222222222222', 'Morning latte', 5.50, 'ca002222-2222-2222-2222-222222222222', NOW(), NOW()),
    
    -- Uber ride
    ('7e003333-3333-3333-3333-333333333333', '7a003333-3333-3333-3333-333333333333', 'Ride to airport', 25.00, 'ca003333-3333-3333-3333-333333333333', NOW(), NOW()),
    
    -- Monthly salary
    ('7e004444-4444-4444-4444-444444444444', '7a004444-4444-4444-4444-444444444444', 'Monthly salary payment', 3500.00, 'ca004444-4444-4444-4444-444444444444', NOW(), NOW()),
    
    -- El Corte Inglés shopping - multiple items
    ('7e005555-5555-5555-5555-555555555555', '7a005555-5555-5555-5555-555555555555', 'Home decoration items', 89.99, 'ca005555-5555-5555-5555-555555555555', NOW(), NOW()),
    ('7e005556-5555-5555-5555-555555555555', '7a005555-5555-5555-5555-555555555555', 'Books and magazines', 25.50, 'ca006666-6666-6666-6666-666666666666', NOW(), NOW()),
    
    -- Netflix subscription
    ('7e006666-6666-6666-6666-666666666666', '7a006666-6666-6666-6666-666666666666', 'Netflix monthly subscription', 12.99, 'ca007777-7777-7777-7777-777777777777', NOW(), NOW()),
    
    -- MediaMarkt electronics
    ('7e007777-7777-7777-7777-777777777777', '7a007777-7777-7777-7777-777777777777', 'Wireless headphones', 149.99, 'ca008888-8888-8888-8888-888888888888', NOW(), NOW()),
    ('7e007778-7777-7777-7777-777777777777', '7a007777-7777-7777-7777-777777777777', 'Phone charging cable', 19.99, 'ca008888-8888-8888-8888-888888888888', NOW(), NOW()),
    
    -- Zara clothing
    ('7e008888-8888-8888-8888-888888888888', '7a008888-8888-8888-8888-888888888888', 'Summer t-shirt', 29.95, 'ca009999-9999-9999-9999-999999999999', NOW(), NOW()),
    ('7e008889-8888-8888-8888-888888888888', '7a008888-8888-8888-8888-888888888888', 'Jeans', 59.95, 'ca009999-9999-9999-9999-999999999999', NOW(), NOW()),
    
    -- Freelance income
    ('7e009999-9999-9999-9999-999999999999', '7a009999-9999-9999-9999-999999999999', 'Freelance project payment', 1200.00, 'ca00aaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', NOW(), NOW()),
    
    -- Investment returns in USD
    ('7eaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', '7aaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'Monthly dividend payment', 250.00, 'ca00bbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', NOW(), NOW()),
    
    -- Travel expense in GBP
    ('7ebbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', '7abbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', 'Hotel accommodation in London', 180.00, 'ca00cccc-cccc-cccc-cccc-cccccccccccc', NOW(), NOW()),
    
    -- Transfer between accounts (movement) - debit entry
    ('7ecccccc-cccc-cccc-cccc-cccccccccccc', '7acccccc-cccc-cccc-cccc-cccccccccccc', 'Transfer to savings account', 500.00, 'ca00dddd-dddd-dddd-dddd-dddddddddddd', NOW(), NOW()),
    -- Transfer between accounts (movement) - credit entry
    ('7ecccccd-cccc-cccc-cccc-cccccccccccc', '7acccccc-cccc-cccc-cccc-cccccccccccc', 'Transfer from checking account', 500.00, 'ca00dddd-dddd-dddd-dddd-dddddddddddd', NOW(), NOW()),
    
    -- Transaction entries for secondary user (99bb3300-0011-2233-4455-667788990022)
    -- Grocery shopping at Dia Supermercado
    ('7e00aa00-0000-1111-2222-333344445555', '7a00aa00-0000-1111-2222-333344445555', 'Weekly grocery shopping', 32.45, 'caa00000-0000-1111-2222-333344445555', NOW(), NOW()),
    ('7e00aa01-0000-1111-2222-333344445555', '7a00aa00-0000-1111-2222-333344445555', 'Personal care items', 15.20, 'ca440000-0000-1111-2222-333344445564', NOW(), NOW()),
    
    -- Coffee at Local Café
    ('7e00bb00-0000-1111-2222-333344445556', '7a00bb00-0000-1111-2222-333344445556', 'Morning coffee and pastry', 6.80, 'caa00000-0000-1111-2222-333344445555', NOW(), NOW()),
    
    -- Part-time job income
    ('7e00cc00-0000-1111-2222-333344445557', '7a00cc00-0000-1111-2222-333344445557', 'Part-time work payment', 450.00, 'ca110000-0000-1111-2222-333344445561', NOW(), NOW()),
    
    -- ATM withdrawal
    ('7e00dd00-0000-1111-2222-333344445558', '7a00dd00-0000-1111-2222-333344445558', 'ATM withdrawal fee', 2.50, 'ca220000-0000-1111-2222-333344445562', NOW(), NOW()),
    
    -- Mixed shopping transaction
    ('7e00ee00-0000-1111-2222-333344445559', '7a00ee00-0000-1111-2222-333344445559', 'Pharmacy medicine', 18.75, 'caff0000-0000-1111-2222-333344445560', NOW(), NOW()),
    ('7e00ee01-0000-1111-2222-333344445559', '7a00ee00-0000-1111-2222-333344445559', 'Educational book', 25.90, 'ca330000-0000-1111-2222-333344445563', NOW(), NOW()),
    ('7e00ee02-0000-1111-2222-333344445559', '7a00ee00-0000-1111-2222-333344445559', 'Phone bill payment', 35.00, 'caee0000-0000-1111-2222-333344445559', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
