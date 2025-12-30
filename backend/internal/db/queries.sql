-- name: InsertBrand :exec
INSERT INTO brand (id, name)
VALUES ($1, $2);

-- name: InsertGuitar :exec
INSERT INTO guitar (id, slug, name, brand_id, model, type, year, description, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);

-- name: InsertSpec :exec
INSERT INTO spec (id, code, label, value_type, unit, filterable, searchable, guitar_type)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8);

-- name: InsertSpecOption :exec
INSERT INTO spec_option (id, spec_id, value, sort_order)
VALUES ($1, $2, $3, $4);

-- name: InsertGuitarSpecValue :exec
INSERT INTO guitar_spec_value (guitar_id, spec_id, value_text, value_number, value_bool, value_option_id, source)
VALUES ($1, $2, $3, $4, $5, $6, $7);

-- name: InsertGuitarMedia :exec
INSERT INTO guitar_media (id, guitar_id, kind, url, sort_order)
VALUES ($1, $2, $3, $4, $5);

-- name: ListGuitars :many
SELECT
  guitar.id,
  guitar.slug,
  guitar.name,
  guitar.model,
  guitar.type,
  guitar.year,
  brand.name AS brand_name,
  COALESCE(
    (
      SELECT url
      FROM guitar_media
      WHERE guitar_media.guitar_id = guitar.id
        AND guitar_media.kind = 'image'
      ORDER BY sort_order ASC
      LIMIT 1
    ),
    ''
  )::text AS thumbnail_url
FROM guitar
JOIN brand ON brand.id = guitar.brand_id
ORDER BY guitar.name ASC;

-- name: GetGuitarBySlug :one
SELECT
  guitar.id,
  guitar.slug,
  guitar.name,
  guitar.model,
  guitar.type,
  guitar.year,
  guitar.description,
  brand.name AS brand_name
FROM guitar
JOIN brand ON brand.id = guitar.brand_id
WHERE guitar.slug = $1;

-- name: GetGuitarSpecs :many
SELECT
  spec.code,
  COALESCE(
    guitar_spec_value.value_text,
    guitar_spec_value.value_number::text,
    guitar_spec_value.value_bool::text,
    spec_option.value,
    ''
  ) AS value
FROM guitar_spec_value
JOIN spec ON spec.id = guitar_spec_value.spec_id
LEFT JOIN spec_option ON spec_option.id = guitar_spec_value.value_option_id
WHERE guitar_spec_value.guitar_id = $1
ORDER BY spec.code ASC;

-- name: GetGuitarMedia :many
SELECT
  kind,
  url
FROM guitar_media
WHERE guitar_id = $1
ORDER BY sort_order ASC;
