# Configuration

Configuration is read once at startup by `config.Load()` (`internal/config/`) from the process
environment, after loading `.env` if present (godotenv). Nothing else reads the environment. Copy
`.env.example` to `.env`; it works as-is apart from the two JWT secrets.

## Variables

| Variable | Default | Required | Notes |
| --- | --- | --- | --- |
| `APP_NAME` | app name | no | shown by `GET /` and the docs page title |
| `APP_PORT` | `3000` | no | HTTP port |
| `APP_ENV` | `development` | no | `development`, `staging`, or `production`. (The TypeScript siblings call this `NODE_ENV`.) |
| `APP_TIMEZONE` | `UTC` | no | used for date filters on list endpoints |
| `LOG_LEVEL` | `info` | no | `info`, `warn`, or `debug` |
| `ENABLE_API_DOCS` | `false` | no | the only switch for `/docs` and `/docs/openapi.json`, independent of `APP_ENV`. `.env.example` sets `true` for local work; keep it off elsewhere |
| `CLIENT_URL` | `http://localhost:3000` | no | base of the links in verification and reset mails |
| `DATABASE_URL` | local placeholder | no | `postgres://user:pass@host:5432/db?sslmode=disable`; used by the server, seed, migrate, and the GORM connection |
| `DB_POOL_MAX_OPEN` | `25` | no | max open connections per process; at least 1 |
| `DB_POOL_MAX_IDLE` | `10` | no | max idle connections kept in the pool; capped at `DB_POOL_MAX_OPEN` |
| `DB_POOL_CONN_MAX_LIFETIME` | `5m` | no | Go duration; a connection is closed after this age. `0` means no limit |
| `DB_POOL_CONN_MAX_IDLE_TIME` | `0` | no | Go duration; an idle connection is closed after this long. `0` means no limit |
| `DB_SLOW_QUERY_MS` | `200` | no | GORM logs a statement slower than this at warn level, with parameters as placeholders (never values); `0` disables. Query errors are logged too, except record-not-found. `APP_ENV=development` with `LOG_LEVEL=debug` still logs every statement |
| `REDIS_HOST` | `localhost` | no | cache and queue |
| `REDIS_PORT` | `6379` | no | |
| `REDIS_PASSWORD` | empty | no | empty means no auth |
| `REDIS_DB` | `0` | no | |
| `MAIL_HOST` | placeholder | no | worker only — the API server never connects to SMTP |
| `MAIL_PORT` | `587` | no | |
| `MAIL_SECURE` | `false` | no | `true` or `false` |
| `MAIL_FROM` | placeholder | no | |
| `MAIL_USER` | placeholder | no | |
| `MAIL_PASS` | placeholder | no | |
| `JWT_SECRET` | — | **yes** | signs access tokens (1 day) |
| `JWT_REFRESH_SECRET` | — | **yes** | signs refresh tokens (7 days); must differ from `JWT_SECRET` |
| `ALLOWED_HOST` | empty | no | CORS origins, comma-separated; empty allows `*` |
| `POSTGRES_DB`, `POSTGRES_USER`, `POSTGRES_PASSWORD` | — | compose only | read by the `postgres` service in `docker-compose.yml`; keep them consistent with `DATABASE_URL` |
| `POSTGRES_SLOW_QUERY_MS` | `100` | compose only | `log_min_duration_statement` of the compose `postgres` service: statements slower than this many ms go to the container log; `-1` disables |

Generate each secret with `openssl rand -hex 32`. A missing required variable stops the process at
boot with the variable's name in the message.

Hardcoded on purpose: token lifetimes (access 1 day, refresh 7 days, mail tokens 1 hour), the cache
TTLs (3600 s for `UserInformation` and the select-option lists), and the rate limit (100 requests / 60 s per IP).

## Compose hostnames

Inside `docker-compose.yml` the app and worker reach Postgres at `postgres` and Redis at `redis`;
compose overrides those hosts for its services. On the host machine `.env` uses `localhost`.

## Removed in the sync with the base contract

`APP_KEY`, `APP_URL`, `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `DB_SSLMODE`,
`DB_TIMEZONE`, `JWT_SECRET_KEY`, `JWT_REFRESH_SECRET_KEY`, `JWT_EXPIRES_IN`,
`JWT_REFRESH_EXPIRES_IN`, `CORS_ALLOWED_HOST`, `CORS_ALLOWED_METHOD`, `MAIL_PASSWORD`. An old `.env`
using them must be updated.
