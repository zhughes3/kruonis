-- +goose Up

CREATE TABLE IF NOT EXISTS group_views (
    id INTEGER NOT NULL,
    time TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT current_timestamp,
    FOREIGN KEY(id) REFERENCES groups(id) ON DELETE CASCADE
);
ALTER TABLE groups ADD views integer DEFAULT 0;

-- +goose Down

ALTER TABLE groups DROP COLUMN views;
DROP TABLE IF EXISTS group_views;
