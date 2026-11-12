# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project

`clean-go-fiber` is a Go 1.26 REST API starter: Fiber v2 (fasthttp), GORM over PostgreSQL with
hand-written golang-migrate SQL migrations (embedded, applied on server boot), Redis for the
`UserInformation` cache (go-redis) and the job queue (asynq), JWT bearer auth (golang-jwt v5, access
1d + refresh 7d), go-playground/validator, zap logging, and swaggo OpenAPI. It ships auth (login,
refresh, register, email verification, forgot/reset password), a self-service profile, and settings
for users, roles, and permissions with RBAC. Four binaries: `cmd/server`, `cmd/worker` (sends mail),
`cmd/seed`, `cmd/migrate`. Its API contract (routes, permission strings, envelope, list query, env
names) follows the sibling `clean-elysia`.

## Commands

Go toolchain only; `make` is the entry point. `make help` lists every target.

```bash
# Development
make dev              # server (air hot reload) + worker
make dev-server       # server only (air)
make dev-worker       # worker only
make build            # go build ./cmd/server and ./cmd/worker
make start            # run the built server

# Quality
make format           # gofmt -w .
make lint             # gofmt -l (must be empty) + go vet ./...
make docs             # regenerate docs/swagger.* with swag

# Tests
make test             # go test ./... — there are no _test.go files yet

# Database
make db-migrate                      # apply pending migrations (go run ./cmd/migrate)
make db-migrate-down                 # roll all migrations back
make db-migrate-create NAME=<name>   # new numbered .up.sql/.down.sql pair
make db-seed                         # idempotent seed (go run ./cmd/seed)
make fresh                           # DESTRUCTIVE: down + migrate + seed
make reset                           # migrate + seed

# Deployment (Docker)
make docker-build docker-up docker-down docker-restart docker-logs docker-ps
make docker-migrate docker-seed      # run the migrator stage against the compose network
```

Seeding is idempotent: a second `make db-seed` changes nothing and exits 0. It creates 15
permissions, the `superuser` and `admin` roles (both with the full catalogue), and
`superuser@example.com` / `admin@example.com` (password `S3crEtP4ssw0rd!`, pre-verified).

## Architecture

| Package | Purpose |
| --- | --- |
| `cmd/{server,worker,seed,migrate}` | entrypoints |
| `internal/bootstrap` | constructor wiring and route registration |
| `internal/config` | env loading; the only reader of the environment (`config.C`) |
| `internal/database/postgres` | GORM connection, embedded migrations, `seed/` |
| `internal/database/redis` | Redis client (pinged on connect) |
| `internal/errors` | `AppError` — the client-facing error vocabulary |
| `internal/guards` | `guards.Permission`, `guards.Role` |
| `internal/middleware` | `Auth`, request logger, slow-request logger |
| `internal/models` | GORM models mirroring the SQL |
| `internal/modules/<name>` | `handler.go`, `routes.go`, `schema.go`, `service.go` per feature |
| `internal/repositories` | interfaces + GORM implementations |
| `internal/libs/*` | cache, mailer, pagination, utils (response, hash, JWT, tokens), validator, txmanager |
| `internal/queue`, `internal/worker` | asynq client/tasks and server/handlers |
| `internal/logging` | zap logger built from `APP_ENV` / `LOG_LEVEL` |

### Request flow

`routes.go` (group + `auth` + guard) → `handler.go` (parse, validate, one service call,
respond) → `service.go` (rules, existence/uniqueness checks, transactions, cache invalidation,
enqueue) → `internal/repositories` (queries only, `WithTx` for transactions) → `internal/models`.
The service opens every transaction; repositories never do.

### The non-obvious bits

- **Repositories return `(nil, nil)` for not found.** The service turns it into 404/422.
- **Auth is default-open.** `middleware.Auth(userRepo, cache)` is built once in bootstrap and passed
  to each module's `RegisterRoutes(router, auth)` for its group; a
  route without a guard is open to every logged-in user. Guards run per route.
- **`superuser` bypasses both guards.** `roles.name` is UNIQUE and the role cannot be renamed or
  deleted. A non-superuser cannot assign or remove it, nor modify, delete, or reset the password of a
  user who holds it (403).
- **Handlers use `utils.ParseBody`, `utils.ParamUUID`, `pagination.Parse`, and
  `middleware.CurrentUser`**; account mails (verification, reset) go through `AccountMail` in
  `internal/modules/auth/account_mail.go`, shared by the auth and users services.
- **Access vs refresh tokens** are signed with different secrets and carry a `type` claim;
  `middleware.Auth()` rejects a refresh token, `/auth/refresh` rejects an access token. The current
  refresh `jti` lives in Redis at `refresh:<id>`; issuing a pair rotates it, and
  `cache.RevokeSessions` (password change/reset, admin reset, delete, non-active status) drops it
  together with `user:<id>`.
- **`UserInformation` is cached** in Redis at `user:<id>` for 3600 s; the DB lookup requires
  `status = 'active' AND deleted_at IS NULL`. Every write that changes a user's status, roles,
  identity, password, or deletion — and every role/permission change, for all holders — drops the
  key.
- **Select-option lists are cached** in Redis at `select_options:permissions` and
  `select_options:roles` for 3600 s; every role or permission create, update, and delete drops both
  (`cache.ForgetSelectOptions`).
- **Pool and slow-query log are env-driven** (`DB_POOL_*`, `DB_SLOW_QUERY_MS`). GORM logs slow
  statements and errors at warn in every environment, with parameters as placeholders; every
  statement only in development with `LOG_LEVEL=debug`.
- **Responses go through `internal/libs/utils/response.go`**; failures are returned as
  `*errors.AppError` and the central `ErrorHandler` builds the envelope. Unknown errors become a
  logged 500 with a fixed message — never `err.Error()`.
- **Env is read only in `internal/config`.** `JWT_SECRET` and `JWT_REFRESH_SECRET` are required and
  must differ; everything else has a default.
- **`/docs` and `/docs/openapi.json` exist only when `ENABLE_API_DOCS=true`** (default `false`).
- **Mail is queued.** Services enqueue after commit; only the worker talks SMTP.
- **One-time tokens are stored as SHA-256 hashes**, single use (`used_at`), 1 h, revoked on reissue.
- **Users are soft-deleted** (`gorm.DeletedAt`); roles and permissions are hard-deleted.
- **Migrations run on server boot.** Never edit an applied migration.

## Conventions

- gofmt (tabs), `go vet` clean; no `any` in signatures where a type exists; always handle errors;
  zap only, never `fmt.Print*` in request paths; no comments unless one block comment is needed.
- Permission strings are `"<group> <action>"`, space-separated (`"user list"`).
- `PATCH`, never `PUT`; list/create paths register `""`; deletes return 200 with `data: null`.
- Uniqueness conflicts are **422**; login and token failures are 400; malformed `:id` is 422.
- Response fields camelCase; status values lowercase (`active`, `inactive`, `suspended`, `blocked`).
- No i18n catalogue — English messages ([.claude/rules/i18n.md](.claude/rules/i18n.md)).
- Global rate limit 100 req / 60 s per IP, hardcoded.
- Flat routes, no version prefix; settings under `/settings`.
- Always in scope: [`contradiction-halt.md`](.claude/rules/contradiction-halt.md) and
  [`documentation.md`](.claude/rules/documentation.md).

## Adding a new feature module

1. Write the migration pair (`make db-migrate-create NAME=create_<name>_table`) and the model in
   `internal/models/`.
2. Add the repository interface to `internal/repositories/interfaces.go` and implement it in
   `internal/repositories/<name>_repository.go` (with `WithTx` and allow-lists if it lists).
3. Create `internal/modules/<name>/schema.go` (request/response structs, validator tags, camelCase
   json).
4. Create `service.go` (interface + `NewService`, owns checks, transactions, cache invalidation).
5. Create `handler.go` with full swag annotations, and `routes.go` with guards.
6. If it needs new permissions, extend the seed in the same change.
7. Wire it in `internal/bootstrap/bootstrap.go`.
8. Update `.claude/rules/routes.md`, `docs/API_DOCUMENTATION.md`, `README.md`; run `make docs`.

## Deeper documentation

- `docs/README.md` — index
- `docs/API_DOCUMENTATION.md` — envelope, auth, route map, list query
- `docs/CONFIGURATION.md` — every env var
- `docs/DEPLOYMENT.md` — Docker and compose
- `docs/ERROR_HANDLING.md` — `AppError`, status mapping
- `docs/SECURITY.md` — tokens, guards, limits, headers
- `docs/MIDDLEWARE.md` — the middleware chain and guards

## Skills, rules, and commands

- `.claude/skills` is a symlink to `.agents/skills/` (vendored, no lock file): the universal
  workflow/craft bundle plus `go-fiber` and `golang-pro`. See `.agents/skills/README.md`.
- `.claude/rules/` holds 24 rule files indexed by `.claude/rules/README.md`.
- Commands: `/commit`, `/audit-flow` (writes `docs/audit-findings.md`), `/update-todo` (reconciles
  `todo.md`).
- `.gemini/GEMINI.md` and `.github/copilot-instructions.md` point other agents at this file and the
  rules.
