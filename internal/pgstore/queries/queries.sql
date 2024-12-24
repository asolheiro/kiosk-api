--- USERS ---
-- name: GetUser :one
SELECT 
    * 
FROM users
WHERE id = $1 AND deleted_at IS NULL
LIMIT 1;

-- name: ListUsers :many
SELECT
    *
FROM 
    users
WHERE deleted_at IS NULL
ORDER BY
    full_name;

-- name: CreateUser :one
INSERT INTO users (
    full_name, email, password
) VALUES (
    $1, $2, $3
)
RETURNING *;

-- name: UpdateUser :one
UPDATE users
SET 
    full_name = $2,
    email = $3,
    password = $4,
    updated_at = CURRENT_TIMESTAMP
WHERE 
    id = $1
RETURNING *;

-- name: SoftDeleteUser :exec
UPDATE users
SET deleted_at = NOW()
WHERE id = $1 AND deleted_at IS NULL;

--- EVENTS --- 
-- name: GetEvent :one
SELECT 
    * 
FROM events
WHERE id = $1 AND deleted_at IS NULL
LIMIT 1;

-- name: ListEvents :many
SELECT
    *
FROM 
    events
WHERE deleted_at IS NULL
ORDER BY
    name;

-- name: CreateEvent :one
INSERT INTO events (
    name, primary_color, logo
) VALUES (
    $1, $2, $3
)
RETURNING *;

-- name: UpdateEvent :one
UPDATE events
SET 
    name = $2,
    primary_color = $3,
    logo = $4,
    updated_at = CURRENT_TIMESTAMP
WHERE 
    id = $1
RETURNING *;

-- name: SoftDeleteEvent :exec
UPDATE events
SET deleted_at = NOW()
WHERE id = $1 AND deleted_at IS NULL;