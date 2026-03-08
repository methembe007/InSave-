-- Create budgets table
CREATE TABLE budgets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    month DATE NOT NULL,
    total_budget DECIMAL(15, 2) NOT NULL CHECK (total_budget >= 0),
    total_spent DECIMAL(15, 2) DEFAULT 0 CHECK (total_spent >= 0),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(user_id, month)
);

-- Create indexes for budgets table
CREATE INDEX idx_budgets_user_month ON budgets(user_id, month DESC);
CREATE INDEX idx_budgets_month ON budgets(month);

-- Add comments
COMMENT ON TABLE budgets IS 'Stores monthly budget plans for users';
COMMENT ON COLUMN budgets.month IS 'First day of the month for the budget period';
COMMENT ON COLUMN budgets.total_budget IS 'Total budget amount for the month';
COMMENT ON COLUMN budgets.total_spent IS 'Total amount spent in the month, updated atomically with spending transactions';
