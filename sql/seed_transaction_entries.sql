-- seed_transaction_entries.sql: Sample transaction entries data for Aurora PostgreSQL Ahorro Transactions Service

INSERT INTO transaction_entry (id, transaction_id, description, amount, category_id, created_at, updated_at) VALUES
    -- Mercadona groceries transaction
    ('e1111111-1111-1111-1111-111111111111', 't1111111-1111-1111-1111-111111111111', 'Weekly groceries at Mercadona', -45.67, '11111111-1111-1111-1111-111111111111', NOW(), NOW()),
    
    -- Starbucks coffee
    ('e2222222-2222-2222-2222-222222222222', 't2222222-2222-2222-2222-222222222222', 'Morning latte', -5.50, '22222222-2222-2222-2222-222222222222', NOW(), NOW()),
    
    -- Uber ride
    ('e3333333-3333-3333-3333-333333333333', 't3333333-3333-3333-3333-333333333333', 'Ride to airport', -25.00, '33333333-3333-3333-3333-333333333333', NOW(), NOW()),
    
    -- Monthly salary
    ('e4444444-4444-4444-4444-444444444444', 't4444444-4444-4444-4444-444444444444', 'Monthly salary payment', 3500.00, '44444444-4444-4444-4444-444444444444', NOW(), NOW()),
    
    -- El Corte Inglés shopping - multiple items
    ('e5555555-5555-5555-5555-555555555555', 't5555555-5555-5555-5555-555555555555', 'Home decoration items', -89.99, '55555555-5555-5555-5555-555555555555', NOW(), NOW()),
    ('e5555556-5555-5555-5555-555555555555', 't5555555-5555-5555-5555-555555555555', 'Books and magazines', -25.50, '66666666-6666-6666-6666-666666666666', NOW(), NOW()),
    
    -- Netflix subscription
    ('e6666666-6666-6666-6666-666666666666', 't6666666-6666-6666-6666-666666666666', 'Netflix monthly subscription', -12.99, '77777777-7777-7777-7777-777777777777', NOW(), NOW()),
    
    -- MediaMarkt electronics
    ('e7777777-7777-7777-7777-777777777777', 't7777777-7777-7777-7777-777777777777', 'Wireless headphones', -149.99, '88888888-8888-8888-8888-888888888888', NOW(), NOW()),
    ('e7777778-7777-7777-7777-777777777777', 't7777777-7777-7777-7777-777777777777', 'Phone charging cable', -19.99, '88888888-8888-8888-8888-888888888888', NOW(), NOW()),
    
    -- Zara clothing
    ('e8888888-8888-8888-8888-888888888888', 't8888888-8888-8888-8888-888888888888', 'Summer t-shirt', -29.95, '99999999-9999-9999-9999-999999999999', NOW(), NOW()),
    ('e8888889-8888-8888-8888-888888888888', 't8888888-8888-8888-8888-888888888888', 'Jeans', -59.95, '99999999-9999-9999-9999-999999999999', NOW(), NOW()),
    
    -- Freelance income
    ('e9999999-9999-9999-9999-999999999999', 't9999999-9999-9999-9999-999999999999', 'Freelance project payment', 1200.00, 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', NOW(), NOW()),
    
    -- Investment returns in USD
    ('eaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'taaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'Monthly dividend payment', 250.00, 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', NOW(), NOW()),
    
    -- Travel expense in GBP
    ('ebbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', 'tbbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', 'Hotel accommodation in London', -180.00, 'cccccccc-cccc-cccc-cccc-cccccccccccc', NOW(), NOW()),
    
    -- Transfer between accounts (movement) - debit entry
    ('ecccccccc-cccc-cccc-cccc-cccccccccccc', 'tcccccccc-cccc-cccc-cccc-cccccccccccc', 'Transfer to savings account', -500.00, 'dddddddd-dddd-dddd-dddd-dddddddddddd', NOW(), NOW()),
    -- Transfer between accounts (movement) - credit entry
    ('ecccccccd-cccc-cccc-cccc-cccccccccccc', 'tcccccccc-cccc-cccc-cccc-cccccccccccc', 'Transfer from checking account', 500.00, 'dddddddd-dddd-dddd-dddd-dddddddddddd', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
