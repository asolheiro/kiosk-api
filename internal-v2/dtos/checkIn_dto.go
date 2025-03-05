package dto

import (
	"time"

	"github.com/asolheiro/kiosk-api/internal-v2/sqlitestore"
)

type CheckIn struct {
	ID        string    `db:"id" json:"id"`
	EventID   string    `db:"event_id" json:"event_id"`
	GuestID   string    `db:"guest_id" json:"guest_id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

func FromSQLCCheckIn(c sqlitestore.Checkin) (CheckIn) {
	return CheckIn{
		ID: c.ID,
		EventID: c.EventID,
		GuestID: c.GuestID,
		CreatedAt: c.CreatedAt,
	}
}

func NewCheckInsList(checkIns []sqlitestore.Checkin) ([]CheckIn) {
	checkInsDto := make([]CheckIn, len(checkIns))

	for i, checkIn := range checkIns {
		checkInsDto[i] = FromSQLCCheckIn(checkIn)
	}
	return checkInsDto
}