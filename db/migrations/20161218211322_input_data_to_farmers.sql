
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
INSERT INTO farmers VALUES(1, "Harish", "belgam", "karnataka", 9091090091);
INSERT INTO farmers VALUES(2, "girish", "raichur", "karnataka", 9091093423);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DELETE FROM farmers WHERE name in ("Harish", "girish");
