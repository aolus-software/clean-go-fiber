# Clean Go Fiber - TODO List

## 🔴 High Priority

### Project Setup

- [x] **Initialize Go module**
  - `go mod init clean-go-fiber`
  - Set up Go version (1.22+)
  - Add `.gitignore` for Go

- [x] **Install core dependencies**
  - `gofiber/fiber/v2` — HTTP framework
  - `gorm.io/gorm` + `gorm.io/driver/postgres` — ORM
  - `redis/go-redis/v9` — Redis client
  - `golang-jwt/jwt/v5` — JWT auth
  - `joho/godotenv` — Env loading
  - `go-playground/validator/v10` — Input validation
  - `uber-go/zap` — Structured logging
  - `hibiken/asynq` — Background job queue (BullMQ equivalent)
  - `gofiber/swagger` + `swaggo/swag` — OpenAPI/Swagger docs

- [x] **Set up project folder structure**
  - Clean architecture layout (see `steps.md`)
  - Separate `cmd/` for server and worker entrypoints

### Configuration Management

- [x] **Create `.env.example`**
  - App: `APP_NAME`, `APP_PORT`, `APP_ENV`, `APP_KEY`, `APP_URL`
  - JWT: `JWT_SECRET`, `JWT_EXPIRY`
  - Database: `DATABASE_URL`, `DB_MAX_OPEN`, `DB_MAX_IDLE`
  - Redis: `REDIS_HOST`, `REDIS_PORT`, `REDIS_PASSWORD`, `REDIS_DB`
  - Mail: `MAIL_HOST`, `MAIL_PORT`, `MAIL_USER`, `MAIL_PASS`, `MAIL_FROM`
  - CORS: `CORS_ORIGIN`, `CORS_METHODS`

- [x] **Implement config loading with validation**
  - `internal/config/` package
  - Validate required env vars at startup (fail fast)
  - Typed config structs (AppConfig, DBConfig, RedisConfig, etc.)

### Database Layer

- [x] **Set up GORM + PostgreSQL connection**
  - Connection pool (MaxOpenConns, MaxIdleConns, ConnMaxLifetime)
  - Auto-migration support
  - Transaction helpers

- [x] **Define database schemas/models**
  - `users` table — id, email, password, name, is_verified, is_superuser
  - `email_verifications` table — id, user_id, token, expires_at
  - `password_reset_tokens` table — id, user_id, token, expires_at
  - `roles` table — id, name, slug, description
  - `permissions` table — id, name, slug, description
  - `role_permissions` table — pivot
  - `user_roles` table — pivot

- [x] **Create repository interfaces + implementations**
  - `UserRepository` — CRUD + find by email
  - `ForgotPasswordRepository` — token management
  - `RoleRepository` — role queries
  - `PermissionRepository` — permission queries

- [x] **Set up database migrations**
  - Migration runner script
  - Initial migration files

- [x] **Create database seeder**
  - Seed default roles: admin, user
  - Seed default permissions
  - Seed superuser account

### Authentication Module

- [x] **Implement auth routes + handlers**
  - `POST /api/auth/login`
  - `POST /api/auth/register`
  - `POST /api/auth/verify-email`
  - `POST /api/auth/resend-verification`
  - `POST /api/auth/forgot-password`
  - `POST /api/auth/reset-password`
  - `POST /api/auth/refresh-token`

- [x] **Implement auth service**
  - Password hashing with bcrypt
  - JWT access + refresh token generation
  - Email verification logic
  - Password reset logic

- [x] **Implement JWT middleware**
  - Extract and validate Bearer token
  - Attach user info to Fiber context
  - Handle token expiry errors

- [ ] **Implement RBAC guards**
  - `RoleGuard` — check user role
  - `PermissionGuard` — check user permission
  - Superuser bypass

---

## 🟡 Medium Priority

### Profile & Settings Modules

- [ ] **Profile module**
  - `GET /api/profile/me` — get current user
  - `PATCH /api/profile/me` — update profile
  - `PATCH /api/profile/me/password` — change password

- [ ] **Settings module**
  - `GET /api/settings/roles` — list roles
  - `POST /api/settings/roles` — create role
  - `GET /api/settings/roles/:id` — get role
  - `PATCH /api/settings/roles/:id` — update role
  - `DELETE /api/settings/roles/:id` — delete role
  - Same CRUD for permissions
  - Assign role/permission endpoints

### API Documentation

- [x] **Set up Swagger/OpenAPI with swag**
  - Install `swaggo/swag` CLI
  - Add doc annotations to all handlers
  - Serve Scalar-style UI or Swagger UI at `/docs`
  - Document request/response schemas
  - Document Bearer auth security scheme

### Middleware

- [ ] **Request ID middleware** — attach unique ID to each request
- [ ] **Structured logging middleware** — log method, path, status, duration (zap)
- [ ] **Performance middleware** — warn on slow requests (>1s)
- [ ] **CORS middleware** — configurable origins
- [ ] **Rate limiting middleware** — global rate limit (Redis-backed)
- [ ] **Security headers middleware** — equivalent of Helmet
- [ ] **Body limit middleware** — max 100KB request body
- [ ] **Recovery middleware** — panic recovery with error logging
- [ ] **Error handler** — centralized error response formatting

### Caching Layer

- [ ] **Redis cache wrapper**
  - `Get(key)`, `Set(key, value, ttl)`, `Delete(key)`, `Flush()`
  - `Remember(key, ttl, fn)` — cache-or-compute helper
  - Key generation helper (`generateKey(...parts)`)
  - Cache constants file

### Background Jobs (Asynq)

- [ ] **Set up Asynq client + server**
  - `cmd/worker/main.go` — separate worker process
  - Queue client for enqueuing tasks
  - Worker server for processing tasks

- [ ] **Send email task**
  - Task type: `email:send`
  - Payload: to, subject, body, template
  - Handler with retry logic and error logging

### Email (Mailer)

- [ ] **SMTP mailer setup**
  - `internal/libs/mailer/` package
  - Send function with template support
  - Email templates: verify-email, reset-password, welcome

### Dependency Injection

- [x] **Implement DI container**
  - Manual DI (constructor injection) or use `uber-go/fx`
  - Bootstrap function to wire all services
  - Inject repositories into services, services into handlers

### Custom Error Types

- [x] **Define error types**
  - `BadRequestError` (400)
  - `UnauthorizedError` (401)
  - `ForbiddenError` (403)
  - `NotFoundError` (404)
  - `ConflictError` (409)
  - `UnprocessableEntityError` (422)
  - `TooManyRequestsError` (429)
  - `InternalServerError` (500)

---

## 🟢 Low Priority

### Developer Experience

- [x] **Create Makefile** (partial — missing build, start, lint, format, test, fresh)
  - `make dev` — run with hot reload (air)
  - `make dev-worker` — run worker
  - `make seed` — seed database
  - `make migrate-*` — migration commands
  - `make docs` — generate swagger docs

- [x] **Hot reload with Air**
  - Install `cosmtrek/air`
  - Create `.air.toml` config
  - Separate configs for server and worker

### Docker & Deployment

- [ ] **Create Dockerfile**
  - Multi-stage build (builder + runner)
  - Alpine base image
  - Build Go binary
  - Non-root user

- [ ] **Create docker-compose.yml**
  - `postgres` service
  - `redis` service
  - `app` service (depends on postgres + redis)
  - Health checks for all services
  - `app_network` bridge network
  - Named volumes

- [ ] **Create `.dockerignore`**

### Code Quality

- [ ] **Set up golangci-lint**
  - `.golangci.yml` config
  - Rules: gofmt, govet, errcheck, staticcheck, gosec

- [ ] **Set up pre-commit hooks**
  - Run lint + format before commit

### Testing

- [ ] **Add test infrastructure**
  - Unit tests for services (`_test.go` files)
  - Integration tests for repository layer
  - E2E tests for API endpoints (using `net/http/httptest`)
  - Test database setup (testcontainers-go)

### Documentation

- [ ] **Write README.md**
  - Project overview and architecture
  - Quick start guide
  - Environment variables reference
  - API documentation link
  - Docker setup instructions
  - Development workflow

---

## 📊 Feature Parity Matrix (vs clean-elysia / clean-hono)

| Feature                    | clean-elysia | clean-hono | clean-go-fiber |
| -------------------------- | :----------: | :--------: | :------------: |
| JWT Authentication         |      ✅      |     ✅     |       ✅       |
| RBAC (Roles + Permissions) |      ✅      |     ✅     |       ⬜       |
| OpenAPI / Swagger Docs     |      ✅      |     ✅     |       ✅       |
| Email Verification Flow    |      ✅      |     ✅     |       ✅       |
| Password Reset Flow        |      ✅      |     ✅     |       ✅       |
| Background Jobs (Queue)    |      ✅      |     ✅     |       ⬜       |
| PostgreSQL (GORM/Drizzle)  |      ✅      |     ✅     |       ✅       |
| Redis Cache                |      ✅      |     ✅     |       ⬜       |
| Rate Limiting              |      ✅      |     ✅     |       ⬜       |
| Security Headers           |      ✅      |     ✅     |       ⬜       |
| CORS                       |      ✅      |     ✅     |       ⬜       |
| Body Limit                 |      ✅      |     ✅     |       ⬜       |
| Request ID                 |      ✅      |     ✅     |       ⬜       |
| Structured Logging         |      ✅      |     ✅     |       ⬜       |
| Performance Monitoring     |      ✅      |     ✅     |       ⬜       |
| DI Container               |      ✅      |     ✅     |       ✅       |
| Docker + docker-compose    |      ✅      |     ✅     |       ⬜       |
| Makefile                   |      ✅      |     ✅     |       🔶       |
| Hot Reload (dev)           |      ✅      |     ✅     |       ✅       |
| Env Validation             |      ✅      |     ✅     |       ✅       |
| Unit Tests                 |      ❌      |     ❌     |       ⬜       |

---

## 🎯 Recommended Implementation Order

1. **Phase 1 — Foundation**
   - Go module init, folder structure, config, GORM setup

2. **Phase 2 — Core Features**
   - Auth module (login, register, verify, reset)
   - JWT middleware, RBAC guards

3. **Phase 3 — Supporting Services**
   - Redis cache, Asynq jobs, Email mailer
   - Profile + Settings modules

4. **Phase 4 — API & Docs**
   - Swagger annotations + swag generate
   - OpenAPI schema documentation

5. **Phase 5 — Infrastructure**
   - Docker, docker-compose, Makefile, Air hot reload

6. **Phase 6 — Quality**
   - Tests, linting, CI/CD
