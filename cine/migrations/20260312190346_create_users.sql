-- +goose Up
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(10) NOT NULL CHECK (role IN ('admin', 'user')),

    createdat TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updatedat TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    token TEXT,
    refreshtoken TEXT
);

CREATE TABLE user_favourite_genres (
    user_id UUID NOT NULL,
    genre_id UUID NOT NULL,

    PRIMARY KEY (user_id, genre_id),

    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (genre_id) REFERENCES genres(genreid) ON DELETE CASCADE
);


-- +goose Down
DROP TABLE IF EXISTS user_favourite_genres;
DROP TABLE IF EXISTS users;