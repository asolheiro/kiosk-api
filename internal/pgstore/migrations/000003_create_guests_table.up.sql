CREATE TABLE IF NOT EXISTS guests (
    id              UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    full_name       VARCHAR(255) NOT NULL UNIQUE,
    email           VARCHAR(255) NULL,
    document_number VARCHAR(255) NOT NULL UNIQUE,
    occupation      VARCHAR(255) NULL,
    profile_picture VARCHAR(255) NULL,
    event_id        UUID NOT NULL REFERENCES events(id),
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at      TIMESTAMP NULL
);