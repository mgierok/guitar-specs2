package seed

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	sqlc "github.com/mgierok/guitar-specs2/backend/internal/db/sqlc"
)

type Dataset struct {
	Brands           []Brand           `json:"brands"`
	Specs            []Spec            `json:"specs"`
	SpecOptions      []SpecOption      `json:"spec_options"`
	Guitars          []Guitar          `json:"guitars"`
	GuitarSpecValues []GuitarSpecValue `json:"guitar_spec_values"`
	GuitarMedia      []GuitarMedia     `json:"guitar_media"`
}

type Brand struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Spec struct {
	ID         string  `json:"id"`
	Code       string  `json:"code"`
	Label      string  `json:"label"`
	ValueType  string  `json:"value_type"`
	Unit       *string `json:"unit"`
	Filterable bool    `json:"filterable"`
	Searchable bool    `json:"searchable"`
	GuitarType *string `json:"guitar_type"`
}

type SpecOption struct {
	ID        string `json:"id"`
	SpecID    string `json:"spec_id"`
	Value     string `json:"value"`
	SortOrder int    `json:"sort_order"`
}

type Guitar struct {
	ID          string    `json:"id"`
	Slug        string    `json:"slug"`
	Name        string    `json:"name"`
	BrandID     string    `json:"brand_id"`
	Model       string    `json:"model"`
	Type        string    `json:"type"`
	Year        *int32    `json:"year"`
	Description *string   `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type GuitarSpecValue struct {
	GuitarID      string   `json:"guitar_id"`
	SpecID        string   `json:"spec_id"`
	ValueText     *string  `json:"value_text"`
	ValueNumber   *float64 `json:"value_number"`
	ValueBool     *bool    `json:"value_bool"`
	ValueOptionID *string  `json:"value_option_id"`
	Source        *string  `json:"source"`
}

type GuitarMedia struct {
	ID        string `json:"id"`
	GuitarID  string `json:"guitar_id"`
	Kind      string `json:"kind"`
	URL       string `json:"url"`
	SortOrder int    `json:"sort_order"`
}

func Load(path string) (Dataset, error) {
	file, err := os.Open(path)
	if err != nil {
		return Dataset{}, err
	}
	defer file.Close()

	return loadFromReader(file)
}

func loadFromReader(reader io.Reader) (Dataset, error) {
	var data Dataset
	if err := json.NewDecoder(reader).Decode(&data); err != nil {
		return Dataset{}, err
	}
	return data, nil
}

func Apply(ctx context.Context, queries *sqlc.Queries, data Dataset) error {
	for _, brand := range data.Brands {
		if err := queries.InsertBrand(ctx, sqlc.InsertBrandParams{
			ID:   brand.ID,
			Name: brand.Name,
		}); err != nil {
			return err
		}
	}

	for _, spec := range data.Specs {
		unit, err := textValue(spec.Unit)
		if err != nil {
			return err
		}
		guitarType := nullGuitarType(spec.GuitarType)

		if err := queries.InsertSpec(ctx, sqlc.InsertSpecParams{
			ID:         spec.ID,
			Code:       spec.Code,
			Label:      spec.Label,
			ValueType:  spec.ValueType,
			Unit:       unit,
			Filterable: spec.Filterable,
			Searchable: spec.Searchable,
			GuitarType: guitarType,
		}); err != nil {
			return err
		}
	}

	for _, option := range data.SpecOptions {
		if err := queries.InsertSpecOption(ctx, sqlc.InsertSpecOptionParams{
			ID:        option.ID,
			SpecID:    option.SpecID,
			Value:     option.Value,
			SortOrder: int32(option.SortOrder),
		}); err != nil {
			return err
		}
	}

	for _, guitar := range data.Guitars {
		year := intValue(guitar.Year)
		description, err := textValue(guitar.Description)
		if err != nil {
			return err
		}

		if err := queries.InsertGuitar(ctx, sqlc.InsertGuitarParams{
			ID:          guitar.ID,
			Slug:        guitar.Slug,
			Name:        guitar.Name,
			BrandID:     guitar.BrandID,
			Model:       guitar.Model,
			Type:        guitar.Type,
			Year:        year,
			Description: description,
			CreatedAt:   guitar.CreatedAt,
			UpdatedAt:   guitar.UpdatedAt,
		}); err != nil {
			return err
		}
	}

	for _, value := range data.GuitarSpecValues {
		textVal, err := textValue(value.ValueText)
		if err != nil {
			return err
		}
		numVal, err := numericValue(value.ValueNumber)
		if err != nil {
			return err
		}
		boolVal := boolValue(value.ValueBool)
		optionID, err := uuidValue(value.ValueOptionID)
		if err != nil {
			return err
		}
		source := nullSpecSource(value.Source)

		if err := queries.InsertGuitarSpecValue(ctx, sqlc.InsertGuitarSpecValueParams{
			GuitarID:      value.GuitarID,
			SpecID:        value.SpecID,
			ValueText:     textVal,
			ValueNumber:   numVal,
			ValueBool:     boolVal,
			ValueOptionID: optionID,
			Source:        source,
		}); err != nil {
			return err
		}
	}

	for _, media := range data.GuitarMedia {
		if err := queries.InsertGuitarMedia(ctx, sqlc.InsertGuitarMediaParams{
			ID:        media.ID,
			GuitarID:  media.GuitarID,
			Kind:      media.Kind,
			Url:       media.URL,
			SortOrder: int32(media.SortOrder),
		}); err != nil {
			return err
		}
	}

	return nil
}

func textValue(value *string) (pgtype.Text, error) {
	if value == nil {
		return pgtype.Text{}, nil
	}
	return pgtype.Text{String: *value, Valid: true}, nil
}

func numericValue(value *float64) (pgtype.Numeric, error) {
	if value == nil {
		return pgtype.Numeric{}, nil
	}
	var numeric pgtype.Numeric
	if err := numeric.Scan(*value); err != nil {
		return pgtype.Numeric{}, err
	}
	return numeric, nil
}

func boolValue(value *bool) pgtype.Bool {
	if value == nil {
		return pgtype.Bool{}
	}
	return pgtype.Bool{Bool: *value, Valid: true}
}

func uuidValue(value *string) (pgtype.UUID, error) {
	if value == nil {
		return pgtype.UUID{}, nil
	}
	var id pgtype.UUID
	if err := id.Scan(*value); err != nil {
		return pgtype.UUID{}, err
	}
	return id, nil
}

func intValue(value *int32) pgtype.Int4 {
	if value == nil {
		return pgtype.Int4{}
	}
	return pgtype.Int4{Int32: *value, Valid: true}
}

func nullGuitarType(value *string) sqlc.NullGuitarType {
	if value == nil {
		return sqlc.NullGuitarType{}
	}
	return sqlc.NullGuitarType{
		GuitarType: sqlc.GuitarType(*value),
		Valid:      true,
	}
}

func nullSpecSource(value *string) sqlc.NullSpecSource {
	if value == nil {
		return sqlc.NullSpecSource{}
	}
	return sqlc.NullSpecSource{
		SpecSource: sqlc.SpecSource(*value),
		Valid:      true,
	}
}
