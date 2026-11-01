# Rule: Rate limiting

One global limiter, Fiber's `limiter` middleware, registered in `cmd/server/main.go`:
**100 requests per 60 seconds per client IP**. Over the limit the caller gets 429 with
`{ status: 429, success: false, message }` (no `data` key).

## Why the numbers are hardcoded

They are not env vars, matching the sibling `clean-elysia`. A configurable limit invites an
environment to loosen it for convenience and ship that to production. Changing them is a code
change, reviewed like one.

## Rules

- Credential endpoints (`/auth/login`, `/auth/register`, `/auth/forgot-password`,
  `/auth/reset-password`, `/auth/refresh`, `/auth/resend-verification`, `/auth/verify-email`) are
  never exempted or loosened. A stricter per-route limiter on them is welcome.
- The limiter is in-memory per process. Behind several replicas the effective limit multiplies —
  say so in deployment notes rather than pretending otherwise.
- The key is `c.IP()`. Behind a reverse proxy, configure Fiber's trusted proxy settings before
  trusting `X-Forwarded-For`; never key on a client-supplied header blindly.
