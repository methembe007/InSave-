.PHONY: build run test clean

build:
	go build -o bin/goal-service cmd/server/main.go

run:
	go run cmd/server/main.go

test:
	go test -v ./...

clean:
	rm -rf bin/

deps:
	go mod download
	go mod tidy

docker-build:
	docker build -t goal-service:latest .

docker-run:
	docker run -p 8005:8005 --env-file .env goal-service:latest
