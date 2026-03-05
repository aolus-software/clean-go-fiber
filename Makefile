dev:
	go run ./cmd/server/main.go
	go run ./cmd/worker/main.go

dev-server:
	go run ./cmd/server/main.go

dev-worker:
	go run ./cmd/worker/main.go

# Create a new migration file pair.
# Usage: make migrate-create NAME=add_users_table
migrate-create:
	@test -n "$(NAME)" || (echo "Error: NAME is required. Usage: make migrate-create NAME=<migration_name>" && exit 1)
	migrate create -ext sql -dir internal/database/postgres/migrations -seq $(NAME)

.PHONY: dev dev-server dev-worker migrate-create