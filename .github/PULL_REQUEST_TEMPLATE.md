## What does this change?

<!-- One or two sentences. What is different after this PR that was not true before? -->

## Why?

<!-- The problem being solved. Link the issue if there is one: Closes #123 -->

## Type of change

- [ ] `feat` — new feature, module, or endpoint
- [ ] `fix` — bug fix
- [ ] `refactor` — restructuring, no behaviour change
- [ ] `docs` — documentation only
- [ ] `chore` — config, dependencies, tooling
- [ ] `db` — schema or migration change

## Checklist

<!-- The repo's standards live in .claude/rules/ and CLAUDE.md. Consult the relevant rule before ticking. -->

- [ ] Commit messages follow [Conventional Commits](https://www.conventionalcommits.org/) (see `.github/CONTRIBUTING.md`).
- [ ] `make lint` (gofmt + go vet) and `go build ./...` pass locally.
- [ ] Layering respected: `handler.go` (HTTP only) → `service.go` (logic, transactions, cache invalidation, enqueue) → `internal/repositories` (queries only). No `*gorm.DB` outside repositories.
- [ ] The service owns existence and uniqueness checks; repositories return `(nil, nil)` for not found; uniqueness is 422.
- [ ] Every error is handled; no `err.Error()` of an unknown error reaches a client; zap only, no `fmt.Print*`.
- [ ] Request structs have validator tags and camelCase json; responses are camelCase structs, never raw models.
- [ ] Successful responses go through the `utils` response helpers; failures return `*errors.AppError`.
- [ ] Protected routes sit under `middleware.Auth()` and carry `guards.Permission(...)` / `guards.Role(...)` with strings the seed produces.
- [ ] Writes that change a user's roles, status, identity, password, or deletion (or a role's permissions) invalidate the cached `UserInformation`.
- [ ] Handlers have complete swag annotations and `make docs` was re-run.
- [ ] Mail and other background work is enqueued after commit, with `asynq.MaxRetry`, and enqueue errors are logged.
- [ ] Any new env var is read in `internal/config` and added to `.env.example` and `docs/CONFIGURATION.md`.
- [ ] Docs my change makes wrong are fixed **in this PR** (`README.md`, `CLAUDE.md`, `docs/`, the matching `.claude/rules/` file).

## Database changes

<!-- Delete this section if the PR touches no schema. -->

- [ ] New numbered pair in `internal/database/postgres/migrations/` (`.up.sql` and a `.down.sql` that reverses it).
- [ ] No already-applied migration was edited.
- [ ] GORM model in `internal/models/` matches the SQL.
- [ ] Migration reviewed for destructive operations (dropped columns, narrowed types, lost data).
- [ ] Soft-deleted users still excluded from every read.
- [ ] Seed in `internal/database/postgres/seed/` updated if the new shape needs it, and still idempotent.

## How was this tested?

<!--
There are no _test.go files yet, so say what you actually exercised:
commands run, endpoints hit, request/response bodies observed. "It builds" is not testing.
-->

## Screenshots or output

<!-- Optional: request/response bodies, /docs screenshots, worker logs. -->
