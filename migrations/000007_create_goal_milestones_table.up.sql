-- Create goal milestones table
CREATE TABLE goal_milestones (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    goal_id UUID NOT NULL REFERENCES goals(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    amount DECIMAL(15, 2) NOT NULL CHECK (amount > 0),
    is_completed BOOLEAN DEFAULT FALSE,
    completed_at TIMESTAMP WITH TIME ZONE,
    "order" INT NOT NULL
);

-- Create indexes for goal milestones
CREATE INDEX idx_milestones_goal ON goal_milestones(goal_id, "order");
CREATE INDEX idx_milestones_completed ON goal_milestones(goal_id, is_completed);

-- Add comments
COMMENT ON TABLE goal_milestones IS 'Stores intermediate milestones for financial goals';
COMMENT ON COLUMN goal_milestones.amount IS 'Milestone target amount';
COMMENT ON COLUMN goal_milestones."order" IS 'Display order of milestones within a goal';
