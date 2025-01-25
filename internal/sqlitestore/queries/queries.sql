--- USERS ---
-- name: GetUser :one
SELECT 
    * 
FROM users
WHERE id = ? AND deleted_at IS NULL
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
    ?, ?, ?
)
RETURNING *;

-- name: UpdateUser :one
UPDATE users
SET 
    full_name = :full_name,
    email = :email,
    password = :password,
    updated_at = CURRENT_TIMESTAMP
WHERE 
    id = :id
RETURNING *;


-- name: SoftDeleteUser :exec
UPDATE users
SET deleted_at = NOW()
WHERE id = ? AND deleted_at IS NULL;


--- EVENTS --- 
-- name: GetEvent :one
SELECT 
    * 
FROM 
    events
WHERE 
    id = ? AND deleted_at IS NULL
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
    ?, ?, ?
)
RETURNING *;

-- name: UpdateEvent :one
UPDATE events
SET 
    name = :name,
    primary_color = :primary_color,
    logo = :logo,
    updated_at = CURRENT_TIMESTAMP
WHERE 
    id = :id
RETURNING *;

-- name: SoftDeleteEvent :exec
UPDATE events
SET deleted_at = NOW()
WHERE id = ? AND deleted_at IS NULL;


--- GUESTS ---
-- name: CreateGuest :one
INSERT INTO guests (
    full_name, email, document_number, occupation, profile_picture, event_id
) VALUES (
    ?, ?, ?, ?, ?, ?
) RETURNING *;


-- name: GetGuest :one
SELECT 
    *
FROM 
    guests
WHERE 
    id = ? AND deleted_at IS NULL
LIMIT 
    1;


-- name: GetGuestByDocumentNumber :one
SELECT 
    *
FROM 
    guests
WHERE 
    document_number = ? AND deleted_at IS NULL
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
    full_name = :full_name,
    email = :email,
    occupation = :occupation,
    profile_picture = :profile_picture,
    document_number = :document_number,
    event_id = :event_id,
    updated_at = CURRENT_TIMESTAMP
WHERE 
    id = :id
RETURNING *;


-- name: SoftDeleteGuest :exec
UPDATE
    guests
SET
    deleted_at = NOW()
WHERE 
    id = ? AND deleted_at IS NULL;


--- CHECKIN ---
-- name: CreateCheckIn :one
INSERT INTO checkins (
    guest_id, event_id
) VALUES (
    ?, ?
) RETURNING *;


-- name: GetCheckIn :one
SELECT 
    *
FROM 
    checkins
WHERE 
    id = ?
LIMIT 
    1;


-- name: ListCheckIns :many
SELECT 
    *
FROM
    checkins
ORDER BY 
    created_at;


-- name: BulkInsertGuests :exec
INSERT INTO guests (
    full_name, email, document_number, occupation, profile_picture, event_id
) VALUES
    -- sqlc will dynamically replace these placeholders with actual data
    (?, ?, ?, ?, ?, ?);


--- CONFIG ---
-- name: GetConfig :one
SELECT 
    * 
FROM 
    config
WHERE 
    id = ?
LIMIT 
    1;

-- name: ListConfigs :many
SELECT
    *
FROM 
    config
ORDER BY
    updated_at;

-- name: CreateConfig :one
INSERT INTO config (
    id, template_image, printer, orientation
) VALUES (
    ?, ?, ?, ?
)
RETURNING *;

-- name: UpdateConfig :one
UPDATE config
SET 
    template_image = :template_image,
    printer = :printer,
    orientation = :orientation,
    updated_at = CURRENT_TIMESTAMP
WHERE 
    id = :id
RETURNING *;

-- name: DeleteConfig :exec
DELETE FROM config
WHERE id = ?;