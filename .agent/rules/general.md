# Rules — clean-go-fiber

> These rules are **non-negotiable**. Follow them on every task, every file, every PR.

---

## Project Structure Rules

- All business logic lives in `internal/modules/<module>/service.go` — nowhere else
- Handlers only: parse request body → call service → return response. No logic in handlers.
- Repository interfaces are defined in `internal/repositories/interfaces.go`
- Concrete GORM implementations go in `internal/repositories/<name>.go`
- Models live only in `internal/models/` — no methods, no business logic on models
- Config is loaded once at startup. Always access it via `config.C` — never `os.Getenv` directly outside of `internal/config/`
- No direct database access outside the repository layer
- Every new module must follow the four-file pattern: `handler.go`, `routes.go`, `schema.go`, `service.go`
- Do not create packages outside of the established folder tree without a strong reason

---

## Go Conventions

- **Always handle errors** — never use `_` to discard errors from external calls (DB, Redis, network)
- Use named return values only when it meaningfully improves readability
- Prefer explicit over implicit — no `init()` functions outside `main` packages, no hidden side effects
- Package names: lowercase, single word, no underscores (e.g., `repositories`, not `repo_layer`)
- File names: `snake_case.go`
- Interface names: suffix with `Repository` or `Service` (e.g., `UserRepository`, `AuthService`)
- All exported types, functions, and methods must have a godoc comment
- Pass `context.Context` through all DB, Redis, and HTTP client calls — never use `context.Background()` inside a handler; propagate from `c.Context()` (Fiber)
- Avoid global mutable state — the only acceptable globals are `config.C`, the DB instance, and the Redis client, all initialized at boot

---

## API Response Rules

- All responses use helpers from `internal/libs/utils/response.go` (e.g., `utils.Success()`, `utils.Error()`)
- Never expose raw Go error messages to the client
- Validation errors return HTTP `422` with per-field details
- Successful resource creation returns HTTP `201`
- Successful deletion returns HTTP `204` with no body
- Pagination responses always include: `data`, `total`, `page`, `limit`, `total_pages`
- HTTP status codes must be semantically correct — do not use `200` for everything

---

## Security Rules

- **Never log** passwords, JWT tokens, secrets, or any PII
- Always hash passwords with `bcrypt` before storing — never store plaintext
- JWT secret must come from `config.C.JWT.Secret` — never hardcode it
- All protected routes must pass through the JWT auth middleware
- Auth endpoints (`/login`, `/register`, `/forgot-password`) must have rate limiting applied
- Input validation is required on every `POST`, `PATCH`, and `PUT` handler before calling the service
- Sanitize and validate email addresses — do not trust user input

---

## Database Rules

- Use **transactions** for any operation that involves more than one write
- Use soft delete (`gorm.DeletedAt`) for all user-facing entities (users, roles, permissions)
- Use hard delete only for token tables (`email_verifications`, `password_reset_tokens`)
- Avoid raw SQL — if GORM cannot express the query, add a comment explaining why raw SQL is needed
- Migrations are **additive only** — never drop or rename columns in a migration; use a new additive migration
- Seed data must be **idempotent** — always use `FirstOrCreate` or `FirstOrUpdate`, never bare `Create`
- Always preload associations explicitly — never rely on GORM's global preload behavior

---

## Background Job Rules

- Task type constants are defined in `internal/queue/tasks/` — never use raw strings for task names
- Task payloads must be JSON-serializable structs
- Every job handler must log both on success and on failure (with error details)
- Email sending **always** goes through the Asynq queue — never call the mailer inline inside an HTTP handler
- All jobs must define an explicit `MaxRetry` in their `asynq.Options`
- Do not spawn raw goroutines inside HTTP handlers — use the queue for all async work

---

## Middleware Rules

- Middleware registration order in `cmd/server/main.go` matters — follow this order:
  1. Recovery
  2. Request ID
  3. Logger
  4. Security headers
  5. CORS
  6. Body limit
  7. Rate limiter
  8. Application routes
- Auth middleware (`middleware.Auth()`) is applied per-route-group, not globally
- Guards (`guards.Role()`, `guards.Permission()`) are applied per-route, after auth middleware

---

## Documentation Rules

- Every Fiber handler must have complete `swaggo` annotations:
  `@Summary`, `@Tags`, `@Accept`, `@Produce`, `@Param`, `@Success`, `@Failure`, `@Security`, `@Router`
- `.env.example` must stay in sync with all fields in `internal/config/`
- `README.md` must reflect the current setup and workflow
- Makefile `help` target must list all available targets with descriptions
- After adding/changing API handlers, always run `make docs` to regenerate Swagger

---

## Testing Rules

- Test files are co-located with source: `<file>_test.go`
- Use **table-driven tests** for all service and utility functions
- Repository tests use `testcontainers-go` with a real PostgreSQL instance — no SQLite mocks
- Hand-write mocks or use `mockery` for interface mocks — no monkey-patching
- No test may depend on global state or share a database transaction with another test
- Each test must clean up after itself

## Code Styles

- IMPORTANT: DO NOT ADD **_ANY_** COMMENTS unless asked
- IMPORTANT: DO NOT ADD **_ANY_** DOCS unless asked

---

## Forbidden Patterns

```
❌ init() outside of main packages
❌ panic() outside of startup in main()
❌ os.Getenv() outside of internal/config/
❌ Goroutines spawned inside HTTP handlers
❌ Inline email sending inside request handlers
❌ Circular imports between internal packages
❌ Raw SQL without an explaining comment
❌ Hardcoded secrets, passwords, or keys
❌ Bare Create() in seed files (must be idempotent)
❌ Skipping input validation on mutating endpoints
❌ Returning raw Go errors to the HTTP client
```
