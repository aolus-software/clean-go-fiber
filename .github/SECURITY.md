# Security Policy

## Supported versions

Only the latest commit on `main` is supported. This is a starter: once you generate a project from
it, keeping that project patched is yours.

## Reporting a vulnerability

Report privately through
[GitHub Security Advisories](https://github.com/aolus-software/clean-go-fiber/security/advisories/new).
**Do not open a public issue.** Include the affected route or package, steps to reproduce, and the
impact you observed. Expect an acknowledgement within a few days.

## Scope

In scope: anything in this repository — auth and token handling, RBAC guards, validation, the
response envelope, configuration defaults, the Docker image, and CI.

Out of scope: vulnerabilities in third-party modules with no exploitable path in this code (report
them upstream), and issues that require a misconfigured deployment the docs warn against (for
example `ENABLE_API_DOCS=true` in production).

## Known hardening notes

- `JWT_SECRET` and `JWT_REFRESH_SECRET` are required, have no defaults, and must differ. Access
  tokens live 1 day, refresh tokens 7 days; each carries a `type` claim that is checked.
- Every authenticated request re-checks `status = 'active' AND deleted_at IS NULL` through a cached
  `UserInformation` (Redis, 3600 s) that is invalidated on relevant writes.
- Verification and reset tokens are stored as SHA-256 hashes, expire after 1 hour, are single use,
  and are revoked when a new one is issued.
- Login uses the same message for unknown email and wrong password and compares against a dummy hash
  when the email is unknown.
- The global rate limit (100 requests / 60 s per IP) is in-memory per process.
- `/docs` is off unless `ENABLE_API_DOCS=true`.
- Seeded accounts use a published password. Change or remove them before any shared deployment.

See [`docs/SECURITY.md`](../docs/SECURITY.md) for the full reference.
