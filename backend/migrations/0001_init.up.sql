CREATE TYPE guitar_type AS ENUM ('electric', 'acoustic', 'classical', 'other');
CREATE TYPE spec_value_type AS ENUM ('text', 'number', 'bool', 'enum');
CREATE TYPE media_kind AS ENUM ('image', 'video');
CREATE TYPE spec_source AS ENUM ('manufacturer', 'editor', 'user');

CREATE TABLE brand (
  id UUID PRIMARY KEY,
  name TEXT NOT NULL UNIQUE
);

CREATE TABLE guitar (
  id UUID PRIMARY KEY,
  slug TEXT NOT NULL UNIQUE,
  name TEXT NOT NULL,
  brand_id UUID NOT NULL REFERENCES brand(id),
  model TEXT NOT NULL,
  type guitar_type NOT NULL,
  year INT,
  description TEXT,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE spec (
  id UUID PRIMARY KEY,
  code TEXT NOT NULL UNIQUE,
  label TEXT NOT NULL,
  value_type spec_value_type NOT NULL,
  unit TEXT,
  filterable BOOL NOT NULL DEFAULT FALSE,
  searchable BOOL NOT NULL DEFAULT FALSE,
  guitar_type guitar_type
);

CREATE TABLE spec_option (
  id UUID PRIMARY KEY,
  spec_id UUID NOT NULL REFERENCES spec(id),
  value TEXT NOT NULL,
  sort_order INT NOT NULL DEFAULT 0,
  UNIQUE (spec_id, value)
);

CREATE TABLE guitar_spec_value (
  guitar_id UUID NOT NULL REFERENCES guitar(id),
  spec_id UUID NOT NULL REFERENCES spec(id),
  value_text TEXT,
  value_number NUMERIC,
  value_bool BOOL,
  value_option_id UUID REFERENCES spec_option(id),
  source spec_source,
  PRIMARY KEY (guitar_id, spec_id)
);

CREATE TABLE guitar_media (
  id UUID PRIMARY KEY,
  guitar_id UUID NOT NULL REFERENCES guitar(id),
  kind media_kind NOT NULL,
  url TEXT NOT NULL,
  sort_order INT NOT NULL DEFAULT 0
);

CREATE INDEX guitar_type_idx ON guitar(type);
CREATE INDEX guitar_brand_idx ON guitar(brand_id);
CREATE INDEX guitar_spec_value_spec_idx ON guitar_spec_value(spec_id);
