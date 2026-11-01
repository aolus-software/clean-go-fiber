# Rule: OpenAPI (swaggo)

The spec is **generated** from swag annotations by `make docs` into `docs/swagger.json`,
`docs/swagger.yaml`, `docs/docs.go`. Never hand-edit those files. The target passes every package
directory under `cmd/server` and `internal` to `swag init -d`; pointing swag at the repo root (which has
no Go files) breaks import-path resolution, so aliased imports and generic types fail to resolve.

## Serving

`GET /docs` (Scalar UI) and `GET /docs/openapi.json` are registered **only when
`ENABLE_API_DOCS=true`**. The flag defaults to `false`, independent of `APP_ENV`, so an environment
that never sets it cannot expose the schema. `.env.example` sets it to `true` for local work. The
Docker image ships `docs/swagger.json` so the flag works in a container.

## Every handler declares

```go
// @Summary     Update role
// @Description Replace a role's name and permissions. Requires 'role edit' permission.
// @Tags        Settings/Roles
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       id    path  string                  true  "Role ID (uuid)"
// @Param       body  body  roles.UpdateRoleRequest true  "Role"
// @Success     200   {object} utils.Response
// @Failure     400   {object} utils.ErrorResponse
// @Failure     401   {object} utils.Response
// @Failure     403   {object} utils.Response
// @Failure     404   {object} utils.Response
// @Failure     422   {object} utils.ErrorResponse
// @Router      /settings/roles/{id} [patch]
```

- `@Tags` mirror the module path: `General`, `Authentication`, `Profile`, `Settings/Users`,
  `Settings/Roles`, `Settings/Permissions`, `Settings/Select Options`.
- `@Security BearerAuth` only on protected routes; public routes omit it.
- List every status the route can return — 403 only if a guard runs, 404 only if an id can miss.
- List routes document `page`, `perPage`, `search`, `sort`, `sortDirection`, and the filter keys
  from the repository's allow-list.
- `@Description` states the permission or role required.

Response type names above illustrate the shape; reference the structs that exist in
`internal/libs/utils` and the module's `schema.go`.

## After changing a handler

Run `make docs` (needs the `swag` CLI: `go install github.com/swaggo/swag/cmd/swag@v1.16.6`) and
commit the regenerated files with the change. A spec that lags the code is an incomplete change.
