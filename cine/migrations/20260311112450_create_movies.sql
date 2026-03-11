-- +goose Up
CREATE TABLE movies (
    id UUID PRIMARY KEY,
    imdbid TEXT NOT NULL,
    title TEXT NOT NULL,
    posterpath TEXT NOT NULL,
    youtubeid TEXT NOT NULL,
    adminreview TEXT NOT NULL
);

-- +goose Down
DROP TABLE movies;
