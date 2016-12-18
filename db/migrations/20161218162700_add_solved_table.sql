
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE solved (
    sovledId INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    friendId INTEGER,
    sovledDate DATE,
    FOREIGN KEY(friendId) REFERENCES friends(friendId));

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE solved;
