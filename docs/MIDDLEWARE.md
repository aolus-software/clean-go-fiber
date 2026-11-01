# Middleware

Fiber middleware is a `fiber.Handler`; it runs in registration order and calls `c.Next()` to pass
control on. This project uses three layers: the global chain, the per-group `auth` handler, and
per-route guards.

## Global chain (`cmd/server/main.go`)

| # | Middleware | Purpose |
| - | --- | --- |
| 1 | recover | turn a panic into a 500 instead of killing the process |
| 2 | requestid | `X-Request-ID` on every response |
| 3 | helmet | security headers |
| 4 | cors | origins from `ALLOWED_HOST` (empty → `*`) |
| 5 | limiter | 100 requests / 60 s per IP, 429 envelope |
| 6 | `middleware.Logger` | one zap line per request: method, path, status, duration, ip |
| 7 | `middleware.Performance` | warns on requests slower than 1 s |

`fiber.Config` also sets `BodyLimit`, read/write timeouts, and the central `ErrorHandler`
([ERROR_HANDLING.md](./ERROR_HANDLING.md)).

## Authentication: `middleware.Auth(userRepo, cache)`

Built once in `internal/bootstrap` and passed as `auth fiber.Handler` to each module's
`RegisterRoutes(router, auth)`, which applies it to its group (`/profile`, `/settings/...`). Public
routes never see it.

1. Requires `Authorization: Bearer <token>`.
2. Verifies the JWT with `JWT_SECRET` and requires `type == "access"`.
3. Loads `UserInformation` cache-aside (Redis `user:<id>`, 3600 s; database fallback requires an
   active, non-deleted user).
4. Stores it for later handlers: `middleware.CurrentUser(c)`.

Any failure is 401.

## Authorization: guards

```go
r := settings.Group("/users", auth)
r.Get("", guards.Permission("user list"), h.List)
r.Post("/:id/reset-password", guards.Role("superuser"), h.ResetPassword)
```

- `guards.Permission(names...)` — the caller must hold every named permission.
- `guards.Role(names...)` — the caller must hold every named role.
- Both bypass for `superuser`; both answer 403 on refusal.

## Adding a middleware

1. Create `internal/middleware/<name>.go` returning `fiber.Handler`.
2. Fail by returning an `*errors.AppError`; never write a response body directly.
3. Register it globally in `cmd/server/main.go` or on the group that needs it.
4. Update this page and `.claude/rules/middleware.md`.
