package api

import (
	"net/http"
	"strconv"
	"strings"
)

type ListParams struct {
	Page     int
	PageSize int
	Sort     string
	Filters  map[string]string
}

func ParseListParams(r *http.Request) ListParams {
	page := parseInt(r.URL.Query().Get("page"), 1)
	pageSize := parseInt(r.URL.Query().Get("pageSize"), 20)

	filters := map[string]string{}
	for _, value := range r.URL.Query()["filter"] {
		parts := strings.SplitN(value, ":", 2)
		if len(parts) == 2 {
			filters[parts[0]] = parts[1]
		}
	}

	sort := r.URL.Query().Get("sort")
	if sort == "" {
		sort = "name:asc"
	}

	return ListParams{
		Page:     page,
		PageSize: pageSize,
		Sort:     sort,
		Filters:  filters,
	}
}

func parseInt(value string, fallback int) int {
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil || parsed <= 0 {
		return fallback
	}
	return parsed
}
