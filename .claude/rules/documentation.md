# Documentation Upkeep Rule

## Principle: docs are part of the change, not an afterthought

When a change makes a documentation file wrong, fixing that doc belongs to the **same change**. Code
and docs are committed together so the repository never carries documentation that contradicts the
code. If nothing a doc covers changed, leave it alone. Do not create new `.md` files nobody asked for.

## Docs that must stay in sync

| Doc | Update it when… |
|---|---|
| `README.md` | setup steps, `make` targets, env vars, the project tree, or the Docker section change |
| `CLAUDE.md` | a package, the boot order, the layering, a convention, a `docs/` guide, or a module changes |
| `Makefile` | a canonical command is added, renamed, or removed — the `.PHONY` line, the `help` target, and the recipe move together, plus `README.md` |
| `docker-compose.yml`, `Dockerfile` | a service, port, volume, or build stage changes — `docker-migrate` / `docker-seed` depend on the `migrator` stage |
| `.env.example` | any env var is added, renamed, or removed — it must be copyable as-is (no empty value the loader rejects) |
| `internal/config/*.go` | a new env var is read. It is read there or nowhere; then `.env.example` and `docs/CONFIGURATION.md` |
| `internal/database/postgres/migrations/` + `internal/models/` | a table, column, enum, or index changes — a new numbered `.up.sql`/`.down.sql` pair and the GORM model move together; the seed moves with them when a fresh database needs different data |
| `docs/API_DOCUMENTATION.md` | a route is added, renamed, re-gated, or removed, or a request/response shape changes |
| `docs/CONFIGURATION.md` | the env var reference drifts from `internal/config` / `.env.example` |
| `docs/DEPLOYMENT.md` | the Docker flow, the migrator stage, or the deploy sequence changes |
| `docs/ERROR_HANDLING.md` | `AppError`, the central error handler, or the envelope changes |
| `docs/MIDDLEWARE.md` | the global middleware chain or a middleware/guard changes |
| `docs/SECURITY.md` | auth, guards, tokens, rate limiting, CORS, headers, or body limits change |
| `docs/README.md` | a file is added to or removed from `docs/` — it is the index |
| `docs/swagger.json` / `swagger.yaml` / `docs.go` | a handler's swag annotations change — run `make docs`; never hand-edit |
| `.claude/rules/*.md` | a coded convention changes, or a new pattern ships with no rule yet |
| `.claude/commands/*.md` | a command's workflow or scope changes |
| `.github/copilot-instructions.md`, `.gemini/GEMINI.md` | the pointers they carry (to `CLAUDE.md` and `.claude/rules/`) go stale |
| `todo.md` | an item it tracks is done or dropped (`/update-todo`) |

## What "up to date" means

- **Exact facts.** Route paths, table and column names, env var names, permission strings, and
  `make` targets match reality. A stale env var name is a documentation bug, not a nitpick.
- **No orphan references.** Rename, move, or delete something → update every doc that names it.
- **Permission strings are quoted verbatim.** They are space-separated (`"user list"`,
  `"role create"`), produced by the seeder as `<group> <action>`.
- **New surfaces get docs.** A new middleware, guard, task/handler pair, repository, or external
  integration that establishes a pattern needs its rule written.
- **The OpenAPI spec is documentation.** A handler without complete swag annotations is an
  incomplete change ([openapi.md](./openapi.md)).

## When a doc is wrong but the current task didn't cause it

Per `contradiction-halt.md`, **report it** — do not silently rewrite unrelated docs.
