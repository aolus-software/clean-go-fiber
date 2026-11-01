# GEMINI.md — clean-go-fiber

You are working on a Go REST API starter (Fiber v2, GORM, golang-migrate, PostgreSQL, Redis, asynq).

## Start here

1. Read `CLAUDE.md` — project overview, commands, architecture, conventions.
2. Read `.claude/rules/README.md` and the rules for the area you are changing. They are binding.
3. Skills live in `.agents/skills/` (`go-fiber`, `golang-pro`, and the general bundle).

## Verification commands

- `make lint` — `gofmt -l` (must be empty) + `go vet ./...`
- `go build ./...`
- `make docs` — regenerate the swagger spec after changing handlers
- `make test` — `go test ./...` (no tests exist yet)
