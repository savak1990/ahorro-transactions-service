-- seed_category_groups.sql: Sample category groups data for Aurora PostgreSQL Ahorro Transactions Service

INSERT INTO category_group (id, name, rank, created_at, updated_at) VALUES
    ('fd123456-7890-1234-5678-901234567890', 'Food & Dining', 1, NOW(), NOW()),
    ('sh123456-7890-1234-5678-901234567890', 'Shopping', 2, NOW(), NOW()),
    ('tr123456-7890-1234-5678-901234567890', 'Transportation', 3, NOW(), NOW()),
    ('en123456-7890-1234-5678-901234567890', 'Entertainment', 4, NOW(), NOW()),
    ('bi123456-7890-1234-5678-901234567890', 'Bills & Utilities', 5, NOW(), NOW()),
    ('he123456-7890-1234-5678-901234567890', 'Health & Wellness', 6, NOW(), NOW()),
    ('in123456-7890-1234-5678-901234567890', 'Income', 7, NOW(), NOW()),
    ('fi123456-7890-1234-5678-901234567890', 'Financial', 8, NOW(), NOW()),
    ('ed123456-7890-1234-5678-901234567890', 'Education', 9, NOW(), NOW()),
    ('mi123456-7890-1234-5678-901234567890', 'Miscellaneous', 10, NOW(), NOW())
ON CONFLICT (id) DO NOTHING;