-- +gooser Up
-- SQL in this section is executed when the migration is applied.

ALTER TABLE groups
    DROP CONSTRAINT user_id_fk;

-- +gooser Down
-- SQL in this section is executed when the migration is rolled back.
ALTER TABLE groups
    ADD CONSTRAINT user_id_fk FOREIGN KEY (user_id) REFERENCES users(id);
