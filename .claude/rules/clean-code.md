# Rule: Clean code

Formatting and vetting are machine-enforced. This rule covers what the machine checks and the
conventions it cannot.

## The commands

| Task | Command | Make target |
| --- | --- | --- |
| Format | `gofmt -w .` | `make format` |
| Lint (format check + vet) | `gofmt -l .` must print nothing, then `go vet ./...` | `make lint` |
| Build | `go build ./...` | `make build` |
| Test | `go test ./...` (no `_test.go` files exist yet) | `make test` |

There is no golangci-lint config in the repo; `go vet` is the linter of record.

## Formatting

gofmt owns it: tabs, its alignment, its import grouping within a block. Don't hand-align, and don't
reformat a file you didn't otherwise touch.

## Types

- Avoid `any` / `interface{}` in signatures where a concrete type exists. It is acceptable where the
  API is genuinely untyped (`json.Marshal` input, the response helper's `data`, cache values).
- Request and response shapes are named structs in the module's `schema.go`, not `fiber.Map`.
- Exported functions declare their return types (Go forces it); return `(T, error)`, never panic.

## Errors

- Always handle errors. Never discard an error from the database, Redis, the queue, or the network
  with `_`.
- Wrap with context when passing an internal error up (`fmt.Errorf("find user: %w", err)`); translate
  to an `*errors.AppError` at the service boundary so the handler never decides a status code.
- `panic` only during startup in `main`. No `init()` outside `main` packages.

## Logging

Use the zap logger, structured fields first-class:

```go
log.Error("failed to enqueue verification email", zap.String("userId", user.ID.String()), zap.Error(err))
```

No `fmt.Print*` or the standard `log` package in request paths (the standard `log` is acceptable in
`cmd/*` before the zap logger exists, and in `internal/config` for fatal startup errors). Never log a
password, a token, a token-bearing URL, or a JWT.

## Context

Pass `context.Context` through every DB, Redis, and queue call. In a handler, propagate
`c.UserContext()` / `c.Context()`; never `context.Background()` inside a request path.

## Comments

None unless a block genuinely needs one — then a single block comment above the block, explaining
what and why. No line-by-line narration, no commented-out code. swag annotations on handlers are
required and are not "comments" for this rule ([openapi.md](./openapi.md)).

## Don't

1. No emoji or decorative symbols in code, comments, log messages, or templates.
2. No `os.Getenv` outside `internal/config`.
3. No goroutines spawned inside HTTP handlers — use the queue.
4. No circular imports between internal packages.
5. No new README or docs file unless asked; keeping an existing doc true is `documentation.md`.
