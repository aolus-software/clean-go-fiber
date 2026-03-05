# AGENTS.md — clean-go-fiber

> **Read this file first before touching any code.**
> Then read `rules/general.md` and any relevant `skills/**/SKILL.md` for your task.
> Also we have some rules

---

## 1. Project Overview

`clean-go-fiber` is a production-ready Go REST API boilerplate built with Fiber v2, GORM, PostgreSQL, Redis, Asynq, JWT authentication, RBAC, and Swagger/OpenAPI docs. It follows clean architecture principles and serves as the Go equivalent of TypeScript boilerplates like `clean-elysia` or `clean-hono` — feature-parity, same patterns, different runtime.

---

## 2. Agent Orientation

Before writing any code:

1. **Read `rules/general.md`** — non-negotiable coding rules for this project
2. **Read the relevant `skills/*.md`** — pick whichever skills apply to your task
3. **Never invent a new folder** — respect the structure in `cmd/` and `internal/`
4. **Entrypoints:**
   - HTTP server → `cmd/server/main.go`
   - Background worker → `cmd/worker/main.go`
5. **New features always go in** → `internal/modules/<module>/`

---

## 3. Tech Stack

| Layer            | Technology                        |
| ---------------- | --------------------------------- |
| HTTP Framework   | Go Fiber v2                       |
| ORM              | GORM                              |
| Database         | PostgreSQL 16                     |
| Cache / Queue    | Redis 7 (go-redis/v9 + Asynq)     |
| Authentication   | golang-jwt/jwt/v5 (Bearer tokens) |
| Validation       | go-playground/validator/v10       |
| Logging          | uber-go/zap                       |
| API Docs         | swaggo/swag + Swagger UI          |
| Email            | gopkg.in/gomail.v2 (SMTP)         |
| Hot Reload       | cosmtrek/air                      |
| Linting          | golangci-lint                     |
| Containerization | Docker + docker-compose           |
| DI               | Manual constructor injection      |

---

## 4. Folder Structure

```
clean-go-fiber/
├── cmd/
│   ├── server/main.go          # HTTP server entrypoint
│   └── worker/main.go          # Asynq worker entrypoint
├── internal/
│   ├── bootstrap/              # DI wiring — wires all layers together
│   ├── config/                 # Env loading, typed config structs
│   ├── database/
│   │   ├── postgres/           # GORM connection, migrations, seeds
│   │   └── redis/              # go-redis client
│   ├── errors/                 # Custom AppError types
│   ├── guards/                 # RoleGuard, PermissionGuard middleware
│   ├── libs/
│   │   ├── cache/              # Redis cache wrapper
│   │   ├── mailer/             # SMTP mailer + HTML templates
│   │   ├── utils/              # hash, token, response, string helpers
│   │   └── validator/          # Validation wrapper
│   ├── middleware/             # Fiber middleware (auth, logger, cors, etc.)
│   ├── models/                 # GORM models only — no business logic
│   ├── modules/                # Feature modules (auth, profile, settings)
│   │   └── <module>/
│   │       ├── handler.go      # HTTP handlers (parse → call service → respond)
│   │       ├── routes.go       # Route registration
│   │       ├── schema.go       # Request/response structs
│   │       └── service.go      # Business logic
│   ├── queue/                  # Asynq client + task definitions
│   ├── repositories/           # Interfaces + GORM implementations
│   └── worker/                 # Asynq server + task handlers
├── docs/                       # swag-generated OpenAPI files (do not edit)
├── .air.toml
├── .env / .env.example
├── docker-compose.yml
├── Dockerfile
├── Makefile
└── go.mod
```

---

## 5. Key Architectural Decisions

| Decision            | Rule                                                                                      |
| ------------------- | ----------------------------------------------------------------------------------------- |
| **Data flow**       | Handler → Service → Repository → Model. Never skip layers.                                |
| **DI**              | Constructor injection only. Wired in `internal/bootstrap/bootstrap.go`.                   |
| **Config**          | Loaded once at startup. Accessed globally via `config.C`.                                 |
| **Async work**      | All async tasks (email, etc.) go through **Asynq**. Never inline in a request handler.    |
| **RBAC**            | `RoleGuard` and `PermissionGuard` are Fiber middleware. Applied per-route in `routes.go`. |
| **Superuser**       | `user.IsSuperuser == true` bypasses all permission checks in guards.                      |
| **Error responses** | All errors go through the centralized Fiber error handler. Never return raw Go errors.    |
| **Validation**      | Every POST/PATCH/PUT handler must validate input via the validator wrapper.               |

---

## 6. Development Workflow

```bash
# 1. Set up environment
cp .env.example .env

# 2. Start infrastructure
docker compose up postgres redis -d

# 3. Install dependencies
go mod tidy

# 4. Generate Swagger docs
make docs

# 5. Run migrations and seed
make migrate
make seed

# 6. Start dev server (hot reload)
make dev

# 7. Start worker (separate terminal)
make dev-worker

# 8. View API docs
open http://localhost:3000/docs
```

---

## 7. Common Task Patterns

### ➕ Adding a new API module

1. Create `internal/modules/<name>/schema.go` — request/response structs with validator tags
2. Create `internal/modules/<name>/service.go` — business logic, accepts repository interface
3. Create `internal/modules/<name>/handler.go` — parse request → call service → return response; add swaggo annotations
4. Create `internal/modules/<name>/routes.go` — register routes, apply guards
5. Add repository interface to `internal/repositories/interfaces.go`
6. Implement repository in `internal/repositories/<name>.go`
7. Wire everything in `internal/bootstrap/bootstrap.go`

### ⚙️ Adding a new background job

1. Define task constant + payload struct in `internal/queue/tasks/<name>.go`
2. Implement handler in `internal/worker/handlers/<name>.go`
3. Register handler in `internal/worker/server.go`
4. Enqueue from service using the Asynq client in `internal/queue/client.go`

### 🔒 Adding a new middleware

1. Create `internal/middleware/<name>.go`
2. Register with `app.Use(middleware.Name())` in `cmd/server/main.go` before routes

### 🗄️ Adding a new model

1. Define struct in `internal/models/<name>.go` with GORM tags
2. Add to the auto-migration list in `internal/database/postgres/migrate.go`
3. Add seed data (idempotent) in `internal/database/postgres/seed/seed.go`

---

## 8. File Reference Map

| What you need         | Where it lives                         |
| --------------------- | -------------------------------------- |
| Config structs        | `internal/config/`                     |
| GORM models           | `internal/models/`                     |
| Repository interfaces | `internal/repositories/interfaces.go`  |
| JWT helpers           | `internal/libs/utils/token.go`         |
| Password hashing      | `internal/libs/utils/hash.go`          |
| API response helpers  | `internal/libs/utils/response.go`      |
| Cache wrapper         | `internal/libs/cache/cache.go`         |
| Cache key constants   | `internal/libs/cache/const.go`         |
| Email sending         | `internal/libs/mailer/mailer.go`       |
| Input validation      | `internal/libs/validator/validator.go` |
| Custom error types    | `internal/errors/errors.go`            |
| DI wiring             | `internal/bootstrap/bootstrap.go`      |
| Task definitions      | `internal/queue/tasks/`                |
| Worker handlers       | `internal/worker/handlers/`            |

---

## 9. Coding Conventions

- **Do not add comments** unless explicitly asked.
- Never assume a library is available — check `go.mod` or existing imports first.
- When creating a new module, look at existing modules (e.g. `internal/modules/auth/`) to match naming, typing, and patterns.
- When editing code, read surrounding context and imports before making changes to stay idiomatic.
- Never introduce code that exposes or logs secrets/keys. Never commit credentials.

---

## 10. Out of Scope

Do **not** implement the following — they are outside this boilerplate's goals:

- GraphQL or gRPC
- WebSockets
- File upload / storage (S3, etc.)
- Multi-tenancy
- OAuth2 / social login
- Frontend / SSR

---

## 11. Agent-Specific Instructions

| Agent          | Config file                       | Notes                |
| -------------- | --------------------------------- | -------------------- |
| Claude         | `CLAUDE.md`                       | References this file |
| Gemini         | `.gemini/GEMINI.md`               | References this file |
| GitHub Copilot | `.github/copilot-instructions.md` | References this file |
| Qwen           | `QWEN.md`                         | References this file |

All agents share the **same rules and skills**. Differences are only in how each agent's config file is written — the actual coding standards are identical.
