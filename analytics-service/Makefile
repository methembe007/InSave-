.PHONY: build run test clean docker-build docker-run

# Build the application
build:
	go build -o bin/analytics-service cmd/server/main.go

# Run the application
run:
	go run cmd/server/main.go

# Run tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -v -cover ./...

# Clean build artifacts
clean:
	rm -rf bin/

# Download dependencies
deps:
	go mod download
	go mod tidy

# Build Docker image
docker-build:
	docker build -t analytics-service:latest .

# Run Docker container
docker-run:
	docker run -p 8008:8008 \
		-e DB_HOST=host.docker.internal \
		-e DB_PASSWORD=postgres \
		-e JWT_SECRET=your-secret-key \
		analytics-service:latest

# Format code
fmt:
	go fmt ./...

# Run linter
lint:
	golangci-lint run

# Run the service with environment variables
dev:
	export DB_HOST=localhost && \
	export DB_PORT=5432 && \
	export DB_USER=postgres && \
	export DB_PASSWORD=postgres && \
	export DB_NAME=insavein && \
	export JWT_SECRET=your-secret-key && \
	export PORT=8008 && \
	go run cmd/server/main.go
