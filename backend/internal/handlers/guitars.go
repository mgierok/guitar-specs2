package handlers

import (
	"context"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/mgierok/guitar-specs2/backend/internal/api"
	db "github.com/mgierok/guitar-specs2/backend/internal/db/sqlc"
)

type GuitarQueries interface {
	ListGuitars(ctx context.Context) ([]db.ListGuitarsRow, error)
	GetGuitarBySlug(ctx context.Context, slug string) (db.GetGuitarBySlugRow, error)
	GetGuitarSpecs(ctx context.Context, guitarID string) ([]db.GetGuitarSpecsRow, error)
	GetGuitarMedia(ctx context.Context, guitarID string) ([]db.GetGuitarMediaRow, error)
}

type GuitarHandler struct {
	Queries GuitarQueries
}

type GuitarListItem struct {
	ID        string `json:"id"`
	Slug      string `json:"slug"`
	Name      string `json:"name"`
	Brand     string `json:"brand"`
	Model     string `json:"model"`
	Type      string `json:"type"`
	Year      *int32 `json:"year,omitempty"`
	Thumbnail string `json:"thumbnail,omitempty"`
}

type GuitarDetail struct {
	ID          string            `json:"id"`
	Slug        string            `json:"slug"`
	Name        string            `json:"name"`
	Brand       string            `json:"brand"`
	Model       string            `json:"model"`
	Type        string            `json:"type"`
	Year        *int32            `json:"year,omitempty"`
	Description *string           `json:"description,omitempty"`
	Specs       map[string]string `json:"specs"`
	Media       []GuitarMedia     `json:"media"`
}

type GuitarMedia struct {
	Kind string `json:"kind"`
	URL  string `json:"url"`
}

func (h *GuitarHandler) List(w http.ResponseWriter, r *http.Request) {
	params := api.ParseListParams(r)

	items, err := h.Queries.ListGuitars(r.Context())
	if err != nil {
		api.WriteError(w, http.StatusInternalServerError, "list_failed", "Failed to list guitars", nil)
		return
	}

	response := api.ListResponse[GuitarListItem]{
		Items:    mapListItems(items),
		Total:    len(items),
		Page:     params.Page,
		PageSize: params.PageSize,
	}

	api.WriteJSON(w, http.StatusOK, response)
}

func (h *GuitarHandler) Detail(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	if strings.TrimSpace(slug) == "" {
		api.WriteError(w, http.StatusBadRequest, "invalid_slug", "Slug is required", nil)
		return
	}

	guitar, err := h.Queries.GetGuitarBySlug(r.Context(), slug)
	if err != nil {
		api.WriteError(w, http.StatusNotFound, "not_found", "Guitar not found", nil)
		return
	}

	specRows, err := h.Queries.GetGuitarSpecs(r.Context(), guitar.ID)
	if err != nil {
		api.WriteError(w, http.StatusInternalServerError, "specs_failed", "Failed to load specs", nil)
		return
	}

	mediaRows, err := h.Queries.GetGuitarMedia(r.Context(), guitar.ID)
	if err != nil {
		api.WriteError(w, http.StatusInternalServerError, "media_failed", "Failed to load media", nil)
		return
	}

	detail := GuitarDetail{
		ID:          guitar.ID,
		Slug:        guitar.Slug,
		Name:        guitar.Name,
		Brand:       guitar.BrandName,
		Model:       guitar.Model,
		Type:        guitar.Type,
		Year:        int32Ptr(guitar.Year),
		Description: stringPtr(guitar.Description),
		Specs:       mapSpecs(specRows),
		Media:       mapMedia(mediaRows),
	}

	api.WriteJSON(w, http.StatusOK, detail)
}

func (h *GuitarHandler) WithContext(ctx context.Context, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		next(w, r.WithContext(ctx))
	}
}

func mapListItems(rows []db.ListGuitarsRow) []GuitarListItem {
	items := make([]GuitarListItem, 0, len(rows))
	for _, row := range rows {
		items = append(items, GuitarListItem{
			ID:        row.ID,
			Slug:      row.Slug,
			Name:      row.Name,
			Brand:     row.BrandName,
			Model:     row.Model,
			Type:      row.Type,
			Year:      int32Ptr(row.Year),
			Thumbnail: row.ThumbnailUrl,
		})
	}
	return items
}

func mapSpecs(rows []db.GetGuitarSpecsRow) map[string]string {
	result := map[string]string{}
	for _, row := range rows {
		result[row.Code] = row.Value
	}
	return result
}

func mapMedia(rows []db.GetGuitarMediaRow) []GuitarMedia {
	media := make([]GuitarMedia, 0, len(rows))
	for _, row := range rows {
		media = append(media, GuitarMedia{
			Kind: row.Kind,
			URL:  row.Url,
		})
	}
	return media
}

func int32Ptr(value pgtype.Int4) *int32 {
	if !value.Valid {
		return nil
	}
	return &value.Int32
}

func stringPtr(value pgtype.Text) *string {
	if !value.Valid {
		return nil
	}
	return &value.String
}
