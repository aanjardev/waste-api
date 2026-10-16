.PHONY: up down build logs restart psql tidy run test migrate-up migrate-down seed

up:
	docker compose up --build

down:
	docker compose down

build:
	docker compose build

logs:
	docker compose logs -f app

restart:
	docker compose restart app

psql:
	docker compose exec db psql -U postgres -d waste_db

tidy:
	go mod tidy

run:
	go run ./cmd/api

test:
	go test ./... -v

migrate-up:
	docker compose exec -T db psql -U postgres -d waste_db < migrations/000001_init.up.sql
	docker compose exec -T db psql -U postgres -d waste_db < migrations/000002_seed.up.sql

migrate-down:
	docker compose exec -T db psql -U postgres -d waste_db < migrations/000001_init.down.sql