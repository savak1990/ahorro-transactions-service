-- seed_merchants.sql: Sample merchants data for Aurora PostgreSQL Ahorro Transactions Service

INSERT INTO merchant (id, group_id, user_id, name, description, rank, created_at, updated_at) VALUES
    -- Spanish Supermarkets & Retail
    ('4e001234-1234-5678-9abc-def012345678', '88aa1100-0011-2233-4455-667788990011', '99bb2200-0011-2233-4455-667788990011', 'Mercadona', 'Spanish supermarket chain', 1, NOW(), NOW()),
    ('4e002345-2345-6789-abcd-ef0123456789', '88aa1100-0011-2233-4455-667788990011', '99bb2200-0011-2233-4455-667788990011', 'El Corte Inglés', 'Spanish department store chain', 2, NOW(), NOW()),
    ('4e003456-3456-789a-bcde-f01234567890', '88aa1100-0011-2233-4455-667788990011', '99bb2200-0011-2233-4455-667788990011', 'Carrefour', 'French multinational retail corporation', 3, NOW(), NOW()),
    ('4e004567-4567-89ab-cdef-012345678901', '88aa1100-0011-2233-4455-667788990011', '99bb2200-0011-2233-4455-667788990011', 'Lidl', 'German discount supermarket chain', 4, NOW(), NOW()),
    ('4e005678-5678-9abc-def0-123456789012', '88aa1100-0011-2233-4455-667788990011', '99bb2200-0011-2233-4455-667788990011', 'Zara', 'Spanish fashion retailer', 5, NOW(), NOW()),
    ('4e006789-6789-abcd-ef01-234567890123', '88aa1100-0011-2233-4455-667788990011', '99bb2200-0011-2233-4455-667788990011', 'MediaMarkt', 'German electronics retailer', 6, NOW(), NOW()),
    
    -- Restaurants & Food
    ('4e007890-789a-bcde-f012-345678901234', '88aa1100-0011-2233-4455-667788990011', '99bb2200-0011-2233-4455-667788990011', 'Starbucks', 'American coffeehouse chain', 7, NOW(), NOW()),
    ('4e008901-89ab-cdef-0123-456789012345', '88aa1100-0011-2233-4455-667788990011', '99bb2200-0011-2233-4455-667788990011', 'McDonald''s', 'American fast food restaurant chain', 8, NOW(), NOW()),
    ('4e009012-9abc-def0-1234-567890123456', '88aa1100-0011-2233-4455-667788990011', '99bb2200-0011-2233-4455-667788990011', 'Telepizza', 'Spanish pizza delivery chain', 9, NOW(), NOW()),
    ('4e00a123-abcd-ef01-2345-678901234567', '88aa1100-0011-2233-4455-667788990011', '99bb2200-0011-2233-4455-667788990011', 'Domino''s Pizza', 'American pizza restaurant chain', 10, NOW(), NOW()),
    ('4e00b0b0-bcde-f012-3456-789012345678', '88aa1100-0011-2233-4455-667788990011', '99bb2200-0011-2233-4455-667788990011', 'Burger King', 'American fast food restaurant chain', 11, NOW(), NOW()),
    
    -- Transportation
    ('4e00c0c0-cdef-0123-4567-890123456789', '88aa1100-0011-2233-4455-667788990011', '99bb2200-0011-2233-4455-667788990011', 'Uber', 'American mobility company', 12, NOW(), NOW()),
    ('4e00d0d0-def0-1234-5678-901234567890', '88aa1100-0011-2233-4455-667788990011', '99bb2200-0011-2233-4455-667788990011', 'Cabify', 'Spanish ride-hailing company', 13, NOW(), NOW()),
    ('4e00e0e0-ef01-2345-6789-012345678901', '88aa1100-0011-2233-4455-667788990011', '99bb2200-0011-2233-4455-667788990011', 'Renfe', 'Spanish railway company', 14, NOW(), NOW()),
    ('4e00f0f0-f012-3456-789a-123456789012', '88aa1100-0011-2233-4455-667788990011', '99bb2200-0011-2233-4455-667788990011', 'Repsol', 'Spanish oil and gas company', 15, NOW(), NOW()),
    ('4e001010-0123-4567-89ab-234567890123', '88aa1100-0011-2233-4455-667788990011', '99bb2200-0011-2233-4455-667788990011', 'Metro Madrid', 'Madrid metro system', 16, NOW(), NOW()),
    
    -- Entertainment & Streaming
    ('4e001111-1234-5678-9abc-345678901234', '88aa1100-0011-2233-4455-667788990011', '99bb2200-0011-2233-4455-667788990011', 'Netflix', 'American streaming service', 17, NOW(), NOW()),
    ('4e001212-2345-6789-abcd-456789012345', '88aa1100-0011-2233-4455-667788990011', '99bb2200-0011-2233-4455-667788990011', 'Amazon Prime', 'Amazon streaming service', 18, NOW(), NOW()),
    ('4e001313-3456-789a-bcde-567890123456', '88aa1100-0011-2233-4455-667788990011', '99bb2200-0011-2233-4455-667788990011', 'Spotify', 'Swedish music streaming service', 19, NOW(), NOW()),
    ('4e001414-4567-89ab-cdef-678901234567', '88aa1100-0011-2233-4455-667788990011', '99bb2200-0011-2233-4455-667788990011', 'Disney+', 'American streaming service', 20, NOW(), NOW()),
    ('4e001515-5678-9abc-def0-789012345678', '88aa1100-0011-2233-4455-667788990011', '99bb2200-0011-2233-4455-667788990011', 'HBO Max', 'American streaming service', 21, NOW(), NOW()),
    
    -- Technology & Online Services
    ('4e001616-6789-abcd-ef01-890123456789', '88aa1100-0011-2233-4455-667788990011', '99bb2200-0011-2233-4455-667788990011', 'Amazon', 'American e-commerce company', 22, NOW(), NOW()),
    ('4e001717-789a-bcde-f012-901234567890', '88aa1100-0011-2233-4455-667788990011', '99bb2200-0011-2233-4455-667788990011', 'Apple', 'American technology company', 23, NOW(), NOW()),
    ('4e001818-89ab-cdef-0123-012345678901', '88aa1100-0011-2233-4455-667788990011', '99bb2200-0011-2233-4455-667788990011', 'Google', 'American technology company', 24, NOW(), NOW()),
    ('4e001919-9abc-def0-1234-123456789012', '88aa1100-0011-2233-4455-667788990011', '99bb2200-0011-2233-4455-667788990011', 'Microsoft', 'American technology company', 25, NOW(), NOW()),
    
    -- Utilities & Services
    ('4e001a1a-abcd-ef01-2345-234567890123', '88aa1100-0011-2233-4455-667788990011', '99bb2200-0011-2233-4455-667788990011', 'Iberdrola', 'Spanish electric utility company', 26, NOW(), NOW()),
    ('4e001b1b-bcde-f012-3456-345678901234', '88aa1100-0011-2233-4455-667788990011', '99bb2200-0011-2233-4455-667788990011', 'Telefónica', 'Spanish telecommunications company', 27, NOW(), NOW()),
    ('4e001c1c-cdef-0123-4567-456789012345', '88aa1100-0011-2233-4455-667788990011', '99bb2200-0011-2233-4455-667788990011', 'Orange', 'French telecommunications company', 28, NOW(), NOW()),
    ('4e001d1d-def0-1234-5678-567890123456', '88aa1100-0011-2233-4455-667788990011', '99bb2200-0011-2233-4455-667788990011', 'Vodafone', 'British telecommunications company', 29, NOW(), NOW()),
    
    -- Banks & Financial
    ('4e001e1e-ef01-2345-6789-678901234567', '88aa1100-0011-2233-4455-667788990011', '99bb2200-0011-2233-4455-667788990011', 'BBVA', 'Spanish multinational bank', 30, NOW(), NOW()),
    ('4e001f1f-f012-3456-789a-789012345678', '88aa1100-0011-2233-4455-667788990011', '99bb2200-0011-2233-4455-667788990011', 'Santander', 'Spanish multinational bank', 31, NOW(), NOW()),
    ('4e002020-0123-4567-89ab-890123456789', '88aa1100-0011-2233-4455-667788990011', '99bb2200-0011-2233-4455-667788990011', 'CaixaBank', 'Spanish bank', 32, NOW(), NOW()),
    ('4e002121-1234-5678-9abc-901234567890', '88aa1100-0011-2233-4455-667788990011', '99bb2200-0011-2233-4455-667788990011', 'ING', 'Dutch multinational bank', 33, NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
