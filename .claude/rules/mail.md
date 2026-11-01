# Rule: Mail

Mail is **always queued**. A service enqueues an `email:send` task (`internal/queue/tasks`) and the
worker (`cmd/worker`) delivers it through `internal/libs/mailer` (gomail over SMTP). The API server
never opens an SMTP connection, so it boots without working `MAIL_*` settings.

## Links

| Mail | Link |
| --- | --- |
| verification | `${CLIENT_URL}/auth/verify-email?token=<raw token>` |
| password reset | `${CLIENT_URL}/auth/reset-password?token=<raw token>` |

The raw token appears only in the mail body; the database stores its SHA-256 hash.

## Rules

- Enqueue **after** the transaction that created the token commits, never inside it.
- An enqueue failure is logged with zap (`zap.Error(err)`, the user id — never the token or link)
  and does not fail a request whose data is committed.
- Every enqueue sets `asynq.MaxRetry`.
- Never log a mail body, a token, or a token-bearing URL.
- `MAIL_FROM`, `MAIL_HOST`, `MAIL_PORT`, `MAIL_SECURE`, `MAIL_USER`, `MAIL_PASS` are read in
  `internal/config/mail.go` and all have defaults.
