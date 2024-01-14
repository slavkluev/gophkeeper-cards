CREATE TABLE IF NOT EXISTS cards
(
    id       INTEGER PRIMARY KEY,
    number   TEXT    NOT NULL,
    cvv      TEXT    NOT NULL,
    month    TEXT    NOT NULL,
    year     TEXT    NOT NULL,
    info     TEXT    NOT NULL DEFAULT '',
    user_uid INTEGER NOT NULL
);
