
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE problems (
    problemId INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    farmerId INTEGER NOT NULL,
    problemDesc VARCHAR(500),
    postedDate DATE,
    isSolved INTEGER,
    isDeleted INTEGER,
    FOREIGN KEY(farmerId) REFERENCES farmers(farmerId));

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE problems;
