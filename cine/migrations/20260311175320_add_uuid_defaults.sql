-- +goose Up
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

ALTER TABLE movies
ALTER COLUMN id SET DEFAULT gen_random_uuid();

ALTER TABLE genres
ALTER COLUMN genreid SET DEFAULT gen_random_uuid();

-- +goose Down
ALTER TABLE movies
ALTER COLUMN id DROP DEFAULT;

ALTER TABLE genres
ALTER COLUMN genreid DROP DEFAULT;