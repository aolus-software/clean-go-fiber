# Contributing

Thanks for helping improve `clean-go-fiber`.

## Before you start

- Read `CLAUDE.md` and the rules in `.claude/rules/` for the area you are changing.
- For anything larger than a bug fix, open an issue or a discussion first so the design can be
  agreed before you write it.

## Local setup

Requirements: Go 1.25+, Docker (for Postgres and Redis), and optionally
[air](https://github.com/air-verse/air) for hot reload and [swag](https://github.com/swaggo/swag) for
`make docs`.

```bash
git clone https://github.com/aolus-software/clean-go-fiber.git
cd clean-go-fiber
cp .env.example .env              # set JWT_SECRET and JWT_REFRESH_SECRET
docker compose up -d postgres redis
go mod download
make db-migrate
make db-seed
make dev
```

### With Docker only

```bash
make docker-build
make docker-up
make docker-migrate
make docker-seed
```

## Coding standards

- `make lint` (gofmt + go vet) and `go build ./...` must pass.
- Follow the layering: handler → service → repository → model.
- Handle every error; log with zap; never send an unknown error's text to a client.
- No comments unless a block genuinely needs one.
- Response fields camelCase; permission strings `"<group> <action>"`.

## Adding a feature module

Follow "Adding a new feature module" in `CLAUDE.md`: migration + model, repository, `schema.go`,
`service.go`, `handler.go` with swag annotations, `routes.go` with guards, wiring in
`internal/bootstrap`, then the route map in `.claude/rules/routes.md` and
`docs/API_DOCUMENTATION.md`.

## Database changes

- `make db-migrate-create NAME=<name>` creates a numbered `.up.sql` / `.down.sql` pair. Write both.
- Never edit a migration that has been applied anywhere.
- Keep the GORM model in step with the SQL, and keep the seed idempotent.

## Commits

[Conventional Commits](https://www.conventionalcommits.org/): `feat:`, `fix:`, `chore:`,
`refactor:`, `docs:`, `remove:`, `test:`. Lowercase summary, imperative, under 72 characters. See
`.claude/rules/commit.md`.

## Pull requests

- Fill in the PR template honestly, including how you tested.
- Keep a PR to one concern; fix docs your change makes wrong in the same PR.
- CI (`.github/workflows/ci.yml`) must be green.

## Reporting bugs

Use the bug report template. Security issues go through
[Security Advisories](https://github.com/aolus-software/clean-go-fiber/security/advisories/new),
never a public issue.

## License

By contributing you agree that your contributions are licensed under the MIT License in `LICENSE`.
