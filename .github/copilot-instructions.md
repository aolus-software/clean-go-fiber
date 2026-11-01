# GitHub Copilot Instructions — clean-go-fiber

Go 1.25 REST API starter: Fiber v2, GORM + golang-migrate, PostgreSQL, Redis (go-redis + asynq), JWT,
go-playground/validator, zap, swaggo.

Read `CLAUDE.md` first, then the rules in `.claude/rules/` for the area you are changing. This file
mirrors their key conventions with examples; when the two drift, the rules win.

## Layering

`routes.go` → `handler.go` → `service.go` → `internal/repositories` → `internal/models`. Handlers
parse, validate, call one service method, and respond. Services own business rules, existence and
uniqueness checks (422), transactions, cache invalidation, and enqueueing. Repositories only query
and return `(nil, nil)` when nothing is found.

## Comments

None unless a block genuinely needs one, then a single block comment above it. swag annotations on
handlers are required.

Good:

```go
func (s *service) Create(ctx context.Context, req *CreateUserRequest) error {
	existing, err := s.userRepo.FindLiveByEmail(ctx, req.Email)
	if err != nil {
		return err
	}
	if existing != nil {
		return apperrors.NewUnprocessableEntityError("Email already registered", apperrors.Field("email", "Email already registered"))
	}
	return s.tx.WithinTransaction(ctx, func(tx *gorm.DB) error {
		return s.userRepo.WithTx(tx).Create(ctx, user, roleIDs)
	})
}
```

Bad: a comment on every line restating the code.

## Responses

Use the helpers in `internal/libs/utils/response.go`; return `*errors.AppError` for failures. The
envelope is `{ status, success, message, data }`; lists put `{ data, meta: { page, limit,
totalCount } }` inside `data`; 400/422 carry `errors: [{ field, message }]`. Never return
`err.Error()` of an unknown error to a client. Response fields are camelCase.

## Routes and RBAC

```go
r := settings.Group("/roles", auth)
r.Get("", guards.Permission("role list"), h.List)
r.Post("", guards.Permission("role create"), h.Create)
r.Patch("/:id", guards.Permission("role edit"), h.Update)
```

Permission strings are `"<group> <action>"` from the seed (`user|role|permission` x
`list|create|detail|edit|delete`). `PATCH`, never `PUT`. Privilege-granting routes use
`guards.Role("superuser")`.

## Queue

```go
task, err := tasks.NewEmailSendTask(payload)
if err != nil {
	return err
}
if _, err := s.queue.EnqueueContext(ctx, task, asynq.MaxRetry(3)); err != nil {
	s.log.Error("failed to enqueue email", zap.String("userId", userID), zap.Error(err))
}
```

Enqueue after the transaction commits. Never log tokens or mail bodies.

## Config

Read env only in `internal/config`; use `config.C` elsewhere.
