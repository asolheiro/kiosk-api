CREATE TABLE IF NOT EXISTS checkins (
    id          TEXT NOT NULL PRIMARY KEY,
    event_id    TEXT NOT NULL REFERENCES events(id),
    guest_id    TEXT NOT NULL REFERENCES guests(id),
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
