CREATE TABLE IF NOT EXISTS users (
    id          UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    full_name   VARCHAR(255) NOT NULL,
    email       VARCHAR(255) NOT NULL UNIQUE,
    password    VARCHAR(255) NOT NULL
);

---- create above / drop below ----

DROP TABLE IF EXISTS users;