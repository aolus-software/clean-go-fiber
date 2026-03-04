# Clean Go Fiber - Step-by-Step Implementation Guide

## Prerequisites

```bash
go 1.22+
docker & docker compose
golangci-lint
swag CLI (go install github.com/swaggo/swag/cmd/swag@latest)
air CLI (go install github.com/cosmtrek/air@latest)
```

---

## Step 1: Initialize Go Module

```bash
mkdir clean-go-fiber && cd clean-go-fiber
go mod init clean-go-fiber
```

### Install all dependencies at once:

```bash
go get github.com/gofiber/fiber/v2
go get github.com/gofiber/swagger
go get github.com/swaggo/swag
go get gorm.io/gorm
go get gorm.io/driver/postgres
go get github.com/redis/go-redis/v9
go get github.com/golang-jwt/jwt/v5
go get github.com/joho/godotenv
go get github.com/go-playground/validator/v10
go get go.uber.org/zap
go get github.com/hibiken/asynq
go get golang.org/x/crypto
go get github.com/google/uuid
go get gopkg.in/gomail.v2
```

---

## Step 2: Project Folder Structure

```
clean-go-fiber/
├── cmd/
│   ├── server/
│   │   └── main.go           # HTTP server entrypoint
│   └── worker/
│       └── main.go           # Asynq worker entrypoint
├── internal/
│   ├── bootstrap/
│   │   └── bootstrap.go      # DI wiring (wire all services)
│   ├── config/
│   │   ├── config.go         # Load & validate env vars
│   │   ├── app.go
│   │   ├── database.go
│   │   ├── redis.go
│   │   ├── jwt.go
│   │   ├── mail.go
│   │   └── cors.go
│   ├── database/
│   │   ├── postgres/
│   │   │   ├── postgres.go   # GORM connection
│   │   │   ├── migrate.go    # Auto-migration runner
│   │   │   └── seed/
│   │   │       └── seed.go
│   │   └── redis/
│   │       └── redis.go      # go-redis client
│   ├── errors/
│   │   └── errors.go         # Custom error types
│   ├── middleware/
│   │   ├── auth.go           # JWT auth middleware
│   │   ├── cors.go
│   │   ├── logger.go
│   │   ├── rate_limit.go
│   │   ├── request_id.go
│   │   ├── security.go       # Security headers
│   │   ├── body_limit.go
│   │   ├── performance.go
│   │   └── recovery.go
│   ├── guards/
│   │   ├── role.go           # Role guard
│   │   └── permission.go     # Permission guard
│   ├── libs/
│   │   ├── cache/
│   │   │   ├── cache.go      # Redis cache wrapper
│   │   │   └── const.go
│   │   ├── mailer/
│   │   │   ├── mailer.go
│   │   │   └── templates/
│   │   │       ├── verify_email.html
│   │   │       ├── reset_password.html
│   │   │       └── welcome.html
│   │   ├── utils/
│   │   │   ├── hash.go       # bcrypt helpers
│   │   │   ├── token.go      # JWT helpers
│   │   │   ├── response.go   # Standardized JSON response
│   │   │   └── string.go
│   │   └── validator/
│   │       └── validator.go  # go-playground/validator wrapper
│   ├── models/               # GORM models
│   │   ├── user.go
│   │   ├── email_verification.go
│   │   ├── password_reset_token.go
│   │   ├── role.go
│   │   └── permission.go
│   ├── repositories/
│   │   ├── interfaces.go     # Repository interfaces
│   │   ├── user.go
│   │   ├── forgot_password.go
│   │   ├── role.go
│   │   └── permission.go
│   ├── modules/
│   │   ├── auth/
│   │   │   ├── handler.go    # Fiber route handlers
│   │   │   ├── routes.go     # Route registration
│   │   │   ├── schema.go     # Request/response structs
│   │   │   └── service.go    # Business logic
│   │   ├── profile/
│   │   │   ├── handler.go
│   │   │   ├── routes.go
│   │   │   ├── schema.go
│   │   │   └── service.go
│   │   └── settings/
│   │       ├── handler.go
│   │       ├── routes.go
│   │       ├── schema.go
│   │       └── service.go
│   ├── queue/
│   │   ├── client.go         # Asynq client (enqueue)
│   │   └── tasks/
│   │       └── send_email.go # Task type + payload
│   └── worker/
│       ├── server.go         # Asynq server setup
│       └── handlers/
│           └── send_email.go # Task handler (process)
├── docs/                     # swag-generated OpenAPI files
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── .air.toml                 # Hot reload config
├── .env
├── .env.example
├── .gitignore
├── .golangci.yml
├── docker-compose.yml
├── Dockerfile
├── Makefile
├── README.md
└── go.mod
```

---

## Step 3: Configuration

### `internal/config/config.go`

```go
package config

import (
    "log"
    "os"
    "strconv"

    "github.com/joho/godotenv"
)

type Config struct {
    App      AppConfig
    Database DatabaseConfig
    Redis    RedisConfig
    JWT      JWTConfig
    Mail     MailConfig
    CORS     CORSConfig
}

var C *Config

func Load() {
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found, using environment variables")
    }

    C = &Config{
        App:      loadApp(),
        Database: loadDatabase(),
        Redis:    loadRedis(),
        JWT:      loadJWT(),
        Mail:     loadMail(),
        CORS:     loadCORS(),
    }

    validate(C)
}

func getEnv(key, fallback string) string {
    if v := os.Getenv(key); v != "" {
        return v
    }
    return fallback
}

func requireEnv(key string) string {
    v := os.Getenv(key)
    if v == "" {
        log.Fatalf("Required environment variable %s is not set", key)
    }
    return v
}

func getEnvInt(key string, fallback int) int {
    v := os.Getenv(key)
    if v == "" {
        return fallback
    }
    n, err := strconv.Atoi(v)
    if err != nil {
        return fallback
    }
    return n
}
```

### `internal/config/app.go`

```go
package config

type AppConfig struct {
    Name    string
    Port    string
    Env     string
    Key     string
    URL     string
}

func loadApp() AppConfig {
    return AppConfig{
        Name: getEnv("APP_NAME", "clean-go-fiber"),
        Port: getEnv("APP_PORT", "3000"),
        Env:  getEnv("APP_ENV", "development"),
        Key:  requireEnv("APP_KEY"),
        URL:  getEnv("APP_URL", "http://localhost:3000"),
    }
}
```

---

## Step 4: Database Models (GORM)

### `internal/models/user.go`

```go
package models

import (
    "time"
    "github.com/google/uuid"
    "gorm.io/gorm"
)

type User struct {
    ID           uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
    Email        string         `gorm:"uniqueIndex;not null"                           json:"email"`
    Password     string         `gorm:"not null"                                       json:"-"`
    Name         string         `gorm:"not null"                                       json:"name"`
    IsVerified   bool           `gorm:"default:false"                                  json:"is_verified"`
    IsSuperuser  bool           `gorm:"default:false"                                  json:"is_superuser"`
    Roles        []Role         `gorm:"many2many:user_roles"                           json:"roles,omitempty"`
    CreatedAt    time.Time                                                             `json:"created_at"`
    UpdatedAt    time.Time                                                             `json:"updated_at"`
    DeletedAt    gorm.DeletedAt `gorm:"index"                                          json:"-"`
}
```

### `internal/models/role.go`

```go
package models

type Role struct {
    ID          uint         `gorm:"primaryKey;autoIncrement" json:"id"`
    Name        string       `gorm:"uniqueIndex;not null"     json:"name"`
    Slug        string       `gorm:"uniqueIndex;not null"     json:"slug"`
    Description string                                        json:"description"`
    Permissions []Permission `gorm:"many2many:role_permissions" json:"permissions,omitempty"`
}

type Permission struct {
    ID          uint   `gorm:"primaryKey;autoIncrement" json:"id"`
    Name        string `gorm:"uniqueIndex;not null"     json:"name"`
    Slug        string `gorm:"uniqueIndex;not null"     json:"slug"`
    Description string                                json:"description"`
}
```

---

## Step 5: Repository Interfaces

### `internal/repositories/interfaces.go`

```go
package repositories

import (
    "clean-go-fiber/internal/models"
    "github.com/google/uuid"
)

type UserRepository interface {
    Create(user *models.User) error
    FindByID(id uuid.UUID) (*models.User, error)
    FindByEmail(email string) (*models.User, error)
    Update(user *models.User) error
    Delete(id uuid.UUID) error
    ExistsByEmail(email string) (bool, error)
}

type ForgotPasswordRepository interface {
    Create(token *models.PasswordResetToken) error
    FindByToken(token string) (*models.PasswordResetToken, error)
    DeleteByUserID(userID uuid.UUID) error
}

type RoleRepository interface {
    FindAll() ([]models.Role, error)
    FindByID(id uint) (*models.Role, error)
    FindBySlug(slug string) (*models.Role, error)
    Create(role *models.Role) error
    Update(role *models.Role) error
    Delete(id uint) error
}
```

---

## Step 6: Custom Error Types

### `internal/errors/errors.go`

```go
package errors

import "net/http"

type AppError struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
}

func (e *AppError) Error() string { return e.Message }

func NewBadRequestError(msg string) *AppError {
    return &AppError{Code: http.StatusBadRequest, Message: msg}
}

func NewUnauthorizedError(msg string) *AppError {
    return &AppError{Code: http.StatusUnauthorized, Message: msg}
}

func NewForbiddenError(msg string) *AppError {
    return &AppError{Code: http.StatusForbidden, Message: msg}
}

func NewNotFoundError(msg string) *AppError {
    return &AppError{Code: http.StatusNotFound, Message: msg}
}

func NewConflictError(msg string) *AppError {
    return &AppError{Code: http.StatusConflict, Message: msg}
}

func NewUnprocessableEntityError(msg string) *AppError {
    return &AppError{Code: http.StatusUnprocessableEntity, Message: msg}
}

func NewTooManyRequestsError(msg string) *AppError {
    return &AppError{Code: http.StatusTooManyRequests, Message: msg}
}

func NewInternalServerError(msg string) *AppError {
    return &AppError{Code: http.StatusInternalServerError, Message: msg}
}
```

---

## Step 7: Standard JSON Response Helper

### `internal/libs/utils/response.go`

```go
package utils

import "github.com/gofiber/fiber/v2"

type Response struct {
    Success bool        `json:"success"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
}

type PaginatedResponse struct {
    Success bool        `json:"success"`
    Message string      `json:"message"`
    Data    interface{} `json:"data"`
    Meta    PaginationMeta `json:"meta"`
}

type PaginationMeta struct {
    Total       int64 `json:"total"`
    Page        int   `json:"page"`
    PerPage     int   `json:"per_page"`
    TotalPages  int   `json:"total_pages"`
}

func OK(c *fiber.Ctx, message string, data interface{}) error {
    return c.Status(fiber.StatusOK).JSON(Response{
        Success: true,
        Message: message,
        Data:    data,
    })
}

func Created(c *fiber.Ctx, message string, data interface{}) error {
    return c.Status(fiber.StatusCreated).JSON(Response{
        Success: true,
        Message: message,
        Data:    data,
    })
}

func ErrorResponse(c *fiber.Ctx, code int, message string) error {
    return c.Status(code).JSON(fiber.Map{
        "success": false,
        "message": message,
    })
}
```

---

## Step 8: Auth Module

### `internal/modules/auth/schema.go`

```go
package auth

type LoginRequest struct {
    Email    string `json:"email"    validate:"required,email"`
    Password string `json:"password" validate:"required,min=8"`
}

type RegisterRequest struct {
    Name     string `json:"name"     validate:"required,min=2,max=100"`
    Email    string `json:"email"    validate:"required,email"`
    Password string `json:"password" validate:"required,min=8"`
}

type ForgotPasswordRequest struct {
    Email string `json:"email" validate:"required,email"`
}

type ResetPasswordRequest struct {
    Token    string `json:"token"    validate:"required"`
    Password string `json:"password" validate:"required,min=8"`
}

type TokenResponse struct {
    AccessToken  string `json:"access_token"`
    RefreshToken string `json:"refresh_token"`
    TokenType    string `json:"token_type"`
    ExpiresIn    int    `json:"expires_in"`
}
```

### `internal/modules/auth/routes.go`

```go
package auth

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(router fiber.Router, handler *Handler) {
    auth := router.Group("/auth")

    auth.Post("/login",               handler.Login)
    auth.Post("/register",            handler.Register)
    auth.Post("/verify-email",        handler.VerifyEmail)
    auth.Post("/resend-verification", handler.ResendVerification)
    auth.Post("/forgot-password",     handler.ForgotPassword)
    auth.Post("/reset-password",      handler.ResetPassword)
    auth.Post("/refresh-token",       handler.RefreshToken)
}
```

### `internal/modules/auth/handler.go`

```go
package auth

import (
    "clean-go-fiber/internal/libs/utils"
    "github.com/gofiber/fiber/v2"
)

type Handler struct {
    service *Service
}

func NewHandler(service *Service) *Handler {
    return &Handler{service: service}
}

// @Summary Login
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body LoginRequest true "Login credentials"
// @Success 200 {object} utils.Response
// @Router /api/auth/login [post]
func (h *Handler) Login(c *fiber.Ctx) error {
    var req LoginRequest
    if err := c.BodyParser(&req); err != nil {
        return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
    }

    if err := validate.Struct(req); err != nil {
        return utils.ErrorResponse(c, fiber.StatusUnprocessableEntity, err.Error())
    }

    result, err := h.service.Login(req)
    if err != nil {
        return utils.ErrorResponse(c, fiber.StatusUnauthorized, err.Error())
    }

    return utils.OK(c, "Login successful", result)
}
```

---

## Step 9: JWT Middleware

### `internal/middleware/auth.go`

```go
package middleware

import (
    "strings"
    "clean-go-fiber/internal/config"
    "clean-go-fiber/internal/libs/utils"
    "github.com/gofiber/fiber/v2"
    "github.com/golang-jwt/jwt/v5"
)

type UserClaims struct {
    UserID      string   `json:"user_id"`
    Email       string   `json:"email"`
    IsSuperuser bool     `json:"is_superuser"`
    Permissions []string `json:"permissions"`
    jwt.RegisteredClaims
}

func JWTAuth() fiber.Handler {
    return func(c *fiber.Ctx) error {
        auth := c.Get("Authorization")
        if !strings.HasPrefix(auth, "Bearer ") {
            return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Missing or invalid token")
        }

        tokenStr := strings.TrimPrefix(auth, "Bearer ")
        claims := &UserClaims{}

        token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
            return []byte(config.C.JWT.Secret), nil
        })

        if err != nil || !token.Valid {
            return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid or expired token")
        }

        c.Locals("user", claims)
        return c.Next()
    }
}
```

---

## Step 10: RBAC Guards

### `internal/guards/permission.go`

```go
package guards

import (
    "clean-go-fiber/internal/libs/utils"
    "clean-go-fiber/internal/middleware"
    "github.com/gofiber/fiber/v2"
)

func RequirePermission(permissions ...string) fiber.Handler {
    return func(c *fiber.Ctx) error {
        user := c.Locals("user").(*middleware.UserClaims)

        // Superuser bypasses all permission checks
        if user.IsSuperuser {
            return c.Next()
        }

        userPerms := make(map[string]bool)
        for _, p := range user.Permissions {
            userPerms[p] = true
        }

        for _, required := range permissions {
            if !userPerms[required] {
                return utils.ErrorResponse(c, fiber.StatusForbidden, "Insufficient permissions")
            }
        }

        return c.Next()
    }
}
```

---

## Step 11: Redis Cache Wrapper

### `internal/libs/cache/cache.go`

```go
package cache

import (
    "context"
    "encoding/json"
    "time"
    "clean-go-fiber/internal/database/redis"
)

const DefaultTTL = 3600 * time.Second

func Get(key string, dest interface{}) error {
    val, err := redis.Client.Get(context.Background(), key).Result()
    if err != nil {
        return err
    }
    return json.Unmarshal([]byte(val), dest)
}

func Set(key string, value interface{}, ttl time.Duration) error {
    data, err := json.Marshal(value)
    if err != nil {
        return err
    }
    return redis.Client.Set(context.Background(), key, data, ttl).Err()
}

func Delete(key string) error {
    return redis.Client.Del(context.Background(), key).Err()
}

func Remember(key string, ttl time.Duration, dest interface{}, fn func() (interface{}, error)) error {
    if err := Get(key, dest); err == nil {
        return nil // cache hit
    }
    result, err := fn()
    if err != nil {
        return err
    }
    _ = Set(key, result, ttl)
    data, _ := json.Marshal(result)
    return json.Unmarshal(data, dest)
}

func GenerateKey(parts ...string) string {
    key := ""
    for i, p := range parts {
        if i > 0 {
            key += ":"
        }
        key += p
    }
    return key
}
```

---

## Step 12: Background Jobs (Asynq)

### `internal/queue/tasks/send_email.go`

```go
package tasks

import (
    "encoding/json"
    "github.com/hibiken/asynq"
)

const TypeSendEmail = "email:send"

type SendEmailPayload struct {
    To       string `json:"to"`
    Subject  string `json:"subject"`
    Template string `json:"template"`
    Data     map[string]interface{} `json:"data"`
}

func NewSendEmailTask(payload SendEmailPayload) (*asynq.Task, error) {
    data, err := json.Marshal(payload)
    if err != nil {
        return nil, err
    }
    return asynq.NewTask(TypeSendEmail, data), nil
}
```

### `internal/worker/handlers/send_email.go`

```go
package handlers

import (
    "context"
    "encoding/json"
    "clean-go-fiber/internal/libs/mailer"
    "clean-go-fiber/internal/queue/tasks"
    "github.com/hibiken/asynq"
    "go.uber.org/zap"
)

func HandleSendEmailTask(logger *zap.Logger, m *mailer.Mailer) asynq.HandlerFunc {
    return func(ctx context.Context, t *asynq.Task) error {
        var payload tasks.SendEmailPayload
        if err := json.Unmarshal(t.Payload(), &payload); err != nil {
            return err
        }

        if err := m.Send(payload.To, payload.Subject, payload.Template, payload.Data); err != nil {
            logger.Error("Failed to send email", zap.Error(err), zap.String("to", payload.To))
            return err
        }

        logger.Info("Email sent successfully", zap.String("to", payload.To))
        return nil
    }
}
```

---

## Step 13: Server Entrypoint

### `cmd/server/main.go`

```go
package main

import (
    "log"
    "clean-go-fiber/internal/bootstrap"
    "clean-go-fiber/internal/config"
    "clean-go-fiber/internal/database/postgres"
    "clean-go-fiber/internal/database/redis"
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/recover"
    _ "clean-go-fiber/docs" // swag docs
)

func main() {
    // Load config
    config.Load()

    // Connect databases
    postgres.Connect()
    postgres.Migrate()
    redis.Connect()

    // Bootstrap DI
    container := bootstrap.NewContainer()

    // Create Fiber app
    app := fiber.New(fiber.Config{
        AppName:      config.C.App.Name,
        ErrorHandler: errorHandler,
    })

    // Global middleware
    app.Use(recover.New())
    // ... register middleware

    // Register routes
    api := app.Group("/api")
    container.RegisterRoutes(api)

    // Swagger docs at /docs
    // app.Get("/docs/*", swagger.HandlerDefault)

    log.Fatal(app.Listen(":" + config.C.App.Port))
}
```

### `cmd/worker/main.go`

```go
package main

import (
    "log"
    "clean-go-fiber/internal/config"
    "clean-go-fiber/internal/worker"
)

func main() {
    config.Load()
    srv := worker.NewServer()
    log.Fatal(srv.Run())
}
```

---

## Step 14: Makefile

```makefile
.PHONY: help dev dev-worker build start lint format test migrate seed fresh docs

help:
	@echo "Available commands:"
	@echo "  dev            - Run server with hot reload (air)"
	@echo "  dev-worker     - Run worker with hot reload"
	@echo "  build          - Build server binary"
	@echo "  start          - Run built binary"
	@echo "  lint           - Run golangci-lint"
	@echo "  format         - Format code (gofmt + goimports)"
	@echo "  test           - Run tests"
	@echo "  migrate        - Run database migrations"
	@echo "  seed           - Seed database"
	@echo "  fresh          - Drop + migrate + seed"
	@echo "  docs           - Generate swagger docs"

dev:
	air -c .air.toml

dev-worker:
	air -c .air.worker.toml

build:
	go build -o bin/server ./cmd/server

start:
	./bin/server

lint:
	golangci-lint run ./...

format:
	gofmt -w .
	goimports -w .

test:
	go test ./... -v -cover

migrate:
	go run ./cmd/server -migrate

seed:
	go run ./cmd/server -seed

fresh: migrate seed
	@echo "Database migrated and seeded!"

docs:
	swag init -g cmd/server/main.go -o docs
```

---

## Step 15: Air Hot Reload Config

### `.air.toml`

```toml
root = "."
tmp_dir = "tmp"

[build]
cmd = "go build -o ./tmp/main ./cmd/server"
bin = "./tmp/main"
include_ext = ["go", "tpl", "tmpl", "html"]
exclude_dir = ["tmp", "vendor", "docs", "node_modules"]
delay = 1000

[log]
time = true

[color]
main = "magenta"
watcher = "cyan"
build = "yellow"
runner = "green"

[misc]
clean_on_exit = true
```

---

## Step 16: Docker Setup

### `Dockerfile`

```dockerfile
# Stage 1: Builder
FROM golang:1.22-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o bin/server ./cmd/server

# Stage 2: Runner
FROM alpine:3.19

WORKDIR /app

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

COPY --from=builder /app/bin/server .
COPY --from=builder /app/.env.example .env

USER appuser

EXPOSE 3000

CMD ["./server"]
```

### `docker-compose.yml`

```yaml
services:
  postgres:
    image: postgres:16-alpine
    container_name: gofiber-postgres
    restart: unless-stopped
    environment:
      POSTGRES_DB: gofiber_db
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
    volumes:
      - postgres_data:/var/lib/postgresql/data
    ports:
      - "5432:5432"
    networks:
      - app_network
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 30s
      timeout: 10s
      retries: 3

  redis:
    image: redis:7-alpine
    container_name: gofiber-redis
    restart: unless-stopped
    command: redis-server --requirepass password123
    volumes:
      - redis_data:/data
    ports:
      - "6379:6379"
    networks:
      - app_network
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 30s
      timeout: 10s
      retries: 3

  app:
    build: .
    container_name: gofiber-app
    restart: unless-stopped
    env_file: .env
    ports:
      - "3000:3000"
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_healthy
    networks:
      - app_network

volumes:
  postgres_data:
  redis_data:

networks:
  app_network:
    driver: bridge
```

---

## Step 17: `.env.example`

```env
# App
APP_NAME=clean-go-fiber
APP_PORT=3000
APP_ENV=development
APP_KEY=change-this-to-a-random-secret-key
APP_URL=http://localhost:3000

# JWT
JWT_SECRET=change-this-jwt-secret
JWT_EXPIRY=3600

# Database
DATABASE_URL=postgres://postgres:postgres@localhost:5432/gofiber_db?sslmode=disable
DB_MAX_OPEN=10
DB_MAX_IDLE=2

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=password123
REDIS_DB=0

# Mail (SMTP)
MAIL_HOST=smtp.mailtrap.io
MAIL_PORT=587
MAIL_USER=your-mailtrap-user
MAIL_PASS=your-mailtrap-pass
MAIL_FROM=noreply@example.com
MAIL_FROM_NAME=Clean Go Fiber

# CORS
CORS_ORIGIN=http://localhost:5173
CORS_METHODS=GET,POST,PUT,PATCH,DELETE,OPTIONS

# Log
LOG_LEVEL=info
```

---

## Step 18: OpenAPI Swagger Annotations Example

Add to `cmd/server/main.go`:

```go
// @title           Clean Go Fiber API
// @version         1.0
// @description     Clean Architecture REST API built with Go Fiber
// @termsOfService  http://example.com/terms/

// @contact.name   API Support
// @contact.email  support@example.com

// @license.name  MIT

// @host      localhost:3000
// @BasePath  /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
```

Generate docs:
```bash
swag init -g cmd/server/main.go -o docs
```

Serve at `/docs`:
```go
import "github.com/gofiber/swagger"
app.Get("/docs/*", swagger.HandlerDefault)
```

---

## Go vs TypeScript Equivalents Reference

| TypeScript (Elysia/Hono)       | Go (Fiber)                          |
|--------------------------------|-------------------------------------|
| Bun / Node.js runtime          | Go runtime                          |
| Elysia / Hono / Fastify        | Go Fiber v2                         |
| Drizzle ORM                    | GORM                                |
| Zod / TypeBox validation       | go-playground/validator             |
| BullMQ (Redis queue)           | Asynq (Redis queue)                 |
| Pino / zap logger (JS)         | uber-go/zap                         |
| @hono/zod-openapi + Scalar     | swaggo/swag + Swagger UI            |
| envalid (env validation)       | requireEnv() fail-fast pattern      |
| tsyringe / manual DI           | Constructor injection (manual DI)   |
| Nodemailer                     | gopkg.in/gomail.v2                  |
| jsonwebtoken                   | golang-jwt/jwt                      |
| bcryptjs                       | golang.org/x/crypto/bcrypt          |
| `process.env.VAR`              | `os.Getenv("VAR")`                  |
| Husky + lint-staged            | pre-commit hooks with golangci-lint |
| `bun run dev`                  | `air` (hot reload)                  |
| `bun run build`                | `go build`                          |

---

## Quick Start (After Implementation)

```bash
# 1. Clone and enter
cd clean-go-fiber

# 2. Copy env
cp .env.example .env

# 3. Start infrastructure
docker compose up postgres redis -d

# 4. Install dependencies
go mod tidy

# 5. Generate swagger docs
make docs

# 6. Run migrations and seed
make migrate
make seed

# 7. Start dev server with hot reload
make dev

# 8. Open API docs
open http://localhost:3000/docs
```
