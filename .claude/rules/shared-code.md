# Rule: Shared code

Where a thing lives once a second module needs it. Feature code stays in `internal/modules/<name>/`;
nothing in one module imports another module.

| Package | Holds |
| --- | --- |
| `internal/config` | env loading and typed config structs; the only place that reads the environment |
| `internal/errors` | `AppError` and its constructors — the client-facing error vocabulary |
| `internal/guards` | `guards.Permission`, `guards.Role` |
| `internal/middleware` | `Auth`, request logger, slow-request logger |
| `internal/models` | GORM models and the `user_status` constants |
| `internal/repositories` | repository interfaces and GORM implementations |
| `internal/database/postgres` | connection, embedded migrations runner, `seed/` |
| `internal/database/redis` | Redis client (pinged at connect) |
| `internal/queue` | asynq client and `tasks/` (task types + payloads) |
| `internal/worker` | asynq server and `handlers/` |
| `internal/bootstrap` | constructor wiring and route registration; builds `middleware.Auth(userRepo, cache)` once |
| `internal/modules/auth/account_mail.go` | `AccountMail` — issues verification/reset tokens and enqueues their mails; shared by the auth and users services (the one sanctioned cross-module use) |
| `internal/logging` | zap logger construction from `APP_ENV` / `LOG_LEVEL`, shared by server and worker |
| `internal/libs/cache` | Redis cache-aside wrapper and key constants (`user:<id>`) |
| `internal/libs/mailer` | SMTP sender used by the worker only |
| `internal/libs/pagination` | `Parse(c)`, `Options` allow-lists, `Apply(db, opts, tz)` |
| `internal/libs/utils` | `OK` / `Created` / `Error`, `PaginatedData`, `ParseBody`, `ParamUUID`, bcrypt, JWT, one-time tokens |
| `internal/libs/validator` | go-playground/validator instance, custom tags, json field names, UUID param parsing |
| `internal/libs/txmanager` | transaction helper used by services |

## Rules

- A helper used by one module stays in that module. The second consumer moves it to the right
  bucket above.
- Do not add a new `internal/libs/<bucket>` without updating this table and `CLAUDE.md`.
- No package outside this table without a strong reason; never a `common` or `helpers` grab-bag.
