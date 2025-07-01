-- seed_category_groups.sql: Sample category groups data for Aurora PostgreSQL Ahorro Transactions Service

INSERT INTO category_group (id, name, rank, created_at, updated_at) VALUES
    ('c9001234-7890-1234-5678-901234567890', 'Food & Dining', 1, NOW(), NOW()),
    ('c9012345-7890-1234-5678-901234567890', 'Shopping', 2, NOW(), NOW()),
    ('c9023456-7890-1234-5678-901234567890', 'Transportation', 3, NOW(), NOW()),
    ('c9034567-7890-1234-5678-901234567890', 'Entertainment', 4, NOW(), NOW()),
    ('c9045678-7890-1234-5678-901234567890', 'Bills & Utilities', 5, NOW(), NOW()),
    ('c9056789-7890-1234-5678-901234567890', 'Health & Wellness', 6, NOW(), NOW()),
    ('c9067890-7890-1234-5678-901234567890', 'Income', 7, NOW(), NOW()),
    ('c9078901-7890-1234-5678-901234567890', 'Financial', 8, NOW(), NOW()),
    ('c9089012-7890-1234-5678-901234567890', 'Education', 9, NOW(), NOW()),
    ('c9090123-7890-1234-5678-901234567890', 'Miscellaneous', 10, NOW(), NOW())
ON CONFLICT (id) DO NOTHING;