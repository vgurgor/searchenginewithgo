.PHONY: up down logs api-logs fe-logs migrate seed build

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


