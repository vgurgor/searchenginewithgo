.PHONY: help build test run clean migrate-up migrate-down seed docker-up docker-down docker-logs test-integration test-coverage lint fmt deps swagger

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

build: ## Build the application
	go build -o bin/api ./cmd/api

test: ## Run unit tests
	go test -v -race -coverprofile=coverage.out ./...

test-integration: ## Run integration tests
	go test -v -tags=integration ./tests/integration/...

test-coverage: test ## Generate coverage report
	go tool cover -html=coverage.out -o coverage.html

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

.DEFAULT_GOAL := help


