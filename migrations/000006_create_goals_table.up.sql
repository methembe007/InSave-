-- Create goals table
CREATE TABLE goals (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
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

-- Create indexes for goals table
CREATE INDEX idx_goals_user_status ON goals(user_id, status);
CREATE INDEX idx_goals_target_date ON goals(target_date);
CREATE INDEX idx_goals_user_id ON goals(user_id);

-- Add comments
COMMENT ON TABLE goals IS 'Stores user financial goals with target amounts and dates';
COMMENT ON COLUMN goals.target_amount IS 'Target amount for the goal';
COMMENT ON COLUMN goals.current_amount IS 'Current progress toward the goal';
COMMENT ON COLUMN goals.status IS 'Goal status: active, completed, or paused';
