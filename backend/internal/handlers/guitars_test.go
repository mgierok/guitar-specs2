package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"

	db "github.com/mgierok/guitar-specs2/backend/internal/db/sqlc"
)

type fakeQueries struct {
	list      []db.ListGuitarsRow
	detail    db.GetGuitarBySlugRow
	specs     []db.GetGuitarSpecsRow
	media     []db.GetGuitarMediaRow
	detailErr error
}

func (f fakeQueries) ListGuitars(_ context.Context) ([]db.ListGuitarsRow, error) {
	return f.list, nil
}

func (f fakeQueries) GetGuitarBySlug(_ context.Context, _ string) (db.GetGuitarBySlugRow, error) {
	return f.detail, f.detailErr
}

func (f fakeQueries) GetGuitarSpecs(_ context.Context, _ string) ([]db.GetGuitarSpecsRow, error) {
	return f.specs, nil
}

func (f fakeQueries) GetGuitarMedia(_ context.Context, _ string) ([]db.GetGuitarMediaRow, error) {
	return f.media, nil
}

func TestListEmpty(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/guitars", nil)
	w := httptest.NewRecorder()

	h := GuitarHandler{Queries: fakeQueries{list: []db.ListGuitarsRow{}}}
	h.List(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var payload struct {
		Items []interface{} `json:"items"`
		Total int           `json:"total"`
	}
	if err := json.NewDecoder(w.Body).Decode(&payload); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if payload.Total != 0 {
		t.Fatalf("expected total 0, got %d", payload.Total)
	}
	if len(payload.Items) != 0 {
		t.Fatalf("expected 0 items, got %d", len(payload.Items))
	}
}

func TestDetailNotFound(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/guitars/missing", nil)
	routeContext := chi.NewRouteContext()
	routeContext.URLParams.Add("slug", "missing")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeContext))
	w := httptest.NewRecorder()

	h := GuitarHandler{Queries: fakeQueries{detailErr: context.Canceled}}
	h.Detail(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestDetailSuccess(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/guitars/fender-player-stratocaster", nil)
	routeContext := chi.NewRouteContext()
	routeContext.URLParams.Add("slug", "fender-player-stratocaster")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeContext))
	w := httptest.NewRecorder()

	h := GuitarHandler{
		Queries: fakeQueries{
			detail: db.GetGuitarBySlugRow{
				ID:        "g1",
				Slug:      "fender-player-stratocaster",
				Name:      "Player Stratocaster",
				Model:     "Stratocaster",
				Type:      "electric",
				BrandName: "Fender",
			},
			specs: []db.GetGuitarSpecsRow{
				{Code: "scale_length", Value: "25.5"},
			},
			media: []db.GetGuitarMediaRow{
				{Kind: "image", Url: "https://example.com/1.jpg"},
			},
		},
	}

	h.Detail(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var payload struct {
		Name  string            `json:"name"`
		Brand string            `json:"brand"`
		Specs map[string]string `json:"specs"`
		Media []struct {
			Kind string `json:"kind"`
			URL  string `json:"url"`
		} `json:"media"`
	}
	if err := json.NewDecoder(w.Body).Decode(&payload); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if payload.Name != "Player Stratocaster" {
		t.Fatalf("unexpected name: %s", payload.Name)
	}
	if payload.Brand != "Fender" {
		t.Fatalf("unexpected brand: %s", payload.Brand)
	}
	if payload.Specs["scale_length"] != "25.5" {
		t.Fatalf("unexpected specs: %v", payload.Specs)
	}
	if len(payload.Media) != 1 || payload.Media[0].Kind != "image" {
		t.Fatalf("unexpected media: %+v", payload.Media)
	}
}
