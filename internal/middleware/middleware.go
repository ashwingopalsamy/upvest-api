package middleware

import (
	"context"
	"net/http"
	"strconv"
)

type PagingParams struct {
	Offset int
	Limit  int
	Sort   string
	Order  string
}

const (
	DefaultOffset = 0
	DefaultLimit  = 100
	MaxLimit      = 1000
	DefaultSort   = "created_at"
	DefaultOrder  = "ASC"
)

func ExtractPagingParams(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()

		// Parse and validate offset
		offset, err := strconv.Atoi(query.Get("offset"))
		if err != nil || offset < 0 {
			offset = DefaultOffset
		}

		// Parse and validate limit
		limit, err := strconv.Atoi(query.Get("limit"))
		if err != nil || limit <= 0 || limit > MaxLimit {
			limit = DefaultLimit
		}

		// Parse sort field
		sort := query.Get("sort")
		if sort != "created_at" && sort != "updated_at" {
			sort = DefaultSort
		}

		// Parse sort order
		order := query.Get("order")
		if order != "ASC" && order != "DESC" {
			order = DefaultOrder
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, "pagingParams", PagingParams{
			Offset: offset,
			Limit:  limit,
			Sort:   sort,
			Order:  order,
		})

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
