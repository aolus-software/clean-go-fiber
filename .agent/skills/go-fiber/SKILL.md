# Skill — Go Fiber v2

## When to use
Any time you are writing or modifying HTTP route handlers, middleware, route registration, or anything that touches the Fiber application instance.

## Key Patterns

- **Handler signature:** always `func(c *fiber.Ctx) error`
- **Parse body:** `c.BodyParser(&req)` — always validate after parsing
- **Route groups:** `api := app.Group("/api")`, then `api.Group("/auth")`
- **Error return:** `return c.Status(fiber.StatusBadRequest).JSON(...)` — never panic
- **Locals (context values):** `c.Locals("user", claims)` to set, `c.Locals("user").(*jwt.Claims)` to get
- **Path params:** `c.Params("id")`
- **Query params:** `c.Query("page", "1")`
- **Middleware registration:** `app.Use(middleware.Name())`
- **Per-route middleware:** `router.Get("/me", authMiddleware, handler.GetMe)`

## Avoid
- Business logic inside handlers — it belongs in the service layer
- Calling `c.Next()` and also returning a response in the same middleware — pick one
- Using `fiber.Map{}` for complex response shapes — use typed response structs

---

# Skill — GORM + PostgreSQL

## When to use
Defining models, writing repository implementations, setting up the database connection, migrations, or seeds.

## Key Patterns

- **UUID primary key:** `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
- **Soft delete:** embed `gorm.DeletedAt` and use `gorm.Model` or add it manually
- **Many-to-many:** `gorm:"many2many:user_roles"` on the struct field
- **Preload:** `db.Preload("Roles").Preload("Roles.Permissions").First(&user, id)`
- **Transaction:** `db.Transaction(func(tx *gorm.DB) error { ... })`
- **Connection pool:**
  ```go
  sqlDB, _ := db.DB()
  sqlDB.SetMaxOpenConns(cfg.MaxOpen)
  sqlDB.SetMaxIdleConns(cfg.MaxIdle)
  sqlDB.SetConnMaxLifetime(time.Hour)
  ```
- **Repository pattern:** accept `*gorm.DB` in constructor, store as field, use in methods
- **Idempotent seed:** `db.FirstOrCreate(&role, models.Role{Slug: "admin"})`

## Avoid
- Raw SQL unless GORM cannot express the query (add a comment explaining why)
- Accessing `db` outside the repository layer
- Dropping or renaming columns in migrations — always additive
- Relying on GORM's auto-preload global setting — always preload explicitly

---

# Skill — Clean Architecture Layout

## When to use
Every task. This is the foundational structure — always respect it.

## Key Patterns

- **Dependency direction:** `handler → service → repository → model`
- **No skipping layers:** handlers never call repositories directly
- **No circular imports:** `models` imports nothing internal; `repositories` imports only `models`; `modules` imports `repositories` and `libs`
- **Bootstrap wiring:**
  ```go
  userRepo := repositories.NewUserRepository(db)
  authService := auth.NewService(userRepo, cfg)
  authHandler := auth.NewHandler(authService)
  authHandler.RegisterRoutes(api)
  ```
- **Schema structs** (request/response) live in `schema.go` per module — never reuse GORM models as request bodies

## Avoid
- Putting business logic in models
- Putting DB calls in handlers
- Sharing schema structs between modules
- Creating packages outside the established tree without a clear reason

---

# Skill — JWT Authentication

## When to use
Writing or modifying the auth middleware, login/register handlers, token generation, or token refresh logic.

## Key Patterns

- **Generate token:**
  ```go
  claims := jwt.MapClaims{
      "sub": user.ID,
      "exp": time.Now().Add(expiry).Unix(),
  }
  token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
  signed, _ := token.SignedString([]byte(config.C.JWT.Secret))
  ```
- **Validate in middleware:** parse `Authorization: Bearer <token>` header, validate signature and expiry, attach claims to `c.Locals("user", claims)`
- **Refresh token:** store refresh token hash in Redis with TTL; on refresh, validate hash, issue new pair, invalidate old
- **Token errors:** distinguish between expired (`jwt.ErrTokenExpired`) and invalid — return `401` for both but with different messages

## Avoid
- Storing full JWT in the database — Redis for refresh tokens is fine
- Hardcoding the JWT secret
- Accepting tokens from query params in production routes — header only

---

# Skill — RBAC Guards

## When to use
Adding role or permission checks to routes.

## Key Patterns

- **RoleGuard:** middleware that reads `c.Locals("user")`, checks `user.Roles` for the required role slug, returns `403` if not found
- **PermissionGuard:** same pattern but checks flattened permissions across all user roles
- **Superuser bypass:** always check `user.IsSuperuser` first — if true, call `c.Next()` immediately
- **Apply per-route in routes.go:**
  ```go
  router.Delete("/:id", authMiddleware,
      guards.Role("admin"),
      handler.DeleteUser,
  )
  ```
- **Multiple roles/permissions:** guards should accept variadic args — user passes if they have ANY of the listed roles/permissions

## Avoid
- Putting RBAC logic inside service layer — it belongs in middleware/guards
- Checking permissions in handlers — use guards
- Forgetting to apply `authMiddleware` before a guard (guard assumes user is already in locals)

---

# Skill — Redis + Asynq

## When to use
Caching data, enqueuing background jobs, or implementing rate limiting.

## Key Patterns

**Cache wrapper:**
```go
// Remember pattern — cache or compute
value, err := cache.Remember(ctx, "users:list", 5*time.Minute, func() (any, error) {
    return repo.FindAll()
})
// Invalidate
cache.Delete(ctx, "users:list")
```

**Key naming:** `<entity>:<id>:<field>` — e.g., `users:abc123:profile`, `roles:list`

**Asynq — enqueue:**
```go
payload, _ := json.Marshal(tasks.SendEmailPayload{To: email, Template: "verify"})
task := asynq.NewTask(tasks.TypeSendEmail, payload, asynq.MaxRetry(3))
client.Enqueue(task)
```

**Asynq — handler:**
```go
func HandleSendEmail(ctx context.Context, t *asynq.Task) error {
    var p tasks.SendEmailPayload
    json.Unmarshal(t.Payload(), &p)
    // process...
    return nil
}
```

## Avoid
- Caching data that must be strongly consistent (e.g., user balances, permissions mid-request)
- Using raw string literals for task type names — always use constants
- Enqueueing without setting `MaxRetry`
- Sending emails synchronously inside HTTP handlers

---

# Skill — Validation

## When to use
Any handler that accepts a request body (POST, PATCH, PUT).

## Key Patterns

- **Struct tags:** `validate:"required,email"`, `validate:"required,min=8"`, `validate:"required,oneof=admin user"`
- **Validate in handler:**
  ```go
  if err := validator.Validate(req); err != nil {
      return c.Status(422).JSON(utils.ValidationError(err))
  }
  ```
- **Validation error response:** field-level map `{ "field": "error message" }` — not a single string
- **Custom messages:** wrap `go-playground/validator` to translate tag names to human-readable messages

## Avoid
- Validating inside service layer — validation is a handler concern
- Returning the raw validator error string to the client
- Skipping validation on "simple" endpoints — always validate

---

# Skill — Structured Logging (zap)

## When to use
Any logging — request logs, error logs, job logs, startup logs.

## Key Patterns

- **Init:** production uses `zap.NewProduction()`, development uses `zap.NewDevelopment()`
- **Request log fields:** `method`, `path`, `status`, `latency`, `request_id`
- **Error logging:** `logger.Error("message", zap.Error(err), zap.String("context", ...))`
- **Job logging:** log task type, payload summary, success/failure outcome
- **Never use `log.Println`** in application code — always use the zap instance

## Avoid
- Logging sensitive fields (passwords, tokens, emails in some contexts)
- Using `zap.Sugar()` in hot paths — use the typed logger
- Creating a new logger per-request — pass it via context or use a global instance initialized at boot

---

# Skill — Swagger / swaggo

## When to use
Adding or modifying any HTTP handler that should appear in the API docs.

## Key Patterns

```go
// @Summary     Login user
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       body  body      auth.LoginRequest  true  "Login credentials"
// @Success     200   {object}  utils.Response{data=auth.TokenResponse}
// @Failure     401   {object}  utils.ErrorResponse
// @Failure     422   {object}  utils.ValidationErrorResponse
// @Router      /auth/login [post]
func (h *Handler) Login(c *fiber.Ctx) error {
```

- **Security annotation** on protected routes: `// @Security BearerAuth`
- **Generate:** `make docs` (runs `swag init -g cmd/server/main.go -o docs`)
- **Serve:** `app.Get("/docs/*", swagger.HandlerDefault)`
- **Main file annotation** (in `cmd/server/main.go`):
  ```
  // @title           Clean Go Fiber API
  // @version         1.0
  // @securityDefinitions.apikey BearerAuth
  // @in header
  // @name Authorization
  ```

## Avoid
- Editing files in `docs/` manually — they are auto-generated
- Forgetting `@Failure` annotations — document at minimum 400/401/422/500
- Using `fiber.Map` as the response type in annotations — always use a typed struct

---

# Skill — Error Handling

## When to use
Defining new error types, returning errors from services, or configuring the global error handler.

## Key Patterns

- **AppError struct:**
  ```go
  type AppError struct {
      Code    int    `json:"code"`
      Message string `json:"message"`
  }
  func (e *AppError) Error() string { return e.Message }
  ```
- **Constructor helpers:** `errors.NewNotFound("user")`, `errors.NewConflict("email already exists")`
- **In service:** `return nil, errors.NewNotFound("user")` — services return `(T, error)`
- **In Fiber error handler:**
  ```go
  app.Use(func(c *fiber.Ctx, err error) error {
      var appErr *errors.AppError
      if ok := errors.As(err, &appErr); ok {
          return c.Status(appErr.Code).JSON(appErr)
      }
      return c.Status(500).JSON(errors.NewInternal())
  })
  ```

## Avoid
- Returning raw `errors.New("...")` from services — always use typed AppErrors
- Sending stack traces or internal error details to the client
- Using `http.StatusInternalServerError` for known, expected errors
