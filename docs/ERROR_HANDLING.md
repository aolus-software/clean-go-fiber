# Error Handling

## How a failure becomes a response

1. A service returns an `*errors.AppError` (`internal/errors`) for anything the client should see.
2. The handler returns it unchanged.
3. The `ErrorHandler` in `fiber.Config` (`cmd/server/main.go`) maps it to the envelope with
   `utils.Error(...)`.

Anything that is not an `AppError` (a driver error, a panic recovered by the recover middleware, an
unexpected Fiber error) becomes a **500 with the fixed message "Internal server error"**. The real
error is logged with zap and never sent to the client.

## Constructors

```go
errors.NewBadRequestError("Invalid credentials", errors.Field("email", "Invalid email or password"))
errors.NewUnprocessableEntityError("Validation error", errors.Field("email", "Email already registered"))
errors.NewUnauthorizedError("Unauthorized")
errors.NewForbiddenError("Forbidden")
errors.NewNotFoundError("User not found")
```

400 and 422 carry `errors: [{ field, message }]`; the others carry `data: null`.

## Response shapes

```jsonc
// 400 / 422
{ "status": 422, "success": false, "message": "Validation error",
  "errors": [ { "field": "email", "message": "must be a valid email" },
              { "field": "password", "message": "must be at least 8 characters" } ] }

// 401 / 403 / 404 / 500
{ "status": 403, "success": false, "message": "Forbidden", "data": null }

// 429
{ "status": 429, "success": false, "message": "Too many requests" }
```

`errors` is always an **array** of `{ field, message }`. Validation fields use the request's json
names (`newPassword`, `roleIds`, `filter[status]`).

## Status mapping

| Source | Status |
| --- | --- |
| Body not parseable (`utils.ParseBody`) | 400, `field: "body"` |
| Validator failure | 422, one entry per field |
| Malformed `:id` (`utils.ParamUUID`) | 422, `field: "id"` |
| `page` / `perPage` / `sort` / `sortDirection` invalid | 422 |
| Unknown filter key, bad enum or date in a filter | 400 |
| Login failure | 400, `field: "email"` |
| Bad verification/reset token | 400, `field: "token"` |
| Uniqueness (service check, or `gorm.ErrDuplicatedKey` from a lost race) | 422 |
| Wrong current password on `PATCH /profile/password` | 422, `field: "currentPassword"` |
| Missing/invalid bearer, wrong token type, inactive/deleted user | 401 |
| Guard refusal, protected superuser target | 403 |
| Missing row, unmatched route | 404 |
| Rate limit | 429 |
| Anything else | 500 |
| `/health` with a failing probe | 503 (health payload, not an error envelope) |

## Writing new code

- Return, don't write: handlers never call `c.Status(...).JSON(...)` for an error.
- Translate third-party errors in the service; never put `err.Error()` into a message.
- Keep secrets, tokens, SQL, and stack traces out of `message` and `errors[]`.
