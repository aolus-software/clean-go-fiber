# Rule: The canonical CRUD module

The five routes every CRUD resource exposes, with the statuses each can return. `users`, `roles`, and
`permissions` under `/settings` follow it.

| Route | Guard | Success | Body / data | Can also return |
| --- | --- | --- | --- | --- |
| `GET /settings/<res>` | `"<group> list"` | 200 | paginated: `data: { data: [...], meta: { page, limit, totalCount } }` | 400 (filter), 401, 403, 422 (page/perPage/sort/sortDirection), 500 |
| `POST /settings/<res>` | `"<group> create"` | **201** | `data: null` | 400, 401, 403, 422 (validation or uniqueness), 500 |
| `GET /settings/<res>/:id` | `"<group> detail"` | 200 | the detail object | 401, 403, 404, 422 (bad id), 500 |
| `PATCH /settings/<res>/:id` | `"<group> edit"` | 200 | `data: null` | 400, 401, 403, 404, 422, 500 |
| `DELETE /settings/<res>/:id` | `"<group> delete"` | 200 | `data: null` | 401, 403, 404, 422 (bad id), 500 |

Conventions:

- `PATCH`, never `PUT`. Updates are full-replace of the documented fields (all required).
- Create and update return `data: null`; the client re-reads the detail if it needs it.
- Delete returns **200 with a body**, never 204. Users are soft-deleted; roles and permissions are
  hard-deleted (join rows cascade).
- Relationship ids travel in the create/update body as camelCase arrays: `roleIds` on users,
  `permissionIds` on roles. There are no separate sync routes.
- `:id` is parsed as a UUID before the service is called; a malformed id is 422, a well-formed but
  missing id is 404.

## routes.go

```go
func (h *Handler) RegisterRoutes(settings fiber.Router, auth fiber.Handler) {
	r := settings.Group("/roles", auth)
	r.Get("", guards.Permission("role list"), h.List)
	r.Post("", guards.Permission("role create"), h.Create)
	r.Get("/:id", guards.Permission("role detail"), h.Detail)
	r.Patch("/:id", guards.Permission("role edit"), h.Update)
	r.Delete("/:id", guards.Permission("role delete"), h.Delete)
}
```

List and create use `""`, not `"/"`, so the path carries no trailing slash. The `auth` handler is
`middleware.Auth(userRepo, cache)`, built once in `internal/bootstrap` and passed to every module's
`RegisterRoutes(router, auth)`; it runs before the guards on the group.
