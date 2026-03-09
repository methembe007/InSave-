.PHONY: build run test clean deps

# Build the application
build:
	go build -o user-service cmd/server/main.go

# Run the application
run:
	go run cmd/server/main.go

# Run tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Install dependencies
deps:
	go mod download
	go mod tidy

# Clean build artifacts
clean:
	rm -f user-service user-service.exe coverage.out coverage.html

# Format code
fmt:
	go fmt ./...

# Run linter
lint:
	golangci-lint run

# Build Docker image
docker-build:
	docker build -t insavein/user-service:latest .

# Run with Docker Compose
docker-up:
	docker-compose up -d

# Stop Docker Compose
docker-down:
	docker-compose down
