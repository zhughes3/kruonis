-- +gooser Up
-- SQL in this section is executed when the migration is applied.

ALTER TABLE groups
    ADD COLUMN private boolean,
    ADD COLUMN user_id integer,
    ADD CONSTRAINT user_id_fk FOREIGN KEY (user_id) REFERENCES users(id);


-- +gooser Down
-- SQL in this section is executed when the migration is rolled back.

ALTER TABLE groups
    DROP CONSTRAINT user_id_fk,
    DROP COLUMN user_id,
    DROP COLUMN private;