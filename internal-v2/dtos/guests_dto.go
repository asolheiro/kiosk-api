package dto

import (
	"database/sql"
	"time"

	"github.com/asolheiro/kiosk-api/internal-v2/sqlitestore"
)

type Guest struct {
	ID             string         `db:"id" json:"id"`
	FullName       string         `db:"full_name" json:"full_name"`
	Email          sql.NullString `db:"email" json:"email"`
	DocumentNumber string         `db:"document_number" json:"document_number"`
	Occupation     sql.NullString `db:"occupation" json:"occupation"`
	ProfilePicture sql.NullString `db:"profile_picture" json:"profile_picture"`
	EventID        string         `db:"event_id" json:"event_id"`
	CreatedAt      time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time      `db:"updated_at" json:"updated_at"`
}

func FromSQLGuest(g sqlitestore.Guest) Guest {
	return Guest{
		ID: g.ID,
		FullName: g.FullName,
		Email: g.Email,
		DocumentNumber: g.DocumentNumber,
		Occupation: g.Occupation,
		ProfilePicture: g.ProfilePicture,
		EventID: g.EventID,
		CreatedAt: g.CreatedAt,
		UpdatedAt: g.UpdatedAt,
	}
}

func NewGuestsList(guests []sqlitestore.Guest) []Guest {
	result := make([]Guest, len(guests))
	for i, g := range guests {
		result[i] = FromSQLGuest(g)
	}
	return result
}