package paginate

import (
	"math"
	"strconv"

	"github.com/labstack/echo"
)

type (
	Paginate struct {
		Data interface{}
		Pages
		Items
	}
	Pages struct {
		Current int
		Prev    int
		HasPrev bool
		Next    int
		HasNext bool
		Total   int
	}
	Items struct {
		Limit int
		Begin int
		End   int
		Total int
	}
)

func Generate(data interface{}, count, page, limit int) Paginate {

	totalPage := math.Ceil(float64(count) / float64(limit))
	begin := ((page * limit) - limit) + 1
	end := page * limit

	result := Paginate{
		Data: data,
		Pages: Pages{
			Current: page,
			Prev:    page - 1,
			HasPrev: (page - 1) != 0,
			Next:    page + 1,
			HasNext: (page + 1) <= int(totalPage),
			Total:   int(totalPage),
		},
		Items: Items{
			Limit: limit,
			Begin: begin,
			End:   end,
			Total: count,
		},
	}

	if begin > count {
		result.Items.Begin = count
	}

	if end > count {
		result.Items.End = count
	}

	return result
}

func HandleQueries(c echo.Context) (int, int) {

	var page, limit = 1, 10
	queryPage := c.QueryParam("page")
	queryLimit := c.QueryParam("limit")

	if queryPage != "" {
		p, err := strconv.Atoi(queryPage)
		if err == nil {
			page = p
		}
	}

	if queryLimit != "" {
		l, err := strconv.Atoi(queryLimit)
		if err == nil {
			limit = l
		}
	}

	return page, limit
}
