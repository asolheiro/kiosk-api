package utils

import (
	"github.com/asolheiro/kiosk-api/internal/api"
	"github.com/go-chi/chi"
)



func UsersRouter(r chi.Router, api api.API) {
	r.Route("/users", func(r chi.Router) {
		r.Post("/", api.PostUser)
		r.Get("/", api.ListUsers)
		r.Get("/{userId}", api.GetUser)
		r.Put("/{userId}", api.PutUser)
	})
}