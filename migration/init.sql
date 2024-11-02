CREATE TABLE users
(
    id       SERIAL PRIMARY KEY UNIQUE,
    name     VARCHAR NOT NULL,
    login    VARCHAR NOT NULL UNIQUE,
    password VARCHAR NOT NULL,
    balance  INT     NOT NULL
);

CREATE TABLE transactions
(
    id      SERIAL                                          NOT NULL UNIQUE,
    user_id INTEGER REFERENCES users (id) ON DELETE CASCADE NOT NULL,
    type    VARCHAR                                         NOT NULL,
    amount  INTEGER                                         NOT NULL
);