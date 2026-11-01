# API Documentation

Base URL in development: `http://localhost:3000`. No global prefix, no version segment. Requests and
responses are `application/json`. Interactive docs: `GET /docs` (spec at `GET /docs/openapi.json`),
available only when `ENABLE_API_DOCS=true`.

## Envelope

`status` is in the body and is also the HTTP status. Response fields are camelCase.

```jsonc
// success
{ "status": 200, "success": true, "message": "User retrieved successfully", "data": { } }

// list
{ "status": 200, "success": true, "message": "…",
  "data": { "data": [ ], "meta": { "page": 1, "limit": 10, "totalCount": 42 } } }

// 400 / 422 — field errors (an array, not an object)
{ "status": 422, "success": false, "message": "Validation error",
  "errors": [ { "field": "email", "message": "must be a valid email" } ] }

// 401 / 403 / 404 / 500
{ "status": 404, "success": false, "message": "User not found", "data": null }

// 429
{ "status": 429, "success": false, "message": "Too many requests" }
```

Deletes return 200 with `data: null`. Create endpoints return 201 with `data: null`.

## Authentication

1. `POST /auth/login` returns `{ user, accessToken, refreshToken }`.
2. Send `Authorization: Bearer <accessToken>` on protected routes.
3. When the access token expires (1 day), `POST /auth/refresh` with `{ "refreshToken": "…" }` returns a
   new `{ accessToken, refreshToken }` and invalidates the refresh token it was given. Refresh tokens
   last 7 days and are rejected as bearer tokens.

A user who is not `active`, or is deleted, gets 401 on every protected route once their cached
identity expires or is invalidated (writes that change it invalidate it immediately).

## Shapes

- `UserInformation` = `{ id, name, email, roles: string[], permissions: string[] }`
- `UserList` = `{ id, name, email, status, remark, roles: string[], createdAt, updatedAt }`
- `UserDetail` = `{ id, name, email, status, remark, emailVerifiedAt, roles: [{ id, name }], createdAt, updatedAt }`
- `RoleList` = `{ id, name, createdAt, updatedAt }`
- `RoleDetail` = `RoleList` + `permissions: [{ group, names: [{ id, name, isAssigned }] }]` — every
  permission in the system, grouped, flagged for this role
- `PermissionList` = `{ id, name, group, createdAt, updatedAt }`

`status` is one of `active`, `inactive`, `suspended`, `blocked`.

## Routes

| Method | Path | Access | Body | Success |
| --- | --- | --- | --- | --- |
| GET | `/` | public | — | 200 `{ appName, appEnv, date }` |
| GET | `/health` | public | — | 200 `{ status: "healthy", timestamp, services: { database, redis } }`; 503 if a probe fails |
| GET | `/live` | public | — | 200 `{ alive: true }` |
| POST | `/auth/login` | public | `email`, `password` | 200 `{ user: UserInformation, accessToken, refreshToken }` |
| POST | `/auth/refresh` | public | `refreshToken` | 200 `{ accessToken, refreshToken }` |
| POST | `/auth/register` | public | `email`, `name` (1-255), `password` (strong) | 201 `null` |
| POST | `/auth/resend-verification` | public | `email` | 200 `null` (always) |
| POST | `/auth/verify-email` | public | `token` | 200 `null` |
| POST | `/auth/forgot-password` | public | `email` | 200 `null` (always) |
| POST | `/auth/reset-password` | public | `token`, `newPassword` (strong) | 200 `null` |
| GET | `/profile` | authenticated | — | 200 `UserInformation` |
| PATCH | `/profile` | authenticated | `name`, `email` | 200 `UserInformation` |
| PATCH | `/profile/password` | authenticated | `currentPassword`, `newPassword` (strong) | 200 `null`; 422 on `currentPassword` if it is wrong |
| GET | `/settings/users` | `user list` | query | 200 paginated `UserList` |
| POST | `/settings/users` | `user create` | `name`, `email`, `password` (strong), `status`, `remark` (string or null), `roleIds` (uuid[]) | 201 `null` |
| GET | `/settings/users/:id` | `user detail` | — | 200 `UserDetail` |
| PATCH | `/settings/users/:id` | `user edit` | `name`, `email`, `status`, `remark`, `roleIds` (all required) | 200 `null` |
| DELETE | `/settings/users/:id` | `user delete` | — | 200 `null` (soft delete) |
| POST | `/settings/users/:id/reset-password` | role `superuser` | `newPassword` (strong) | 200 `null` |
| POST | `/settings/users/:id/send-verification-email` | `user create` | — | 200 `null`; 400 if already verified |
| POST | `/settings/users/:id/send-reset-password-email` | `user create` | — | 200 `null` |
| GET | `/settings/roles` | `role list` | query | 200 paginated `RoleList` |
| POST | `/settings/roles` | `role create` | `name`, `permissionIds` (uuid[]) | 201 `null` |
| GET | `/settings/roles/:id` | `role detail` | — | 200 `RoleDetail` |
| PATCH | `/settings/roles/:id` | `role edit` | `name`, `permissionIds` | 200 `null` |
| DELETE | `/settings/roles/:id` | `role delete` | — | 200 `null` |
| GET | `/settings/permissions` | `permission list` | query | 200 paginated `PermissionList` |
| POST | `/settings/permissions` | `permission create` | `name` (string[] — the actions), `group` (1-100); stores `"<group> <name>"` per entry | 201 `null` |
| GET | `/settings/permissions/:id` | `permission detail` | — | 200 `PermissionList` |
| PATCH | `/settings/permissions/:id` | `permission edit` | `name` (full `"<group> <action>"`, stored verbatim), `group` | 200 `null` |
| DELETE | `/settings/permissions/:id` | `permission delete` | — | 200 `null` |
| GET | `/settings/select-options/permissions` | role `superuser` | — | 200 `[{ group, permissions: [{ id, name, group }] }]` |
| GET | `/settings/select-options/roles` | role `superuser` | — | 200 `[{ id, name }]`, excluding `superuser`, ordered by name |

Strong password: at least 8 characters with a lowercase letter, an uppercase letter, a digit, and a
symbol. `superuser` bypasses every permission and role check. The `superuser` role cannot be renamed
or deleted, and only a superuser can assign it or act on a user who holds it.

Refresh tokens rotate: each `/auth/refresh` call invalidates the token it was given. A password
change, password reset, admin password reset, user delete, or suspension revokes the user's refresh
token, so the next `/auth/refresh` returns 401.

Mail links: verification `${CLIENT_URL}/auth/verify-email?token=…`, reset
`${CLIENT_URL}/auth/reset-password?token=…`. Tokens last 1 hour and are single use.

## List query parameters

| Param | Default | Rule | Invalid |
| --- | --- | --- | --- |
| `page` | 1 | integer ≥ 1 | 422 |
| `perPage` | 10 | integer 1..100 | 422 |
| `search` | — | case-insensitive substring over the resource's search columns | — |
| `sort` | `createdAt` | key from the allow-list below | 422 |
| `sortDirection` | `desc` | `asc` or `desc` | 422 |
| `filter[<key>]` | — | one flat query param per key | unknown key, bad enum, bad date: 400 |

Filter semantics: enum and id values accept a comma list (`IN`); dates accept `start,end` (inclusive
whole days in `APP_TIMEZONE`; a single date is that day); other keys match a substring.

| Resource | Sort | Filter | Search |
| --- | --- | --- | --- |
| users | `id`, `name`, `email`, `status`, `createdAt`, `updatedAt` | `status`, `name`, `email`, `roleId`, `createdAt`, `updatedAt` | `name`, `email`, `status` |
| roles | `id`, `name`, `createdAt`, `updatedAt` | `name`, `createdAt`, `updatedAt` | `name` |
| permissions | `id`, `name`, `group`, `createdAt`, `updatedAt` | `name`, `group`, `createdAt`, `updatedAt` | `name`, `group` |

Example: `GET /settings/users?page=2&perPage=20&sort=name&sortDirection=asc&filter[status]=active,inactive&filter[createdAt]=2026-01-01,2026-01-31`

## Status codes

| Case | Status |
| --- | --- |
| Body not parseable | 400 (`field: "body"`) |
| Unknown filter key / bad filter value | 400 |
| Login failure (unknown email, wrong password, unverified, not active) | 400 (`field: "email"`) |
| Invalid, expired, or used verification/reset token | 400 (`field: "token"`) |
| Missing / invalid / expired bearer, refresh token used as bearer, inactive or deleted user | 401 |
| Guard refuses | 403 |
| Id not found | 404 |
| Validation, malformed `:id`, bad pagination/sort | 422 |
| Duplicate email, role name, or permission name | 422 |
| Rate limit (100 / 60 s per IP) | 429 |
| Unexpected error | 500 |

## Example

```bash
curl -s -X POST localhost:3000/auth/login -H 'Content-Type: application/json' \
  -d '{"email":"admin@example.com","password":"S3crEtP4ssw0rd!"}'

curl -s localhost:3000/settings/users -H "Authorization: Bearer $ACCESS_TOKEN"
```
