# Rule: Middleware

Applies to `internal/middleware/`, `internal/guards/`, and the global chain in `cmd/server/main.go`.

## Global chain (order matters)

Registered with `app.Use(...)` in `cmd/server/main.go`, before `bootstrap.RegisterRoutes`:

1. recover
2. request id
3. helmet (security headers)
4. CORS — origins from `ALLOWED_HOST` (empty → `*`)
5. rate limiter — 100 requests / 60 s per IP, 429 envelope ([rate-limiting.md](./rate-limiting.md))
6. request logger (zap)
7. slow-request logger

Body size is capped by `fiber.Config.BodyLimit`. The central `ErrorHandler` in `fiber.Config` turns
every returned error into the envelope ([errors-and-responses.md](./errors-and-responses.md)).

## Per-group and per-route

- `middleware.Auth(userRepo, cache)` is built once in `internal/bootstrap` and passed as
  `auth fiber.Handler` into each module's `RegisterRoutes(router, auth)`, which applies it to its
  group (`/profile`, `/settings/...`), never globally — the
  public routes (`/`, `/health`, `/live`, `/auth/*`, `/docs`) must stay reachable.
- Guards (`guards.Permission(...)`, `guards.Role(...)`) are applied per route, **after**
  `middleware.Auth()`. See [rbac.md](./rbac.md).

## What `middleware.Auth()` does

1. Reads `Authorization: Bearer <token>`; missing or malformed → 401.
2. Verifies the JWT with `JWT_SECRET` and requires claim `type == "access"` — a refresh token is
   rejected here.
3. Loads `UserInformation` (`id, name, email, roles[], permissions[]`) cache-aside from Redis key
   `user:<id>` (TTL 3600 s). The database lookup requires `status = 'active' AND deleted_at IS NULL`;
   otherwise 401.
4. Stores it in `c.Locals`; guards read it there and handlers use `middleware.CurrentUser(c)`.

Because of step 3, a status change, role change, deletion, or password reset only takes effect once
the cache key is dropped — services own that invalidation ([services.md](./services.md)).

## Writing a new middleware

1. `internal/middleware/<name>.go`, returning `fiber.Handler`.
2. Return an `*errors.AppError` to fail a request — never write a response body directly.
3. Register it in the chain above (or on the group that needs it) and update this rule and
   `docs/MIDDLEWARE.md` in the same change.
