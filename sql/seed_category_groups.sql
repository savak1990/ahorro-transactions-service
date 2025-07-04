-- seed_category_groups.sql: Sample category groups data for Aurora PostgreSQL Ahorro Transactions Service
-- Note: Higher rank numbers indicate more important category groups
-- Rank uses steps of 1000 to allow easy insertion of new category groups between existing ones

INSERT INTO category_group (id, name, rank, image_url, created_at, updated_at) VALUES
    ('c9001234-7890-1234-5678-901234567890', 'Food & Dining', 10000, 'https://images.unsplash.com/photo-1567620905732-2d1ec7ab7445?w=100&h=100&fit=crop', NOW(), NOW()),
    ('c9012345-7890-1234-5678-901234567890', 'Shopping', 9000, 'https://images.unsplash.com/photo-1472851294608-062f824d29cc?w=100&h=100&fit=crop', NOW(), NOW()),
    ('c9023456-7890-1234-5678-901234567890', 'Transportation', 8000, 'https://images.unsplash.com/photo-1544620347-c4fd4a3d5957?w=100&h=100&fit=crop', NOW(), NOW()),
    ('c9034567-7890-1234-5678-901234567890', 'Entertainment', 7000, 'https://images.unsplash.com/photo-1489599856069-42b7dbc8b963?w=100&h=100&fit=crop', NOW(), NOW()),
    ('c9045678-7890-1234-5678-901234567890', 'Bills & Utilities', 6000, 'https://images.unsplash.com/photo-1554224155-6726b3ff858f?w=100&h=100&fit=crop', NOW(), NOW()),
    ('c9056789-7890-1234-5678-901234567890', 'Health & Wellness', 5000, 'https://images.unsplash.com/photo-1559757148-5c350d0d3c56?w=100&h=100&fit=crop', NOW(), NOW()),
    ('c9067890-7890-1234-5678-901234567890', 'Income', 4000, 'https://images.unsplash.com/photo-1579621970563-ebec7560ff3e?w=100&h=100&fit=crop', NOW(), NOW()),
    ('c9078901-7890-1234-5678-901234567890', 'Financial', 3000, 'https://images.unsplash.com/photo-1611974789855-9c2a0a7236a3?w=100&h=100&fit=crop', NOW(), NOW()),
    ('c9089012-7890-1234-5678-901234567890', 'Education', 2000, 'https://images.unsplash.com/photo-1481627834876-b7833e8f5570?w=100&h=100&fit=crop', NOW(), NOW()),
    ('c9090123-7890-1234-5678-901234567890', 'Miscellaneous', 1000, 'https://images.unsplash.com/photo-1588345921523-c2dcdb7f1dcd?w=100&h=100&fit=crop', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;