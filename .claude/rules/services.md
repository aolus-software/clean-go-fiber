# Rule: Services

Applies to `internal/modules/<name>/service.go`.

## Shape

A service is an exported interface plus an unexported struct built by `NewService(...)`, taking its
repositories, the cache, the queue client, the logger, and config as constructor arguments. Every
method takes `ctx context.Context` first.

## Owns

- **Business rules and the checks behind them** — existence (404), uniqueness (422), account state.
- **Transactions.** Any operation with more than one write runs inside one transaction opened by the
  service (`txmanager.WithinTransaction` / `db.Transaction`), passing `tx` to repositories via
  `repo.WithTx(tx)`. Repositories never open a transaction.
- **Cache invalidation.** A write that changes what `UserInformation` contains deletes the
  `user:<id>` key after commit:
  - user name, email, status, roles, password, or deletion → that user's key;
  - password change, password reset, admin password reset, deletion, or a status other than
    `active` → `cache.RevokeSessions`, which also drops the stored refresh `jti` (`refresh:<id>`);
  - role permissions, rename, or delete → every holder of the role;
  - permission edit or delete → every holder of a role that has it.
  A stale permission cache is a security defect, not a performance one.
- **Enqueueing** — after the transaction commits, never inside it. An enqueue failure is logged with
  zap and does not fail a request whose data was already committed (the user can ask for a resend).

## Errors

Return `*errors.AppError` for anything the client should see (400, 401, 403, 404, 422). Wrap any
other failure and let the central handler turn it into a 500 — never pass `err.Error()` of a driver,
bcrypt, or Redis error into a client-facing message. A repository's `(nil, nil)` "not found" is
turned into a 404 or 422 here, not in the repository.

## Enumeration-safe paths

`/auth/forgot-password` and `/auth/resend-verification` always succeed (200, `data: null`) whether or
not the account exists or is already verified. Login answers unknown email and wrong password with
the same message and compares against a dummy hash when the email is unknown, then requires a
verified email and `status = active`.

## One-time tokens

Verification and reset tokens: 32 random bytes hex, 1 h lifetime computed at issue time, stored as a
SHA-256 hash, single use via `used_at`. Issuing a new token stamps `used_at` on the user's earlier
unused tokens of the same kind (revoke-on-issue). Consuming one stamps every unused token of that
kind for the user, in the same transaction as the user update.
