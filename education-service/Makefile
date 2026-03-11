.PHONY: build run test clean docker-build docker-run deps

# Build the application
build:
	go build -o bin/education-service cmd/server/main.go

# Run the application
run:
	go run cmd/server/main.go

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -rf bin/

# Install dependencies
deps:
	go mod download
	go mod tidy

# Build Docker image
docker-build:
	docker build -t education-service:latest .

# Run Docker container
docker-run:
	docker run -p 8085:8085 --env-file .env education-service:latest

# Run with hot reload (requires air)
dev:
	air
