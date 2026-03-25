-- Sample Goals Seed Data
-- Creates financial goals with milestones for testing

-- John Doe's goals (User 1) - Multiple goals at different stages
INSERT INTO goals (id, user_id, title, description, target_amount, current_amount, currency, target_date, status, created_at, updated_at)
VALUES
    ('g1111111-1111-1111-1111-111111111111', '11111111-1111-1111-1111-111111111111', 'Emergency Fund', 'Build a 6-month emergency fund for financial security', 10000.00, 4500.00, 'USD', CURRENT_DATE + INTERVAL '180 days', 'active', NOW() - INTERVAL '60 days', NOW()),
    ('g1111112-1111-1111-1111-111111111111', '11111111-1111-1111-1111-111111111111', 'Vacation to Europe', 'Save for a 2-week trip to Europe next summer', 5000.00, 2800.00, 'USD', CURRENT_DATE + INTERVAL '240 days', 'active', NOW() - INTERVAL '45 days', NOW()),
    ('g1111113-1111-1111-1111-111111111111', '11111111-1111-1111-1111-111111111111', 'New Laptop', 'Save for a MacBook Pro for work', 2500.00, 2500.00, 'USD', CURRENT_DATE - INTERVAL '5 days', 'completed', NOW() - INTERVAL '90 days', NOW() - INTERVAL '5 days');

-- John's goal milestones
INSERT INTO goal_milestones (id, goal_id, title, amount, is_completed, completed_at, "order")
VALUES
    -- Emergency Fund milestones
    (gen_random_uuid(), 'g1111111-1111-1111-1111-111111111111', 'First $1,000', 1000.00, true, NOW() - INTERVAL '55 days', 1),
    (gen_random_uuid(), 'g1111111-1111-1111-1111-111111111111', 'Quarter Way ($2,500)', 2500.00, true, NOW() - INTERVAL '40 days', 2),
    (gen_random_uuid(), 'g1111111-1111-1111-1111-111111111111', 'Halfway Point ($5,000)', 5000.00, false, NULL, 3),
    (gen_random_uuid(), 'g1111111-1111-1111-1111-111111111111', 'Three Quarters ($7,500)', 7500.00, false, NULL, 4),
    (gen_random_uuid(), 'g1111111-1111-1111-1111-111111111111', 'Goal Complete!', 10000.00, false, NULL, 5),
    
    -- Vacation milestones
    (gen_random_uuid(), 'g1111112-1111-1111-1111-111111111111', 'Flight Money ($1,500)', 1500.00, true, NOW() - INTERVAL '30 days', 1),
    (gen_random_uuid(), 'g1111112-1111-1111-1111-111111111111', 'Accommodation ($3,000)', 3000.00, false, NULL, 2),
    (gen_random_uuid(), 'g1111112-1111-1111-1111-111111111111', 'Spending Money ($5,000)', 5000.00, false, NULL, 3),
    
    -- Laptop milestones (all completed)
    (gen_random_uuid(), 'g1111113-1111-1111-1111-111111111111', 'First $500', 500.00, true, NOW() - INTERVAL '80 days', 1),
    (gen_random_uuid(), 'g1111113-1111-1111-1111-111111111111', 'Halfway ($1,250)', 1250.00, true, NOW() - INTERVAL '50 days', 2),
    (gen_random_uuid(), 'g1111113-1111-1111-1111-111111111111', 'Almost There ($2,000)', 2000.00, true, NOW() - INTERVAL '20 days', 3),
    (gen_random_uuid(), 'g1111113-1111-1111-1111-111111111111', 'Goal Achieved!', 2500.00, true, NOW() - INTERVAL '5 days', 4);

-- Jane Smith's goals (User 2) - Early stage goals
INSERT INTO goals (id, user_id, title, description, target_amount, current_amount, currency, target_date, status, created_at, updated_at)
VALUES
    ('g2222221-2222-2222-2222-222222222222', '22222222-2222-2222-2222-222222222222', 'Down Payment for Car', 'Save $8,000 for a reliable used car', 8000.00, 1200.00, 'USD', CURRENT_DATE + INTERVAL '365 days', 'active', NOW() - INTERVAL '30 days', NOW()),
    ('g2222222-2222-2222-2222-222222222222', '22222222-2222-2222-2222-222222222222', 'Pay Off Credit Card', 'Eliminate credit card debt', 3500.00, 800.00, 'USD', CURRENT_DATE + INTERVAL '120 days', 'active', NOW() - INTERVAL '20 days', NOW());

-- Jane's goal milestones
INSERT INTO goal_milestones (id, goal_id, title, amount, is_completed, completed_at, "order")
VALUES
    -- Car down payment milestones
    (gen_random_uuid(), 'g2222221-2222-2222-2222-222222222222', 'First $1,000', 1000.00, true, NOW() - INTERVAL '15 days', 1),
    (gen_random_uuid(), 'g2222221-2222-2222-2222-222222222222', 'Quarter Way ($2,000)', 2000.00, false, NULL, 2),
    (gen_random_uuid(), 'g2222221-2222-2222-2222-222222222222', 'Halfway ($4,000)', 4000.00, false, NULL, 3),
    (gen_random_uuid(), 'g2222221-2222-2222-2222-222222222222', 'Three Quarters ($6,000)', 6000.00, false, NULL, 4),
    (gen_random_uuid(), 'g2222221-2222-2222-2222-222222222222', 'Goal Complete!', 8000.00, false, NULL, 5),
    
    -- Credit card payoff milestones
    (gen_random_uuid(), 'g2222222-2222-2222-2222-222222222222', 'First $500', 500.00, true, NOW() - INTERVAL '12 days', 1),
    (gen_random_uuid(), 'g2222222-2222-2222-2222-222222222222', 'Quarter Paid ($875)', 875.00, false, NULL, 2),
    (gen_random_uuid(), 'g2222222-2222-2222-2222-222222222222', 'Halfway ($1,750)', 1750.00, false, NULL, 3),
    (gen_random_uuid(), 'g2222222-2222-2222-2222-222222222222', 'Almost Done ($2,625)', 2625.00, false, NULL, 4),
    (gen_random_uuid(), 'g2222222-2222-2222-2222-222222222222', 'Debt Free!', 3500.00, false, NULL, 5);

-- Mike Johnson's goals (User 3) - Just starting out
INSERT INTO goals (id, user_id, title, description, target_amount, current_amount, currency, target_date, status, created_at, updated_at)
VALUES
    ('g3333331-3333-3333-3333-333333333333', '33333333-3333-3333-3333-333333333333', 'First $1,000 Saved', 'Build initial savings cushion', 1000.00, 85.00, 'USD', CURRENT_DATE + INTERVAL '90 days', 'active', NOW() - INTERVAL '10 days', NOW());

-- Mike's goal milestones
INSERT INTO goal_milestones (id, goal_id, title, amount, is_completed, completed_at, "order")
VALUES
    (gen_random_uuid(), 'g3333331-3333-3333-3333-333333333333', 'First $100', 100.00, false, NULL, 1),
    (gen_random_uuid(), 'g3333331-3333-3333-3333-333333333333', 'Quarter Way ($250)', 250.00, false, NULL, 2),
    (gen_random_uuid(), 'g3333331-3333-3333-3333-333333333333', 'Halfway ($500)', 500.00, false, NULL, 3),
    (gen_random_uuid(), 'g3333331-3333-3333-3333-333333333333', 'Three Quarters ($750)', 750.00, false, NULL, 4),
    (gen_random_uuid(), 'g3333331-3333-3333-3333-333333333333', 'Goal Achieved!', 1000.00, false, NULL, 5);

-- Sarah Williams's goals (User 4) - Ambitious saver
INSERT INTO goals (id, user_id, title, description, target_amount, current_amount, currency, target_date, status, created_at, updated_at)
VALUES
    ('g4444441-4444-4444-4444-444444444444', '44444444-4444-4444-4444-444444444444', 'House Down Payment', 'Save $30,000 for first home down payment', 30000.00, 8500.00, 'USD', CURRENT_DATE + INTERVAL '540 days', 'active', NOW() - INTERVAL '75 days', NOW()),
    ('g4444442-4444-4444-4444-444444444444', '44444444-4444-4444-4444-444444444444', 'Investment Portfolio', 'Start investing with $5,000', 5000.00, 3200.00, 'USD', CURRENT_DATE + INTERVAL '150 days', 'active', NOW() - INTERVAL '50 days', NOW()),
    ('g4444443-4444-4444-4444-444444444444', '44444444-4444-4444-4444-444444444444', 'Wedding Fund', 'Save for dream wedding', 15000.00, 6800.00, 'USD', CURRENT_DATE + INTERVAL '450 days', 'active', NOW() - INTERVAL '60 days', NOW());

-- Sarah's goal milestones
INSERT INTO goal_milestones (id, goal_id, title, amount, is_completed, completed_at, "order")
VALUES
    -- House down payment milestones
    (gen_random_uuid(), 'g4444441-4444-4444-4444-444444444444', 'First $5,000', 5000.00, true, NOW() - INTERVAL '60 days', 1),
    (gen_random_uuid(), 'g4444441-4444-4444-4444-444444444444', 'Quarter Way ($7,500)', 7500.00, true, NOW() - INTERVAL '30 days', 2),
    (gen_random_uuid(), 'g4444441-4444-4444-4444-444444444444', 'Halfway ($15,000)', 15000.00, false, NULL, 3),
    (gen_random_uuid(), 'g4444441-4444-4444-4444-444444444444', 'Three Quarters ($22,500)', 22500.00, false, NULL, 4),
    (gen_random_uuid(), 'g4444441-4444-4444-4444-444444444444', 'Goal Complete!', 30000.00, false, NULL, 5),
    
    -- Investment portfolio milestones
    (gen_random_uuid(), 'g4444442-4444-4444-4444-444444444444', 'First $1,000', 1000.00, true, NOW() - INTERVAL '45 days', 1),
    (gen_random_uuid(), 'g4444442-4444-4444-4444-444444444444', 'Halfway ($2,500)', 2500.00, true, NOW() - INTERVAL '25 days', 2),
    (gen_random_uuid(), 'g4444442-4444-4444-4444-444444444444', 'Almost There ($4,000)', 4000.00, false, NULL, 3),
    (gen_random_uuid(), 'g4444442-4444-4444-4444-444444444444', 'Ready to Invest!', 5000.00, false, NULL, 4),
    
    -- Wedding fund milestones
    (gen_random_uuid(), 'g4444443-4444-4444-4444-444444444444', 'First $3,000', 3000.00, true, NOW() - INTERVAL '50 days', 1),
    (gen_random_uuid(), 'g4444443-4444-4444-4444-444444444444', 'Quarter Way ($3,750)', 3750.00, true, NOW() - INTERVAL '35 days', 2),
    (gen_random_uuid(), 'g4444443-4444-4444-4444-444444444444', 'Halfway ($7,500)', 7500.00, false, NULL, 3),
    (gen_random_uuid(), 'g4444443-4444-4444-4444-444444444444', 'Three Quarters ($11,250)', 11250.00, false, NULL, 4),
    (gen_random_uuid(), 'g4444443-4444-4444-4444-444444444444', 'Dream Wedding!', 15000.00, false, NULL, 5);

-- Alex Brown's goals (User 5) - New user with simple goal
INSERT INTO goals (id, user_id, title, description, target_amount, current_amount, currency, target_date, status, created_at, updated_at)
VALUES
    ('g5555551-5555-5555-5555-555555555555', '55555555-5555-5555-5555-555555555555', 'New Phone', 'Save for iPhone 15 Pro', 1200.00, 245.00, 'USD', CURRENT_DATE + INTERVAL '60 days', 'active', NOW() - INTERVAL '14 days', NOW());

-- Alex's goal milestones
INSERT INTO goal_milestones (id, goal_id, title, amount, is_completed, completed_at, "order")
VALUES
    (gen_random_uuid(), 'g5555551-5555-5555-5555-555555555555', 'First $200', 200.00, true, NOW() - INTERVAL '7 days', 1),
    (gen_random_uuid(), 'g5555551-5555-5555-5555-555555555555', 'Quarter Way ($300)', 300.00, false, NULL, 2),
    (gen_random_uuid(), 'g5555551-5555-5555-5555-555555555555', 'Halfway ($600)', 600.00, false, NULL, 3),
    (gen_random_uuid(), 'g5555551-5555-5555-5555-555555555555', 'Almost There ($900)', 900.00, false, NULL, 4),
    (gen_random_uuid(), 'g5555551-5555-5555-5555-555555555555', 'New Phone!', 1200.00, false, NULL, 5);
