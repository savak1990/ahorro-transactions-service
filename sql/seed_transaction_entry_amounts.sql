-- seed_transaction_entry_amounts.sql: Multi-currency amounts for transaction entries
-- Creates TransactionEntryAmount records for supported currencies (EUR, USD, GBP)
-- Uses realistic exchange rates as of August 2025

-- Exchange rates used (as of August 2025):
-- EUR/USD: 1.09 (1 EUR = 1.09 USD)
-- EUR/GBP: 0.86 (1 EUR = 0.86 GBP)
-- USD/EUR: 0.92 (1 USD = 0.92 EUR)
-- USD/GBP: 0.79 (1 USD = 0.79 GBP)
-- GBP/EUR: 1.16 (1 GBP = 1.16 EUR)
-- GBP/USD: 1.27 (1 GBP = 1.27 USD)

INSERT INTO transaction_entry_amount (transaction_entry_id, currency, amount, exchange_rate, created_at, updated_at) VALUES
    -- Mercadona groceries (base: EUR 45.67)
    ('7e001111-1111-1111-1111-111111111111', 'EUR', 4567, 1.0, NOW(), NOW()),
    ('7e001111-1111-1111-1111-111111111111', 'USD', 4978, 1.09, NOW(), NOW()),
    ('7e001111-1111-1111-1111-111111111111', 'GBP', 3928, 0.86, NOW(), NOW()),
    
    -- Starbucks coffee (base: EUR 5.50)
    ('7e002222-2222-2222-2222-222222222222', 'EUR', 550, 1.0, NOW(), NOW()),
    ('7e002222-2222-2222-2222-222222222222', 'USD', 600, 1.09, NOW(), NOW()),
    ('7e002222-2222-2222-2222-222222222222', 'GBP', 473, 0.86, NOW(), NOW()),
    
    -- Uber ride (base: EUR 25.00)
    ('7e003333-3333-3333-3333-333333333333', 'EUR', 2500, 1.0, NOW(), NOW()),
    ('7e003333-3333-3333-3333-333333333333', 'USD', 2725, 1.09, NOW(), NOW()),
    ('7e003333-3333-3333-3333-333333333333', 'GBP', 2150, 0.86, NOW(), NOW()),
    
    -- Monthly salary (base: EUR 3500.00)
    ('7e004444-4444-4444-4444-444444444444', 'EUR', 350000, 1.0, NOW(), NOW()),
    ('7e004444-4444-4444-4444-444444444444', 'USD', 381500, 1.09, NOW(), NOW()),
    ('7e004444-4444-4444-4444-444444444444', 'GBP', 301000, 0.86, NOW(), NOW()),
    
    -- El Corte Inglés - home decoration (base: EUR 89.99)
    ('7e005555-5555-5555-5555-555555555555', 'EUR', 8999, 1.0, NOW(), NOW()),
    ('7e005555-5555-5555-5555-555555555555', 'USD', 9809, 1.09, NOW(), NOW()),
    ('7e005555-5555-5555-5555-555555555555', 'GBP', 7739, 0.86, NOW(), NOW()),
    
    -- El Corte Inglés - books (base: EUR 25.50)
    ('7e005556-5555-5555-5555-555555555555', 'EUR', 2550, 1.0, NOW(), NOW()),
    ('7e005556-5555-5555-5555-555555555555', 'USD', 2780, 1.09, NOW(), NOW()),
    ('7e005556-5555-5555-5555-555555555555', 'GBP', 2193, 0.86, NOW(), NOW()),
    
    -- Netflix subscription (base: EUR 12.99)
    ('7e006666-6666-6666-6666-666666666666', 'EUR', 1299, 1.0, NOW(), NOW()),
    ('7e006666-6666-6666-6666-666666666666', 'USD', 1416, 1.09, NOW(), NOW()),
    ('7e006666-6666-6666-6666-666666666666', 'GBP', 1117, 0.86, NOW(), NOW()),
    
    -- MediaMarkt headphones (base: EUR 149.99)
    ('7e007777-7777-7777-7777-777777777777', 'EUR', 14999, 1.0, NOW(), NOW()),
    ('7e007777-7777-7777-7777-777777777777', 'USD', 16349, 1.09, NOW(), NOW()),
    ('7e007777-7777-7777-7777-777777777777', 'GBP', 12899, 0.86, NOW(), NOW()),
    
    -- MediaMarkt cable (base: EUR 19.99)
    ('7e007778-7777-7777-7777-777777777777', 'EUR', 1999, 1.0, NOW(), NOW()),
    ('7e007778-7777-7777-7777-777777777777', 'USD', 2179, 1.09, NOW(), NOW()),
    ('7e007778-7777-7777-7777-777777777777', 'GBP', 1719, 0.86, NOW(), NOW()),
    
    -- Zara t-shirt (base: EUR 29.95)
    ('7e008888-8888-8888-8888-888888888888', 'EUR', 2995, 1.0, NOW(), NOW()),
    ('7e008888-8888-8888-8888-888888888888', 'USD', 3264, 1.09, NOW(), NOW()),
    ('7e008888-8888-8888-8888-888888888888', 'GBP', 2576, 0.86, NOW(), NOW()),
    
    -- Zara jeans (base: EUR 59.95)
    ('7e008889-8888-8888-8888-888888888888', 'EUR', 5995, 1.0, NOW(), NOW()),
    ('7e008889-8888-8888-8888-888888888888', 'USD', 6535, 1.09, NOW(), NOW()),
    ('7e008889-8888-8888-8888-888888888888', 'GBP', 5156, 0.86, NOW(), NOW()),
    
    -- Freelance income (base: EUR 1200.00)
    ('7e009999-9999-9999-9999-999999999999', 'EUR', 120000, 1.0, NOW(), NOW()),
    ('7e009999-9999-9999-9999-999999999999', 'USD', 130800, 1.09, NOW(), NOW()),
    ('7e009999-9999-9999-9999-999999999999', 'GBP', 103200, 0.86, NOW(), NOW()),
    
    -- Investment returns - USD base (base: USD 250.00)
    ('7eaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'USD', 25000, 1.0, NOW(), NOW()),
    ('7eaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'EUR', 23000, 0.92, NOW(), NOW()),
    ('7eaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'GBP', 19750, 0.79, NOW(), NOW()),
    
    -- Travel expense - GBP base (base: GBP 180.00)
    ('7ebbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', 'GBP', 18000, 1.0, NOW(), NOW()),
    ('7ebbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', 'EUR', 20880, 1.16, NOW(), NOW()),
    ('7ebbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', 'USD', 22860, 1.27, NOW(), NOW()),
    
    -- Transfer to savings - debit (base: EUR 500.00)
    ('7ecccccc-cccc-cccc-cccc-cccccccccccc', 'EUR', 50000, 1.0, NOW(), NOW()),
    ('7ecccccc-cccc-cccc-cccc-cccccccccccc', 'USD', 54500, 1.09, NOW(), NOW()),
    ('7ecccccc-cccc-cccc-cccc-cccccccccccc', 'GBP', 43000, 0.86, NOW(), NOW()),
    
    -- Transfer from checking - credit (base: EUR 500.00)
    ('7ecccccd-cccc-cccc-cccc-cccccccccccc', 'EUR', 50000, 1.0, NOW(), NOW()),
    ('7ecccccd-cccc-cccc-cccc-cccccccccccc', 'USD', 54500, 1.09, NOW(), NOW()),
    ('7ecccccd-cccc-cccc-cccc-cccccccccccc', 'GBP', 43000, 0.86, NOW(), NOW()),
    
    -- Secondary user transactions
    -- Grocery shopping at Dia (base: EUR 32.45)
    ('7e00aa00-0000-1111-2222-333344445555', 'EUR', 3245, 1.0, NOW(), NOW()),
    ('7e00aa00-0000-1111-2222-333344445555', 'USD', 3537, 1.09, NOW(), NOW()),
    ('7e00aa00-0000-1111-2222-333344445555', 'GBP', 2791, 0.86, NOW(), NOW()),
    
    -- Personal care items (base: EUR 15.20)
    ('7e00aa01-0000-1111-2222-333344445555', 'EUR', 1520, 1.0, NOW(), NOW()),
    ('7e00aa01-0000-1111-2222-333344445555', 'USD', 1657, 1.09, NOW(), NOW()),
    ('7e00aa01-0000-1111-2222-333344445555', 'GBP', 1307, 0.86, NOW(), NOW()),
    
    -- Coffee at Local Café (base: EUR 6.80)
    ('7e00bb00-0000-1111-2222-333344445556', 'EUR', 680, 1.0, NOW(), NOW()),
    ('7e00bb00-0000-1111-2222-333344445556', 'USD', 741, 1.09, NOW(), NOW()),
    ('7e00bb00-0000-1111-2222-333344445556', 'GBP', 585, 0.86, NOW(), NOW()),
    
    -- Part-time work payment (base: EUR 450.00)
    ('7e00cc00-0000-1111-2222-333344445557', 'EUR', 45000, 1.0, NOW(), NOW()),
    ('7e00cc00-0000-1111-2222-333344445557', 'USD', 49050, 1.09, NOW(), NOW()),
    ('7e00cc00-0000-1111-2222-333344445557', 'GBP', 38700, 0.86, NOW(), NOW()),
    
    -- ATM withdrawal fee (base: EUR 2.50)
    ('7e00dd00-0000-1111-2222-333344445558', 'EUR', 250, 1.0, NOW(), NOW()),
    ('7e00dd00-0000-1111-2222-333344445558', 'USD', 273, 1.09, NOW(), NOW()),
    ('7e00dd00-0000-1111-2222-333344445558', 'GBP', 215, 0.86, NOW(), NOW()),
    
    -- Pharmacy medicine (base: EUR 18.75)
    ('7e00ee00-0000-1111-2222-333344445559', 'EUR', 1875, 1.0, NOW(), NOW()),
    ('7e00ee00-0000-1111-2222-333344445559', 'USD', 2044, 1.09, NOW(), NOW()),
    ('7e00ee00-0000-1111-2222-333344445559', 'GBP', 1613, 0.86, NOW(), NOW()),
    
    -- Educational book (base: EUR 25.90)
    ('7e00ee01-0000-1111-2222-333344445559', 'EUR', 2590, 1.0, NOW(), NOW()),
    ('7e00ee01-0000-1111-2222-333344445559', 'USD', 2823, 1.09, NOW(), NOW()),
    ('7e00ee01-0000-1111-2222-333344445559', 'GBP', 2227, 0.86, NOW(), NOW()),
    
    -- Phone bill payment (base: EUR 35.00)
    ('7e00ee02-0000-1111-2222-333344445559', 'EUR', 3500, 1.0, NOW(), NOW()),
    ('7e00ee02-0000-1111-2222-333344445559', 'USD', 3815, 1.09, NOW(), NOW()),
    ('7e00ee02-0000-1111-2222-333344445559', 'GBP', 3010, 0.86, NOW(), NOW())
ON CONFLICT (transaction_entry_id, currency) DO NOTHING;
