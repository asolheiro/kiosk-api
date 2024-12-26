CREATE TABLE IF NOT EXISTS checkins (
    id          UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    event_id    UUID NOT NULL REFERENCES events(id),
    guest_id    UUID NOT NULL REFERENCES guests(id),
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
