# Rule: Fiber application structure

Applies to the whole tree. Framework: Fiber v2 on fasthttp; ORM: GORM; migrations: golang-migrate.

## Entrypoints

| Binary | Path | Does |
| --- | --- | --- |
| server | `cmd/server/main.go` | load config, connect Postgres, run migrations, connect Redis (pinged), build the Fiber app and global middleware, `bootstrap.RegisterRoutes`, graceful shutdown |
| worker | `cmd/worker/main.go` | asynq server; the only process that sends mail |
| seed | `cmd/seed/main.go` | idempotent seeder; exits non-zero on failure |
| migrate | `cmd/migrate/main.go` | `up` (default) or `down` against `DATABASE_URL` |

## Layering

`handler.go` → `service.go` → `internal/repositories/*` → `internal/models/*`. Never skip a layer:

- Handlers parse, validate, call one service method, and respond. No business logic, no GORM.
- Services own business rules, existence/uniqueness checks, transactions, cache invalidation, and
  enqueueing. They return `*errors.AppError` for anything a client should see.
- Repositories own every query. No `*gorm.DB` outside `internal/repositories`, `internal/database`,
  and the seed.
- Models are plain GORM structs mirroring the SQL. No methods carrying business logic.

## The four-file module

Every feature under `internal/modules/<name>/` has `handler.go`, `routes.go`, `schema.go`,
`service.go`. Settings modules (`users`, `roles`, `permissions`, `selectoptions`) mount under the
`/settings` group. `health` and `docs` have no service.

## Dependency injection

Constructor injection only, wired by hand in `internal/bootstrap/bootstrap.go`: build repositories,
pass them into `NewService(...)`, pass the service into `NewHandler(...)`, call
`RegisterRoutes(...)`. Services and repositories are interfaces so they can be faked in tests. There
is no container — do not add one.

## Config

`config.Load()` runs once at startup and populates `config.C`. Read configuration through `config.C`
(or a `*config.Config` passed in); never `os.Getenv` outside `internal/config`. A required variable
missing at boot is a fatal error with the variable's name in the message.

## Globals

The only acceptable globals are `config.C`, the GORM handle, and the Redis client — all initialised
at boot. Everything else is passed in.

## Reuse before you build

| Need | Use |
| --- | --- |
| response envelope | `internal/libs/utils/response.go` |
| client-facing error | `internal/errors` |
| request validation | `internal/libs/validator` |
| list query parsing | `internal/libs/pagination` |
| password hash / JWT / one-time token | `internal/libs/utils` |
| cache-aside | `internal/libs/cache` |
| transaction | `internal/libs/txmanager` |
| mail | enqueue a task from `internal/queue/tasks` |

## Out of scope for this boilerplate

GraphQL, gRPC, WebSockets, file storage, multi-tenancy, OAuth/social login, SSR.
