package dto

import (
	"time"

	"github.com/asolheiro/kiosk-api/internal-v2/sqlitestore"
)

type Event struct {
	ID           string    `db:"id" json:"id"`
	Name         string    `db:"name" json:"name"`
	PrimaryColor string    `db:"primary_color" json:"primary_color"`
	Logo         string    `db:"logo" json:"logo"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}

func FromSQLCEvent(e sqlitestore.Event) Event {
	return Event{
		ID:           e.ID,
		Name:         e.Name,
		PrimaryColor: e.PrimaryColor,
		Logo:         e.Logo,
		CreatedAt:    e.CreatedAt,
		UpdatedAt:    e.UpdatedAt,
	}
}

func NewEventResponseList(events []sqlitestore.Event) []Event {
	result := make([]Event, len(events))
	for i, u := range events {
		result[i] = FromSQLCEvent(u)
	}
	return result
}

type CreateEventRequest struct {
	Name         string `db:"name" json:"name"`
	PrimaryColor string `db:"primary_color" json:"primary_color"`
	Logo         string `db:"logo" json:"logo"`
}

