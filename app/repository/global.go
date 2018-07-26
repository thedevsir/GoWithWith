package repository

import (
	"math"
)

var (
	String string
)

type (
	Pagination struct {
		Data  interface{}
		Pages PaginationPages
		Items PaginationItems
	}
	PaginationPages struct {
		Current int
		Prev    int
		HasPrev bool
		Next    int
		HasNext bool
		Total   int
	}
	PaginationItems struct {
		Limit int
		Begin int
		End   int
		Total int
	}
)

func GeneratePagination(data interface{}, count, page, limit int) Pagination {

	totalPage := math.Ceil(float64(count) / float64(limit))
	begin := ((page * limit) - limit) + 1
	end := page * limit

	result := Pagination{
		Data: data,
		Pages: PaginationPages{
			Current: page,
			Prev:    page - 1,
			HasPrev: (page - 1) != 0,
			Next:    page + 1,
			HasNext: (page + 1) <= int(totalPage),
			Total:   int(totalPage),
		},
		Items: PaginationItems{
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
