# clean-go-fiber

A production-ready Go REST API boilerplate following clean architecture principles. Built with Fiber v2, GORM, PostgreSQL, Redis, Asynq, JWT authentication, RBAC, and Swagger/OpenAPI docs.

---

## Tech Stack

| Layer            | Technology                        |
| ---------------- | --------------------------------- |
| HTTP Framework   | Go Fiber v2                       |
| ORM              | GORM                              |
| Database         | PostgreSQL 16                     |
| Migrations       | golang-migrate/migrate/v4         |
| Cache / Queue    | Redis 7 (go-redis/v9 + Asynq)     |
| Authentication   | golang-jwt/jwt/v5 (Bearer tokens) |
| Validation       | go-playground/validator/v10       |
| Logging          | uber-go/zap                       |
| API Docs         | swaggo/swag + Swagger UI          |
| Email            | gopkg.in/gomail.v2 (SMTP)         |
| Hot Reload       | cosmtrek/air                      |
| Containerization | Docker + docker-compose           |

---

## Prerequisites

- Go 1.23+
- Docker & docker-compose
- [`golang-migrate` CLI](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate) (for creating new migrations)

```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

---

## Getting Started

```bash
# 1. Clone and enter the project
git clone https://github.com/aolus-software/clean-go-fiber.git
cd clean-go-fiber

# 2. Copy environment file and fill in values
cp .env.example .env

# 3. Start infrastructure
docker compose up postgres redis -d

# 4. Install dependencies
go mod tidy

# 5. Run the server (migrations run automatically on startup)
make dev-server

# 6. (Optional) Start the background worker in a separate terminal
make dev-worker
```

Swagger UI is available at `http://localhost:3000/docs` once the server is running.

---

## Database Migrations

Migrations live in `internal/database/postgres/migrations/` as plain SQL files and are embedded into the binary at compile time. They run automatically when the server starts — it is safe to restart; already-applied migrations are skipped.

### Create a new migration

```bash
make migrate-create NAME=add_sessions_table
```

This generates a sequential pair:

```
internal/database/postgres/migrations/
  000002_add_sessions_table.up.sql
  000002_add_sessions_table.down.sql
```

Write your `CREATE TABLE` / `ALTER TABLE` SQL in the `.up.sql` file and the corresponding rollback SQL in `.down.sql`.

---

## Project Structure

```
clean-go-fiber/
├── cmd/
│   ├── server/main.go          # HTTP server entrypoint
│   └── worker/main.go          # Asynq worker entrypoint
├── internal/
│   ├── bootstrap/              # DI wiring
│   ├── config/                 # Env loading, typed config structs
│   ├── database/
│   │   ├── postgres/           # GORM connection + migrate runner
│   │   │   └── migrations/     # SQL migration files (embedded)
│   │   └── redis/              # go-redis client
│   ├── errors/                 # Custom AppError types
│   ├── guards/                 # RoleGuard, PermissionGuard middleware
│   ├── libs/
│   │   ├── cache/              # Redis cache wrapper
│   │   ├── mailer/             # SMTP mailer + HTML templates
│   │   ├── utils/              # hash, token, response, string helpers
│   │   └── validator/          # Validation wrapper
│   ├── middleware/             # Fiber middleware (auth, logger, cors)
│   ├── models/                 # GORM models — no business logic
│   ├── modules/                # Feature modules (auth, profile, users, roles, permissions)
│   ├── queue/                  # Asynq client + task definitions
│   ├── repositories/           # Interfaces + GORM implementations
│   └── worker/                 # Asynq server + task handlers
└── docs/                       # swag-generated OpenAPI files
```

---

## Architecture

All features follow a strict layering: **Handler → Service → Repository → Model**. No layer may skip another.

- **Handlers** — parse request, call service, return response. No business logic.
- **Services** — business logic only. Accept repository interfaces.
- **Repositories** — all database access. Return domain models.
- **Models** — GORM model structs. No methods or logic.

---

## API Endpoints

All routes (except auth) require a `Bearer` token in the `Authorization` header.

### Auth — `/auth`
| Method | Path | Description |
|--------|------|-------------|
| POST | `/auth/register` | Register a new account |
| POST | `/auth/login` | Login and get tokens |
| POST | `/auth/refresh-token` | Refresh access token |
| POST | `/auth/verify-email` | Verify email address |
| POST | `/auth/resend-verification` | Resend verification email |
| POST | `/auth/forgot-password` | Request password reset |
| POST | `/auth/reset-password` | Reset password with token |

### Profile — `/profile`
| Method | Path | Description |
|--------|------|-------------|
| GET | `/profile/me` | Get current user profile |
| PATCH | `/profile/me` | Update profile |
| PATCH | `/profile/me/password` | Change password |

### Users — `/users`
| Method | Path | Description |
|--------|------|-------------|
| GET | `/users` | List users |
| POST | `/users` | Create user |
| GET | `/users/:id` | Get user by ID |
| PATCH | `/users/:id` | Update user |
| DELETE | `/users/:id` | Delete user |
| POST | `/users/:id/roles` | Sync roles to user |

### Roles — `/roles`
| Method | Path | Description |
|--------|------|-------------|
| GET | `/roles` | List roles |
| POST | `/roles` | Create role |
| GET | `/roles/:id` | Get role by ID |
| PATCH | `/roles/:id` | Update role |
| DELETE | `/roles/:id` | Delete role |
| POST | `/roles/:id/permissions` | Sync permissions to role |

### Permissions — `/permissions`
| Method | Path | Description |
|--------|------|-------------|
| GET | `/permissions` | List permissions |
| POST | `/permissions` | Create permission |
| GET | `/permissions/:id` | Get permission by ID |
| PATCH | `/permissions/:id` | Update permission |
| DELETE | `/permissions/:id` | Delete permission |

---

## API Documentation

Docs are generated with [swaggo/swag](https://github.com/swaggo/swag). After modifying handlers run:

```bash
make docs
```
