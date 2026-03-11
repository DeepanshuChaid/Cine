-- +goose Up
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE movies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    imdbid TEXT NOT NULL,
    title TEXT NOT NULL,
    posterpath TEXT NOT NULL,
    youtubeid TEXT NOT NULL,
    adminreview TEXT NOT NULL
);

-- +goose Down
DROP TABLE movies;