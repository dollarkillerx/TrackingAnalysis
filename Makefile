.PHONY: dev build docker up down

# Run server locally
dev:
	go run cmd/server/main.go

# Build Go binary for linux amd64
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o bin/tracking-server ./cmd/server

# Build Docker image
docker:
	@echo "Building Docker image..."
	@docker build -t tracking-server:latest .

# Start all services
up:
	docker compose up -d

# Stop all services
down:
	docker compose down
