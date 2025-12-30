DO $$
BEGIN
  CREATE TYPE guitar_type AS ENUM ('electric', 'acoustic', 'classical', 'other');
EXCEPTION
  WHEN duplicate_object THEN NULL;
END $$;

DO $$
BEGIN
  CREATE TYPE spec_value_type AS ENUM ('text', 'number', 'bool', 'enum');
EXCEPTION
  WHEN duplicate_object THEN NULL;
END $$;

DO $$
BEGIN
  CREATE TYPE media_kind AS ENUM ('image', 'video');
EXCEPTION
  WHEN duplicate_object THEN NULL;
END $$;

DO $$
BEGIN
  CREATE TYPE spec_source AS ENUM ('manufacturer', 'editor', 'user');
EXCEPTION
  WHEN duplicate_object THEN NULL;
END $$;

CREATE TABLE IF NOT EXISTS brand (
  id UUID PRIMARY KEY,
  name TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS guitar (
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

CREATE TABLE IF NOT EXISTS spec (
  id UUID PRIMARY KEY,
  code TEXT NOT NULL UNIQUE,
  label TEXT NOT NULL,
  value_type spec_value_type NOT NULL,
  unit TEXT,
  filterable BOOL NOT NULL DEFAULT FALSE,
  searchable BOOL NOT NULL DEFAULT FALSE,
  guitar_type guitar_type
);

CREATE TABLE IF NOT EXISTS spec_option (
  id UUID PRIMARY KEY,
  spec_id UUID NOT NULL REFERENCES spec(id),
  value TEXT NOT NULL,
  sort_order INT NOT NULL DEFAULT 0,
  UNIQUE (spec_id, value)
);

CREATE TABLE IF NOT EXISTS guitar_spec_value (
  guitar_id UUID NOT NULL REFERENCES guitar(id),
  spec_id UUID NOT NULL REFERENCES spec(id),
  value_text TEXT,
  value_number NUMERIC,
  value_bool BOOL,
  value_option_id UUID REFERENCES spec_option(id),
  source spec_source,
  PRIMARY KEY (guitar_id, spec_id)
);

CREATE TABLE IF NOT EXISTS guitar_media (
  id UUID PRIMARY KEY,
  guitar_id UUID NOT NULL REFERENCES guitar(id),
  kind media_kind NOT NULL,
  url TEXT NOT NULL,
  sort_order INT NOT NULL DEFAULT 0
);

CREATE INDEX IF NOT EXISTS guitar_type_idx ON guitar(type);
CREATE INDEX IF NOT EXISTS guitar_brand_idx ON guitar(brand_id);
CREATE INDEX IF NOT EXISTS guitar_spec_value_spec_idx ON guitar_spec_value(spec_id);
