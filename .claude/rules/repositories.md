# Rule: Repositories

Applies to `internal/repositories/`.

## Shape

- The interface lives in `internal/repositories/interfaces.go`; the GORM implementation in
  `internal/repositories/<name>_repository.go` (`user`, `role`, `permission`, `email_verification`,
  `password_reset_token`), built by `New<Name>Repository(db)`.
- Every repository has `WithTx(tx *gorm.DB) <Name>Repository` so a service can run it inside its
  transaction. Repositories never begin, commit, or roll back.
- Every method takes `ctx` and calls `r.db.WithContext(ctx)`.

## Not found is `(nil, nil)`

A lookup that finds nothing returns `(nil, nil)` — translate `gorm.ErrRecordNotFound` inside the
repository. Repositories never return `*errors.AppError`; the service owns 404/422
([services-crud.md](./services-crud.md)).

## Soft delete

`users` is soft-deleted (`gorm.DeletedAt`, GORM adds `deleted_at IS NULL` to every query on the
model). Do not call `Unscoped()` on users outside the seed. Raw SQL or `Joins` touching `users` must
filter `users.deleted_at IS NULL` explicitly. Roles, permissions, and token rows are hard-deleted.

## Lists: allow-lists, one source of truth

Each list repository exports its allow-list once as a `pagination.Options` value —
`repositories.UserListOptions`, `RoleListOptions`, `PermissionListOptions` — and applies the parsed
query with `q.Apply(db, opts, tz)`:

| Resource | Sortable (API key → column) | Filterable | Search columns |
| --- | --- | --- | --- |
| users | `id`, `name`, `email`, `status`, `createdAt`, `updatedAt` | `status` (enum), `name`, `email`, `roleId` (id), `createdAt` (date), `updatedAt` (date) | `name`, `email`, `status::text` |
| roles | `id`, `name`, `createdAt`, `updatedAt` | `name`, `createdAt` (date), `updatedAt` (date) | `name` |
| permissions | `id`, `name`, `group`, `createdAt`, `updatedAt` | `name`, `group`, `createdAt` (date), `updatedAt` (date) | `name`, `"group"` |

- Sort keys are camelCase API names mapped to snake_case columns; a raw query value is never
  concatenated into `ORDER BY`.
- Unknown filter key, out-of-range enum value, or unparseable date → 400 from the pagination package;
  the service passes that error through unchanged.
- `perPage` is capped at 100.
- `"group"` is a reserved word — always quote it in SQL.

## Lookups the services need

`FindLiveByEmail(email, excludeID)` (users, deleted excluded), `FindByName(name, excludeID)` (roles),
`FindExistingNames(names, excludeID)` (permissions), `FindIDsByRoleIDs` / `FindIDsByPermissionID`
(cache invalidation), `FindUserInformation` (active and live only), `FindUsable(tokenHash, now)` (token tables,
unused and unexpired only).
