.PHONY: help build test run clean migrate-up migrate-down migrate-test seed docker-up docker-down docker-logs test-integration test-coverage test-unit lint fmt deps swagger test-env-up test-env-down check-go-version fix-go-version

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

build: ## Build the application
	go build -o bin/api ./cmd/api

test-unit: ## Run unit tests
	go test -v -race -coverprofile=coverage.out ./... -short

test-integration: test-env-up ## Run integration tests
	@echo "Running integration tests..."
	@sleep 5 # Wait for containers to be ready
	go test -v -tags=integration ./tests/integration/... -timeout 30s

test: test-unit test-integration ## Run all tests

test-coverage: test-unit ## Generate coverage report
	go tool cover -html=coverage.out -o coverage.html

test-env-up: ## Start test environment (PostgreSQL + Redis)
	docker compose -f ../../docker-compose.test.yml up -d
	@echo "Waiting for test containers to be ready..."
	@sleep 10

test-env-down: ## Stop test environment
	docker compose -f ../../docker-compose.test.yml down -v

run: ## Run the application
	go run ./cmd/api

docker-up: ## Start Docker services
	docker compose up -d

docker-down: ## Stop Docker services
	docker compose down

docker-logs: ## Show Docker logs
	docker compose logs -f

migrate-up: ## Run database migrations
	docker compose run --rm migrate up

migrate-down: ## Rollback database migrations
	docker compose run --rm migrate down

migrate-test: ## Run migrations for test environment (handled by testcontainers)
	@echo "✅ Test migrations are handled by testcontainers in integration tests"

seed: ## Seed database with test data
	bash scripts/seed.sh

clean: ## Clean build artifacts
	rm -rf bin/ coverage.out coverage.html

lint: ## Run linter (requires golangci-lint)
	golangci-lint run

fmt: ## Format code
	go fmt ./...

deps: ## Download dependencies
	go mod download
	go mod tidy

swagger: ## Generate Swagger docs (requires swag)
	swag init -g cmd/api/main.go -o docs

check-go-version: ## Check if go.mod has correct Go version (1.22)
	@if grep -q "^go 1.24" go.mod; then \
		echo "❌ ERROR: go.mod has Go 1.24.0 which is not released yet!"; \
		echo ""; \
		echo "Dockerfile uses golang:1.22-alpine"; \
		echo "Please change go.mod to: go 1.22"; \
		echo ""; \
		echo "Run: make fix-go-version"; \
		exit 1; \
	else \
		echo "✅ go.mod Go version is correct (1.22)"; \
	fi

fix-go-version: ## Fix go.mod to use Go 1.22
	@echo "Fixing go.mod to use Go 1.22..."
	@go mod edit -go=1.22
	@echo "✅ Fixed! go.mod now uses Go 1.22"
	@echo ""
	@grep "^go " go.mod

.DEFAULT_GOAL := help


