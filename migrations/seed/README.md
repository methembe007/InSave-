# Database Seed Scripts

This directory contains SQL scripts to populate the InSavein database with sample data for testing and development.

## Overview

The seed scripts create realistic sample data including:
- 5 sample users with different usage patterns
- Savings transactions with various streak scenarios
- Monthly budgets with spending categories
- Spending transactions across different categories
- Financial goals with milestones at various completion stages
- Educational lessons with user progress tracking

## Running the Seeds

### Using PostgreSQL CLI

From the `migrations/seed` directory:

```bash
psql -U postgres -d insavein -f run_seeds.sql
```

### Using Docker

If running PostgreSQL in Docker:

```bash
docker exec -i insavein-postgres psql -U postgres -d insavein < run_seeds.sql
```

### Windows (PowerShell)

```powershell
Get-Content run_seeds.sql | docker exec -i insavein-postgres psql -U postgres -d insavein
```

## Sample Users

All users have the password: `password123` (bcrypt hashed with cost factor 12)

| Email | Name | History | Streak | Profile |
|-------|------|---------|--------|---------|
| john.doe@example.com | John Doe | 90 days | 7 days | Consistent saver, well-managed budget |
| jane.smith@example.com | Jane Smith | 60 days | Broken | Irregular saver, some overspending |
| mike.johnson@example.com | Mike Johnson | 45 days | 3 days | New saver, minimal spending |
| sarah.williams@example.com | Sarah Williams | 30 days | 14 days | Improving pattern, near budget limits |
| alex.brown@example.com | Alex Brown | 15 days | 3 days | Multiple daily saves (streak testing) |

## Data Scenarios

### Savings Streaks
- **John Doe**: Current 7-day streak (consistent daily saver)
- **Jane Smith**: Broken streak (last save 3 days ago, had 5-day streak)
- **Mike Johnson**: New user with 3-day streak
- **Sarah Williams**: Strong 14-day streak (improving pattern)
- **Alex Brown**: Multiple transactions per day (tests same-day counting)

### Budget Alerts
- **John Doe**: Well under budget in all categories
- **Jane Smith**: Over budget in Groceries, Transportation, and Dining Out
- **Mike Johnson**: Minimal spending, well under budget
- **Sarah Williams**: Multiple categories at 80%+ (warning alerts) and 100%+ (critical alerts)
- **Alex Brown**: New user, light spending

### Goals Progress
- **John Doe**: 3 goals (2 active, 1 completed)
  - Emergency Fund: 45% complete
  - Europe Vacation: 56% complete
  - New Laptop: 100% complete ✓
- **Jane Smith**: 2 active goals (early stage)
  - Car Down Payment: 15% complete
  - Credit Card Payoff: 23% complete
- **Mike Johnson**: 1 goal (just starting)
  - First $1,000: 8.5% complete
- **Sarah Williams**: 3 ambitious goals
  - House Down Payment: 28% complete
  - Investment Portfolio: 64% complete
  - Wedding Fund: 45% complete
- **Alex Brown**: 1 simple goal
  - New Phone: 20% complete

### Education Progress
- **John Doe**: 5/8 lessons completed (62.5%)
- **Jane Smith**: 3/8 lessons completed (37.5%)
- **Mike Johnson**: 1/8 lessons completed (12.5%)
- **Sarah Williams**: 8/8 lessons completed (100%) ✓
- **Alex Brown**: 2/8 lessons completed (25%)

## Individual Seed Files

1. **001_sample_users.sql** - Creates 5 users with different preferences
2. **002_sample_savings_transactions.sql** - Savings history with various date patterns
3. **003_sample_budgets.sql** - Monthly budgets with category allocations
4. **004_sample_spending_transactions.sql** - Spending records linked to budgets
5. **005_sample_goals.sql** - Financial goals with milestone tracking
6. **006_sample_education_lessons.sql** - Educational content and user progress

## Clearing Seed Data

To remove all seed data and start fresh:

```sql
-- Delete in reverse order to respect foreign keys
DELETE FROM education_progress WHERE user_id IN (
    '11111111-1111-1111-1111-111111111111',
    '22222222-2222-2222-2222-222222222222',
    '33333333-3333-3333-3333-333333333333',
    '44444444-4444-4444-4444-444444444444',
    '55555555-5555-5555-5555-555555555555'
);

DELETE FROM goal_milestones WHERE goal_id IN (
    SELECT id FROM goals WHERE user_id IN (
        '11111111-1111-1111-1111-111111111111',
        '22222222-2222-2222-2222-222222222222',
        '33333333-3333-3333-3333-333333333333',
        '44444444-4444-4444-4444-444444444444',
        '55555555-5555-5555-5555-555555555555'
    )
);

DELETE FROM goals WHERE user_id IN (
    '11111111-1111-1111-1111-111111111111',
    '22222222-2222-2222-2222-222222222222',
    '33333333-3333-3333-3333-333333333333',
    '44444444-4444-4444-4444-444444444444',
    '55555555-5555-5555-5555-555555555555'
);

DELETE FROM spending_transactions WHERE user_id IN (
    '11111111-1111-1111-1111-111111111111',
    '22222222-2222-2222-2222-222222222222',
    '33333333-3333-3333-3333-333333333333',
    '44444444-4444-4444-4444-444444444444',
    '55555555-5555-5555-5555-555555555555'
);

DELETE FROM budget_categories WHERE budget_id IN (
    SELECT id FROM budgets WHERE user_id IN (
        '11111111-1111-1111-1111-111111111111',
        '22222222-2222-2222-2222-222222222222',
        '33333333-3333-3333-3333-333333333333',
        '44444444-4444-4444-4444-444444444444',
        '55555555-5555-5555-5555-555555555555'
    )
);

DELETE FROM budgets WHERE user_id IN (
    '11111111-1111-1111-1111-111111111111',
    '22222222-2222-2222-2222-222222222222',
    '33333333-3333-3333-3333-333333333333',
    '44444444-4444-4444-4444-444444444444',
    '55555555-5555-5555-5555-555555555555'
);

DELETE FROM savings_transactions WHERE user_id IN (
    '11111111-1111-1111-1111-111111111111',
    '22222222-2222-2222-2222-222222222222',
    '33333333-3333-3333-3333-333333333333',
    '44444444-4444-4444-4444-444444444444',
    '55555555-5555-5555-5555-555555555555'
);

DELETE FROM users WHERE id IN (
    '11111111-1111-1111-1111-111111111111',
    '22222222-2222-2222-2222-222222222222',
    '33333333-3333-3333-3333-333333333333',
    '44444444-4444-4444-4444-444444444444',
    '55555555-5555-5555-5555-555555555555'
);

-- Also delete sample lessons if needed
DELETE FROM lessons WHERE id LIKE 'l0000%';
```

## Testing Scenarios

### Streak Calculation Testing
- Use Alex Brown's account to test multiple transactions on the same day
- Use Jane Smith's account to test broken streaks
- Use Sarah Williams's account to test long streaks

### Budget Alert Testing
- Use Sarah Williams's account to test warning (80-99%) and critical (100%+) alerts
- Use Jane Smith's account to test over-budget scenarios

### Goal Progress Testing
- Use John Doe's completed laptop goal to test completion flow
- Use Sarah Williams's multiple goals to test concurrent goal tracking
- Use Mike Johnson's new goal to test early-stage progress

### Education Progress Testing
- Use Sarah Williams's account to test 100% completion
- Use other accounts to test partial completion scenarios

## Notes

- All timestamps are relative to NOW() for consistency
- UUIDs are deterministic for users and budgets to enable foreign key relationships
- Transaction IDs use gen_random_uuid() for uniqueness
- Partitioned tables (savings_transactions, spending_transactions) will automatically route to appropriate partitions
