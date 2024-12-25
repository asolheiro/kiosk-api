package utils

import (
	"github.com/asolheiro/kiosk-api/internal/api"
	"github.com/go-chi/chi"
)



func UsersRouter(r chi.Router, api api.API) {
	r.Route("/user", func(r chi.Router) {
		r.Post("/", api.PostUser)
		r.Get("/", api.ListUsers)
		r.Get("/{userId}", api.GetUser)
		r.Put("/{userId}", api.PutUser)
		r.Delete("/{userId}", api.DeleteUser)
	})
}

func EventsRouter(r chi.Router, api api.API) {
	r.Route("/event", func(r chi.Router) {
		r.Post("/", api.PostEvent)
		r.Get("/", api.ListEvents)
		r.Get("/{eventId}", api.GetEvent)
		r.Put("/{eventId}", api.PutEvent)
		r.Delete("/{eventId}", api.DeleteEvent)
	})
}

func GuestsRouter(r chi.Router, api api.API) {
	r.Route("/guest", func(r chi.Router) {
		r.Post("/", api.PostGuest)
		r.Get("/", api.ListGuests)
		r.Get("/{guestId}", api.GetGuest)
		r.Get("/{documentNumber}/document", api.GetGuestByDocument)
		r.Put("/{guestId}", api.PutGuest)
		r.Delete("/{guestId}", api.DeleteGuest)
	})
}