-- +goose Up
ALTER TABLE groups ADD views integer DEFAULT 0;

-- +goose Down
ALTER TABLE groups DROP COLUMN views;


