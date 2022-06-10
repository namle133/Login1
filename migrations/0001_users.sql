CREATE TABLE users (
    username varchar(255)  PRIMARY KEY,
    password bytea  NOT NULL,
    email varchar(255)  NOT NULL
);

