# Clean Go Fiber - TODO List

> Updated: 2026-10-31 (synced with the base contract: routes under `/settings`, `"<group> <action>"` permissions, E envelope, E env names)

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
  - App: `APP_NAME`, `APP_PORT`, `APP_ENV`, `APP_TIMEZONE`, `LOG_LEVEL`, `ENABLE_API_DOCS`, `CLIENT_URL`
  - JWT: `JWT_SECRET`, `JWT_REFRESH_SECRET` (required; lifetimes hardcoded 1d / 7d)
  - Database: `DATABASE_URL`
  - Redis: `REDIS_HOST`, `REDIS_PORT`, `REDIS_PASSWORD`, `REDIS_DB`
  - Mail: `MAIL_HOST`, `MAIL_PORT`, `MAIL_SECURE`, `MAIL_USER`, `MAIL_PASS`, `MAIL_FROM`
  - CORS: `ALLOWED_HOST`

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
  - `POST /auth/login`
  - `POST /auth/register`
  - `POST /auth/verify-email`
  - `POST /auth/resend-verification`
  - `POST /auth/forgot-password`
  - `POST /auth/reset-password`
  - `POST /auth/refresh`

- [x] **Implement auth service**
  - Password hashing with bcrypt
  - JWT access + refresh token generation
  - Email verification logic
  - Password reset logic

- [x] **Implement JWT middleware**
  - Extract and validate Bearer token
  - Attach user info to Fiber context
  - Handle token expiry errors

- [x] **Implement RBAC guards**
  - `RoleGuard` — check user role
  - `PermissionGuard` — check user permission
  - Superuser bypass

---

## 🟡 Medium Priority

### Profile & Settings Modules

- [x] **Profile module**
  - `GET /profile` — get current user
  - `PATCH /profile` — update name and email

- [x] **Permissions module** (`/settings/permissions`)
  - list, create, detail, update, delete

- [x] **Roles module** (`/settings/roles`)
  - list, create, detail, update, delete (`permissionIds` in the body)

- [x] **Users module** (`/settings/users`)
  - list, create, detail, update, delete (`roleIds` in the body)
  - `reset-password` (superuser), `send-verification-email`, `send-reset-password-email`

- [x] **Select options** (`/settings/select-options/permissions`, `/roles`)

### API Documentation

- [x] **Set up Swagger/OpenAPI with swag**
  - Install `swaggo/swag` CLI
  - Add doc annotations to all handlers
  - Serve Scalar UI at `/docs`, gated by `ENABLE_API_DOCS`
  - Document request/response schemas
  - Document Bearer auth security scheme

### Middleware

- [x] **Request ID middleware** — attach unique ID to each request
- [x] **Structured logging middleware** — log method, path, status, duration (zap)
- [x] **Performance middleware** — warn on slow requests (>1s)
- [x] **CORS middleware** — configurable origins
- [x] **Rate limiting middleware** — global rate limit (memory-backed, 100 req/min)
- [x] **Security headers middleware** — equivalent of Helmet
- [x] **Body limit middleware** — max 100KB request body
- [x] **Recovery middleware** — panic recovery with error logging
- [x] **Error handler** — centralized error response formatting

### Caching Layer

- [x] **Redis cache wrapper**
  - `Get(key)`, `Set(key, value, ttl)`, `Delete(key)`, `Flush()`
  - `Remember(key, ttl, fn)` — cache-or-compute helper
  - Key generation helper (`generateKey(...parts)`)
  - Cache constants file

### Background Jobs (Asynq)

- [x] **Set up Asynq client + server**
  - `cmd/worker/main.go` — separate worker process
  - Queue client for enqueuing tasks
  - Worker server for processing tasks

- [x] **Send email task**
  - Task type: `email:send`
  - Payload: to, subject, body, template
  - Handler with retry logic and error logging

### Email (Mailer)

- [x] **SMTP mailer setup**
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

- [x] **Create Makefile**
  - `make dev` / `dev-server` / `dev-worker`, `build`, `start`, `lint`, `format`, `test`, `docs`
  - `make db-migrate` / `db-migrate-down` / `db-migrate-create` / `db-seed` / `fresh` / `reset`
  - `make docker-*` suite

- [x] **Hot reload with Air**
  - Install `cosmtrek/air`
  - Create `.air.toml` config
  - Separate configs for server and worker

### Docker & Deployment

- [x] **Create Dockerfile**
  - Multi-stage build (builder + runner)
  - Alpine base image
  - Build Go binary
  - Non-root user

- [x] **Create docker-compose.yml**
  - `postgres` service
  - `redis` service
  - `app` service (depends on postgres + redis)
  - Health checks for all services
  - `app_network` bridge network
  - Named volumes

- [x] **Create `.dockerignore`**

### Code Quality

- [ ] **Set up golangci-lint**
  - `.golangci.yml` config
  - Rules: gofmt, govet, errcheck, staticcheck, gosec

- [x] **CI workflow** (`.github/workflows/ci.yml`: lint, build, database)

- [ ] **Set up pre-commit hooks**
  - Run lint + format before commit

### Testing

- [ ] **Add test infrastructure**
  - Unit tests for services (`_test.go` files)
  - Integration tests for repository layer
  - E2E tests for API endpoints (using `net/http/httptest`)
  - Test database setup (testcontainers-go)

### Documentation

- [x] **Write README.md**
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
| RBAC (Roles + Permissions) |      ✅      |     ✅     |       ✅       |
| OpenAPI / Swagger Docs     |      ✅      |     ✅     |       ✅       |
| Email Verification Flow    |      ✅      |     ✅     |       ✅       |
| Password Reset Flow        |      ✅      |     ✅     |       ✅       |
| Background Jobs (Queue)    |      ✅      |     ✅     |       ✅       |
| PostgreSQL (GORM/Drizzle)  |      ✅      |     ✅     |       ✅       |
| Redis Cache                |      ✅      |     ✅     |       ✅       |
| Rate Limiting              |      ✅      |     ✅     |       ✅       |
| Security Headers           |      ✅      |     ✅     |       ✅       |
| CORS                       |      ✅      |     ✅     |       ✅       |
| Body Limit                 |      ✅      |     ✅     |       ✅       |
| Request ID                 |      ✅      |     ✅     |       ✅       |
| Structured Logging         |      ✅      |     ✅     |       ✅       |
| Performance Monitoring     |      ✅      |     ✅     |       ✅       |
| DI Container               |      ✅      |     ✅     |  manual (constructors)  |
| Docker + docker-compose    |      ✅      |     ✅     |       ✅       |
| Makefile                   |      ✅      |     ✅     |       ✅       |
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
