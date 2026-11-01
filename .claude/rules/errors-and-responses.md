# Rule: Errors and responses

The app speaks **one** success envelope and **one** error vocabulary. Handlers build successes with
`utils.OK(c, message, data)` and `utils.Created(c, message, nil)` (`internal/libs/utils/response.go`);
list data is `utils.PaginatedData[T]{Data, Meta{Page, Limit, TotalCount}}`. Failures are
**returned** as `*errors.AppError` and the central `ErrorHandler` in `cmd/server/main.go` turns them
into HTTP with `utils.Error(...)` — the only caller of that helper.

## Constructors (`internal/errors`)

```go
errors.NewBadRequestError("Invalid credentials", errors.Field("email", "Invalid email or password"))
errors.NewUnprocessableEntityError("Validation error", errors.Field("email", "Email already registered"))
errors.NewUnauthorizedError("Unauthorized")
errors.NewForbiddenError("Forbidden")
errors.NewNotFoundError("User not found")
```

## The envelope

`status` is inside the body **and** is the HTTP status.

```jsonc
// success
{ "status": 200, "success": true, "message": "…", "data": { } }

// list
{ "status": 200, "success": true, "message": "…",
  "data": { "data": [], "meta": { "page": 1, "limit": 10, "totalCount": 42 } } }

// 400 / 422 — field errors
{ "status": 422, "success": false, "message": "…", "errors": [ { "field": "email", "message": "…" } ] }

// 401 / 403 / 404 / 500 — no field detail
{ "status": 404, "success": false, "message": "…", "data": null }

// 429 — rate limiter, no data key
{ "status": 429, "success": false, "message": "…" }
```

- `meta` is exactly `{ page, limit, totalCount }` — the request param is `perPage`, the meta key is
  `limit`. No `totalPages`.
- `errors` is an **array** of `{ field, message }`, never an object of arrays.
- Pass `nil` data when there is nothing to return; the key stays (`"data": null`).
- Response fields are **camelCase** (`createdAt`, `emailVerifiedAt`, `isAssigned`) — models are
  mapped to response structs in the service; GORM models are never serialised directly.

## Status mapping

| Case | Status | Body |
| --- | --- | --- |
| Body not parseable | 400 | `errors[{ field: "body" }]` |
| Unknown filter key / bad filter value | 400 | `errors[{ field: "filter[<key>]" }]` |
| Login failure (unknown email, wrong password, unverified, not active) | 400 | `errors[{ field: "email" }]` |
| Invalid / expired / used verification or reset token | 400 | `errors[{ field: "token" }]` |
| Missing / invalid / expired bearer, wrong token type, user inactive or deleted | 401 | `data: null` |
| Guard refuses | 403 | `data: null` |
| `:id` not found | 404 | `data: null` |
| Unmatched route | 404 | `data: null` |
| Validation, malformed `:id`, bad page/perPage/sort/sortDirection | 422 | `errors[...]` |
| Uniqueness (email, role name, permission name) | 422 | `errors[...]` |
| Rate limit | 429 | no `data` |
| Anything else | 500 | `"Internal server error"`, logged with zap |
| `/health` with a failed probe | 503 | health payload |

## Rules

- **Never send `err.Error()` of an unknown error to the client.** The central handler logs it and
  answers 500 with a fixed message. Only an `*errors.AppError`'s message reaches a client.
- **Never hand-build the envelope** with `fiber.Map` in a handler.
- **Never `c.Status(4xx).JSON(...)` in a handler** — return the error.
- **Translate driver errors in the service**, never leak SQL or bcrypt text. GORM runs with
  `TranslateError: true`, so a unique violation surfaces as `gorm.ErrDuplicatedKey` and the service
  maps it to the same 422 as its own uniqueness check.
- **No secrets in `message` or `errors[]`.**
- **Deletes return 200 with a body**, never 204.
- Enumeration-leaky paths (forgot-password, resend-verification) return 200 rather than revealing
  whether an account exists.
