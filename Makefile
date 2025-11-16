.PHONY: up down logs api-logs fe-logs migrate seed build install-hooks check-go-version fix-go-version

up:
	docker compose up --build -d

down:
	docker compose down

logs:
	docker compose logs -f

api-logs:
	docker compose logs -f api

fe-logs:
	docker compose logs -f frontend

migrate:
	bash backend/scripts/migrate.sh

seed:
	bash backend/scripts/seed.sh

build:
	docker compose build

install-hooks: ## Install git hooks
	@echo "📥 Installing git hooks..."
	@cp .github/hooks/pre-commit .git/hooks/pre-commit
	@chmod +x .git/hooks/pre-commit
	@echo "✅ Git hooks installed!"
	@echo "   - pre-commit: Checks Go version"

check-go-version: ## Check backend/go.mod Go version
	@cd backend && make check-go-version

fix-go-version: ## Fix backend/go.mod to use Go 1.22
	@cd backend && make fix-go-version


