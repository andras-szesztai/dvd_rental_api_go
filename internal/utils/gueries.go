package utils

import (
	"net/http"
	"strconv"
)

type MovieQuery struct {
	Limit  int    `json:"limit" validate:"gte=1,lte=20" default:"20"`
	Offset int    `json:"offset" validate:"gte=0" default:"0"`
	Sort   string `json:"sort" validate:"oneof=asc desc" default:"desc"`
	Search string `json:"search" validate:"max=100" default:""`
}

func (q *MovieQuery) Parse(r *http.Request) (MovieQuery, error) {
	query := r.URL.Query()
	limit, _ := strconv.Atoi(query.Get("limit"))
	if limit == 0 {
		limit = q.Limit
	}
	offset, _ := strconv.Atoi(query.Get("offset"))
	if offset == 0 {
		offset = q.Offset
	}
	sort := query.Get("sort")
	if sort == "" {
		sort = q.Sort
	}
	search := query.Get("search")
	return MovieQuery{
		Limit:  limit,
		Offset: offset,
		Sort:   sort,
		Search: search,
	}, nil
}
