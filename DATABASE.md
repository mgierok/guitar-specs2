# Database Reference

This document is the single source of truth for data entities. Keep backend models, migrations, and frontend types aligned with it.

## Goals
- Support predefined, normalized specs for reliable filtering.
- Allow unlimited custom parameters per guitar without schema churn.
- Keep search fast on common filters and full-text on descriptions.

## Core Entities

### Guitar
Base record for every instrument.

- `id` (uuid, pk)
- `slug` (text, unique)
- `name` (text)
- `brand_id` (uuid, fk -> `brand.id`)
- `model` (text)
- `type` (enum) — `electric`, `acoustic`, `classical`, `other`
- `year` (int, nullable)
- `description` (text, nullable)
- `created_at` (timestamptz)
- `updated_at` (timestamptz)

### Brand
- `id` (uuid, pk)
- `name` (text, unique)

### Spec
Normalized definition of a parameter. Most parameters live here.

- `id` (uuid, pk)
- `code` (text, unique) — stable key (e.g., `scale_length`)
- `label` (text) — UI label (e.g., "Scale length")
- `value_type` (enum) — `text`, `number`, `bool`, `enum`
- `unit` (text, nullable) — e.g., `in`, `mm`
- `filterable` (bool) — enable faceted filtering
- `searchable` (bool) — include in free-text search
- `guitar_type` (enum, nullable) — restrict to a type

### SpecOption
For enum specs (pickup configs, colors, finishes, etc.).

- `id` (uuid, pk)
- `spec_id` (uuid, fk -> `spec.id`)
- `value` (text)
- `sort_order` (int)

### GuitarSpecValue
Links a guitar to a spec and stores the actual value.

- `guitar_id` (uuid, fk -> `guitar.id`)
- `spec_id` (uuid, fk -> `spec.id`)
- `value_text` (text, nullable)
- `value_number` (numeric, nullable)
- `value_bool` (bool, nullable)
- `value_option_id` (uuid, fk -> `spec_option.id`, nullable)
- `source` (enum, nullable) — `manufacturer`, `editor`, `user`

Composite unique: (`guitar_id`, `spec_id`)

### GuitarMedia
- `id` (uuid, pk)
- `guitar_id` (uuid, fk -> `guitar.id`)
- `kind` (enum) — `image`, `video`
- `url` (text)
- `sort_order` (int)

## Relations
- `guitar` N:1 `brand` (many guitars per brand).
- `guitar` 1:N `guitar_spec_value` (one guitar has many spec values).
- `spec` 1:N `spec_option` (one spec can have many allowed options).
- `spec` 1:N `guitar_spec_value` (one spec used by many guitars).
- `spec_option` 1:N `guitar_spec_value` (option chosen by many guitars).
- `guitar` 1:N `guitar_media` (gallery items per guitar).

## Searching and Filtering
- Use indexes on `guitar.type`, `brand_id`, and `spec_id`.
- For filters, query `GuitarSpecValue` by `spec_id` and the appropriate value column.
- For text search, combine `guitar.name`, `brand.name`, and selected `spec` values marked `searchable`.

## Notes
- Add new predefined specs by inserting into `spec` and optional `spec_option`.
- This model allows unlimited parameters per guitar while keeping normalized, filterable specs.
- Canonical test dataset: `data/guitars.json` (update alongside schema changes).
