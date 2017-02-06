
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE friends (
    friendId INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    name VARCHAR(50),
    phoneNumber VARCHAR(10),
    isDeleted INTEGER );

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE friends;
