# Integration Tests - Quick Start Guide

## 🚀 Run Tests in 3 Steps

### 1. Navigate to integration tests directory
```bash
cd integration-tests
```

### 2. Start test environment
```bash
# Linux/Mac
./run-tests.sh

# Windows
run-tests.bat

# Or using Make
make full-test
```

### 3. That's it! ✅
The script will:
- Start all services in Docker
- Wait for services to be healthy
- Run all integration tests
- Show results

## 📋 Common Commands

### Run All Tests
```bash
make test
```

### Run Specific Test Suite
```bash
make test-user      # User registration tests
make test-savings   # Savings flow tests
make test-budget    # Budget alert tests
make test-goal      # Goal progress tests
```

### View Logs
```bash
make logs           # All services
make logs-auth      # Auth service only
make logs-savings   # Savings service only
```

### Generate Coverage Report
```bash
make test-coverage
```

### Cleanup
```bash
make clean          # Stop and remove everything
```

## 🔍 Troubleshooting

### Services not starting?
```bash
# Check Docker is running
docker info

# View service status
make ps

# View logs
make logs
```

### Tests failing?
```bash
# Restart environment
make restart

# Run tests individually
go test -v -run TestUserRegistration
go test -v -run TestSavings
```

### Need fresh start?
```bash
# Complete cleanup and restart
make clean
make start
make test
```

## 📊 Test Coverage

- **26 test cases** across 4 major workflows
- **13 requirements** validated
- **100%** critical path coverage

## 🎯 What Gets Tested

✅ User registration and authentication
✅ Profile and preferences management
✅ Savings transactions and streaks
✅ Budget alerts and thresholds
✅ Goal progress and milestones
✅ Cross-service communication
✅ Data consistency
✅ Input validation
✅ Error handling

## 📚 More Information

- See `README.md` for detailed documentation
- See `IMPLEMENTATION_SUMMARY.md` for technical details
- See `TASK_27_INTEGRATION_TESTING_COMPLETE.md` for completion status

## 💡 Tips

- Tests run in isolated environment (won't affect dev/prod)
- Test data is automatically seeded
- Each test run starts with clean state
- Services run on different ports (18080-18086)
- Test database on port 5435

## 🆘 Need Help?

Check the logs:
```bash
make logs
```

View service status:
```bash
make ps
```

Read the full documentation:
```bash
cat README.md
```
