# Rule: Imports and naming

## Imports

Module path: `github.com/aolus-software/clean-go-fiber`. Two groups, gofmt-sorted within each:
standard library, then everything else (internal and third-party together, as the existing files
do).

```go
import (
	"context"
	"time"

	apperrors "github.com/aolus-software/clean-go-fiber/internal/errors"
	"github.com/aolus-software/clean-go-fiber/internal/repositories"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)
```

Aliases only to resolve a clash or a misleading name: `apperrors` for `internal/errors` (clashes
with the standard `errors`), `redisdb` for `internal/database/redis`.

## Packages and files

- Package names: lowercase, one word, no underscores (`selectoptions`, `repositories`).
- File names: `snake_case.go` (`user_repository.go`, `permission_guard.go`).
- Module files are always `handler.go`, `routes.go`, `schema.go`, `service.go`.
- Migrations: `NNNNNN_snake_case_name.up.sql` / `.down.sql`.

## Symbols

- Interfaces end in their role: `UserRepository`, `Service` (inside a module package).
- Constructors: `New<Thing>(...)`. Unexported implementation structs (`userRepository`, `service`).
- Request structs: `<Action><Entity>Request` (`CreateUserRequest`); response structs:
  `<Entity><Shape>` (`UserDetail`, `RoleList`, `UserInformation`).
- JSON tags camelCase; DB columns snake_case (GORM `column:` tag when the default mapping differs).
- Initialisms keep their case: `ID`, `URL`, `JWT` (`UserID`, `FindByID`).
