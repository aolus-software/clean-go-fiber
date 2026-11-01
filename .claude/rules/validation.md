# Rule: Validation

Applies to `internal/modules/<name>/schema.go`, `internal/libs/validator`, and
`internal/libs/pagination`.

## Request structs

- One named struct per request body, in the module's `schema.go`, with `json` tags in **camelCase**
  (`newPassword`, `roleIds`, `permissionIds`, `refreshToken`) and `validate` tags from
  go-playground/validator.
- Every POST/PATCH handler calls `utils.ParseBody(c, &req)`, which parses and validates in one step,
  before calling the service.
- Validation errors are 422 with `errors: [{ field, message }]`, where `field` is the **json name**
  (the validator is configured to report json tag names), never the Go field name.

Common tags:

| Need | Tag |
| --- | --- |
| required string, length | `validate:"required,min=1,max=255"` |
| email | `validate:"required,email,max=255"` |
| strong password | `validate:"required,min=8,strongpassword"` — lower, upper, digit, symbol, 8+ |
| status enum | `validate:"required,oneof=active inactive suspended blocked"` |
| uuid array | `validate:"required,dive,uuid"` |
| nullable string, key required | pointer field (`*string`) with `validate:"omitempty,max=255"` |

Register a custom tag once in `internal/libs/validator`; never re-implement a check in a service.

## Path ids

`:id` is parsed with `utils.ParamUUID(c, "id")` in the handler. Malformed → 422 `errors: [{ field: "id", ... }]`. The
service receives a `uuid.UUID`, never a raw string.

## List query

Parsed by `pagination.Parse(c)` and applied by the repository with `q.Apply(db, opts, tz)` against
its exported allow-list (`repositories.UserListOptions`, `RoleListOptions`, `PermissionListOptions`):

| Param | Rule | Invalid → |
| --- | --- | --- |
| `page` | integer ≥ 1, default 1 | 422 |
| `perPage` | integer 1..100, default 10 | 422 |
| `search` | free text, `ILIKE %s%` over the resource's search columns | — |
| `sort` | camelCase key from the allow-list, default `createdAt` | 422 |
| `sortDirection` | `asc` \| `desc`, default `desc` | 422 |
| `filter[<key>]=<value>` | flat query params, one per key | unknown key 400; bad enum 400; bad date 400 |

Filter semantics by kind: enum and id — comma means several values (`IN`); date — `start,end`
inclusive whole days in `APP_TIMEZONE`, a single date is that whole day; plain — `ILIKE %value%`.

## What never appears in a response

Password hashes, token values (except the JWTs returned by login/refresh), `deletedAt`.
