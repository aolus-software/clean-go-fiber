# Documentation Index

| Guide | Covers |
| --- | --- |
| [API_DOCUMENTATION.md](./API_DOCUMENTATION.md) | Envelope, authentication, the full route map, request bodies, list query parameters |
| [CONFIGURATION.md](./CONFIGURATION.md) | Every environment variable, its default, and who reads it |
| [DEPLOYMENT.md](./DEPLOYMENT.md) | Docker image, compose, migrations, seeding |
| [ERROR_HANDLING.md](./ERROR_HANDLING.md) | `AppError`, the central error handler, status mapping |
| [SECURITY.md](./SECURITY.md) | Tokens, guards, caching, rate limiting, CORS, headers |
| [MIDDLEWARE.md](./MIDDLEWARE.md) | The global middleware chain, `middleware.Auth`, guards |

Generated files (do not edit by hand; run `make docs`):

- `swagger.json`, `swagger.yaml`, `docs.go` — the OpenAPI spec produced by swag from handler
  annotations. `swagger.json` is served at `/docs/openapi.json` when `ENABLE_API_DOCS=true`.

`audit-findings.md` appears here once `/audit-flow` has been run.

Conventions for contributors live in [`../CLAUDE.md`](../CLAUDE.md) and
[`../.claude/rules/`](../.claude/rules/README.md).
