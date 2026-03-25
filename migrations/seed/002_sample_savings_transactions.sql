-- Sample Savings Transactions Seed Data
-- Creates transactions with various dates for streak testing

-- John Doe's savings (User 1) - Consistent saver with current 7-day streak
INSERT INTO savings_transactions (id, user_id, amount, currency, description, category, created_at)
VALUES
    -- Current streak (last 7 days)
    (gen_random_uuid(), '11111111-1111-1111-1111-111111111111', 25.00, 'USD', 'Daily coffee savings', 'Daily Savings', NOW() - INTERVAL '0 days'),
    (gen_random_uuid(), '11111111-1111-1111-1111-111111111111', 30.00, 'USD', 'Lunch money saved', 'Daily Savings', NOW() - INTERVAL '1 days'),
    (gen_random_uuid(), '11111111-1111-1111-1111-111111111111', 20.00, 'USD', 'Skipped takeout', 'Food', NOW() - INTERVAL '2 days'),
    (gen_random_uuid(), '11111111-1111-1111-1111-111111111111', 15.00, 'USD', 'Morning coffee saved', 'Daily Savings', NOW() - INTERVAL '3 days'),
    (gen_random_uuid(), '11111111-1111-1111-1111-111111111111', 40.00, 'USD', 'Entertainment savings', 'Entertainment', NOW() - INTERVAL '4 days'),
    (gen_random_uuid(), '11111111-1111-1111-1111-111111111111', 35.00, 'USD', 'Transportation savings', 'Transport', NOW() - INTERVAL '5 days'),
    (gen_random_uuid(), '11111111-1111-1111-1111-111111111111', 50.00, 'USD', 'Weekly bonus save', 'Bonus', NOW() - INTERVAL '6 days'),
    
    -- Historical savings (older transactions)
    (gen_random_uuid(), '11111111-1111-1111-1111-111111111111', 100.00, 'USD', 'Paycheck savings', 'Income', NOW() - INTERVAL '15 days'),
    (gen_random_uuid(), '11111111-1111-1111-1111-111111111111', 75.00, 'USD', 'Side gig earnings', 'Income', NOW() - INTERVAL '20 days'),
    (gen_random_uuid(), '11111111-1111-1111-1111-111111111111', 45.00, 'USD', 'Refund saved', 'Other', NOW() - INTERVAL '25 days'),
    (gen_random_uuid(), '11111111-1111-1111-1111-111111111111', 60.00, 'USD', 'Monthly challenge', 'Challenge', NOW() - INTERVAL '30 days'),
    (gen_random_uuid(), '11111111-1111-1111-1111-111111111111', 80.00, 'USD', 'Bonus from work', 'Income', NOW() - INTERVAL '45 days'),
    (gen_random_uuid(), '11111111-1111-1111-1111-111111111111', 55.00, 'USD', 'Gift money saved', 'Gift', NOW() - INTERVAL '60 days');

-- Jane Smith's savings (User 2) - Irregular saver, broke streak
INSERT INTO savings_transactions (id, user_id, amount, currency, description, category, created_at)
VALUES
    -- Recent save (streak broken - more than 1 day ago)
    (gen_random_uuid(), '22222222-2222-2222-2222-222222222222', 50.00, 'USD', 'Weekly savings', 'Weekly', NOW() - INTERVAL '3 days'),
    
    -- Previous streak (5 days, ended 10 days ago)
    (gen_random_uuid(), '22222222-2222-2222-2222-222222222222', 20.00, 'USD', 'Coffee money', 'Daily Savings', NOW() - INTERVAL '10 days'),
    (gen_random_uuid(), '22222222-2222-2222-2222-222222222222', 25.00, 'USD', 'Lunch savings', 'Food', NOW() - INTERVAL '11 days'),
    (gen_random_uuid(), '22222222-2222-2222-2222-222222222222', 30.00, 'USD', 'Transport saved', 'Transport', NOW() - INTERVAL '12 days'),
    (gen_random_uuid(), '22222222-2222-2222-2222-222222222222', 15.00, 'USD', 'Snack money', 'Food', NOW() - INTERVAL '13 days'),
    (gen_random_uuid(), '22222222-2222-2222-2222-222222222222', 40.00, 'USD', 'Entertainment skip', 'Entertainment', NOW() - INTERVAL '14 days'),
    
    -- Older irregular saves
    (gen_random_uuid(), '22222222-2222-2222-2222-222222222222', 100.00, 'USD', 'Paycheck portion', 'Income', NOW() - INTERVAL '25 days'),
    (gen_random_uuid(), '22222222-2222-2222-2222-222222222222', 75.00, 'USD', 'Tax refund', 'Other', NOW() - INTERVAL '40 days'),
    (gen_random_uuid(), '22222222-2222-2222-2222-222222222222', 60.00, 'USD', 'Birthday money', 'Gift', NOW() - INTERVAL '55 days');

-- Mike Johnson's savings (User 3) - New saver with 3-day streak
INSERT INTO savings_transactions (id, user_id, amount, currency, description, category, created_at)
VALUES
    (gen_random_uuid(), '33333333-3333-3333-3333-333333333333', 10.00, 'USD', 'First savings!', 'Daily Savings', NOW() - INTERVAL '0 days'),
    (gen_random_uuid(), '33333333-3333-3333-3333-333333333333', 15.00, 'USD', 'Day 2 savings', 'Daily Savings', NOW() - INTERVAL '1 days'),
    (gen_random_uuid(), '33333333-3333-3333-3333-333333333333', 20.00, 'USD', 'Day 3 savings', 'Daily Savings', NOW() - INTERVAL '2 days'),
    (gen_random_uuid(), '33333333-3333-3333-3333-333333333333', 50.00, 'USD', 'Initial deposit', 'Income', NOW() - INTERVAL '10 days');

-- Sarah Williams's savings (User 4) - Improving pattern with 14-day streak
INSERT INTO savings_transactions (id, user_id, amount, currency, description, category, created_at)
VALUES
    -- Current 14-day streak
    (gen_random_uuid(), '44444444-4444-4444-4444-444444444444', 30.00, 'USD', 'Daily discipline', 'Daily Savings', NOW() - INTERVAL '0 days'),
    (gen_random_uuid(), '44444444-4444-4444-4444-444444444444', 25.00, 'USD', 'Morning routine', 'Daily Savings', NOW() - INTERVAL '1 days'),
    (gen_random_uuid(), '44444444-4444-4444-4444-444444444444', 35.00, 'USD', 'Lunch prep savings', 'Food', NOW() - INTERVAL '2 days'),
    (gen_random_uuid(), '44444444-4444-4444-4444-444444444444', 20.00, 'USD', 'Coffee skip', 'Daily Savings', NOW() - INTERVAL '3 days'),
    (gen_random_uuid(), '44444444-4444-4444-4444-444444444444', 40.00, 'USD', 'Bike to work', 'Transport', NOW() - INTERVAL '4 days'),
    (gen_random_uuid(), '44444444-4444-4444-4444-444444444444', 30.00, 'USD', 'Home cooking', 'Food', NOW() - INTERVAL '5 days'),
    (gen_random_uuid(), '44444444-4444-4444-4444-444444444444', 45.00, 'USD', 'Weekend savings', 'Daily Savings', NOW() - INTERVAL '6 days'),
    (gen_random_uuid(), '44444444-4444-4444-4444-444444444444', 25.00, 'USD', 'Daily goal', 'Daily Savings', NOW() - INTERVAL '7 days'),
    (gen_random_uuid(), '44444444-4444-4444-4444-444444444444', 30.00, 'USD', 'Consistent save', 'Daily Savings', NOW() - INTERVAL '8 days'),
    (gen_random_uuid(), '44444444-4444-4444-4444-444444444444', 35.00, 'USD', 'Morning discipline', 'Daily Savings', NOW() - INTERVAL '9 days'),
    (gen_random_uuid(), '44444444-4444-4444-4444-444444444444', 40.00, 'USD', 'Double day', 'Daily Savings', NOW() - INTERVAL '10 days'),
    (gen_random_uuid(), '44444444-4444-4444-4444-444444444444', 20.00, 'USD', 'Quick save', 'Daily Savings', NOW() - INTERVAL '11 days'),
    (gen_random_uuid(), '44444444-4444-4444-4444-444444444444', 50.00, 'USD', 'Milestone day', 'Bonus', NOW() - INTERVAL '12 days'),
    (gen_random_uuid(), '44444444-4444-4444-4444-444444444444', 30.00, 'USD', 'Two weeks strong', 'Daily Savings', NOW() - INTERVAL '13 days');

-- Alex Brown's savings (User 5) - Multiple transactions same day (should count as 1 day in streak)
INSERT INTO savings_transactions (id, user_id, amount, currency, description, category, created_at)
VALUES
    -- Today - multiple transactions
    (gen_random_uuid(), '55555555-5555-5555-5555-555555555555', 10.00, 'USD', 'Morning save', 'Daily Savings', NOW() - INTERVAL '0 days' + INTERVAL '8 hours'),
    (gen_random_uuid(), '55555555-5555-5555-5555-555555555555', 15.00, 'USD', 'Afternoon save', 'Daily Savings', NOW() - INTERVAL '0 days' + INTERVAL '14 hours'),
    (gen_random_uuid(), '55555555-5555-5555-5555-555555555555', 20.00, 'USD', 'Evening save', 'Daily Savings', NOW() - INTERVAL '0 days' + INTERVAL '20 hours'),
    
    -- Yesterday - multiple transactions
    (gen_random_uuid(), '55555555-5555-5555-5555-555555555555', 12.00, 'USD', 'Morning coffee', 'Daily Savings', NOW() - INTERVAL '1 days' + INTERVAL '9 hours'),
    (gen_random_uuid(), '55555555-5555-5555-5555-555555555555', 18.00, 'USD', 'Lunch money', 'Food', NOW() - INTERVAL '1 days' + INTERVAL '13 hours'),
    
    -- 2 days ago
    (gen_random_uuid(), '55555555-5555-5555-5555-555555555555', 25.00, 'USD', 'Daily goal', 'Daily Savings', NOW() - INTERVAL '2 days'),
    
    -- Older saves
    (gen_random_uuid(), '55555555-5555-5555-5555-555555555555', 50.00, 'USD', 'First week bonus', 'Bonus', NOW() - INTERVAL '7 days'),
    (gen_random_uuid(), '55555555-5555-5555-5555-555555555555', 100.00, 'USD', 'Initial deposit', 'Income', NOW() - INTERVAL '14 days');
