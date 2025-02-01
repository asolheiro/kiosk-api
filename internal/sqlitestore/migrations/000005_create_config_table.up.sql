CREATE TABLE IF NOT EXISTS config (
    id          TEXT NOT NULL PRIMARY KEY,
    template_image    TEXT,
    printer    TEXT,
    orientation    TEXT,
    updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
