# Security

How authentication, authorization, and request hardening work. To report a vulnerability, see
[`.github/SECURITY.md`](../.github/SECURITY.md).

## Tokens

| Token | Lifetime | Signed / stored | Notes |
| --- | --- | --- | --- |
| Access JWT | 1 day | HS256 with `JWT_SECRET` | claims `sub` (user id), `type: "access"` |
| Refresh JWT | 7 days | HS256 with `JWT_REFRESH_SECRET` | claims `sub`, `type: "refresh"`; accepted only by `/auth/refresh` |
| Email verification | 1 hour | SHA-256 hash in `email_verifications.token` (unique) | single use (`used_at`), revoked when a new one is issued |
| Password reset | 1 hour | SHA-256 hash in `password_reset_tokens.token` (unique) | same rules |

Both JWT secrets are required, have no defaults, and must differ. `middleware.Auth` rejects a token
whose `type` is not `access`; `/auth/refresh` rejects one whose `type` is not `refresh` and re-checks
that the user is active and not deleted. Raw verification and reset tokens exist only in the mail
link.

Each refresh token carries a random `jti`; the current one is stored in Redis at `refresh:<id>` for
the refresh lifetime. Login and `/auth/refresh` issue a new `jti`, so the previous refresh token is
rejected from then on (rotation). `cache.RevokeSessions` deletes `refresh:<id>` and `user:<id>` on
password change (`PATCH /profile/password`), password reset (`/auth/reset-password`), admin password
reset, user delete, and a user update that leaves the status anything but `active`.
Access tokens are not tracked: one already issued stays valid until it expires (1 day), but the auth
middleware re-reads `UserInformation` once `user:<id>` is dropped, so a deleted or suspended user is
refused on the next request.

## Login

The password is compared first — against a fixed dummy bcrypt hash when the email is unknown, so
response time does not reveal registration — and unknown email and wrong password share one message.
Only then are a verified email and `status = active` required. All failures are 400.

`/auth/forgot-password` and `/auth/resend-verification` always answer 200, so they cannot be used to
discover accounts.

## Authenticated requests

`middleware.Auth` loads `UserInformation` (`id, name, email, roles, permissions`) from Redis key
`user:<id>` (TTL 3600 s), falling back to the database with `status = 'active' AND deleted_at IS
NULL`. Services delete the key whenever a user's status, roles, identity, password, or deletion
changes, and for every holder when a role's permissions change or a role or permission is deleted.

## Authorization

- `guards.Permission(...)` and `guards.Role(...)` run per route after `middleware.Auth`; multiple
  arguments mean AND.
- The `superuser` role bypasses both. `roles.name` is UNIQUE, the role cannot be renamed or deleted,
  and only a superuser can assign it, remove it, or modify, delete, or reset the password of a user
  who holds it.
- Permission strings come from the seed: `user|role|permission` x `list|create|detail|edit|delete`.
- Admin password reset and the select-options endpoints require the `superuser` role.

## Passwords

bcrypt (default cost). New passwords must have 8+ characters with lower, upper, digit, and symbol.

## Request hardening

- Global rate limit: 100 requests / 60 s per client IP (in-memory, per process), 429 on excess.
- Helmet security headers, request ids, recover middleware.
- CORS origins from `ALLOWED_HOST` (empty → `*`; set it in any shared environment).
- Body size capped by Fiber's `BodyLimit`.
- Unknown errors return a fixed 500 message; details only in logs.
- `/docs` is served only when `ENABLE_API_DOCS=true`.

## Data

- Users are soft-deleted; email uniqueness applies to live users only (partial unique index), so a
  deleted user's email can register again.
- Passwords, token hashes, and `deletedAt` never appear in responses or logs.
