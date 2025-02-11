CREATE TABLE IF NOT EXISTS events (
    id              TEXT NOT NULL PRIMARY KEY,
    name            TEXT NOT NULL UNIQUE,
    primary_color   TEXT NOT NULL,
    logo            TEXT NOT NULL,
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at      TIMESTAMP
);