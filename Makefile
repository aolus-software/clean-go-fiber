dev:
	$(MAKE) dev-server & $(MAKE) dev-worker

dev-server:
	air

dev-worker:
	go run ./cmd/worker/main.go

# ---------------------------------------------------------------------------
# Database migrations (requires `migrate` CLI: go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest)
# ---------------------------------------------------------------------------
DB_USER     := $(shell grep '^DB_USER='     .env 2>/dev/null | cut -d= -f2)
DB_PASSWORD := $(shell grep '^DB_PASSWORD=' .env 2>/dev/null | cut -d= -f2)
DB_HOST     := $(shell grep '^DB_HOST='     .env 2>/dev/null | cut -d= -f2)
DB_PORT     := $(shell grep '^DB_PORT='     .env 2>/dev/null | cut -d= -f2)
DB_NAME     := $(shell grep '^DB_NAME='     .env 2>/dev/null | cut -d= -f2)
DB_SSLMODE  := $(shell grep '^DB_SSLMODE='  .env 2>/dev/null | cut -d= -f2)

MIGRATE_URL := postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)
MIGRATE_DIR := internal/database/postgres/migrations

# Create a new migration file pair.
# Usage: make migrate-create NAME=add_users_table
migrate-create:
	@test -n "$(NAME)" || (echo "Error: NAME is required. Usage: make migrate-create NAME=<migration_name>" && exit 1)
	migrate create -ext sql -dir $(MIGRATE_DIR) -seq $(NAME)

# Apply all pending migrations.
migrate-up:
	migrate -path $(MIGRATE_DIR) -database "$(MIGRATE_URL)" up

# Roll back the last N migrations (default: 1).
# Usage: make migrate-down        (rolls back 1)
#        make migrate-down N=3    (rolls back 3)
migrate-down:
	migrate -path $(MIGRATE_DIR) -database "$(MIGRATE_URL)" down $(if $(N),$(N),1)

# Show the current migration version.
migrate-status:
	migrate -path $(MIGRATE_DIR) -database "$(MIGRATE_URL)" version

# Force-set the migration version (use after manually fixing a dirty state).
# Usage: make migrate-force VERSION=1
migrate-force:
	@test -n "$(VERSION)" || (echo "Error: VERSION is required. Usage: make migrate-force VERSION=<n>" && exit 1)
	migrate -path $(MIGRATE_DIR) -database "$(MIGRATE_URL)" force $(VERSION)

.PHONY: dev dev-server dev-worker migrate-create migrate-up migrate-down migrate-status migrate-force