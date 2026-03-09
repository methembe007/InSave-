.PHONY: build run test clean docker-build docker-run

# Build the application
build:
	go build -o bin/auth-service cmd/server/main.go

# Run the application
run:
	go run cmd/server/main.go

# Run tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -cover -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html

# Download dependencies
deps:
	go mod download
	go mod tidy

# Build Docker image
docker-build:
	docker build -t insavein/auth-service:latest .

# Run Docker container
docker-run:
	docker run -p 8080:8080 \
		-e DB_HOST=localhost \
		-e DB_PORT=5432 \
		-e DB_USER=insavein_user \
		-e DB_PASSWORD=insavein_password \
		-e DB_NAME=insavein_db \
		-e JWT_SECRET=your-secret-key \
		insavein/auth-service:latest

# Format code
fmt:
	go fmt ./...

# Lint code
lint:
	golangci-lint run

# Run the service with local database
dev:
	DB_HOST=localhost \
	DB_PORT=5432 \
	DB_USER=insavein_user \
	DB_PASSWORD=insavein_password \
	DB_NAME=insavein_db \
	JWT_SECRET=dev-secret-key-change-in-production \
	go run cmd/server/main.go
