# Rule: Commits

## Message format

Conventional Commits: `feat`, `fix`, `chore`, `refactor`, `docs`, `remove`, `test`. The summary is
lowercase after the colon, imperative mood, under 72 characters, no trailing period.

```
feat: add refresh token rotation
fix: return 422 for a duplicate email on admin create
chore: regenerate swagger spec
```

The full workflow lives in `.claude/commands/commit.md`; it is the authority when the two drift.

## Gates

There is **no pre-commit hook** in this repo (no Husky in a Go project). The gates are run by hand
before every commit:

```
make lint        # gofmt -l . (must print nothing) + go vet ./...
go build ./...
```

If you changed a handler's swag annotations, run `make docs` and commit the regenerated
`docs/swagger.*` and `docs/docs.go` with the change.

## Never commit

1. `.env` — it is gitignored; do not `git add -f` it.
2. Secrets, keys, tokens, or a real `DATABASE_URL` / SMTP credential in a default or a fixture.
3. `tmp/` (air output), compiled `server` / `worker` binaries, `build-errors.log`.
4. `skills-lock.json` — gitignored; the skills bundle is vendored.
5. A hand-edited `docs/swagger.json` / `docs/docs.go`.
6. An edit to an already-applied migration. A schema change is a new numbered pair.
7. gofmt churn in files the change did not otherwise touch — that is its own `chore:` commit.

## Staging

Stage by explicit path. `git add -A` and `git add .` are forbidden by `.claude/commands/commit.md`.
