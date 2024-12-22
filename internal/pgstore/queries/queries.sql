-- name: GetUser :one
SELECT 
    * 
FROM users
WHERE id = $1
LIMIT 1;

-- name: ListUsers :many
SELECT
    *
FROM 
    users
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
UPDATE USERS
set 
    full_name = $2,
    email = $3,
    password = $4
WHERE 
    id = $1
RETURNING *;