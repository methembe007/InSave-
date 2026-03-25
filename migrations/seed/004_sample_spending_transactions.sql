-- Sample Spending Transactions Seed Data
-- Creates spending records linked to budgets and categories

-- Get category IDs for John Doe's budget
DO $$
DECLARE
    v_budget_id UUID := 'b1111111-1111-1111-1111-111111111111';
    v_user_id UUID := '11111111-1111-1111-1111-111111111111';
    v_groceries_id UUID;
    v_transport_id UUID;
    v_entertainment_id UUID;
    v_dining_id UUID;
    v_shopping_id UUID;
    v_utilities_id UUID;
BEGIN
    -- Get category IDs
    SELECT id INTO v_groceries_id FROM budget_categories WHERE budget_id = v_budget_id AND name = 'Groceries';
    SELECT id INTO v_transport_id FROM budget_categories WHERE budget_id = v_budget_id AND name = 'Transportation';
    SELECT id INTO v_entertainment_id FROM budget_categories WHERE budget_id = v_budget_id AND name = 'Entertainment';
    SELECT id INTO v_dining_id FROM budget_categories WHERE budget_id = v_budget_id AND name = 'Dining Out';
    SELECT id INTO v_shopping_id FROM budget_categories WHERE budget_id = v_budget_id AND name = 'Shopping';
    SELECT id INTO v_utilities_id FROM budget_categories WHERE budget_id = v_budget_id AND name = 'Utilities';

    -- Insert spending transactions
    INSERT INTO spending_transactions (id, user_id, budget_id, category_id, amount, description, merchant, transaction_date, created_at)
    VALUES
        (gen_random_uuid(), v_user_id, v_budget_id, v_groceries_id, 85.50, 'Weekly groceries', 'Whole Foods', CURRENT_DATE - 2, NOW()),
        (gen_random_uuid(), v_user_id, v_budget_id, v_groceries_id, 120.00, 'Monthly stock up', 'Costco', CURRENT_DATE - 5, NOW()),
        (gen_random_uuid(), v_user_id, v_budget_id, v_groceries_id, 45.50, 'Fresh produce', 'Farmers Market', CURRENT_DATE - 8, NOW()),
        (gen_random_uuid(), v_user_id, v_budget_id, v_groceries_id, 69.00, 'Grocery run', 'Trader Joes', CURRENT_DATE - 12, NOW()),
        
        (gen_random_uuid(), v_user_id, v_budget_id, v_transport_id, 50.00, 'Gas fill up', 'Shell Station', CURRENT_DATE - 3, NOW()),
        (gen_random_uuid(), v_user_id, v_budget_id, v_transport_id, 60.00, 'Gas', 'Chevron', CURRENT_DATE - 10, NOW()),
        (gen_random_uuid(), v_user_id, v_budget_id, v_transport_id, 40.00, 'Parking fees', 'City Parking', CURRENT_DATE - 15, NOW()),
        
        (gen_random_uuid(), v_user_id, v_budget_id, v_entertainment_id, 45.00, 'Movie tickets', 'AMC Theaters', CURRENT_DATE - 6, NOW()),
        (gen_random_uuid(), v_user_id, v_budget_id, v_entertainment_id, 35.50, 'Concert tickets', 'Ticketmaster', CURRENT_DATE - 14, NOW()),
        (gen_random_uuid(), v_user_id, v_budget_id, v_entertainment_id, 40.00, 'Streaming services', 'Netflix', CURRENT_DATE - 1, NOW()),
        
        (gen_random_uuid(), v_user_id, v_budget_id, v_dining_id, 65.00, 'Dinner with friends', 'Italian Restaurant', CURRENT_DATE - 4, NOW()),
        (gen_random_uuid(), v_user_id, v_budget_id, v_dining_id, 85.00, 'Date night', 'Steakhouse', CURRENT_DATE - 9, NOW()),
        (gen_random_uuid(), v_user_id, v_budget_id, v_dining_id, 45.00, 'Lunch meeting', 'Cafe Downtown', CURRENT_DATE - 11, NOW()),
        (gen_random_uuid(), v_user_id, v_budget_id, v_dining_id, 85.00, 'Weekend brunch', 'Brunch Spot', CURRENT_DATE - 13, NOW()),
        
        (gen_random_uuid(), v_user_id, v_budget_id, v_shopping_id, 120.00, 'New shoes', 'Nike Store', CURRENT_DATE - 7, NOW()),
        (gen_random_uuid(), v_user_id, v_budget_id, v_shopping_id, 60.00, 'Clothing', 'H&M', CURRENT_DATE - 16, NOW()),
        
        (gen_random_uuid(), v_user_id, v_budget_id, v_utilities_id, 120.00, 'Electric bill', 'Power Company', CURRENT_DATE - 5, NOW()),
        (gen_random_uuid(), v_user_id, v_budget_id, v_utilities_id, 80.00, 'Internet', 'ISP Provider', CURRENT_DATE - 10, NOW());
END $$;

-- Jane Smith's spending (User 2) - Some over budget
DO $$
DECLARE
    v_budget_id UUID := 'b2222222-2222-2222-2222-222222222222';
    v_user_id UUID := '22222222-2222-2222-2222-222222222222';
    v_groceries_id UUID;
    v_transport_id UUID;
    v_entertainment_id UUID;
    v_dining_id UUID;
    v_shopping_id UUID;
    v_utilities_id UUID;
    v_subscriptions_id UUID;
BEGIN
    SELECT id INTO v_groceries_id FROM budget_categories WHERE budget_id = v_budget_id AND name = 'Groceries';
    SELECT id INTO v_transport_id FROM budget_categories WHERE budget_id = v_budget_id AND name = 'Transportation';
    SELECT id INTO v_entertainment_id FROM budget_categories WHERE budget_id = v_budget_id AND name = 'Entertainment';
    SELECT id INTO v_dining_id FROM budget_categories WHERE budget_id = v_budget_id AND name = 'Dining Out';
    SELECT id INTO v_shopping_id FROM budget_categories WHERE budget_id = v_budget_id AND name = 'Shopping';
    SELECT id INTO v_utilities_id FROM budget_categories WHERE budget_id = v_budget_id AND name = 'Utilities';
    SELECT id INTO v_subscriptions_id FROM budget_categories WHERE budget_id = v_budget_id AND name = 'Subscriptions';

    INSERT INTO spending_transactions (id, user_id, budget_id, category_id, amount, description, merchant, transaction_date, created_at)
    VALUES
        (gen_random_uuid(), v_user_id, v_budget_id, v_groceries_id, 95.00, 'Weekly shopping', 'Safeway', CURRENT_DATE - 3, NOW()),
        (gen_random_uuid(), v_user_id, v_budget_id, v_groceries_id, 135.00, 'Big grocery haul', 'Whole Foods', CURRENT_DATE - 8, NOW()),
        (gen_random_uuid(), v_user_id, v_budget_id, v_groceries_id, 150.00, 'Monthly restock', 'Costco', CURRENT_DATE - 15, NOW()),
        
        (gen_random_uuid(), v_user_id, v_budget_id, v_transport_id, 75.50, 'Gas and car wash', 'Shell', CURRENT_DATE - 4, NOW()),
        (gen_random_uuid(), v_user_id, v_budget_id, v_transport_id, 120.00, 'Uber rides', 'Uber', CURRENT_DATE - 12, NOW()),
        
        (gen_random_uuid(), v_user_id, v_budget_id, v_entertainment_id, 85.00, 'Concert', 'Live Nation', CURRENT_DATE - 6, NOW()),
        (gen_random_uuid(), v_user_id, v_budget_id, v_entertainment_id, 100.00, 'Theme park', 'Six Flags', CURRENT_DATE - 14, NOW()),
        
        (gen_random_uuid(), v_user_id, v_budget_id, v_dining_id, 110.25, 'Birthday dinner', 'Fine Dining', CURRENT_DATE - 5, NOW()),
        (gen_random_uuid(), v_user_id, v_budget_id, v_dining_id, 90.00, 'Weekend meals', 'Various', CURRENT_DATE - 10, NOW()),
        (gen_random_uuid(), v_user_id, v_budget_id, v_dining_id, 90.00, 'Takeout week', 'DoorDash', CURRENT_DATE - 16, NOW()),
        
        (gen_random_uuid(), v_user_id, v_budget_id, v_shopping_id, 150.00, 'New outfit', 'Zara', CURRENT_DATE - 7, NOW()),
        (gen_random_uuid(), v_user_id, v_budget_id, v_shopping_id, 100.00, 'Online shopping', 'Amazon', CURRENT_DATE - 13, NOW()),
        
        (gen_random_uuid(), v_user_id, v_budget_id, v_utilities_id, 200.00, 'Rent utilities', 'Landlord', CURRENT_DATE - 2, NOW()),
        (gen_random_uuid(), v_user_id, v_budget_id, v_utilities_id, 120.00, 'Electric & water', 'Utility Co', CURRENT_DATE - 11, NOW()),
        
        (gen_random_uuid(), v_user_id, v_budget_id, v_subscriptions_id, 15.00, 'Spotify Premium', 'Spotify', CURRENT_DATE - 1, NOW()),
        (gen_random_uuid(), v_user_id, v_budget_id, v_subscriptions_id, 15.00, 'Netflix', 'Netflix', CURRENT_DATE - 1, NOW());
END $$;

-- Mike Johnson's spending (User 3) - Minimal spending
DO $$
DECLARE
    v_budget_id UUID := 'b3333333-3333-3333-3333-333333333333';
    v_user_id UUID := '33333333-3333-3333-3333-333333333333';
    v_groceries_id UUID;
    v_transport_id UUID;
    v_entertainment_id UUID;
    v_dining_id UUID;
BEGIN
    SELECT id INTO v_groceries_id FROM budget_categories WHERE budget_id = v_budget_id AND name = 'Groceries';
    SELECT id INTO v_transport_id FROM budget_categories WHERE budget_id = v_budget_id AND name = 'Transportation';
    SELECT id INTO v_entertainment_id FROM budget_categories WHERE budget_id = v_budget_id AND name = 'Entertainment';
    SELECT id INTO v_dining_id FROM budget_categories WHERE budget_id = v_budget_id AND name = 'Dining Out';

    INSERT INTO spending_transactions (id, user_id, budget_id, category_id, amount, description, merchant, transaction_date, created_at)
    VALUES
        (gen_random_uuid(), v_user_id, v_budget_id, v_groceries_id, 120.00, 'Groceries', 'Aldi', CURRENT_DATE - 5, NOW()),
        (gen_random_uuid(), v_user_id, v_budget_id, v_groceries_id, 80.00, 'Food shopping', 'Walmart', CURRENT_DATE - 12, NOW()),
        
        (gen_random_uuid(), v_user_id, v_budget_id, v_transport_id, 60.00, 'Gas', 'Gas Station', CURRENT_DATE - 8, NOW()),
        (gen_random_uuid(), v_user_id, v_budget_id, v_transport_id, 40.00, 'Bus pass', 'Transit', CURRENT_DATE - 15, NOW()),
        
        (gen_random_uuid(), v_user_id, v_budget_id, v_entertainment_id, 50.00, 'Video games', 'Steam', CURRENT_DATE - 10, NOW()),
        
        (gen_random_uuid(), v_user_id, v_budget_id, v_dining_id, 100.00, 'Fast food', 'Various', CURRENT_DATE - 7, NOW());
END $$;

-- Sarah Williams's spending (User 4) - Near budget limits
DO $$
DECLARE
    v_budget_id UUID := 'b4444444-4444-4444-4444-444444444444';
    v_user_id UUID := '44444444-4444-4444-4444-444444444444';
    v_groceries_id UUID;
    v_transport_id UUID;
    v_entertainment_id UUID;
    v_dining_id UUID;
    v_shopping_id UUID;
    v_utilities_id UUID;
    v_healthcare_id UUID;
    v_misc_id UUID;
BEGIN
    SELECT id INTO v_groceries_id FROM budget_categories WHERE budget_id = v_budget_id AND name = 'Groceries';
    SELECT id INTO v_transport_id FROM budget_categories WHERE budget_id = v_budget_id AND name = 'Transportation';
    SELECT id INTO v_entertainment_id FROM budget_categories WHERE budget_id = v_budget_id AND name = 'Entertainment';
    SELECT id INTO v_dining_id FROM budget_categories WHERE budget_id = v_budget_id AND name = 'Dining Out';
    SELECT id INTO v_shopping_id FROM budget_categories WHERE budget_id = v_budget_id AND name = 'Shopping';
    SELECT id INTO v_utilities_id FROM budget_categories WHERE budget_id = v_budget_id AND name = 'Utilities';
    SELECT id INTO v_healthcare_id FROM budget_categories WHERE budget_id = v_budget_id AND name = 'Healthcare';
    SELECT id INTO v_misc_id FROM budget_categories WHERE budget_id = v_budget_id AND name = 'Miscellaneous';

    INSERT INTO spending_transactions (id, user_id, budget_id, category_id, amount, description, merchant, transaction_date, created_at)
    VALUES
        (gen_random_uuid(), v_user_id, v_budget_id, v_groceries_id, 130.00, 'Weekly groceries', 'Whole Foods', CURRENT_DATE - 2, NOW()),
        (gen_random_uuid(), v_user_id, v_budget_id, v_groceries_id, 140.00, 'Grocery shopping', 'Trader Joes', CURRENT_DATE - 9, NOW()),
        (gen_random_uuid(), v_user_id, v_budget_id, v_groceries_id, 120.00, 'Food supplies', 'Safeway', CURRENT_DATE - 16, NOW()),
        
        (gen_random_uuid(), v_user_id, v_budget_id, v_transport_id, 90.00, 'Gas', 'Shell', CURRENT_DATE - 4, NOW()),
        (gen_random_uuid(), v_user_id, v_budget_id, v_transport_id, 85.00, 'Gas', 'Chevron', CURRENT_DATE - 11, NOW()),
        (gen_random_uuid(), v_user_id, v_budget_id, v_transport_id, 85.00, 'Fuel', 'BP', CURRENT_DATE - 18, NOW()),
        
        (gen_random_uuid(), v_user_id, v_budget_id, v_entertainment_id, 85.00, 'Spa day', 'Spa & Wellness', CURRENT_DATE - 6, NOW()),
        (gen_random_uuid(), v_user_id, v_budget_id, v_entertainment_id, 90.00, 'Concert tickets', 'Ticketmaster', CURRENT_DATE - 13, NOW()),
        
        (gen_random_uuid(), v_user_id, v_budget_id, v_dining_id, 95.00, 'Restaurant', 'Fine Dining', CURRENT_DATE - 3, NOW()),
        (gen_random_uuid(), v_user_id, v_budget_id, v_dining_id, 100.00, 'Dinner out', 'Steakhouse', CURRENT_DATE - 10, NOW()),
        (gen_random_uuid(), v_user_id, v_budget_id, v_dining_id, 100.00, 'Weekend meals', 'Various', CURRENT_DATE - 17, NOW()),
        
        (gen_random_uuid(), v_user_id, v_budget_id, v_shopping_id, 170.00, 'New clothes', 'Nordstrom', CURRENT_DATE - 5, NOW()),
        (gen_random_uuid(), v_user_id, v_budget_id, v_shopping_id, 170.00, 'Shopping spree', 'Mall', CURRENT_DATE - 14, NOW()),
        
        (gen_random_uuid(), v_user_id, v_budget_id, v_utilities_id, 180.00, 'Electric bill', 'Power Co', CURRENT_DATE - 7, NOW()),
        (gen_random_uuid(), v_user_id, v_budget_id, v_utilities_id, 140.00, 'Internet & phone', 'Telecom', CURRENT_DATE - 15, NOW()),
        
        (gen_random_uuid(), v_user_id, v_budget_id, v_healthcare_id, 150.00, 'Doctor visit', 'Medical Center', CURRENT_DATE - 12, NOW()),
        
        (gen_random_uuid(), v_user_id, v_budget_id, v_misc_id, 50.00, 'Miscellaneous', 'Various', CURRENT_DATE - 8, NOW());
END $$;

-- Alex Brown's spending (User 5) - New user, light spending
DO $$
DECLARE
    v_budget_id UUID := 'b5555555-5555-5555-5555-555555555555';
    v_user_id UUID := '55555555-5555-5555-5555-555555555555';
    v_groceries_id UUID;
    v_transport_id UUID;
    v_entertainment_id UUID;
BEGIN
    SELECT id INTO v_groceries_id FROM budget_categories WHERE budget_id = v_budget_id AND name = 'Groceries';
    SELECT id INTO v_transport_id FROM budget_categories WHERE budget_id = v_budget_id AND name = 'Transportation';
    SELECT id INTO v_entertainment_id FROM budget_categories WHERE budget_id = v_budget_id AND name = 'Entertainment';

    INSERT INTO spending_transactions (id, user_id, budget_id, category_id, amount, description, merchant, transaction_date, created_at)
    VALUES
        (gen_random_uuid(), v_user_id, v_budget_id, v_groceries_id, 80.00, 'Groceries', 'Kroger', CURRENT_DATE - 4, NOW()),
        (gen_random_uuid(), v_user_id, v_budget_id, v_groceries_id, 70.00, 'Food shopping', 'Aldi', CURRENT_DATE - 11, NOW()),
        
        (gen_random_uuid(), v_user_id, v_budget_id, v_transport_id, 50.00, 'Gas', 'Gas Station', CURRENT_DATE - 6, NOW()),
        (gen_random_uuid(), v_user_id, v_budget_id, v_transport_id, 30.00, 'Parking', 'Parking Lot', CURRENT_DATE - 13, NOW()),
        
        (gen_random_uuid(), v_user_id, v_budget_id, v_entertainment_id, 50.00, 'Movies', 'Cinema', CURRENT_DATE - 9, NOW());
END $$;
