CREATE TABLE IF NOT EXISTS guests (
    id              TEXT NOT NULL PRIMARY KEY,
    full_name       TEXT NOT NULL UNIQUE,
    email           TEXT,
    document_number TEXT NOT NULL UNIQUE,
    occupation      TEXT,
    profile_picture TEXT,
    event_id        TEXT NOT NULL REFERENCES events(id),
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at      TIMESTAMP
);