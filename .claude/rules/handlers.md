# Rule: Handlers

Applies to `internal/modules/<name>/handler.go` and `routes.go`.

## What a handler may do

1. Parse and validate the body with `utils.ParseBody(c, &req)` (unparseable → 400
   `errors[{field:"body"}]`; invalid → 422 with per-field `errors`). Parse `:id` with
   `utils.ParamUUID(c, "id")` (malformed → 422 `errors[{field:"id"}]`). Lists use
   `pagination.Parse(c)`.
2. Read the caller's `UserInformation` with `middleware.CurrentUser(c)` when the route acts on the
   caller.
3. Call **one** service method.
4. Return a response helper from `internal/libs/utils/response.go`, or return the service's error.

A handler never touches GORM, Redis, or the queue; never decides a status code for a failure (it
returns the service's `*errors.AppError` and the central handler maps it); never builds an envelope
by hand with `fiber.Map`.

## Four things every route declares

1. **Path and verb** in `routes.go` ([routes.md](./routes.md)).
2. **Auth and guard** — `middleware.Auth()` on the group, then `guards.Permission(...)` or
   `guards.Role(...)` on the route, unless the route is public or acts only on the caller's identity
   ([rbac.md](./rbac.md)).
3. **Request struct** with validator tags in `schema.go` ([validation.md](./validation.md)).
4. **swag annotations** above the handler: `@Summary`, `@Description` (stating the permission),
   `@Tags`, `@Accept`, `@Produce`, `@Security BearerAuth` (protected routes only), `@Param`,
   `@Success`, every `@Failure` the route can return, `@Router` ([openapi.md](./openapi.md)).

## Shape

```go
func (h *Handler) Update(c *fiber.Ctx) error {
	id, err := utils.ParamUUID(c, "id")
	if err != nil {
		return err
	}
	var req UpdateRoleRequest
	if err := utils.ParseBody(c, &req); err != nil {
		return err
	}
	if err := h.service.Update(c.UserContext(), id, &req); err != nil {
		return err
	}
	return utils.OK(c, "Role updated successfully", nil)
}
```

Use these helpers rather than adding parallel ones. `utils.Error(...)` is reserved for the central
`ErrorHandler`; handlers never call it.
