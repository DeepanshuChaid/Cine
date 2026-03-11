-- +goose Up
ALTER TABLE genres
ADD CONSTRAINT unique_genrename UNIQUE (genrename);

-- +goose Down
ALTER TABLE genres
DROP CONSTRAINT unique_genrename;