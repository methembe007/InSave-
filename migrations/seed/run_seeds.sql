-- Master Seed Script
-- Run this script to populate the database with sample data for testing and development

-- Note: This script assumes all migrations have been run and tables exist

\echo 'Starting database seeding...'
\echo ''

\echo '1. Seeding sample users...'
\i 001_sample_users.sql
\echo '   ✓ Sample users created'
\echo ''

\echo '2. Seeding savings transactions...'
\i 002_sample_savings_transactions.sql
\echo '   ✓ Savings transactions created'
\echo ''

\echo '3. Seeding budgets and categories...'
\i 003_sample_budgets.sql
\echo '   ✓ Budgets and categories created'
\echo ''

\echo '4. Seeding spending transactions...'
\i 004_sample_spending_transactions.sql
\echo '   ✓ Spending transactions created'
\echo ''

\echo '5. Seeding goals and milestones...'
\i 005_sample_goals.sql
\echo '   ✓ Goals and milestones created'
\echo ''

\echo '6. Seeding education lessons and progress...'
\i 006_sample_education_lessons.sql
\echo '   ✓ Education content created'
\echo ''

\echo 'Database seeding complete!'
\echo ''
\echo 'Sample Users (all passwords: "password123"):'
\echo '  - john.doe@example.com (90 days history, 7-day streak)'
\echo '  - jane.smith@example.com (60 days history, broken streak)'
\echo '  - mike.johnson@example.com (45 days history, 3-day streak)'
\echo '  - sarah.williams@example.com (30 days history, 14-day streak)'
\echo '  - alex.brown@example.com (15 days history, 3-day streak with multiple daily saves)'
