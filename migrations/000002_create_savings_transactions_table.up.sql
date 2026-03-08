-- Create savings transactions table with partitioning by month
CREATE TABLE savings_transactions (
    id UUID DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    amount DECIMAL(15, 2) NOT NULL CHECK (amount > 0),
    currency VARCHAR(3) NOT NULL DEFAULT 'USD',
    description TEXT,
    category VARCHAR(50),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    PRIMARY KEY (id, created_at)
) PARTITION BY RANGE (created_at);

-- Create indexes for savings transactions
CREATE INDEX idx_savings_user_date ON savings_transactions(user_id, created_at DESC);
CREATE INDEX idx_savings_category ON savings_transactions(category);
CREATE INDEX idx_savings_created_at ON savings_transactions(created_at);

-- Create partitions for 2024-2027 (can be extended as needed)
-- 2024 partitions
CREATE TABLE savings_transactions_2024_01 PARTITION OF savings_transactions
    FOR VALUES FROM ('2024-01-01') TO ('2024-02-01');
CREATE TABLE savings_transactions_2024_02 PARTITION OF savings_transactions
    FOR VALUES FROM ('2024-02-01') TO ('2024-03-01');
CREATE TABLE savings_transactions_2024_03 PARTITION OF savings_transactions
    FOR VALUES FROM ('2024-03-01') TO ('2024-04-01');
CREATE TABLE savings_transactions_2024_04 PARTITION OF savings_transactions
    FOR VALUES FROM ('2024-04-01') TO ('2024-05-01');
CREATE TABLE savings_transactions_2024_05 PARTITION OF savings_transactions
    FOR VALUES FROM ('2024-05-01') TO ('2024-06-01');
CREATE TABLE savings_transactions_2024_06 PARTITION OF savings_transactions
    FOR VALUES FROM ('2024-06-01') TO ('2024-07-01');
CREATE TABLE savings_transactions_2024_07 PARTITION OF savings_transactions
    FOR VALUES FROM ('2024-07-01') TO ('2024-08-01');
CREATE TABLE savings_transactions_2024_08 PARTITION OF savings_transactions
    FOR VALUES FROM ('2024-08-01') TO ('2024-09-01');
CREATE TABLE savings_transactions_2024_09 PARTITION OF savings_transactions
    FOR VALUES FROM ('2024-09-01') TO ('2024-10-01');
CREATE TABLE savings_transactions_2024_10 PARTITION OF savings_transactions
    FOR VALUES FROM ('2024-10-01') TO ('2024-11-01');
CREATE TABLE savings_transactions_2024_11 PARTITION OF savings_transactions
    FOR VALUES FROM ('2024-11-01') TO ('2024-12-01');
CREATE TABLE savings_transactions_2024_12 PARTITION OF savings_transactions
    FOR VALUES FROM ('2024-12-01') TO ('2025-01-01');

-- 2025 partitions
CREATE TABLE savings_transactions_2025_01 PARTITION OF savings_transactions
    FOR VALUES FROM ('2025-01-01') TO ('2025-02-01');
CREATE TABLE savings_transactions_2025_02 PARTITION OF savings_transactions
    FOR VALUES FROM ('2025-02-01') TO ('2025-03-01');
CREATE TABLE savings_transactions_2025_03 PARTITION OF savings_transactions
    FOR VALUES FROM ('2025-03-01') TO ('2025-04-01');
CREATE TABLE savings_transactions_2025_04 PARTITION OF savings_transactions
    FOR VALUES FROM ('2025-04-01') TO ('2025-05-01');
CREATE TABLE savings_transactions_2025_05 PARTITION OF savings_transactions
    FOR VALUES FROM ('2025-05-01') TO ('2025-06-01');
CREATE TABLE savings_transactions_2025_06 PARTITION OF savings_transactions
    FOR VALUES FROM ('2025-06-01') TO ('2025-07-01');
CREATE TABLE savings_transactions_2025_07 PARTITION OF savings_transactions
    FOR VALUES FROM ('2025-07-01') TO ('2025-08-01');
CREATE TABLE savings_transactions_2025_08 PARTITION OF savings_transactions
    FOR VALUES FROM ('2025-08-01') TO ('2025-09-01');
CREATE TABLE savings_transactions_2025_09 PARTITION OF savings_transactions
    FOR VALUES FROM ('2025-09-01') TO ('2025-10-01');
CREATE TABLE savings_transactions_2025_10 PARTITION OF savings_transactions
    FOR VALUES FROM ('2025-10-01') TO ('2025-11-01');
CREATE TABLE savings_transactions_2025_11 PARTITION OF savings_transactions
    FOR VALUES FROM ('2025-11-01') TO ('2025-12-01');
CREATE TABLE savings_transactions_2025_12 PARTITION OF savings_transactions
    FOR VALUES FROM ('2025-12-01') TO ('2026-01-01');

-- 2026 partitions
CREATE TABLE savings_transactions_2026_01 PARTITION OF savings_transactions
    FOR VALUES FROM ('2026-01-01') TO ('2026-02-01');
CREATE TABLE savings_transactions_2026_02 PARTITION OF savings_transactions
    FOR VALUES FROM ('2026-02-01') TO ('2026-03-01');
CREATE TABLE savings_transactions_2026_03 PARTITION OF savings_transactions
    FOR VALUES FROM ('2026-03-01') TO ('2026-04-01');
CREATE TABLE savings_transactions_2026_04 PARTITION OF savings_transactions
    FOR VALUES FROM ('2026-04-01') TO ('2026-05-01');
CREATE TABLE savings_transactions_2026_05 PARTITION OF savings_transactions
    FOR VALUES FROM ('2026-05-01') TO ('2026-06-01');
CREATE TABLE savings_transactions_2026_06 PARTITION OF savings_transactions
    FOR VALUES FROM ('2026-06-01') TO ('2026-07-01');
CREATE TABLE savings_transactions_2026_07 PARTITION OF savings_transactions
    FOR VALUES FROM ('2026-07-01') TO ('2026-08-01');
CREATE TABLE savings_transactions_2026_08 PARTITION OF savings_transactions
    FOR VALUES FROM ('2026-08-01') TO ('2026-09-01');
CREATE TABLE savings_transactions_2026_09 PARTITION OF savings_transactions
    FOR VALUES FROM ('2026-09-01') TO ('2026-10-01');
CREATE TABLE savings_transactions_2026_10 PARTITION OF savings_transactions
    FOR VALUES FROM ('2026-10-01') TO ('2026-11-01');
CREATE TABLE savings_transactions_2026_11 PARTITION OF savings_transactions
    FOR VALUES FROM ('2026-11-01') TO ('2026-12-01');
CREATE TABLE savings_transactions_2026_12 PARTITION OF savings_transactions
    FOR VALUES FROM ('2026-12-01') TO ('2027-01-01');

-- 2027 partitions
CREATE TABLE savings_transactions_2027_01 PARTITION OF savings_transactions
    FOR VALUES FROM ('2027-01-01') TO ('2027-02-01');
CREATE TABLE savings_transactions_2027_02 PARTITION OF savings_transactions
    FOR VALUES FROM ('2027-02-01') TO ('2027-03-01');
CREATE TABLE savings_transactions_2027_03 PARTITION OF savings_transactions
    FOR VALUES FROM ('2027-03-01') TO ('2027-04-01');
CREATE TABLE savings_transactions_2027_04 PARTITION OF savings_transactions
    FOR VALUES FROM ('2027-04-01') TO ('2027-05-01');
CREATE TABLE savings_transactions_2027_05 PARTITION OF savings_transactions
    FOR VALUES FROM ('2027-05-01') TO ('2027-06-01');
CREATE TABLE savings_transactions_2027_06 PARTITION OF savings_transactions
    FOR VALUES FROM ('2027-06-01') TO ('2027-07-01');
CREATE TABLE savings_transactions_2027_07 PARTITION OF savings_transactions
    FOR VALUES FROM ('2027-07-01') TO ('2027-08-01');
CREATE TABLE savings_transactions_2027_08 PARTITION OF savings_transactions
    FOR VALUES FROM ('2027-08-01') TO ('2027-09-01');
CREATE TABLE savings_transactions_2027_09 PARTITION OF savings_transactions
    FOR VALUES FROM ('2027-09-01') TO ('2027-10-01');
CREATE TABLE savings_transactions_2027_10 PARTITION OF savings_transactions
    FOR VALUES FROM ('2027-10-01') TO ('2027-11-01');
CREATE TABLE savings_transactions_2027_11 PARTITION OF savings_transactions
    FOR VALUES FROM ('2027-11-01') TO ('2027-12-01');
CREATE TABLE savings_transactions_2027_12 PARTITION OF savings_transactions
    FOR VALUES FROM ('2027-12-01') TO ('2028-01-01');

-- Add comments
COMMENT ON TABLE savings_transactions IS 'Stores user savings transactions, partitioned by month for performance';
COMMENT ON COLUMN savings_transactions.amount IS 'Savings amount stored as DECIMAL(15,2) with 2 decimal places precision';
