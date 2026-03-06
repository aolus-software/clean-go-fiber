package pagination

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func ParseFromFiber(c *fiber.Ctx) *QueryParamDatatable {
	p := &QueryParamDatatable{
		Search:        c.Query("search"),
		OrderBy:       c.Query("order_by"),
		SortDirection: c.Query("sort_direction"),
		Filters:       make(map[string]string),
	}

	if perPage, err := strconv.Atoi(c.Query("per_page")); err == nil {
		p.PerPage = perPage
	}
	if page, err := strconv.Atoi(c.Query("page")); err == nil {
		p.Page = page
	}

	c.Request().URI().QueryArgs().VisitAll(func(key, value []byte) {
		k := string(key)
		if strings.HasPrefix(k, "filter[") && strings.HasSuffix(k, "]") {
			field := k[7 : len(k)-1]
			p.Filters[field] = string(value)
		}
	})

	return p
}
