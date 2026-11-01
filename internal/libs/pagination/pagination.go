package pagination

import (
	"slices"
	"strconv"
	"strings"
	"time"

	apperrors "github.com/aolus-software/clean-go-fiber/internal/errors"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	DefaultPage          = 1
	DefaultPerPage       = 10
	MaxPerPage           = 100
	DefaultSort          = "createdAt"
	DefaultSortDirection = "desc"
)

type FilterKind int

const (
	FilterText FilterKind = iota
	FilterEnum
	FilterID
	FilterDate
)

type Filter struct {
	Column string
	Kind   FilterKind
	Enum   []string
	Expr   string
}

type Options struct {
	Sortable      map[string]string
	Filterable    map[string]Filter
	SearchColumns []string
}

type Query struct {
	Page          int
	PerPage       int
	Search        string
	Sort          string
	SortDirection string
	Filters       map[string]string
}

func Parse(c *fiber.Ctx) (*Query, error) {
	q := &Query{
		Page:          DefaultPage,
		PerPage:       DefaultPerPage,
		Search:        strings.TrimSpace(c.Query("search")),
		Sort:          c.Query("sort", DefaultSort),
		SortDirection: c.Query("sortDirection", DefaultSortDirection),
		Filters:       map[string]string{},
	}

	var fields []apperrors.FieldError
	if raw := c.Query("page"); raw != "" {
		n, err := strconv.Atoi(raw)
		if err != nil || n < 1 {
			fields = append(fields, apperrors.Field("page", "must be an integer greater than or equal to 1"))
		} else {
			q.Page = n
		}
	}
	if raw := c.Query("perPage"); raw != "" {
		n, err := strconv.Atoi(raw)
		if err != nil || n < 1 || n > MaxPerPage {
			fields = append(fields, apperrors.Field("perPage", "must be an integer between 1 and "+strconv.Itoa(MaxPerPage)))
		} else {
			q.PerPage = n
		}
	}
	if q.SortDirection != "asc" && q.SortDirection != "desc" {
		fields = append(fields, apperrors.Field("sortDirection", "must be one of: asc, desc"))
	}
	if len(fields) > 0 {
		return nil, apperrors.NewUnprocessableEntityError("Validation error", fields...)
	}

	c.Context().QueryArgs().VisitAll(func(key, value []byte) {
		k := string(key)
		if strings.HasPrefix(k, "filter[") && strings.HasSuffix(k, "]") {
			q.Filters[k[len("filter["):len(k)-1]] = string(value)
		}
	})

	return q, nil
}

func (q *Query) Offset() int {
	return (q.Page - 1) * q.PerPage
}

func (q *Query) Apply(db *gorm.DB, opts Options, tz *time.Location) (*gorm.DB, error) {
	column, ok := opts.Sortable[q.Sort]
	if !ok {
		keys := make([]string, 0, len(opts.Sortable))
		for k := range opts.Sortable {
			keys = append(keys, k)
		}
		slices.Sort(keys)
		return nil, apperrors.NewUnprocessableEntityError("Validation error", apperrors.Field("sort", "must be one of: "+strings.Join(keys, ", ")))
	}

	if q.Search != "" && len(opts.SearchColumns) > 0 {
		like := "%" + escapeLike(q.Search) + "%"
		conds := make([]string, 0, len(opts.SearchColumns))
		args := make([]any, 0, len(opts.SearchColumns))
		for _, col := range opts.SearchColumns {
			conds = append(conds, col+" ILIKE ?")
			args = append(args, like)
		}
		db = db.Where("("+strings.Join(conds, " OR ")+")", args...)
	}

	for key, value := range q.Filters {
		def, ok := opts.Filterable[key]
		if !ok {
			return nil, apperrors.NewBadRequestError("Invalid filter", apperrors.Field("filter["+key+"]", "is not a filterable field"))
		}
		var err error
		db, err = applyFilter(db, key, value, def, tz)
		if err != nil {
			return nil, err
		}
	}

	return db.Order(column + " " + q.SortDirection), nil
}

func applyFilter(db *gorm.DB, key, value string, def Filter, tz *time.Location) (*gorm.DB, error) {
	field := "filter[" + key + "]"
	switch def.Kind {
	case FilterEnum:
		values := splitList(value)
		for _, v := range values {
			if !slices.Contains(def.Enum, v) {
				return nil, apperrors.NewBadRequestError("Invalid filter", apperrors.Field(field, "must be one of: "+strings.Join(def.Enum, ", ")))
			}
		}
		return db.Where(def.Column+" IN ?", values), nil
	case FilterID:
		values := splitList(value)
		ids := make([]uuid.UUID, 0, len(values))
		for _, v := range values {
			id, err := uuid.Parse(v)
			if err != nil {
				return nil, apperrors.NewBadRequestError("Invalid filter", apperrors.Field(field, "must be a comma-separated list of uuids"))
			}
			ids = append(ids, id)
		}
		if def.Expr != "" {
			return db.Where(def.Expr, ids), nil
		}
		return db.Where(def.Column+" IN ?", ids), nil
	case FilterDate:
		start, end, err := parseDateRange(value, tz)
		if err != nil {
			return nil, apperrors.NewBadRequestError("Invalid filter", apperrors.Field(field, "must be a date (YYYY-MM-DD) or a range start,end"))
		}
		return db.Where(def.Column+" >= ? AND "+def.Column+" < ?", start, end), nil
	default:
		return db.Where(def.Column+" ILIKE ?", "%"+escapeLike(value)+"%"), nil
	}
}

func parseDateRange(value string, tz *time.Location) (time.Time, time.Time, error) {
	parts := strings.Split(value, ",")
	if len(parts) > 2 {
		return time.Time{}, time.Time{}, apperrors.NewBadRequestError("invalid date range")
	}
	first, err := time.ParseInLocation(time.DateOnly, strings.TrimSpace(parts[0]), tz)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	last := first
	if len(parts) == 2 {
		last, err = time.ParseInLocation(time.DateOnly, strings.TrimSpace(parts[1]), tz)
		if err != nil {
			return time.Time{}, time.Time{}, err
		}
	}
	if last.Before(first) {
		return time.Time{}, time.Time{}, apperrors.NewBadRequestError("invalid date range")
	}
	return first.UTC(), last.AddDate(0, 0, 1).UTC(), nil
}

func splitList(value string) []string {
	parts := strings.Split(value, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	return out
}

func escapeLike(s string) string {
	return strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`).Replace(s)
}
