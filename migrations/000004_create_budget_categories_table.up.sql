-- Create budget categories table
CREATE TABLE budget_categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    budget_id UUID NOT NULL REFERENCES budgets(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    allocated_amount DECIMAL(15, 2) NOT NULL CHECK (allocated_amount >= 0),
    spent_amount DECIMAL(15, 2) DEFAULT 0 CHECK (spent_amount >= 0),
    color VARCHAR(7) DEFAULT '#000000'
);

-- Create indexes for budget categories
CREATE INDEX idx_budget_categories_budget ON budget_categories(budget_id);
CREATE INDEX idx_budget_categories_name ON budget_categories(name);

-- Add comments
COMMENT ON TABLE budget_categories IS 'Stores spending categories within monthly budgets';
COMMENT ON COLUMN budget_categories.allocated_amount IS 'Amount allocated to this category';
COMMENT ON COLUMN budget_categories.spent_amount IS 'Amount spent in this category, updated atomically';
COMMENT ON COLUMN budget_categories.color IS 'Hex color code for UI visualization';
