-- seed_categories.sql: Sample categories data for Aurora PostgreSQL Ahorro Transactions Service
-- Note: Higher rank numbers indicate more important categories within each group
-- Rank ranges: Food & Dining (10000-9000), Shopping (9000-8000), Transportation (8000-7000), 
-- Entertainment (7000-6000), Bills & Utilities (6000-5000), Health & Wellness (5000-4000),
-- Income (4000-3000), Financial (3000-2000), Education (2000-1000), Miscellaneous (1000-0)
-- Each category uses steps of 100 within its group range

INSERT INTO category (id, user_id, group_id, category_group_id, name, "group", description, rank, image_url, created_at, updated_at) VALUES
    -- Food & Dining Categories (10000-9000 range)
    ('ca001111-1111-1111-1111-111111111111', '99bb2200-0011-2233-4455-667788990011', '88aa1100-0011-2233-4455-667788990011', 'c9001234-7890-1234-5678-901234567890', 'Groceries', 'Food & Dining', 'Grocery shopping and food supplies', 9900, 'https://images.unsplash.com/photo-1542838132-92c53300491e?w=100&h=100&fit=crop', NOW(), NOW()),
    ('ca002222-2222-2222-2222-222222222222', '99bb2200-0011-2233-4455-667788990011', '88aa1100-0011-2233-4455-667788990011', 'c9001234-7890-1234-5678-901234567890', 'Restaurants', 'Food & Dining', 'Dining out at restaurants', 9800, 'https://images.unsplash.com/photo-1517248135467-4c7edcad34c4?w=100&h=100&fit=crop', NOW(), NOW()),
    ('ca003333-3333-3333-3333-333333333333', '99bb2200-0011-2233-4455-667788990011', '88aa1100-0011-2233-4455-667788990011', 'c9001234-7890-1234-5678-901234567890', 'Coffee & Tea', 'Food & Dining', 'Coffee shops and tea houses', 9700, 'https://images.unsplash.com/photo-1495474472287-4d71bcdd2085?w=100&h=100&fit=crop', NOW(), NOW()),
    ('ca004444-4444-4444-4444-444444444444', '99bb2200-0011-2233-4455-667788990011', '88aa1100-0011-2233-4455-667788990011', 'c9001234-7890-1234-5678-901234567890', 'Fast Food', 'Food & Dining', 'Quick service restaurants', 9600, 'https://images.unsplash.com/photo-1568901346375-23c9450c58cd?w=100&h=100&fit=crop', NOW(), NOW()),
    
    -- Shopping Categories (9000-8000 range)
    ('ca005555-5555-5555-5555-555555555555', '99bb2200-0011-2233-4455-667788990011', '88aa1100-0011-2233-4455-667788990011', 'c9012345-7890-1234-5678-901234567890', 'Clothing', 'Shopping', 'Clothing and fashion', 8900, 'https://images.unsplash.com/photo-1445205170230-053b83016050?w=100&h=100&fit=crop', NOW(), NOW()),
    ('ca006666-6666-6666-6666-666666666666', '99bb2200-0011-2233-4455-667788990011', '88aa1100-0011-2233-4455-667788990011', 'c9012345-7890-1234-5678-901234567890', 'Electronics', 'Shopping', 'Electronic devices and gadgets', 8800, 'https://images.unsplash.com/photo-1498049794561-7780e7231661?w=100&h=100&fit=crop', NOW(), NOW()),
    ('ca007777-7777-7777-7777-777777777777', '99bb2200-0011-2233-4455-667788990011', '88aa1100-0011-2233-4455-667788990011', 'c9012345-7890-1234-5678-901234567890', 'Home & Garden', 'Shopping', 'Home improvement and garden supplies', 8700, 'https://images.unsplash.com/photo-1416879595882-3373a0480b5b?w=100&h=100&fit=crop', NOW(), NOW()),
    ('ca008888-8888-8888-8888-888888888888', '99bb2200-0011-2233-4455-667788990011', '88aa1100-0011-2233-4455-667788990011', 'c9012345-7890-1234-5678-901234567890', 'Books & Magazines', 'Shopping', 'Books, magazines, and reading materials', 8600, 'https://images.unsplash.com/photo-1481627834876-b7833e8f5570?w=100&h=100&fit=crop', NOW(), NOW()),
    
    -- Transportation Categories (8000-7000 range)
    ('ca009999-9999-9999-9999-999999999999', '99bb2200-0011-2233-4455-667788990011', '88aa1100-0011-2233-4455-667788990011', 'c9023456-7890-1234-5678-901234567890', 'Gas & Fuel', 'Transportation', 'Vehicle fuel and gas stations', 7900, 'https://images.unsplash.com/photo-1545262810-77515befe149?w=100&h=100&fit=crop', NOW(), NOW()),
    ('ca00aaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', '99bb2200-0011-2233-4455-667788990011', '88aa1100-0011-2233-4455-667788990011', 'c9023456-7890-1234-5678-901234567890', 'Public Transit', 'Transportation', 'Bus, train, and subway fares', 7800, 'https://images.unsplash.com/photo-1544620347-c4fd4a3d5957?w=100&h=100&fit=crop', NOW(), NOW()),
    
    -- Entertainment Categories (7000-6000 range)
    ('ca00bbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', '99bb2200-0011-2233-4455-667788990011', '88aa1100-0011-2233-4455-667788990011', 'c9034567-7890-1234-5678-901234567890', 'Streaming Services', 'Entertainment', 'Netflix, Amazon Prime, etc.', 6900, 'https://images.unsplash.com/photo-1522869635100-9f4c5e86aa37?w=100&h=100&fit=crop', NOW(), NOW()),
    ('ca00cccc-cccc-cccc-cccc-cccccccccccc', '99bb2200-0011-2233-4455-667788990011', '88aa1100-0011-2233-4455-667788990011', 'c9034567-7890-1234-5678-901234567890', 'Movies & Cinema', 'Entertainment', 'Movie tickets and cinema', 6800, 'https://images.unsplash.com/photo-1489599856069-42b7dbc8b963?w=100&h=100&fit=crop', NOW(), NOW()),
    
    -- Bills & Utilities Categories (6000-5000 range)
    ('ca00dddd-dddd-dddd-dddd-dddddddddddd', '99bb2200-0011-2233-4455-667788990011', '88aa1100-0011-2233-4455-667788990011', 'c9045678-7890-1234-5678-901234567890', 'Electricity', 'Bills & Utilities', 'Electric utility bills', 5900, 'https://images.unsplash.com/photo-1558618666-fcd25c85cd64?w=100&h=100&fit=crop', NOW(), NOW()),
    ('ca00eeee-eeee-eeee-eeee-eeeeeeeeeeee', '99bb2200-0011-2233-4455-667788990011', '88aa1100-0011-2233-4455-667788990011', 'c9045678-7890-1234-5678-901234567890', 'Internet & Phone', 'Bills & Utilities', 'Internet and phone bills', 5800, 'https://images.unsplash.com/photo-1554224155-6726b3ff858f?w=100&h=100&fit=crop', NOW(), NOW()),
    
    -- Health & Wellness Categories (5000-4000 range)
    ('ca00ffff-ffff-ffff-ffff-ffffffffffff', '99bb2200-0011-2233-4455-667788990011', '88aa1100-0011-2233-4455-667788990011', 'c9056789-7890-1234-5678-901234567890', 'Medical', 'Health & Wellness', 'Doctor visits and medical expenses', 4900, 'https://images.unsplash.com/photo-1559757175-0eb30cd8c063?w=100&h=100&fit=crop', NOW(), NOW()),
    ('ca011111-1111-1111-1111-111111111111', '99bb2200-0011-2233-4455-667788990011', '88aa1100-0011-2233-4455-667788990011', 'c9056789-7890-1234-5678-901234567890', 'Fitness & Gym', 'Health & Wellness', 'Gym memberships and fitness', 4800, 'https://images.unsplash.com/photo-1559757148-5c350d0d3c56?w=100&h=100&fit=crop', NOW(), NOW()),
    
    -- Income Categories (4000-3000 range)
    ('ca012222-2222-2222-2222-222222222222', '99bb2200-0011-2233-4455-667788990011', '88aa1100-0011-2233-4455-667788990011', 'c9067890-7890-1234-5678-901234567890', 'Salary', 'Income', 'Regular salary income', 3900, 'https://images.unsplash.com/photo-1579621970563-ebec7560ff3e?w=100&h=100&fit=crop', NOW(), NOW()),
    ('ca013333-3333-3333-3333-333333333333', '99bb2200-0011-2233-4455-667788990011', '88aa1100-0011-2233-4455-667788990011', 'c9067890-7890-1234-5678-901234567890', 'Freelance', 'Income', 'Freelance work and consulting', 3800, 'https://images.unsplash.com/photo-1460925895917-afdab827c52f?w=100&h=100&fit=crop', NOW(), NOW()),
    
    -- Financial Categories (3000-2000 range)
    ('ca014444-4444-4444-4444-444444444444', '99bb2200-0011-2233-4455-667788990011', '88aa1100-0011-2233-4455-667788990011', 'c9078901-7890-1234-5678-901234567890', 'Bank Fees', 'Financial', 'Banking fees and charges', 2900, 'https://images.unsplash.com/photo-1611974789855-9c2a0a7236a3?w=100&h=100&fit=crop', NOW(), NOW()),
    ('ca015555-5555-5555-5555-555555555555', '99bb2200-0011-2233-4455-667788990011', '88aa1100-0011-2233-4455-667788990011', 'c9078901-7890-1234-5678-901234567890', 'Investments', 'Financial', 'Investment purchases and trades', 2800, 'https://images.unsplash.com/photo-1590283603385-17ffb3a7f29f?w=100&h=100&fit=crop', NOW(), NOW()),
    
    -- Education Categories (2000-1000 range)
    ('ca016666-6666-6666-6666-666666666666', '99bb2200-0011-2233-4455-667788990011', '88aa1100-0011-2233-4455-667788990011', 'c9089012-7890-1234-5678-901234567890', 'Online Courses', 'Education', 'Online learning and courses', 1900, 'https://images.unsplash.com/photo-1481627834876-b7833e8f5570?w=100&h=100&fit=crop', NOW(), NOW()),
    ('ca017777-7777-7777-7777-777777777777', '99bb2200-0011-2233-4455-667788990011', '88aa1100-0011-2233-4455-667788990011', 'c9089012-7890-1234-5678-901234567890', 'School Supplies', 'Education', 'Educational materials and supplies', 1800, 'https://images.unsplash.com/photo-1456735190827-d1262f71b8a3?w=100&h=100&fit=crop', NOW(), NOW()),
    
    -- Miscellaneous Categories (1000-0 range)
    ('ca018888-8888-8888-8888-888888888888', '99bb2200-0011-2233-4455-667788990011', '88aa1100-0011-2233-4455-667788990011', 'c9090123-7890-1234-5678-901234567890', 'Other', 'Miscellaneous', 'Uncategorized expenses', 900, 'https://images.unsplash.com/photo-1588345921523-c2dcdb7f1dcd?w=100&h=100&fit=crop', NOW(), NOW()),
    ('ca019999-9999-9999-9999-999999999999', '99bb2200-0011-2233-4455-667788990011', '88aa1100-0011-2233-4455-667788990011', 'c9090123-7890-1234-5678-901234567890', 'Gifts & Donations', 'Miscellaneous', 'Gifts and charitable donations', 800, 'https://images.unsplash.com/photo-1549301889-9c60e2a8cde6?w=100&h=100&fit=crop', NOW(), NOW()),
    
    -- Categories for secondary user (99bb3300-0011-2233-4455-667788990022) - one for each category group
    ('caa00000-0000-1111-2222-333344445555', '99bb3300-0011-2233-4455-667788990022', '88aa1100-0011-2233-4455-667788990011', 'c9001234-7890-1234-5678-901234567890', 'Household Groceries', 'Food & Dining', 'Daily grocery shopping', 9500, 'https://images.unsplash.com/photo-1542838132-92c53300491e?w=100&h=100&fit=crop', NOW(), NOW()),
    ('cabb0000-0000-1111-2222-333344445556', '99bb3300-0011-2233-4455-667788990022', '88aa1100-0011-2233-4455-667788990011', 'c9012345-7890-1234-5678-901234567890', 'Personal Items', 'Shopping', 'Personal shopping and items', 8500, 'https://images.unsplash.com/photo-1445205170230-053b83016050?w=100&h=100&fit=crop', NOW(), NOW()),
    ('cacc0000-0000-1111-2222-333344445557', '99bb3300-0011-2233-4455-667788990022', '88aa1100-0011-2233-4455-667788990011', 'c9023456-7890-1234-5678-901234567890', 'Commute', 'Transportation', 'Daily commute transportation', 7500, 'https://images.unsplash.com/photo-1544620347-c4fd4a3d5957?w=100&h=100&fit=crop', NOW(), NOW()),
    ('cadd0000-0000-1111-2222-333344445558', '99bb3300-0011-2233-4455-667788990022', '88aa1100-0011-2233-4455-667788990011', 'c9034567-7890-1234-5678-901234567890', 'Movies', 'Entertainment', 'Cinema and movie entertainment', 6500, 'https://images.unsplash.com/photo-1489599856069-42b7dbc8b963?w=100&h=100&fit=crop', NOW(), NOW()),
    ('caee0000-0000-1111-2222-333344445559', '99bb3300-0011-2233-4455-667788990022', '88aa1100-0011-2233-4455-667788990011', 'c9045678-7890-1234-5678-901234567890', 'Phone Bill', 'Bills & Utilities', 'Monthly phone bill', 5500, 'https://images.unsplash.com/photo-1554224155-6726b3ff858f?w=100&h=100&fit=crop', NOW(), NOW()),
    ('caff0000-0000-1111-2222-333344445560', '99bb3300-0011-2233-4455-667788990022', '88aa1100-0011-2233-4455-667788990011', 'c9056789-7890-1234-5678-901234567890', 'Pharmacy', 'Health & Wellness', 'Pharmacy and medical purchases', 4500, 'https://images.unsplash.com/photo-1559757175-0eb30cd8c063?w=100&h=100&fit=crop', NOW(), NOW()),
    ('ca110000-0000-1111-2222-333344445561', '99bb3300-0011-2233-4455-667788990022', '88aa1100-0011-2233-4455-667788990011', 'c9067890-7890-1234-5678-901234567890', 'Part-time Job', 'Income', 'Part-time work income', 3500, 'https://images.unsplash.com/photo-1579621970563-ebec7560ff3e?w=100&h=100&fit=crop', NOW(), NOW()),
    ('ca220000-0000-1111-2222-333344445562', '99bb3300-0011-2233-4455-667788990022', '88aa1100-0011-2233-4455-667788990011', 'c9078901-7890-1234-5678-901234567890', 'ATM Fees', 'Financial', 'Bank and ATM fees', 2500, 'https://images.unsplash.com/photo-1611974789855-9c2a0a7236a3?w=100&h=100&fit=crop', NOW(), NOW()),
    ('ca330000-0000-1111-2222-333344445563', '99bb3300-0011-2233-4455-667788990022', '88aa1100-0011-2233-4455-667788990011', 'c9089012-7890-1234-5678-901234567890', 'Books', 'Education', 'Educational books and materials', 1500, 'https://images.unsplash.com/photo-1481627834876-b7833e8f5570?w=100&h=100&fit=crop', NOW(), NOW()),
    ('ca440000-0000-1111-2222-333344445564', '99bb3300-0011-2233-4455-667788990022', '88aa1100-0011-2233-4455-667788990011', 'c9090123-7890-1234-5678-901234567890', 'Personal Care', 'Miscellaneous', 'Personal care and misc items', 500, 'https://images.unsplash.com/photo-1588345921523-c2dcdb7f1dcd?w=100&h=100&fit=crop', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
