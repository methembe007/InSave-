-- Integration Test Database Initialization Script
-- This script sets up the test database schema and seeds initial test data

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Users table
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    date_of_birth DATE NOT NULL,
    profile_image_url TEXT,
    preferences JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_created_at ON users(created_at);

-- Savings transactions table
CREATE TABLE IF NOT EXISTS savings_transactions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    amount DECIMAL(15, 2) NOT NULL CHECK (amount > 0),
    currency VARCHAR(3) NOT NULL DEFAULT 'USD',
    description TEXT,
    category VARCHAR(50),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_savings_user_date ON savings_transactions(user_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_savings_category ON savings_transactions(category);

-- Budgets table
CREATE TABLE IF NOT EXISTS budgets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    month DATE NOT NULL,
    total_budget DECIMAL(15, 2) NOT NULL CHECK (total_budget >= 0),
    total_spent DECIMAL(15, 2) DEFAULT 0 CHECK (total_spent >= 0),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(user_id, month)
);

CREATE INDEX IF NOT EXISTS idx_budgets_user_month ON budgets(user_id, month DESC);

-- Budget categories table
CREATE TABLE IF NOT EXISTS budget_categories (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    budget_id UUID NOT NULL REFERENCES budgets(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    allocated_amount DECIMAL(15, 2) NOT NULL CHECK (allocated_amount >= 0),
    spent_amount DECIMAL(15, 2) DEFAULT 0 CHECK (spent_amount >= 0),
    color VARCHAR(7) DEFAULT '#000000'
);

CREATE INDEX IF NOT EXISTS idx_budget_categories_budget ON budget_categories(budget_id);

-- Spending transactions table
CREATE TABLE IF NOT EXISTS spending_transactions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    budget_id UUID NOT NULL REFERENCES budgets(id) ON DELETE CASCADE,
    category_id UUID REFERENCES budget_categories(id) ON DELETE SET NULL,
    amount DECIMAL(15, 2) NOT NULL CHECK (amount > 0),
    description TEXT,
    merchant VARCHAR(255),
    transaction_date DATE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_spending_user_date ON spending_transactions(user_id, transaction_date DESC);
CREATE INDEX IF NOT EXISTS idx_spending_category ON spending_transactions(category_id);

-- Goals table
CREATE TABLE IF NOT EXISTS goals (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    target_amount DECIMAL(15, 2) NOT NULL CHECK (target_amount > 0),
    current_amount DECIMAL(15, 2) DEFAULT 0 CHECK (current_amount >= 0),
    currency VARCHAR(3) NOT NULL DEFAULT 'USD',
    target_date DATE NOT NULL,
    status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active', 'completed', 'paused')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_goals_user_status ON goals(user_id, status);
CREATE INDEX IF NOT EXISTS idx_goals_target_date ON goals(target_date);

-- Goal milestones table
CREATE TABLE IF NOT EXISTS goal_milestones (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    goal_id UUID NOT NULL REFERENCES goals(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    amount DECIMAL(15, 2) NOT NULL CHECK (amount > 0),
    is_completed BOOLEAN DEFAULT FALSE,
    completed_at TIMESTAMP WITH TIME ZONE,
    "order" INT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_milestones_goal ON goal_milestones(goal_id, "order");

-- Notifications table
CREATE TABLE IF NOT EXISTS notifications (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL,
    title VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    is_read BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_notifications_user_read ON notifications(user_id, is_read, created_at DESC);

-- Seed test data
-- Note: Password hash is for "TestPassword123!" using bcrypt cost 12
INSERT INTO users (id, email, password_hash, first_name, last_name, date_of_birth, preferences)
VALUES 
    ('11111111-1111-1111-1111-111111111111', 'test1@example.com', '$2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewY5GyYIpSRelgyG', 'Test', 'User1', '1995-01-15', '{"currency": "USD", "notifications_enabled": true}'),
    ('22222222-2222-2222-2222-222222222222', 'test2@example.com', '$2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewY5GyYIpSRelgyG', 'Test', 'User2', '1998-05-20', '{"currency": "USD", "notifications_enabled": true}')
ON CONFLICT (email) DO NOTHING;

-- Seed some savings transactions for test1
INSERT INTO savings_transactions (user_id, amount, currency, description, category, created_at)
VALUES 
    ('11111111-1111-1111-1111-111111111111', 50.00, 'USD', 'Weekly savings', 'general', NOW() - INTERVAL '2 days'),
    ('11111111-1111-1111-1111-111111111111', 25.00, 'USD', 'Coffee savings', 'food', NOW() - INTERVAL '1 day'),
    ('11111111-1111-1111-1111-111111111111', 100.00, 'USD', 'Monthly bonus', 'general', NOW())
ON CONFLICT DO NOTHING;

-- Seed a budget for test1
INSERT INTO budgets (id, user_id, month, total_budget, total_spent)
VALUES 
    ('33333333-3333-3333-3333-333333333333', '11111111-1111-1111-1111-111111111111', DATE_TRUNC('month', CURRENT_DATE), 1000.00, 0.00)
ON CONFLICT (user_id, month) DO NOTHING;

-- Seed budget categories
INSERT INTO budget_categories (budget_id, name, allocated_amount, spent_amount, color)
VALUES 
    ('33333333-3333-3333-3333-333333333333', 'Food', 300.00, 0.00, '#FF6B6B'),
    ('33333333-3333-3333-3333-333333333333', 'Transport', 200.00, 0.00, '#4ECDC4'),
    ('33333333-3333-3333-3333-333333333333', 'Entertainment', 150.00, 0.00, '#95E1D3'),
    ('33333333-3333-3333-3333-333333333333', 'Utilities', 250.00, 0.00, '#F38181')
ON CONFLICT DO NOTHING;

-- Seed a goal for test1
INSERT INTO goals (id, user_id, title, description, target_amount, current_amount, currency, target_date, status)
VALUES 
    ('44444444-4444-4444-4444-444444444444', '11111111-1111-1111-1111-111111111111', 'Emergency Fund', 'Build 6 months emergency fund', 5000.00, 0.00, 'USD', CURRENT_DATE + INTERVAL '12 months', 'active')
ON CONFLICT DO NOTHING;

-- Seed goal milestones
INSERT INTO goal_milestones (goal_id, title, amount, is_completed, "order")
VALUES 
    ('44444444-4444-4444-4444-444444444444', 'First $1000', 1000.00, false, 1),
    ('44444444-4444-4444-4444-444444444444', 'Halfway Point', 2500.00, false, 2),
    ('44444444-4444-4444-4444-444444444444', 'Almost There', 4000.00, false, 3),
    ('44444444-4444-4444-4444-444444444444', 'Goal Complete', 5000.00, false, 4)
ON CONFLICT DO NOTHING;

-- Grant permissions (if needed)
-- Note: In test environment, we're using the postgres superuser
