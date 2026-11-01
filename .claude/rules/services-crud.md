# Rule: The canonical CRUD service

The service **owns existence and uniqueness checks** — settled, do not move them into repositories.
Repositories return `(nil, nil)` for a missing row; the service decides what that means.

## Five methods

| Method | Steps | Errors |
| --- | --- | --- |
| `List(ctx, q)` | delegate to `repo.FindAll(ctx, q)`; map models to camelCase response structs | allow-list violations pass through as 400/422 — **never** rewrap them as 500 |
| `Create(ctx, req)` | uniqueness check (live rows only) → in a transaction: insert + join rows | 422 `errors[{field:"email"\|"name"}]` on a duplicate; 422 on unknown related ids |
| `Detail(ctx, id)` | `repo.FindByID`; nil → 404 | 404 |
| `Update(ctx, id, req)` | `FindByID` → 404; uniqueness excluding this id → 422; transaction: update + replace join rows; invalidate cache | 404, 422 |
| `Delete(ctx, id)` | `FindByID` → 404; delete (soft for users, hard for roles/permissions); invalidate cache | 404, 422 for a protected row |

## Uniqueness is 422 everywhere

Duplicate user email (admin create/update, profile update, **and** `/auth/register`), duplicate role
name, duplicate permission name — all 422 with the offending field in `errors`. The database also
enforces it (`users_email_unique_live`, `roles_name_unique`, `permissions_name_unique`); if the
insert still loses a race, map the unique-violation to the same 422 rather than letting it become a
500.

## Related ids

`roleIds` / `permissionIds` are validated as UUIDs by the validator, then the service loads them and
rejects any id that does not exist with 422 (`field: "roleIds"` / `"permissionIds"`).

## Protected rows

- The `superuser` role cannot be renamed or deleted — 422.
- An actor who is not a superuser cannot assign or remove the `superuser` role, nor edit or delete a
  user who holds it — 403. (Admin password reset is already superuser-only by its guard.)
- Permissions created through the API store `"<group> <action>"` per entry in `name[]`; `PATCH`
  stores `name` verbatim, so an update must send the full `"<group> <action>"` string.

## Example

```go
func (s *service) Create(ctx context.Context, req *CreateRoleRequest) error {
	existing, err := s.roleRepo.FindByName(ctx, req.Name, nil)
	if err != nil {
		return err
	}
	if existing != nil {
		return apperrors.NewUnprocessableEntityError("Role name already exists", apperrors.Field("name", "Role name already exists"))
	}
	return s.tx.WithinTransaction(ctx, func(tx *gorm.DB) error {
		return s.roleRepo.WithTx(tx).Create(ctx, req.Name, permissionIDs)
	})
}
```

`gorm.ErrDuplicatedKey` from a lost race maps to the same 422.
