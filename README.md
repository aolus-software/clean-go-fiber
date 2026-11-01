# clean-go-fiber

A Go REST API starter following clean architecture: Fiber v2, GORM + golang-migrate, PostgreSQL,
Redis (cache + asynq queue), JWT authentication with refresh tokens, RBAC, and swaggo OpenAPI docs.
Its API contract — routes, permission strings, response envelope, list query, env names — matches
the sibling TypeScript starters (`clean-elysia`).

## Tech stack

| Layer | Technology |
| --- | --- |
| HTTP framework | Fiber v2 |
| ORM / migrations | GORM / golang-migrate (SQL pairs, embedded) |
| Database | PostgreSQL 18 |
| Cache / queue | Redis (go-redis v9 + asynq) |
| Auth | golang-jwt v5 — access 1 day, refresh 7 days, separate secrets |
| Validation | go-playground/validator v10 |
| Logging | zap |
| API docs | swaggo/swag + Scalar UI (gated by `ENABLE_API_DOCS`) |
| Mail | gomail (SMTP), sent by the worker only |
| Hot reload | air |

## Prerequisites

- Go 1.26+
- Docker and docker compose (for Postgres and Redis)
- Optional: [air](https://github.com/air-verse/air) for `make dev`, and
  [swag](https://github.com/swaggo/swag) (`go install github.com/swaggo/swag/cmd/swag@v1.16.6`) for
  `make docs`

## Getting started

```bash
git clone https://github.com/aolus-software/clean-go-fiber.git
cd clean-go-fiber

cp .env.example .env          # set JWT_SECRET and JWT_REFRESH_SECRET (openssl rand -hex 32, twice)
docker compose up -d postgres redis
go mod download

make db-migrate               # migrations also run automatically when the server starts
make db-seed                  # idempotent
make dev                      # server (air) + worker
```

The API listens on `http://localhost:3000`. With `ENABLE_API_DOCS=true` (the `.env.example` value)
the docs are at `/docs` and the spec at `/docs/openapi.json`.

Seeded accounts (password `S3crEtP4ssw0rd!`, already verified): `superuser@example.com` (role
`superuser`, bypasses every guard) and `admin@example.com` (role `admin`, all 15 permissions). Change
or remove them before any shared deployment.

## Makefile commands

| Target | Does |
| --- | --- |
| `make help` | list every target |
| `make dev` / `dev-server` / `dev-worker` | server (air) and worker; or one of them |
| `make build` / `make start` | build the binaries; run the built server |
| `make lint` | `gofmt -l` must be empty, then `go vet ./...` |
| `make format` | `gofmt -w .` |
| `make test` | `go test ./...` — there are no `_test.go` files yet |
| `make docs` | regenerate `docs/swagger.*` with swag |
| `make db-migrate` / `db-migrate-down` | apply / roll back migrations (`go run ./cmd/migrate [down]`) |
| `make db-migrate-create NAME=<name>` | new numbered `.up.sql` / `.down.sql` pair |
| `make db-seed` | idempotent seed |
| `make fresh` | **destructive**: roll back, migrate, seed |
| `make reset` | migrate and seed |
| `make docker-build` / `docker-up` / `docker-down` / `docker-restart` / `docker-logs` / `docker-ps` | compose lifecycle |
| `make docker-migrate` / `docker-seed` | run migrations / seed from the `migrator` image on the compose network |

## Configuration

All variables are read in `internal/config`. Only `JWT_SECRET` and `JWT_REFRESH_SECRET` are
required (and must differ); everything else has a default. Full reference:
[docs/CONFIGURATION.md](docs/CONFIGURATION.md).

`APP_NAME`, `APP_PORT`, `APP_ENV`, `APP_TIMEZONE`, `LOG_LEVEL`, `ENABLE_API_DOCS`, `CLIENT_URL`,
`DATABASE_URL`, `REDIS_HOST`, `REDIS_PORT`, `REDIS_PASSWORD`, `REDIS_DB`, `MAIL_HOST`, `MAIL_PORT`,
`MAIL_SECURE`, `MAIL_FROM`, `MAIL_USER`, `MAIL_PASS`, `JWT_SECRET`, `JWT_REFRESH_SECRET`,
`ALLOWED_HOST`; compose also reads `POSTGRES_DB`, `POSTGRES_USER`, `POSTGRES_PASSWORD`.

## API

Envelope: `{ status, success, message, data }`; lists return `data: { data: [...], meta: { page,
limit, totalCount } }`; 400/422 return `errors: [{ field, message }]`. Response fields are
camelCase. Full reference: [docs/API_DOCUMENTATION.md](docs/API_DOCUMENTATION.md).

| Method | Path | Access |
| --- | --- | --- |
| GET | `/`, `/health`, `/live` | public |
| POST | `/auth/login`, `/auth/refresh`, `/auth/register`, `/auth/resend-verification`, `/auth/verify-email`, `/auth/forgot-password`, `/auth/reset-password` | public |
| GET, PATCH | `/profile` | authenticated (own identity) |
| PATCH | `/profile/password` | authenticated (own identity) |
| GET, POST | `/settings/users` | `user list`, `user create` |
| GET, PATCH, DELETE | `/settings/users/:id` | `user detail`, `user edit`, `user delete` |
| POST | `/settings/users/:id/reset-password` | role `superuser` |
| POST | `/settings/users/:id/send-verification-email`, `/send-reset-password-email` | `user create` |
| GET, POST | `/settings/roles` | `role list`, `role create` |
| GET, PATCH, DELETE | `/settings/roles/:id` | `role detail`, `role edit`, `role delete` |
| GET, POST | `/settings/permissions` | `permission list`, `permission create` |
| GET, PATCH, DELETE | `/settings/permissions/:id` | `permission detail`, `permission edit`, `permission delete` |
| GET | `/settings/select-options/permissions`, `/settings/select-options/roles` | role `superuser` |
| GET | `/docs`, `/docs/openapi.json` | public, only when `ENABLE_API_DOCS=true` |

List endpoints accept `page`, `perPage` (max 100), `search`, `sort`, `sortDirection`, and
`filter[<key>]=<value>`.

## Project structure

```
clean-go-fiber/
├── cmd/
│   ├── server/          # HTTP server (runs migrations on boot)
│   ├── worker/          # asynq worker — sends mail
│   ├── seed/            # idempotent seeder
│   └── migrate/         # migrations up / down
├── internal/
│   ├── bootstrap/       # constructor wiring, route registration
│   ├── config/          # env loading (the only reader of the environment)
│   ├── database/
│   │   ├── postgres/    # GORM connection, migrations/ (embedded), seed/
│   │   └── redis/       # Redis client
│   ├── errors/          # AppError
│   ├── guards/          # guards.Permission, guards.Role
│   ├── libs/            # cache, mailer, pagination, txmanager, utils, validator
│   ├── middleware/      # Auth, logging
│   ├── models/          # GORM models
│   ├── modules/         # auth, profile, users, roles, permissions, selectoptions, health, docs
│   ├── queue/           # asynq client + tasks
│   ├── repositories/    # interfaces + GORM implementations
│   └── worker/          # asynq server + handlers
├── docs/                # guides + swag-generated spec
└── .claude/             # rules, commands, skills (symlink to .agents/skills)
```

Layering: **handler → service → repository → model**. Handlers parse, validate, and respond;
services own business rules, transactions, cache invalidation, and enqueueing; repositories own
every query.

## Documentation

- [docs/README.md](docs/README.md) — index of the guides
- [CLAUDE.md](CLAUDE.md) and [.claude/rules/](.claude/rules/README.md) — conventions
- [.github/CONTRIBUTING.md](.github/CONTRIBUTING.md) — how to contribute

## License

[MIT](LICENSE)
