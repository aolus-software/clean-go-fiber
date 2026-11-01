# Deployment

The service ships as two processes from one image: the **server** (`./server`, HTTP) and the
**worker** (`./worker`, asynq consumer that sends mail). Scale them independently — each is one
process per container.

## Docker image

`Dockerfile` is multi-stage:

- **builder** — `golang:1.25-alpine`, builds `server`, `worker`, `seed`, and `migrate`.
- **migrator** — runs migrations and the seeder; used by `make docker-migrate` / `make docker-seed`.
- **runtime** — minimal Alpine, non-root user, the binaries plus `docs/swagger.json` (so the docs
  work when `ENABLE_API_DOCS=true`).

```bash
make docker-build
```

## Compose

`docker-compose.yml` runs `postgres`, `redis`, `app`, and `worker`. The app and worker read `.env`
and reach the databases by service name.

```bash
cp .env.example .env     # set JWT_SECRET and JWT_REFRESH_SECRET
make docker-up
make docker-migrate      # optional — the server also migrates on boot
make docker-seed         # idempotent
make docker-logs
make docker-down
```

## Migrations

Migrations are SQL pairs embedded into the binaries. The server applies pending ones on startup;
`go run ./cmd/migrate` (or the migrator image) applies them explicitly, and `go run ./cmd/migrate
down` rolls everything back. Run migrations once per release before scaling out servers, so several
replicas do not race at boot.

## Checklist for a shared environment

- Unique `JWT_SECRET` and `JWT_REFRESH_SECRET`, at least 32 random bytes each.
- `ENABLE_API_DOCS=false` unless the docs are meant to be public.
- `ALLOWED_HOST` set to your front-end origins.
- Real SMTP settings on the **worker**.
- Seeded accounts removed or their passwords changed.
- The rate limiter is per process; with several replicas, enforce limits at the proxy too.
