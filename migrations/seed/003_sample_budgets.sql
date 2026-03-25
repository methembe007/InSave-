-- Sample Budgets Seed Data
-- Creates monthly budgets with categories for testing

-- John Doe's budget (User 1) - Current month, well-managed
INSERT INTO budgets (id, user_id, month, total_budget, total_spent, created_at, updated_at)
VALUES
    ('b1111111-1111-1111-1111-111111111111', '11111111-1111-1111-1111-111111111111', DATE_TRUNC('month', CURRENT_DATE), 2000.00, 1250.50, NOW(), NOW());

INSERT INTO budget_categories (id, budget_id, name, allocated_amount, spent_amount, color)
VALUES
    (gen_random_uuid(), 'b1111111-1111-1111-1111-111111111111', 'Groceries', 400.00, 320.00, '#10B981'),
    (gen_random_uuid(), 'b1111111-1111-1111-1111-111111111111', 'Transportation', 200.00, 150.00, '#3B82F6'),
    (gen_random_uuid(), 'b1111111-1111-1111-1111-111111111111', 'Entertainment', 150.00, 120.50, '#8B5CF6'),
    (gen_random_uuid(), 'b1111111-1111-1111-1111-111111111111', 'Dining Out', 300.00, 280.00, '#F59E0B'),
    (gen_random_uuid(), 'b1111111-1111-1111-1111-111111111111', 'Shopping', 250.00, 180.00, '#EC4899'),
    (gen_random_uuid(), 'b1111111-1111-1111-1111-111111111111', 'Utilities', 350.00, 200.00, '#6366F1'),
    (gen_random_uuid(), 'b1111111-1111-1111-1111-111111111111', 'Healthcare', 200.00, 0.00, '#14B8A6'),
    (gen_random_uuid(), 'b1111111-1111-1111-1111-111111111111', 'Miscellaneous', 150.00, 0.00, '#64748B');

-- Jane Smith's budget (User 2) - Current month, some categories over budget
INSERT INTO budgets (id, user_id, month, total_budget, total_spent, created_at, updated_at)
VALUES
    ('b2222222-2222-2222-2222-222222222222', '22222222-2222-2222-2222-222222222222', DATE_TRUNC('month', CURRENT_DATE), 1800.00, 1650.75, NOW(), NOW());

INSERT INTO budget_categories (id, budget_id, name, allocated_amount, spent_amount, color)
VALUES
    (gen_random_uuid(), 'b2222222-2222-2222-2222-222222222222', 'Groceries', 350.00, 380.00, '#10B981'),
    (gen_random_uuid(), 'b2222222-2222-2222-2222-222222222222', 'Transportation', 180.00, 195.50, '#3B82F6'),
    (gen_random_uuid(), 'b2222222-2222-2222-2222-222222222222', 'Entertainment', 200.00, 185.00, '#8B5CF6'),
    (gen_random_uuid(), 'b2222222-2222-2222-2222-222222222222', 'Dining Out', 250.00, 290.25, '#F59E0B'),
    (gen_random_uuid(), 'b2222222-2222-2222-2222-222222222222', 'Shopping', 300.00, 250.00, '#EC4899'),
    (gen_random_uuid(), 'b2222222-2222-2222-2222-222222222222', 'Utilities', 320.00, 320.00, '#6366F1'),
    (gen_random_uuid(), 'b2222222-2222-2222-2222-222222222222', 'Subscriptions', 100.00, 30.00, '#F97316'),
    (gen_random_uuid(), 'b2222222-2222-2222-2222-222222222222', 'Miscellaneous', 100.00, 0.00, '#64748B');

-- Mike Johnson's budget (User 3) - Current month, minimal spending
INSERT INTO budgets (id, user_id, month, total_budget, total_spent, created_at, updated_at)
VALUES
    ('b3333333-3333-3333-3333-333333333333', '33333333-3333-3333-3333-333333333333', DATE_TRUNC('month', CURRENT_DATE), 1500.00, 450.00, NOW(), NOW());

INSERT INTO budget_categories (id, budget_id, name, allocated_amount, spent_amount, color)
VALUES
    (gen_random_uuid(), 'b3333333-3333-3333-3333-333333333333', 'Groceries', 300.00, 200.00, '#10B981'),
    (gen_random_uuid(), 'b3333333-3333-3333-3333-333333333333', 'Transportation', 150.00, 100.00, '#3B82F6'),
    (gen_random_uuid(), 'b3333333-3333-3333-3333-333333333333', 'Entertainment', 100.00, 50.00, '#8B5CF6'),
    (gen_random_uuid(), 'b3333333-3333-3333-3333-333333333333', 'Dining Out', 200.00, 100.00, '#F59E0B'),
    (gen_random_uuid(), 'b3333333-3333-3333-3333-333333333333', 'Utilities', 400.00, 0.00, '#6366F1'),
    (gen_random_uuid(), 'b3333333-3333-3333-3333-333333333333', 'Miscellaneous', 350.00, 0.00, '#64748B');

-- Sarah Williams's budget (User 4) - Current month, critical alerts
INSERT INTO budgets (id, user_id, month, total_budget, total_spent, created_at, updated_at)
VALUES
    ('b4444444-4444-4444-4444-444444444444', '44444444-4444-4444-4444-444444444444', DATE_TRUNC('month', CURRENT_DATE), 2200.00, 1980.00, NOW(), NOW());

INSERT INTO budget_categories (id, budget_id, name, allocated_amount, spent_amount, color)
VALUES
    (gen_random_uuid(), 'b4444444-4444-4444-4444-444444444444', 'Groceries', 400.00, 390.00, '#10B981'),
    (gen_random_uuid(), 'b4444444-4444-4444-4444-444444444444', 'Transportation', 250.00, 260.00, '#3B82F6'),
    (gen_random_uuid(), 'b4444444-4444-4444-4444-444444444444', 'Entertainment', 180.00, 175.00, '#8B5CF6'),
    (gen_random_uuid(), 'b4444444-4444-4444-4444-444444444444', 'Dining Out', 300.00, 295.00, '#F59E0B'),
    (gen_random_uuid(), 'b4444444-4444-4444-4444-444444444444', 'Shopping', 350.00, 340.00, '#EC4899'),
    (gen_random_uuid(), 'b4444444-4444-4444-4444-444444444444', 'Utilities', 400.00, 320.00, '#6366F1'),
    (gen_random_uuid(), 'b4444444-4444-4444-4444-444444444444', 'Healthcare', 200.00, 150.00, '#14B8A6'),
    (gen_random_uuid(), 'b4444444-4444-4444-4444-444444444444', 'Miscellaneous', 120.00, 50.00, '#64748B');

-- Alex Brown's budget (User 5) - Current month, new user
INSERT INTO budgets (id, user_id, month, total_budget, total_spent, created_at, updated_at)
VALUES
    ('b5555555-5555-5555-5555-555555555555', '55555555-5555-5555-5555-555555555555', DATE_TRUNC('month', CURRENT_DATE), 1600.00, 280.00, NOW(), NOW());

INSERT INTO budget_categories (id, budget_id, name, allocated_amount, spent_amount, color)
VALUES
    (gen_random_uuid(), 'b5555555-5555-5555-5555-555555555555', 'Groceries', 350.00, 150.00, '#10B981'),
    (gen_random_uuid(), 'b5555555-5555-5555-5555-555555555555', 'Transportation', 200.00, 80.00, '#3B82F6'),
    (gen_random_uuid(), 'b5555555-5555-5555-5555-555555555555', 'Entertainment', 150.00, 50.00, '#8B5CF6'),
    (gen_random_uuid(), 'b5555555-5555-5555-5555-555555555555', 'Dining Out', 250.00, 0.00, '#F59E0B'),
    (gen_random_uuid(), 'b5555555-5555-5555-5555-555555555555', 'Utilities', 400.00, 0.00, '#6366F1'),
    (gen_random_uuid(), 'b5555555-5555-5555-5555-555555555555', 'Miscellaneous', 250.00, 0.00, '#64748B');
