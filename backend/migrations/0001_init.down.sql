DROP INDEX IF EXISTS guitar_spec_value_spec_idx;
DROP INDEX IF EXISTS guitar_brand_idx;
DROP INDEX IF EXISTS guitar_type_idx;

DROP TABLE IF EXISTS guitar_media;
DROP TABLE IF EXISTS guitar_spec_value;
DROP TABLE IF EXISTS spec_option;
DROP TABLE IF EXISTS spec;
DROP TABLE IF EXISTS guitar;
DROP TABLE IF EXISTS brand;

DROP TYPE IF EXISTS spec_source;
DROP TYPE IF EXISTS media_kind;
DROP TYPE IF EXISTS spec_value_type;
DROP TYPE IF EXISTS guitar_type;
