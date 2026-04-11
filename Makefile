.PHONY: build build-alpine run dev swag test lint migrate-up migrate-down tools help

APP_NAME = boilerplate
BINARY   = bin/$(APP_NAME)

LDFLAGS = -ldflags "\
	-X github.com/sule/go-boilerplate/version.GitCommit=$(shell git rev-parse --short HEAD 2>/dev/null || echo 'unknown') \
	-X github.com/sule/go-boilerplate/version.BuildDate=$(shell date -u +%Y-%m-%dT%H:%M:%SZ) \
	-X github.com/sule/go-boilerplate/version.Version=0.1.0"

build: ## Build the binary
	go build $(LDFLAGS) -o $(BINARY) ./main.go

build-alpine: ## Build a static binary for Linux/Alpine
	CGO_ENABLED=0 GOOS=linux go build $(LDFLAGS) -a -installsuffix cgo -o $(BINARY) ./main.go

run: ## Run the server
	go run main.go serve

dev: ## Run with air live-reload
	air

swag: ## Generate Swagger docs
	swag init -g main.go -o docs

test: ## Run tests
	go test ./...

lint: ## Run linter
	golangci-lint run

migrate-up: ## Apply all pending migrations
	go run main.go db migrate-up

migrate-down: ## Roll back the last migration
	go run main.go db migrate-down

tools: ## Install development tools
	go install github.com/air-verse/air@latest
	go install github.com/swaggo/swag/cmd/swag@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
