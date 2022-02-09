-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE IF NOT EXISTS hzads_ads (
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(1024), 
    type VARCHAR(64),
    position VARCHAR(64),
    content VARCHAR(2048),
    status INT,
    sort INT,
    created BIGINT,
    updated BIGINT
);

CREATE TABLE hzads_users (
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(1024), 
    username VARCHAR(1024),
    password VARCHAR(1024),
    last_login BIGINT
);



-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE hzads__ads;
DROP TABLE hzads__users;