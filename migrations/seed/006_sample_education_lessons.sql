-- Sample Education Lessons Seed Data
-- Creates financial education content for testing

INSERT INTO lessons (id, title, description, category, duration_minutes, difficulty, content, video_url, resources, "order")
VALUES
    -- Beginner Lessons
    ('l0000001-0000-0000-0000-000000000001', 
     'Introduction to Personal Finance', 
     'Learn the basics of managing your money, including budgeting, saving, and spending wisely.',
     'Basics',
     15,
     'beginner',
     '# Introduction to Personal Finance

## What is Personal Finance?

Personal finance is the management of your money and financial decisions. It includes:
- **Budgeting**: Planning how to spend your money
- **Saving**: Setting aside money for future needs
- **Investing**: Growing your wealth over time
- **Debt Management**: Handling loans and credit responsibly

## Why is it Important?

Good personal finance habits help you:
1. Achieve your financial goals
2. Reduce financial stress
3. Build wealth over time
4. Prepare for emergencies
5. Secure your future

## Getting Started

The first step is understanding where your money goes. Track your spending for a month to see your patterns.

## Key Takeaways

- Personal finance is about making smart decisions with your money
- Everyone can learn to manage money better
- Small changes can make a big difference over time',
     'https://www.youtube.com/watch?v=example1',
     '{"articles": ["https://www.investopedia.com/personal-finance-4427760"], "tools": ["Budget calculator", "Expense tracker"]}',
     1),

    ('l0000002-0000-0000-0000-000000000002',
     'Creating Your First Budget',
     'Step-by-step guide to creating a realistic budget that works for your lifestyle.',
     'Budgeting',
     20,
     'beginner',
     '# Creating Your First Budget

## The 50/30/20 Rule

A simple budgeting framework:
- **50%** for Needs (rent, food, utilities)
- **30%** for Wants (entertainment, dining out)
- **20%** for Savings and debt repayment

## Steps to Create Your Budget

1. **Calculate Your Income**: Add up all money coming in
2. **List Your Expenses**: Track everything you spend
3. **Categorize Spending**: Group expenses into needs, wants, and savings
4. **Set Limits**: Decide how much to spend in each category
5. **Track and Adjust**: Monitor your spending and make changes as needed

## Common Budget Categories

- Housing (rent/mortgage)
- Transportation
- Food (groceries and dining out)
- Utilities
- Insurance
- Entertainment
- Savings

## Tips for Success

- Be realistic about your spending
- Leave room for unexpected expenses
- Review your budget monthly
- Use budgeting apps to track automatically',
     'https://www.youtube.com/watch?v=example2',
     '{"articles": ["https://www.nerdwallet.com/article/finance/how-to-budget"], "templates": ["Monthly budget template", "Expense tracker spreadsheet"]}',
     2),

    ('l0000003-0000-0000-0000-000000000003',
     'The Power of Saving',
     'Understand why saving is crucial and learn strategies to build your savings habit.',
     'Saving',
     18,
     'beginner',
     '# The Power of Saving

## Why Save Money?

Saving provides:
- **Emergency Fund**: Protection against unexpected expenses
- **Financial Goals**: Money for things you want
- **Peace of Mind**: Less stress about money
- **Future Security**: Preparation for retirement

## How Much Should You Save?

Start with these goals:
1. **Emergency Fund**: 3-6 months of expenses
2. **Short-term Goals**: Vacations, purchases (1-3 years)
3. **Long-term Goals**: House, retirement (5+ years)

## Saving Strategies

### Pay Yourself First
Set up automatic transfers to savings on payday.

### The 24-Hour Rule
Wait 24 hours before making non-essential purchases.

### Round-Up Savings
Round up purchases to the nearest dollar and save the difference.

### Challenge Yourself
Try a 30-day no-spend challenge or save $1 more each day.

## Making Saving Easier

- Automate your savings
- Start small and increase gradually
- Celebrate milestones
- Keep savings in a separate account',
     'https://www.youtube.com/watch?v=example3',
     '{"articles": ["https://www.bankrate.com/banking/savings/how-to-save-money/"], "calculators": ["Savings goal calculator", "Emergency fund calculator"]}',
     3),

    -- Intermediate Lessons
    ('l0000004-0000-0000-0000-000000000004',
     'Understanding Credit and Debt',
     'Learn how credit works, how to build good credit, and strategies for managing debt.',
     'Credit',
     25,
     'intermediate',
     '# Understanding Credit and Debt

## What is Credit?

Credit is borrowed money that you promise to pay back. It includes:
- Credit cards
- Personal loans
- Student loans
- Auto loans
- Mortgages

## Credit Scores

Your credit score (300-850) affects:
- Loan approval
- Interest rates
- Rental applications
- Insurance premiums

### Factors Affecting Your Score

1. **Payment History (35%)**: Pay on time
2. **Credit Utilization (30%)**: Keep balances low
3. **Credit History Length (15%)**: Keep old accounts open
4. **Credit Mix (10%)**: Different types of credit
5. **New Credit (10%)**: Limit new applications

## Managing Debt

### The Debt Snowball Method
Pay off smallest debts first for quick wins.

### The Debt Avalanche Method
Pay off highest interest debts first to save money.

## Best Practices

- Pay more than the minimum
- Avoid new debt while paying off existing debt
- Negotiate lower interest rates
- Consider debt consolidation if appropriate',
     'https://www.youtube.com/watch?v=example4',
     '{"articles": ["https://www.experian.com/blogs/ask-experian/credit-education/"], "tools": ["Credit score simulator", "Debt payoff calculator"]}',
     4),

    ('l0000005-0000-0000-0000-000000000005',
     'Smart Shopping and Spending',
     'Techniques to make better purchasing decisions and avoid impulse buying.',
     'Spending',
     22,
     'intermediate',
     '# Smart Shopping and Spending

## The True Cost of Purchases

Consider:
- **Opportunity Cost**: What else could you do with that money?
- **Cost Per Use**: How often will you use it?
- **Maintenance Costs**: Ongoing expenses

## Avoiding Impulse Purchases

### The 30-Day Rule
Wait 30 days before buying non-essentials. If you still want it, then buy it.

### The 10-10 Rule
Will this matter in 10 days? 10 months? 10 years?

### Unsubscribe from Marketing
Reduce temptation by unsubscribing from promotional emails.

## Smart Shopping Strategies

1. **Make a List**: Stick to it when shopping
2. **Compare Prices**: Use price comparison tools
3. **Buy Generic**: Often same quality, lower price
4. **Use Cash**: Harder to overspend
5. **Shop Sales**: But only for things you need

## Subscription Audit

Review all subscriptions monthly:
- Are you using them?
- Can you get a better deal?
- Can you share with family?

## Quality vs. Price

Sometimes paying more upfront saves money long-term. Consider:
- Durability
- Warranty
- Energy efficiency
- Maintenance costs',
     'https://www.youtube.com/watch?v=example5',
     '{"articles": ["https://www.consumerreports.org/money/"], "apps": ["Price comparison apps", "Coupon apps"]}',
     5),

    ('l0000006-0000-0000-0000-000000000006',
     'Building an Emergency Fund',
     'Step-by-step guide to building and maintaining an emergency fund.',
     'Saving',
     20,
     'intermediate',
     '# Building an Emergency Fund

## What is an Emergency Fund?

Money set aside for unexpected expenses:
- Medical emergencies
- Car repairs
- Job loss
- Home repairs
- Unexpected travel

## How Much Do You Need?

### Starter Emergency Fund
$1,000 for immediate emergencies

### Full Emergency Fund
3-6 months of essential expenses

### Calculate Your Target
Add up monthly costs for:
- Rent/mortgage
- Utilities
- Food
- Transportation
- Insurance
- Minimum debt payments

Multiply by 3-6 months.

## Where to Keep It

- **High-yield savings account**: Easy access, earns interest
- **Money market account**: Slightly higher rates
- **NOT in**: Checking account (too tempting) or investments (too risky)

## Building Your Fund

1. **Start Small**: Even $25/month adds up
2. **Automate**: Set up automatic transfers
3. **Use Windfalls**: Tax refunds, bonuses, gifts
4. **Reduce Expenses**: Find areas to cut temporarily
5. **Increase Income**: Side gigs, overtime

## Maintaining Your Fund

- Only use for true emergencies
- Replenish after using
- Adjust as life changes
- Review annually',
     'https://www.youtube.com/watch?v=example6',
     '{"articles": ["https://www.nerdwallet.com/article/banking/emergency-fund"], "calculators": ["Emergency fund calculator"]}',
     6),

    -- Advanced Lessons
    ('l0000007-0000-0000-0000-000000000007',
     'Introduction to Investing',
     'Learn the basics of investing, including stocks, bonds, and retirement accounts.',
     'Investing',
     30,
     'advanced',
     '# Introduction to Investing

## Why Invest?

Investing helps your money grow faster than savings accounts through:
- **Compound Interest**: Earning returns on your returns
- **Inflation Protection**: Keeping pace with rising costs
- **Wealth Building**: Growing your net worth over time

## Types of Investments

### Stocks
Ownership in companies. Higher risk, higher potential return.

### Bonds
Loans to companies or governments. Lower risk, lower return.

### Mutual Funds
Baskets of stocks and bonds managed by professionals.

### Index Funds
Funds that track market indexes. Low fees, diversified.

### ETFs
Exchange-traded funds. Similar to index funds, trade like stocks.

## Investment Accounts

### 401(k)
Employer retirement account. Often includes employer match.

### IRA (Individual Retirement Account)
Personal retirement account. Tax advantages.

### Roth IRA
After-tax contributions, tax-free withdrawals in retirement.

### Taxable Brokerage Account
No tax advantages, but no restrictions on withdrawals.

## Investment Principles

1. **Start Early**: Time is your biggest advantage
2. **Diversify**: Don''t put all eggs in one basket
3. **Think Long-term**: Ignore short-term volatility
4. **Keep Costs Low**: Fees eat into returns
5. **Stay Consistent**: Invest regularly

## Risk vs. Return

Higher potential returns come with higher risk. Your risk tolerance depends on:
- Age
- Financial goals
- Time horizon
- Personal comfort level

## Getting Started

1. Pay off high-interest debt first
2. Build emergency fund
3. Take advantage of employer 401(k) match
4. Open an IRA
5. Start with index funds',
     'https://www.youtube.com/watch?v=example7',
     '{"articles": ["https://www.investopedia.com/investing-4427685"], "tools": ["Investment calculator", "Retirement calculator"]}',
     7),

    ('l0000008-0000-0000-0000-000000000008',
     'Financial Goal Setting',
     'Learn how to set SMART financial goals and create action plans to achieve them.',
     'Planning',
     25,
     'advanced',
     '# Financial Goal Setting

## SMART Goals

Make goals:
- **Specific**: Clear and well-defined
- **Measurable**: Track progress with numbers
- **Achievable**: Realistic given your situation
- **Relevant**: Aligned with your values
- **Time-bound**: Have a deadline

## Types of Financial Goals

### Short-term (< 1 year)
- Build $1,000 emergency fund
- Pay off credit card
- Save for vacation

### Medium-term (1-5 years)
- Save for car down payment
- Build 6-month emergency fund
- Pay off student loans

### Long-term (5+ years)
- Save for house down payment
- Build retirement savings
- Fund children''s education

## Creating Your Action Plan

1. **Define Your Goal**: Be specific about what you want
2. **Calculate the Cost**: How much money do you need?
3. **Set a Timeline**: When do you want to achieve it?
4. **Break It Down**: Monthly or weekly savings needed
5. **Identify Obstacles**: What might get in the way?
6. **Create Solutions**: How will you overcome obstacles?
7. **Track Progress**: Regular check-ins

## Prioritizing Goals

Use this framework:
1. **Essential**: Emergency fund, debt payoff
2. **Important**: Retirement, major purchases
3. **Nice to Have**: Vacations, upgrades

## Staying Motivated

- Visualize your goal
- Track progress visually
- Celebrate milestones
- Share goals with accountability partner
- Review and adjust regularly

## When Life Changes

Revisit goals when:
- Income changes
- Family situation changes
- Unexpected expenses arise
- Goals are achieved',
     'https://www.youtube.com/watch?v=example8',
     '{"articles": ["https://www.thebalancemoney.com/financial-goals-453822"], "templates": ["Goal planning worksheet", "Progress tracker"]}',
     8);

-- Sample education progress for users
INSERT INTO education_progress (id, user_id, lesson_id, is_completed, completed_at)
VALUES
    -- John Doe (User 1) - Completed 5 lessons
    (gen_random_uuid(), '11111111-1111-1111-1111-111111111111', 'l0000001-0000-0000-0000-000000000001', true, NOW() - INTERVAL '80 days'),
    (gen_random_uuid(), '11111111-1111-1111-1111-111111111111', 'l0000002-0000-0000-0000-000000000002', true, NOW() - INTERVAL '75 days'),
    (gen_random_uuid(), '11111111-1111-1111-1111-111111111111', 'l0000003-0000-0000-0000-000000000003', true, NOW() - INTERVAL '70 days'),
    (gen_random_uuid(), '11111111-1111-1111-1111-111111111111', 'l0000004-0000-0000-0000-000000000004', true, NOW() - INTERVAL '50 days'),
    (gen_random_uuid(), '11111111-1111-1111-1111-111111111111', 'l0000005-0000-0000-0000-000000000005', true, NOW() - INTERVAL '30 days'),
    
    -- Jane Smith (User 2) - Completed 3 lessons
    (gen_random_uuid(), '22222222-2222-2222-2222-222222222222', 'l0000001-0000-0000-0000-000000000001', true, NOW() - INTERVAL '55 days'),
    (gen_random_uuid(), '22222222-2222-2222-2222-222222222222', 'l0000002-0000-0000-0000-000000000002', true, NOW() - INTERVAL '45 days'),
    (gen_random_uuid(), '22222222-2222-2222-2222-222222222222', 'l0000003-0000-0000-0000-000000000003', true, NOW() - INTERVAL '35 days'),
    
    -- Mike Johnson (User 3) - Completed 1 lesson
    (gen_random_uuid(), '33333333-3333-3333-3333-333333333333', 'l0000001-0000-0000-0000-000000000001', true, NOW() - INTERVAL '40 days'),
    
    -- Sarah Williams (User 4) - Completed all 8 lessons
    (gen_random_uuid(), '44444444-4444-4444-4444-444444444444', 'l0000001-0000-0000-0000-000000000001', true, NOW() - INTERVAL '65 days'),
    (gen_random_uuid(), '44444444-4444-4444-4444-444444444444', 'l0000002-0000-0000-0000-000000000002', true, NOW() - INTERVAL '60 days'),
    (gen_random_uuid(), '44444444-4444-4444-4444-444444444444', 'l0000003-0000-0000-0000-000000000003', true, NOW() - INTERVAL '55 days'),
    (gen_random_uuid(), '44444444-4444-4444-4444-444444444444', 'l0000004-0000-0000-0000-000000000004', true, NOW() - INTERVAL '45 days'),
    (gen_random_uuid(), '44444444-4444-4444-4444-444444444444', 'l0000005-0000-0000-0000-000000000005', true, NOW() - INTERVAL '35 days'),
    (gen_random_uuid(), '44444444-4444-4444-4444-444444444444', 'l0000006-0000-0000-0000-000000000006', true, NOW() - INTERVAL '25 days'),
    (gen_random_uuid(), '44444444-4444-4444-4444-444444444444', 'l0000007-0000-0000-0000-000000000007', true, NOW() - INTERVAL '15 days'),
    (gen_random_uuid(), '44444444-4444-4444-4444-444444444444', 'l0000008-0000-0000-0000-000000000008', true, NOW() - INTERVAL '5 days'),
    
    -- Alex Brown (User 5) - Completed 2 lessons
    (gen_random_uuid(), '55555555-5555-5555-5555-555555555555', 'l0000001-0000-0000-0000-000000000001', true, NOW() - INTERVAL '12 days'),
    (gen_random_uuid(), '55555555-5555-5555-5555-555555555555', 'l0000002-0000-0000-0000-000000000002', true, NOW() - INTERVAL '8 days');
