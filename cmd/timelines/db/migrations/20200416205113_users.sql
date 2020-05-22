-- +gooser Up
-- SQL in this section is executed when the migration is applied.

CREATE TABLE IF NOT EXISTS users (
   id SERIAL,
   email text NOT NULL,
   hash text NOT NULL,
   created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
   updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
   is_admin boolean DEFAULT false,
   PRIMARY KEY (id)
);
-- +gooser Down
-- SQL in this section is executed when the migration is rolled back.

DROP TABLE IF EXISTS users;
