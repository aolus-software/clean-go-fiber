# Rule: Route structure and the live route map

## Convention: flat, resource-named, composed by group

No global prefix, no URL versioning. A route's full path is the concatenation of the Fiber groups it
passes through: the `/settings` group composes the module group `/roles` (which carries
`middleware.Auth()`), so the list route
is `GET /settings/roles`.

Resource segments are plural and kebab-case (`/users`, `/roles`, `/permissions`,
`/select-options`). Access is enforced by the per-route guard, never by the path — `/settings/**` is
not a security boundary ([rbac.md](./rbac.md)).

## Method conventions

| Operation | Verb + path | Status |
| --------- | ----------- | ------ |
| list      | `GET /<resource>` | 200 |
| create    | `POST /<resource>` | 201 |
| detail    | `GET /<resource>/:id` | 200 |
| update    | `PATCH /<resource>/:id` | 200 |
| delete    | `DELETE /<resource>/:id` | 200 |

**`PATCH`, never `PUT`.** A sub-resource action gets a trailing segment on the id
(`POST /settings/users/:id/reset-password`). List and create register `""`, not `"/"`.

## Current route map

"—" means no guard: legitimate only for public routes and routes that act on the caller's identity.

```
# health (internal/modules/health) — public
GET    /                                             —   public   { appName, appEnv, date }
GET    /health                                       —   public   200 healthy / 503 degraded
GET    /live                                         —   public   { alive: true }

# docs (internal/modules/docs) — public, only when ENABLE_API_DOCS=true
GET    /docs                                         —   public
GET    /docs/openapi.json                            —   public

# auth (internal/modules/auth, group /auth) — public
POST   /auth/login                                   —   { user, accessToken, refreshToken }
POST   /auth/refresh                                 —   { accessToken, refreshToken }
POST   /auth/register                                —   201, data null
POST   /auth/resend-verification                     —   200 always
POST   /auth/verify-email                            —
POST   /auth/forgot-password                         —   200 always
POST   /auth/reset-password                          —

# profile (internal/modules/profile, group /profile) — middleware.Auth()
GET    /profile                                      —   own identity
PATCH  /profile                                      —   own identity
PATCH  /profile/password                             —   own identity; { currentPassword, newPassword }, 422 on currentPassword, data null

# settings/users (group /settings/users) — middleware.Auth() on each settings module group
GET    /settings/users                               guards.Permission("user list")
POST   /settings/users                               guards.Permission("user create")
GET    /settings/users/:id                           guards.Permission("user detail")
PATCH  /settings/users/:id                           guards.Permission("user edit")
DELETE /settings/users/:id                           guards.Permission("user delete")
POST   /settings/users/:id/reset-password            guards.Role("superuser")
POST   /settings/users/:id/send-verification-email   guards.Permission("user create")
POST   /settings/users/:id/send-reset-password-email guards.Permission("user create")

# settings/roles (group /settings/roles)
GET    /settings/roles                               guards.Permission("role list")
POST   /settings/roles                               guards.Permission("role create")
GET    /settings/roles/:id                           guards.Permission("role detail")
PATCH  /settings/roles/:id                           guards.Permission("role edit")
DELETE /settings/roles/:id                           guards.Permission("role delete")

# settings/permissions (group /settings/permissions)
GET    /settings/permissions                         guards.Permission("permission list")
POST   /settings/permissions                         guards.Permission("permission create")
GET    /settings/permissions/:id                     guards.Permission("permission detail")
PATCH  /settings/permissions/:id                     guards.Permission("permission edit")
DELETE /settings/permissions/:id                     guards.Permission("permission delete")

# settings/select-options (internal/modules/selectoptions, group /settings/select-options)
GET    /settings/select-options/permissions          guards.Role("superuser")
GET    /settings/select-options/roles                guards.Role("superuser")
```

**Keep this map current.** Adding, renaming, or re-gating a route updates this table, `README.md`,
and `docs/API_DOCUMENTATION.md` in the same change.

## Known sharp edges

- **Guard strings are checked against the seed.** The update action is `"role edit"`, not
  `"role update"`. A string the seed never produces fails **closed**: only `superuser` can reach the
  route and an `admin` holding all 15 permissions gets 403.
- `send-verification-email` and `send-reset-password-email` reuse `"user create"` because no
  dedicated permission is seeded. Deliberate.
- `PATCH /settings/permissions/:id` stores `name` verbatim while `POST` prefixes the group per entry,
  so an update must send the full `"<group> <action>"` string.

## Registering a route

1. Add the handler with swag annotations ([handlers.md](./handlers.md), [openapi.md](./openapi.md)).
2. Register it in the module's `routes.go` with its guard.
3. Make sure `internal/bootstrap/bootstrap.go` calls the module's `RegisterRoutes` on the right group.
4. Add it to the map above.
