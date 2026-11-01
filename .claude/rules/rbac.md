# Rule: Authorization (RBAC)

Authentication and authorization are separate steps. `middleware.Auth()` proves **who** the caller
is; a guard decides **what** they may do. A route with the first and not the second is open to every
logged-in user. Nothing scans routes for missing guards — that is the failure mode this rule exists
to prevent.

## The two guards

Both live in `internal/guards/`, run after `middleware.Auth()`, and read the `UserInformation` it put
in `c.Locals` (no database query):

```go
guards.Permission("user list")          // requires every listed permission (AND)
guards.Role("superuser")                // requires every listed role (AND)
```

- Both **bypass for the `superuser` role**. `roles.name` is UNIQUE, so the bypass cannot be claimed
  by a second role with the same name, and the `superuser` role cannot be renamed or deleted.
- Both compare **full permission names** (`"user list"`), so holding `"user list"` never grants
  `"role list"`.
- Non-superusers cannot assign or remove the `superuser` role, or modify, delete, or reset the
  password of a user who holds it — the users service answers 403 even when the guard passed.
- A refusal is 403 with the error envelope. There is no OR — a route that needs "either" is a sign
  the catalogue is wrong.

## Permission names come from the seed

`internal/database/postgres/seed/` builds the catalogue as `"<group> <action>"` over groups `user`,
`role`, `permission` and actions `list`, `create`, `detail`, `edit`, `delete`: **15 permissions,
space-separated** — `"user list"`, not `user:list`. Seeded roles: `superuser`, `admin` (both hold
the full catalogue).

A guard naming a string the seed does not produce fails closed. Grep the seed before writing one;
extend the seed in the same change if a route genuinely needs a new permission.

## Privilege-granting routes are gated on the role

`POST /settings/users/:id/reset-password` and both `select-options` routes use
`guards.Role("superuser")`. Gating password reset on `"user edit"` would let anyone with an edit
permission take over any account, including a superuser's.

## Freshness

Guards read cached `UserInformation` (Redis `user:<id>`, 3600 s). Services invalidate it on every
write that changes a user's roles, status, or deletion, and for every holder when a role's
permissions change or a role or permission is deleted ([services.md](./services.md)). Without that,
a revoked permission keeps working for up to an hour.

## Checklist for a new protected route

- [ ] The route is under a group with `middleware.Auth()`.
- [ ] It has `guards.Permission(...)` or `guards.Role(...)`.
- [ ] Every permission string exists in the seed.
- [ ] Privilege-granting or account-takeover routes use `guards.Role("superuser")`.
- [ ] swag lists `@Failure 401` and `@Failure 403`, and `@Description` names the requirement.
- [ ] The route is in the map in [routes.md](./routes.md).

## Don't

- Don't rely on the URL prefix for authorization.
- Don't inline permission checks in a handler or service — use the guards.
- Don't soften a 403 into an empty list or a 200.
