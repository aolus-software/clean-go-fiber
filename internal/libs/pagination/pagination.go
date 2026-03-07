package pagination

import (
	"fmt"

	apperrors "github.com/aolus-software/clean-go-fiber/internal/errors"
)

type QueryParamDatatable struct {
	Search        string
	Filters       map[string]string
	OrderBy       string
	SortDirection string
	PerPage       int
	Page          int
}

type PaginationResponse[T any] struct {
	Total       int64 `json:"total"`
	PerPage     int   `json:"per_page"`
	CurrentPage int   `json:"current_page"`
	LastPage    int   `json:"last_page"`
	Data        []T   `json:"data"`
}

type DatatableOptions struct {
	AllowableFilter []string
	AllowableSort   []string
	DefaultSort     string
	DefaultSortDir  string
	DefaultPerPage  int
}

func (q *QueryParamDatatable) ApplyDefaults(options DatatableOptions) {
	if q.Page <= 0 {
		q.Page = 1
	}
	if q.PerPage <= 0 {
		if options.DefaultPerPage > 0 {
			q.PerPage = options.DefaultPerPage
		} else {
			q.PerPage = 10
		}
	}
	if q.OrderBy == "" && options.DefaultSort != "" {
		q.OrderBy = options.DefaultSort
	}
	if q.SortDirection == "" && options.DefaultSortDir != "" {
		q.SortDirection = options.DefaultSortDir
	}
}

func (q *QueryParamDatatable) Validate(options DatatableOptions) error {
	if q.OrderBy != "" && !isAllowableSort(q.OrderBy, options) {
		return apperrors.NewBadRequestError("The order_by field is invalid")
	}
	if q.SortDirection != "" && q.SortDirection != "asc" && q.SortDirection != "desc" {
		return apperrors.NewBadRequestError("Invalid sort_direction field")
	}
	for filter := range q.Filters {
		if !isAllowableFilter(filter, options) {
			return apperrors.NewBadRequestError(fmt.Sprintf("Invalid filter field: %s", filter))
		}
	}
	return nil
}

func CalcLastPage(total int64, perPage int) int {
	if perPage <= 0 {
		return 0
	}
	return int((total + int64(perPage) - 1) / int64(perPage))
}

func isAllowableSort(field string, options DatatableOptions) bool {
	for _, f := range options.AllowableSort {
		if f == field {
			return true
		}
	}
	return false
}

func isAllowableFilter(field string, options DatatableOptions) bool {
	for _, f := range options.AllowableFilter {
		if f == field {
			return true
		}
	}
	return false
}
