-- +goose Up
CREATE TABLE movie_genres (
    movie_id UUID REFERENCES movies(id) ON DELETE CASCADE,
    genre_id UUID REFERENCES genres(genreid) ON DELETE CASCADE,
    PRIMARY KEY (movie_id, genre_id)
);

-- +goose Down
DROP TABLE movie_genres;