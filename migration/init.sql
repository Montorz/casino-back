CREATE TABLE users
(
    id            SERIAL PRIMARY KEY,
    name          VARCHAR NOT NULL,
    login         VARCHAR NOT NULL,
    password_hash VARCHAR NOT NULL,
    balance       INT     NOT NULL
);