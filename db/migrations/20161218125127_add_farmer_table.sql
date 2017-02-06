
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE farmers (
    farmerId INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    name VARCHAR(50),
    district VARCHAR(100),
    state VARCHAR(30),
    phoneNumber VARCHAR(10),
    isDeleted INTEGER);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE farmers;
