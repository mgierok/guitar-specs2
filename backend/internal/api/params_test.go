package api

import (
	"net/http"
	"testing"
)

func TestParseListParamsDefaults(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "https://example.com/api/v1/guitars", nil)
	if err != nil {
		t.Fatalf("request: %v", err)
	}

	params := ParseListParams(req)
	if params.Page != 1 {
		t.Fatalf("expected page 1, got %d", params.Page)
	}
	if params.PageSize != 20 {
		t.Fatalf("expected pageSize 20, got %d", params.PageSize)
	}
	if params.Sort != "name:asc" {
		t.Fatalf("expected sort name:asc, got %s", params.Sort)
	}
	if len(params.Filters) != 0 {
		t.Fatalf("expected no filters, got %d", len(params.Filters))
	}
}

func TestParseListParamsFilters(t *testing.T) {
	req, err := http.NewRequest(
		http.MethodGet,
		"https://example.com/api/v1/guitars?filter=type:electric&filter=brand:Fender&filter=brand:Gibson",
		nil,
	)
	if err != nil {
		t.Fatalf("request: %v", err)
	}

	params := ParseListParams(req)
	if params.Filters["type"] != "electric" {
		t.Fatalf("expected type filter electric, got %s", params.Filters["type"])
	}
	if params.Filters["brand"] != "Fender" {
		t.Fatalf("expected brand filter Fender (no duplicates), got %s", params.Filters["brand"])
	}
}

func TestParseListParamsInvalidNumbers(t *testing.T) {
	req, err := http.NewRequest(
		http.MethodGet,
		"https://example.com/api/v1/guitars?page=-1&pageSize=0&sort=year:desc",
		nil,
	)
	if err != nil {
		t.Fatalf("request: %v", err)
	}

	params := ParseListParams(req)
	if params.Page != 1 {
		t.Fatalf("expected fallback page 1, got %d", params.Page)
	}
	if params.PageSize != 20 {
		t.Fatalf("expected fallback pageSize 20, got %d", params.PageSize)
	}
	if params.Sort != "year:desc" {
		t.Fatalf("expected sort year:desc, got %s", params.Sort)
	}
}
