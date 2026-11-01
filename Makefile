MIGRATE_DIR := internal/database/postgres/migrations
COMPOSE     := docker compose

.PHONY: help dev dev-server dev-worker build start lint format test docs \
	db-migrate db-migrate-down db-migrate-create db-seed fresh reset \
	docker-build docker-up docker-down docker-restart docker-logs docker-ps docker-migrate docker-seed

help: ## List every target
	@grep -E '^[a-zA-Z_-]+:.*?## ' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'

dev: ## Run the server (air hot reload) and the worker
	$(MAKE) -j2 dev-server dev-worker

dev-server: ## Run the server with air hot reload
	air

dev-worker: ## Run the asynq worker
	go run ./cmd/worker

build: ## Build the server, worker, migrate and seed binaries into bin/
	go build -o bin/server ./cmd/server
	go build -o bin/worker ./cmd/worker
	go build -o bin/migrate ./cmd/migrate
	go build -o bin/seed ./cmd/seed

start: ## Run the built server
	./bin/server

lint: ## gofmt check (must print nothing) and go vet
	@test -z "$$(gofmt -l .)" || (echo "gofmt needed on:"; gofmt -l .; exit 1)
	go vet ./...

format: ## Format every Go file with gofmt
	gofmt -w .

test: ## Run go test (no _test.go files exist yet)
	go test ./...

docs: ## Regenerate docs/swagger.* (needs the swag CLI)
	swag init -d "$$(go list -f '{{.Dir}}' ./cmd/server ./internal/... | paste -sd, -)" -g main.go -o docs

db-migrate: ## Apply pending migrations
	go run ./cmd/migrate up

db-migrate-down: ## Roll back every migration (destructive)
	go run ./cmd/migrate down

db-migrate-create: ## Create a migration pair: make db-migrate-create NAME=<name>
	@test -n "$(NAME)" || (echo "Usage: make db-migrate-create NAME=<name>" && exit 1)
	@last=$$(ls $(MIGRATE_DIR)/*.up.sql 2>/dev/null | sed 's#.*/##' | cut -d_ -f1 | sort -n | tail -1); \
	next=$$(printf "%06d" $$((10#$${last:-0} + 1))); \
	touch $(MIGRATE_DIR)/$${next}_$(NAME).up.sql $(MIGRATE_DIR)/$${next}_$(NAME).down.sql; \
	echo "Created $(MIGRATE_DIR)/$${next}_$(NAME).{up,down}.sql"

db-seed: ## Seed permissions, roles and users (idempotent)
	go run ./cmd/seed

fresh: ## DESTRUCTIVE: roll back everything, migrate, seed
	go run ./cmd/migrate down
	go run ./cmd/migrate up
	go run ./cmd/seed

reset: ## Migrate and seed
	go run ./cmd/migrate up
	go run ./cmd/seed

docker-build: ## Build the compose images
	$(COMPOSE) build

docker-up: ## Start every compose service
	$(COMPOSE) up -d

docker-down: ## Stop every compose service
	$(COMPOSE) down

docker-restart: ## Restart every compose service
	$(COMPOSE) restart

docker-logs: ## Follow compose logs
	$(COMPOSE) logs -f

docker-ps: ## List compose services
	$(COMPOSE) ps

docker-migrate: ## Run migrations from the migrator stage
	$(COMPOSE) run --rm migrator ./migrate up

docker-seed: ## Seed from the migrator stage
	$(COMPOSE) run --rm migrator ./seed
