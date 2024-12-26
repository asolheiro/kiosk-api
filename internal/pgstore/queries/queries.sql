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
FROM 
    events
WHERE 
    id = $1 AND deleted_at IS NULL
LIMIT 
    1;

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


--- GUESTS ---
-- name: CreateGuest :one
INSERT INTO guests (
    full_name, email, document_number, occupation, profile_picture, event_id
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;


-- name: GetGuest :one
SELECT 
    *
FROM 
    guests
WHERE 
    id = $1 AND deleted_at IS NULL
LIMIT 
    1;


-- name: GetGuestByDocumentNumber :one
SELECT 
    *
FROM 
    guests
WHERE 
    document_number = $1 AND deleted_at IS NULL
LIMIT 
    1;

-- name: ListGuests :many
SELECT 
    *
FROM
    guests
WHERE
    deleted_at IS NULL
ORDER BY 
    full_name;


-- name: UpdateGuest :one
UPDATE 
    guests
SET 
    full_name = $2,
    email = $3,
    occupation = $4,
    profile_picture = $5,
    document_number = $6,
    event_id = $7,
    updated_at = CURRENT_TIMESTAMP
WHERE 
    id = $1
RETURNING *;


-- name: SoftDeleteGuest :exec
UPDATE
    guests
SET
    deleted_at = NOW()
WHERE 
    id = $1 AND deleted_at IS NULL;


--- CHECKIN ---
-- name: CreateCheckIn :one
INSERT INTO checkins (
    guest_id, event_id
) VALUES (
    $1, $2
) RETURNING *;


-- name: GetCheckIn :one
SELECT 
    *
FROM 
    checkins
WHERE 
    id = $1
LIMIT 
    1;


-- name: ListCheckIns :many
SELECT 
    *
FROM
    checkins
ORDER BY 
    created_at;
