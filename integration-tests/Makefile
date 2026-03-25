.PHONY: help setup start stop test test-coverage test-user test-savings test-budget test-goal clean logs

help: ## Show this help message
	@echo "InSavein Integration Tests - Available Commands:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

setup: ## Install Go dependencies
	@echo "Installing Go dependencies..."
	go mod download
	@echo "Dependencies installed!"

start: ## Start test environment
	@echo "Starting test environment..."
	docker-compose -f docker-compose.test.yml up -d
	@echo "Waiting for services to be ready..."
	@sleep 30
	@echo "Test environment is ready!"

stop: ## Stop test environment
	@echo "Stopping test environment..."
	docker-compose -f docker-compose.test.yml down
	@echo "Test environment stopped!"

clean: ## Stop test environment and remove volumes
	@echo "Cleaning up test environment..."
	docker-compose -f docker-compose.test.yml down -v
	@echo "Cleanup complete!"

test: ## Run all integration tests
	@echo "Running all integration tests..."
	go test -v ./...

test-coverage: ## Run tests with coverage report
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

test-user: ## Run user registration flow tests
	@echo "Running user registration tests..."
	go test -v -run TestUserRegistration ./...

test-savings: ## Run savings flow tests
	@echo "Running savings flow tests..."
	go test -v -run TestSavings ./...

test-budget: ## Run budget alert flow tests
	@echo "Running budget alert tests..."
	go test -v -run TestBudget ./...

test-goal: ## Run goal progress flow tests
	@echo "Running goal progress tests..."
	go test -v -run TestGoal ./...

logs: ## Show logs from all test services
	docker-compose -f docker-compose.test.yml logs -f

logs-auth: ## Show auth service logs
	docker-compose -f docker-compose.test.yml logs -f auth-service-test

logs-user: ## Show user service logs
	docker-compose -f docker-compose.test.yml logs -f user-service-test

logs-savings: ## Show savings service logs
	docker-compose -f docker-compose.test.yml logs -f savings-service-test

logs-budget: ## Show budget service logs
	docker-compose -f docker-compose.test.yml logs -f budget-service-test

logs-goal: ## Show goal service logs
	docker-compose -f docker-compose.test.yml logs -f goal-service-test

logs-db: ## Show database logs
	docker-compose -f docker-compose.test.yml logs -f postgres-test

ps: ## Show status of test services
	docker-compose -f docker-compose.test.yml ps

restart: stop start ## Restart test environment

full-test: start test clean ## Run complete test cycle (start, test, cleanup)
