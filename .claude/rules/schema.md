# Rule: Schema and migrations

Applies to `internal/database/postgres/migrations/`, `internal/models/`, and
`internal/database/postgres/seed/`.

## Migrations

- golang-migrate, hand-written SQL pairs: `NNNNNN_<name>.up.sql` and `.down.sql`, sequential
  (`make db-migrate-create NAME=<name>`). They are embedded into the server binary and applied on
  boot, and by `go run ./cmd/migrate` / `make db-migrate`.
- **Never edit an applied migration.** A change is a new pair. The down file reverses the up file.
- A migration that adds a NOT NULL column to a token table deletes the pre-existing rows rather than
  inventing a default — an old token row has no meaningful value for the new column.
- The GORM model in `internal/models/` mirrors the SQL (column names, nullability). GORM's
  `AutoMigrate` is never used.
- `000003_sync_with_base.down.sql` deliberately **fails** when a soft-deleted user shares an email
  with a live user: it cannot restore the old full `users_email_key` unique over duplicates, and a
  down migration never deletes user rows to make one fit. Resolve the duplicates by hand first. The
  down is also lossy by design (roles stay renamed, tokens are cleared).

## Naming

snake_case plural tables, snake_case columns, `uuid` primary keys defaulting to `gen_random_uuid()`,
`created_at` / `updated_at` `timestamp NOT NULL DEFAULT now()` on entity tables. Foreign keys state
`ON DELETE CASCADE` explicitly. Indexes are named `<table>_<cols>_index`, unique constraints
`<table>_<col>_unique`.

## Tables

| Table | Notes |
| --- | --- |
| `users` | `name`, `email`, `password` NOT NULL; `status user_status NOT NULL DEFAULT 'active'`; `remark varchar(255)` nullable; `email_verified_at`, `deleted_at` nullable. **No plain UNIQUE on email**: partial unique `users_email_unique_live (email) WHERE deleted_at IS NULL`, plus `users_email_deleted_at_status_index`; partial `users_created_at_live_index (created_at) WHERE deleted_at IS NULL` backs the default list sort |
| `roles` | `name varchar(100)` UNIQUE (`roles_name_unique`); hard delete |
| `permissions` | `name varchar(255)` UNIQUE (`permissions_name_unique`), stored as `"<group> <action>"`; `"group" varchar(100)`; hard delete |
| `role_permissions` | PK `(role_id, permission_id)`, both FKs CASCADE; `role_permissions_permission_id_index` for holder lookups by permission |
| `user_roles` | PK `(user_id, role_id)`, both FKs CASCADE, `assigned_at timestamp NOT NULL DEFAULT now()`; `user_roles_role_id_index` for holder lookups by role |
| `email_verifications`, `password_reset_tokens` | `user_id` FK CASCADE; `token varchar(255)` UNIQUE (SHA-256 hex of the emailed token); `expired_at` NOT NULL; `used_at` nullable; index `(user_id, used_at)` |

## Enum

`user_status`: `active`, `inactive`, `suspended`, `blocked` — lowercase on the wire. Go mirrors it
with typed constants in `internal/models`; the validator's `oneof` and the filter allow-list use the
same four values.

## Seed

`go run ./cmd/seed` / `make db-seed`. Must be idempotent (a second run changes nothing and exits 0)
and exit non-zero on failure. Order: permissions → roles + grants → users + role assignment.

- 15 permissions: `user`, `role`, `permission` × `list`, `create`, `detail`, `edit`, `delete`.
- Roles `superuser` and `admin`, both granted the full catalogue. The grant reads permissions scoped
  to the catalogue names, so an operator-added permission is never auto-granted.
- Users `superuser@example.com` and `admin@example.com` (password `S3crEtP4ssw0rd!`), active,
  pre-verified, each holding the role of the same name. Lookups are live-only; an existing user is
  skipped.
- Use `FirstOrCreate` / `ON CONFLICT DO NOTHING`; never a bare `Create`.
