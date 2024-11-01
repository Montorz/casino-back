CREATE TABLE users
(
    id            SERIAL PRIMARY KEY UNIQUE,
    name          VARCHAR NOT NULL,
    login         VARCHAR NOT NULL UNIQUE,
    password_hash VARCHAR NOT NULL,
    balance       INT     NOT NULL
);