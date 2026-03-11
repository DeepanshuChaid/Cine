-- +goose Up
CREATE TABLE rankings (
    movie_id UUID PRIMARY KEY REFERENCES movies(id) ON DELETE CASCADE,
    rankingvalue INTEGER NOT NULL,
    rankingname TEXT NOT NULL CHECK (
        rankingname IN ('Excellent','Good','Okay','Bad','Terrible')
    )
);

-- +goose Down
DROP TABLE rankings;