// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: queries.sql

package pgstore

import (
	"context"

	"github.com/google/uuid"
)

const createEvent = `-- name: CreateEvent :one
INSERT INTO events (
    name, primary_color, logo
) VALUES (
    $1, $2, $3
)
RETURNING id, name, primary_color, logo, created_at, updated_at, deleted_at
`

type CreateEventParams struct {
	Name         string `db:"name" json:"name"`
	PrimaryColor string `db:"primary_color" json:"primary_color"`
	Logo         string `db:"logo" json:"logo"`
}

func (q *Queries) CreateEvent(ctx context.Context, arg CreateEventParams) (Event, error) {
	row := q.db.QueryRow(ctx, createEvent, arg.Name, arg.PrimaryColor, arg.Logo)
	var i Event
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.PrimaryColor,
		&i.Logo,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	return i, err
}

const createUser = `-- name: CreateUser :one
INSERT INTO users (
    full_name, email, password
) VALUES (
    $1, $2, $3
)
RETURNING id, full_name, email, password, created_at, updated_at, deleted_at
`

type CreateUserParams struct {
	FullName string `db:"full_name" json:"full_name"`
	Email    string `db:"email" json:"email"`
	Password string `db:"password" json:"password"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRow(ctx, createUser, arg.FullName, arg.Email, arg.Password)
	var i User
	err := row.Scan(
		&i.ID,
		&i.FullName,
		&i.Email,
		&i.Password,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	return i, err
}

const getEvent = `-- name: GetEvent :one
SELECT 
    id, name, primary_color, logo, created_at, updated_at, deleted_at 
FROM events
WHERE id = $1 AND deleted_at IS NULL
LIMIT 1
`

// - EVENTS ---
func (q *Queries) GetEvent(ctx context.Context, id uuid.UUID) (Event, error) {
	row := q.db.QueryRow(ctx, getEvent, id)
	var i Event
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.PrimaryColor,
		&i.Logo,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	return i, err
}

const getUser = `-- name: GetUser :one
SELECT 
    id, full_name, email, password, created_at, updated_at, deleted_at 
FROM users
WHERE id = $1 AND deleted_at IS NULL
LIMIT 1
`

// - USERS ---
func (q *Queries) GetUser(ctx context.Context, id uuid.UUID) (User, error) {
	row := q.db.QueryRow(ctx, getUser, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.FullName,
		&i.Email,
		&i.Password,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	return i, err
}

const listEvents = `-- name: ListEvents :many
SELECT
    id, name, primary_color, logo, created_at, updated_at, deleted_at
FROM 
    events
WHERE deleted_at IS NULL
ORDER BY
    name
`

func (q *Queries) ListEvents(ctx context.Context) ([]Event, error) {
	rows, err := q.db.Query(ctx, listEvents)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Event
	for rows.Next() {
		var i Event
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.PrimaryColor,
			&i.Logo,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.DeletedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listUsers = `-- name: ListUsers :many
SELECT
    id, full_name, email, password, created_at, updated_at, deleted_at
FROM 
    users
WHERE deleted_at IS NULL
ORDER BY
    full_name
`

func (q *Queries) ListUsers(ctx context.Context) ([]User, error) {
	rows, err := q.db.Query(ctx, listUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.FullName,
			&i.Email,
			&i.Password,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.DeletedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const softDeleteEvent = `-- name: SoftDeleteEvent :exec
UPDATE events
SET deleted_at = NOW()
WHERE id = $1 AND deleted_at IS NULL
`

func (q *Queries) SoftDeleteEvent(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.Exec(ctx, softDeleteEvent, id)
	return err
}

const softDeleteUser = `-- name: SoftDeleteUser :exec
UPDATE users
SET deleted_at = NOW()
WHERE id = $1 AND deleted_at IS NULL
`

func (q *Queries) SoftDeleteUser(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.Exec(ctx, softDeleteUser, id)
	return err
}

const updateEvent = `-- name: UpdateEvent :one
UPDATE events
SET 
    name = $2,
    primary_color = $3,
    logo = $4,
    updated_at = CURRENT_TIMESTAMP
WHERE 
    id = $1
RETURNING id, name, primary_color, logo, created_at, updated_at, deleted_at
`

type UpdateEventParams struct {
	ID           uuid.UUID `db:"id" json:"id"`
	Name         string    `db:"name" json:"name"`
	PrimaryColor string    `db:"primary_color" json:"primary_color"`
	Logo         string    `db:"logo" json:"logo"`
}

func (q *Queries) UpdateEvent(ctx context.Context, arg UpdateEventParams) (Event, error) {
	row := q.db.QueryRow(ctx, updateEvent,
		arg.ID,
		arg.Name,
		arg.PrimaryColor,
		arg.Logo,
	)
	var i Event
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.PrimaryColor,
		&i.Logo,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	return i, err
}

const updateUser = `-- name: UpdateUser :one
UPDATE users
SET 
    full_name = $2,
    email = $3,
    password = $4,
    updated_at = CURRENT_TIMESTAMP
WHERE 
    id = $1
RETURNING id, full_name, email, password, created_at, updated_at, deleted_at
`

type UpdateUserParams struct {
	ID       uuid.UUID `db:"id" json:"id"`
	FullName string    `db:"full_name" json:"full_name"`
	Email    string    `db:"email" json:"email"`
	Password string    `db:"password" json:"password"`
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error) {
	row := q.db.QueryRow(ctx, updateUser,
		arg.ID,
		arg.FullName,
		arg.Email,
		arg.Password,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.FullName,
		&i.Email,
		&i.Password,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	return i, err
}
