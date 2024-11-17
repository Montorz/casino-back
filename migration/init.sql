CREATE TABLE users
(
    id         SERIAL  NOT NULL UNIQUE,
    name       VARCHAR NOT NULL,
    login      VARCHAR NOT NULL UNIQUE,
    password   VARCHAR NOT NULL,
    balance    INTEGER NOT NULL,
    avatar_url VARCHAR
);

CREATE TABLE transactions
(
    id           SERIAL                                          NOT NULL UNIQUE,
    user_id      INTEGER REFERENCES users (id) ON DELETE CASCADE NOT NULL,
    type         VARCHAR                                         NOT NULL,
    amount       INTEGER                                         NOT NULL,
    created_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP             NOT NULL
);

CREATE TABLE slots
(
    id      SERIAL  NOT NULL UNIQUE,
    name    VARCHAR NOT NULL UNIQUE,
    min_bet INTEGER NOT NULL,
    max_bet INTEGER NOT NULL
);

CREATE TABLE games
(
    id           SERIAL                                          NOT NULL UNIQUE,
    user_id      INTEGER REFERENCES users (id) ON DELETE CASCADE NOT NULL,
    slot_id      INTEGER REFERENCES slots (id) ON DELETE CASCADE NOT NULL,
    name         VARCHAR                                         NOT NULL,
    bet_amount   INTEGER                                         NOT NULL,
    coefficient  FLOAT                                           NOT NULL,
    win_amount   INTEGER                                         NOT NULL,
    created_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP             NOT NULL
);