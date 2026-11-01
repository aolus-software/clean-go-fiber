# Project Rules

Coding rules for the `clean-go-fiber` codebase (Go 1.25 + Fiber v2 + GORM + golang-migrate +
PostgreSQL + Redis + asynq). Each file is a focused, enforceable contract — read the relevant one
before writing code in that area.

Rules here carry **no** frontmatter; the scope is stated in the table below and in each file's
opening line. They replace the former `.agent/AGENTS.md` and `.agent/rules/general.md`, whose
still-valid content was folded in here.

## Always in scope

| Rule | Scope |
| ---- | ----- |
| [contradiction-halt.md](./contradiction-halt.md) | A request that contradicts a rule, the architecture, or a security invariant is reported and halted — never silently implemented. Lists the invariants |
| [documentation.md](./documentation.md) | A doc your change makes wrong is fixed in the **same** change; lists every doc that must stay in sync |
| [audit-findings.md](./audit-findings.md) | How an audit finding is written — the writing contract for [`/audit-flow`](../commands/audit-flow.md) |
| [clean-code.md](./clean-code.md) | gofmt, `go vet`, explicit types, error handling, zap-only logging, comment density |
| [fiber.md](./fiber.md) | Layering (`handler → service → repository`), constructor injection in `internal/bootstrap`, config, the four-file module layout |

## Layer rules

Ordered outside-in, the way a request travels.

| Rule | Applies to |
| ---- | ---------- |
| [middleware.md](./middleware.md) | `internal/middleware/`, the global chain in `cmd/server/main.go` |
| [handlers.md](./handlers.md) | `internal/modules/<name>/handler.go` + `routes.go` — what a handler may do |
| [handlers-crud.md](./handlers-crud.md) | The canonical five-route CRUD module and its status codes |
| [validation.md](./validation.md) | `internal/modules/<name>/schema.go` — validator tags, `:id` parsing, list query |
| [services.md](./services.md) | `internal/modules/<name>/service.go` — transactions, cache invalidation, error vocabulary |
| [services-crud.md](./services-crud.md) | The canonical CRUD service, which **owns existence and uniqueness checks** (422) |
| [repositories.md](./repositories.md) | `internal/repositories/` — interfaces, `WithTx`, soft-delete filter, allow-lists, `(nil, nil)` for not found |
| [schema.md](./schema.md) | `internal/database/postgres/migrations/` + `internal/models/` — SQL pairs, enum, token tables |
| [shared-code.md](./shared-code.md) | `internal/libs/<bucket>/` and the other shared packages |

## Cross-cutting concerns

| Rule | Applies to |
| ---- | ---------- |
| [rbac.md](./rbac.md) | `guards.Permission` / `guards.Role`, the seeded vocabulary, privilege-granting routes |
| [routes.md](./routes.md) | Path composition, verb conventions, the **live route map** |
| [errors-and-responses.md](./errors-and-responses.md) | The envelope, `AppError`, the central Fiber error handler |
| [openapi.md](./openapi.md) | swaggo annotations, `make docs`, the `ENABLE_API_DOCS` gate |
| [i18n.md](./i18n.md) | There is no catalogue — English only, and where strings live |
| [mail.md](./mail.md) | Mail goes through the asynq queue, never inline |
| [rate-limiting.md](./rate-limiting.md) | The global Fiber limiter and why its numbers are hardcoded |
| [queue.md](./queue.md) | `internal/queue/` + `internal/worker/` — asynq tasks and handlers |
| [imports-and-naming.md](./imports-and-naming.md) | Import grouping, package, file, and symbol naming |
| [commit.md](./commit.md) | Conventional Commits, the gates to run, what never gets committed |

## How to use

- These rules complement `CLAUDE.md`. When a rule conflicts with `CLAUDE.md`, the **rule file wins**.
- Don't introduce a new pattern without updating the relevant rule first, and add a new rule to this
  index ([documentation.md](./documentation.md)).
- Slash commands live in [`../commands/`](../commands/): `/commit`, `/audit-flow`, `/update-todo`.

## Settled decisions (do not reopen as drive-bys)

- **Uniqueness conflicts are 422** (email, role name, permission name), including `/auth/register`.
- **The service owns existence and uniqueness checks.** Repositories return `(nil, nil)` for a
  missing row; the service turns that into 404 or 422.
- **The API contract follows the sibling `clean-elysia`** (routes, permission strings, envelope, list
  query parameters, env names), with camelCase response fields. Go keeps `APP_ENV` where Elysia uses
  `NODE_ENV`, and keeps `/auth/refresh` (separate secret, `type` claim checked).
