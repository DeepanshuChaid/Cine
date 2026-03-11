-- +goose Up
CREATE TABLE genres (
    genreid UUID PRIMARY KEY,
    genrename TEXT NOT NULL
);

-- +goose Down
DROP TABLE genres;