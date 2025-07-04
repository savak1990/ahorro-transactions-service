-- seed_merchants.sql: Sample merchants data for Aurora PostgreSQL Ahorro Transactions Service
-- Note: Higher rank numbers indicate more popular merchants for the user
-- Rank uses steps of 100 to allow easy insertion of new merchants between existing ones

INSERT INTO merchant (id, group_id, user_id, name, description, rank, image_url, created_at, updated_at) VALUES
    -- Spanish Supermarkets & Retail (Most Popular)
    ('4e001234-1234-5678-9abc-def012345678', '88aa1100-0011-2233-4455-667788990011', '99bb2200-0011-2233-4455-667788990011', 'Mercadona', 'Spanish supermarket chain', 3300, 'https://images.unsplash.com/photo-1556909114-f6e7ad7d3136?w=100&h=100&fit=crop', NOW(), NOW()),
    ('4e003456-3456-789a-bcde-f01234567890', '88aa1100-0011-2233-4455-667788990011', '99bb2200-0011-2233-4455-667788990011', 'Carrefour', 'French multinational retail corporation', 3200, 'https://images.unsplash.com/photo-1604719312566-8912e9227c6a?w=100&h=100&fit=crop', NOW(), NOW()),
    ('4e004567-4567-89ab-cdef-012345678901', '88aa1100-0011-2233-4455-667788990011', '99bb2200-0011-2233-4455-667788990011', 'Lidl', 'German discount supermarket chain', 3100, 'https://images.unsplash.com/photo-1578662996442-48f60103fc96?w=100&h=100&fit=crop', NOW(), NOW()),
    ('4e002345-2345-6789-abcd-ef0123456789', '88aa1100-0011-2233-4455-667788990011', '99bb2200-0011-2233-4455-667788990011', 'El Corte Inglés', 'Spanish department store chain', 3000, 'https://images.unsplash.com/photo-1441986300917-64674bd600d8?w=100&h=100&fit=crop', NOW(), NOW()),
    ('4e005678-5678-9abc-def0-123456789012', '88aa1100-0011-2233-4455-667788990011', '99bb2200-0011-2233-4455-667788990011', 'Zara', 'Spanish fashion retailer', 2900, 'https://images.unsplash.com/photo-1445205170230-053b83016050?w=100&h=100&fit=crop', NOW(), NOW()),
    ('4e006789-6789-abcd-ef01-234567890123', '88aa1100-0011-2233-4455-667788990011', '99bb2200-0011-2233-4455-667788990011', 'MediaMarkt', 'German electronics retailer', 2800, 'https://images.unsplash.com/photo-1498049794561-7780e7231661?w=100&h=100&fit=crop', NOW(), NOW()),
    
    -- Restaurants & Food (Very Popular)
    ('4e007890-789a-bcde-f012-345678901234', '88aa1100-0011-2233-4455-667788990011', '99bb2200-0011-2233-4455-667788990011', 'Starbucks', 'American coffeehouse chain', 2700, 'https://images.unsplash.com/photo-1495474472287-4d71bcdd2085?w=100&h=100&fit=crop', NOW(), NOW()),
    ('4e008901-89ab-cdef-0123-456789012345', '88aa1100-0011-2233-4455-667788990011', '99bb2200-0011-2233-4455-667788990011', 'McDonald''s', 'American fast food restaurant chain', 2600, 'https://images.unsplash.com/photo-1568901346375-23c9450c58cd?w=100&h=100&fit=crop', NOW(), NOW()),
    ('4e009012-9abc-def0-1234-567890123456', '88aa1100-0011-2233-4455-667788990011', '99bb2200-0011-2233-4455-667788990011', 'Telepizza', 'Spanish pizza delivery chain', 2500, 'https://images.unsplash.com/photo-1571407970349-bc81e7e96d47?w=100&h=100&fit=crop', NOW(), NOW()),
    ('4e00a123-abcd-ef01-2345-678901234567', '88aa1100-0011-2233-4455-667788990011', '99bb2200-0011-2233-4455-667788990011', 'Domino''s Pizza', 'American pizza restaurant chain', 2400, 'https://images.unsplash.com/photo-1513104890138-7c749659a591?w=100&h=100&fit=crop', NOW(), NOW()),
    ('4e00b0b0-bcde-f012-3456-789012345678', '88aa1100-0011-2233-4455-667788990011', '99bb2200-0011-2233-4455-667788990011', 'Burger King', 'American fast food restaurant chain', 2300, 'https://images.unsplash.com/photo-1571091718767-18b5b1457add?w=100&h=100&fit=crop', NOW(), NOW()),
    
    -- Transportation (Popular)
    ('4e00c0c0-cdef-0123-4567-890123456789', '88aa1100-0011-2233-4455-667788990011', '99bb2200-0011-2233-4455-667788990011', 'Uber', 'American mobility company', 2200, 'https://images.unsplash.com/photo-1544620347-c4fd4a3d5957?w=100&h=100&fit=crop', NOW(), NOW()),
    ('4e00f0f0-f012-3456-789a-123456789012', '88aa1100-0011-2233-4455-667788990011', '99bb2200-0011-2233-4455-667788990011', 'Repsol', 'Spanish oil and gas company', 2100, 'https://images.unsplash.com/photo-1545262810-77515befe149?w=100&h=100&fit=crop', NOW(), NOW()),
    ('4e00d0d0-def0-1234-5678-901234567890', '88aa1100-0011-2233-4455-667788990011', '99bb2200-0011-2233-4455-667788990011', 'Cabify', 'Spanish ride-hailing company', 2000, 'https://images.unsplash.com/photo-1551632811-561732d1e306?w=100&h=100&fit=crop', NOW(), NOW()),
    ('4e001010-0123-4567-89ab-234567890123', '88aa1100-0011-2233-4455-667788990011', '99bb2200-0011-2233-4455-667788990011', 'Metro Madrid', 'Madrid metro system', 1900, 'https://images.unsplash.com/photo-1544620347-c4fd4a3d5957?w=100&h=100&fit=crop', NOW(), NOW()),
    ('4e00e0e0-ef01-2345-6789-012345678901', '88aa1100-0011-2233-4455-667788990011', '99bb2200-0011-2233-4455-667788990011', 'Renfe', 'Spanish railway company', 1800, 'https://images.unsplash.com/photo-1474487548417-781cb71495f3?w=100&h=100&fit=crop', NOW(), NOW()),
    
    -- Entertainment & Streaming (Moderately Popular)
    ('4e001111-1234-5678-9abc-345678901234', '88aa1100-0011-2233-4455-667788990011', '99bb2200-0011-2233-4455-667788990011', 'Netflix', 'American streaming service', 1700, 'https://images.unsplash.com/photo-1522869635100-9f4c5e86aa37?w=100&h=100&fit=crop', NOW(), NOW()),
    ('4e001313-3456-789a-bcde-567890123456', '88aa1100-0011-2233-4455-667788990011', '99bb2200-0011-2233-4455-667788990011', 'Spotify', 'Swedish music streaming service', 1600, 'https://images.unsplash.com/photo-1470225620780-dba8ba36b745?w=100&h=100&fit=crop', NOW(), NOW()),
    ('4e001212-2345-6789-abcd-456789012345', '88aa1100-0011-2233-4455-667788990011', '99bb2200-0011-2233-4455-667788990011', 'Amazon Prime', 'Amazon streaming service', 1500, 'https://images.unsplash.com/photo-1489599856069-42b7dbc8b963?w=100&h=100&fit=crop', NOW(), NOW()),
    ('4e001414-4567-89ab-cdef-678901234567', '88aa1100-0011-2233-4455-667788990011', '99bb2200-0011-2233-4455-667788990011', 'Disney+', 'American streaming service', 1400, 'https://images.unsplash.com/photo-1578662996442-48f60103fc96?w=100&h=100&fit=crop', NOW(), NOW()),
    ('4e001515-5678-9abc-def0-789012345678', '88aa1100-0011-2233-4455-667788990011', '99bb2200-0011-2233-4455-667788990011', 'HBO Max', 'American streaming service', 1300, 'https://images.unsplash.com/photo-1489599856069-42b7dbc8b963?w=100&h=100&fit=crop', NOW(), NOW()),
    
    -- Technology & Online Services (Popular)
    ('4e001616-6789-abcd-ef01-890123456789', '88aa1100-0011-2233-4455-667788990011', '99bb2200-0011-2233-4455-667788990011', 'Amazon', 'American e-commerce company', 1200, 'https://images.unsplash.com/photo-1523474253046-8cd2748b5fd2?w=100&h=100&fit=crop', NOW(), NOW()),
    ('4e001717-789a-bcde-f012-901234567890', '88aa1100-0011-2233-4455-667788990011', '99bb2200-0011-2233-4455-667788990011', 'Apple', 'American technology company', 1100, 'https://images.unsplash.com/photo-1472851294608-062f824d29cc?w=100&h=100&fit=crop', NOW(), NOW()),
    ('4e001818-89ab-cdef-0123-012345678901', '88aa1100-0011-2233-4455-667788990011', '99bb2200-0011-2233-4455-667788990011', 'Google', 'American technology company', 1000, 'https://images.unsplash.com/photo-1573804633927-bfcbcd909acd?w=100&h=100&fit=crop', NOW(), NOW()),
    ('4e001919-9abc-def0-1234-123456789012', '88aa1100-0011-2233-4455-667788990011', '99bb2200-0011-2233-4455-667788990011', 'Microsoft', 'American technology company', 900, 'https://images.unsplash.com/photo-1498049794561-7780e7231661?w=100&h=100&fit=crop', NOW(), NOW()),
    
    -- Utilities & Services (Regular Usage)
    ('4e001a1a-abcd-ef01-2345-234567890123', '88aa1100-0011-2233-4455-667788990011', '99bb2200-0011-2233-4455-667788990011', 'Iberdrola', 'Spanish electric utility company', 800, 'https://images.unsplash.com/photo-1558618666-fcd25c85cd64?w=100&h=100&fit=crop', NOW(), NOW()),
    ('4e001b1b-bcde-f012-3456-345678901234', '88aa1100-0011-2233-4455-667788990011', '99bb2200-0011-2233-4455-667788990011', 'Telefónica', 'Spanish telecommunications company', 700, 'https://images.unsplash.com/photo-1554224155-6726b3ff858f?w=100&h=100&fit=crop', NOW(), NOW()),
    ('4e001c1c-cdef-0123-4567-456789012345', '88aa1100-0011-2233-4455-667788990011', '99bb2200-0011-2233-4455-667788990011', 'Orange', 'French telecommunications company', 600, 'https://images.unsplash.com/photo-1512941937669-90a1b58e7e9c?w=100&h=100&fit=crop', NOW(), NOW()),
    ('4e001d1d-def0-1234-5678-567890123456', '88aa1100-0011-2233-4455-667788990011', '99bb2200-0011-2233-4455-667788990011', 'Vodafone', 'British telecommunications company', 500, 'https://images.unsplash.com/photo-1560472354-b33ff0c44a43?w=100&h=100&fit=crop', NOW(), NOW()),
    
    -- Banks & Financial (Less Frequent)
    ('4e001e1e-ef01-2345-6789-678901234567', '88aa1100-0011-2233-4455-667788990011', '99bb2200-0011-2233-4455-667788990011', 'BBVA', 'Spanish multinational bank', 400, 'https://images.unsplash.com/photo-1611974789855-9c2a0a7236a3?w=100&h=100&fit=crop', NOW(), NOW()),
    ('4e001f1f-f012-3456-789a-789012345678', '88aa1100-0011-2233-4455-667788990011', '99bb2200-0011-2233-4455-667788990011', 'Santander', 'Spanish multinational bank', 300, 'https://images.unsplash.com/photo-1541354329998-f4d9a9f9297f?w=100&h=100&fit=crop', NOW(), NOW()),
    ('4e002020-0123-4567-89ab-890123456789', '88aa1100-0011-2233-4455-667788990011', '99bb2200-0011-2233-4455-667788990011', 'CaixaBank', 'Spanish bank', 200, 'https://images.unsplash.com/photo-1554224154-22dec7ec8818?w=100&h=100&fit=crop', NOW(), NOW()),
    ('4e002121-1234-5678-9abc-901234567890', '88aa1100-0011-2233-4455-667788990011', '99bb2200-0011-2233-4455-667788990011', 'ING', 'Dutch multinational bank', 100, 'https://images.unsplash.com/photo-1601597111158-2fceff292cdc?w=100&h=100&fit=crop', NOW(), NOW()),
    
    -- Merchants for secondary user (99bb3300-0011-2233-4455-667788990022)
    ('4e00aa00-0000-1111-2222-333344445555', '88aa1100-0011-2233-4455-667788990011', '99bb3300-0011-2233-4455-667788990022', 'Dia Supermercado', 'Spanish discount supermarket chain', 3100, 'https://images.unsplash.com/photo-1604719312566-8912e9227c6a?w=100&h=100&fit=crop', NOW(), NOW()),
    ('4e00bb00-0000-1111-2222-333344445566', '88aa1100-0011-2233-4455-667788990011', '99bb3300-0011-2233-4455-667788990022', 'Local Café', 'Neighborhood coffee shop', 2500, 'https://images.unsplash.com/photo-1559925393-8be0ec4767c8?w=100&h=100&fit=crop', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
