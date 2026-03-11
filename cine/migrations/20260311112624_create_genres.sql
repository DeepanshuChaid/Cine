-- +goose Up
CREATE TABLE genres (
    genreid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    genrename TEXT NOT NULL
);

-- +goose Down
DROP TABLE genres;