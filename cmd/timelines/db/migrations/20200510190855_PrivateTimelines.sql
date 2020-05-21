-- +gooser Up
-- SQL in this section is executed when the migration is applied.

ALTER TABLE groups
    ADD COLUMN private boolean DEFAULT false,
    ADD COLUMN user_id integer,
    ADD COLUMN uuid uuid,
    ADD CONSTRAINT user_id_fk FOREIGN KEY (user_id) REFERENCES users(id);

ALTER TABLE users
    ADD COLUMN is_admin boolean DEFAULT false;

-- +gooser Down
-- SQL in this section is executed when the migration is rolled back.

ALTER TABLE users
    DROP COLUMN is_admin;

ALTER TABLE groups
    DROP CONSTRAINT user_id_fk,
    DROP COLUMN uuid,
    DROP COLUMN user_id,
    DROP COLUMN private;

