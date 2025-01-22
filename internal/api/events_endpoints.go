package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/asolheiro/kiosk-api/internal/sqlitestore"
	"github.com/go-chi/chi"
)

// Create a new event
// (POST /events)
func (api API) PostEvent(w http.ResponseWriter, r *http.Request) {
	var body sqlitestore.CreateEventParams

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "error deconding JSON", http.StatusBadRequest)
		return
	}

	event, err := api.repo.CreateEvent(r.Context(), body)
	if err != nil {
		http.Error(w, "error creating event", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(event)
}

// Get an event
// (GET /events/{eventId})
func (api API) GetEvent(w http.ResponseWriter, r *http.Request) {
	stringId := chi.URLParam(r, "eventId")
	eventId := strings.TrimSpace(stringId)

	event, err := api.repo.GetEvent(r.Context(), eventId)
	if err != nil {
		http.Error(w, "event not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(event); err != nil {
		http.Error(w, "error encoding response", http.StatusInternalServerError)
		return
	}
}

// List events
// (GET /events
func (api API) ListEvents(w http.ResponseWriter, r *http.Request) {
	events, err := api.repo.ListEvents(r.Context())
	if err != nil {
		http.Error(w, "error finding events", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(events); err != nil {
		http.Error(w, "error encoding response", http.StatusInternalServerError)
		return
	}
}

// Update an event
// (PUT /events/{eventId})
func (api API) PutEvent(w http.ResponseWriter, r *http.Request) {
	stringId := chi.URLParam(r, "eventId")
	eventId := strings.TrimSpace(stringId)

	_, err := api.repo.GetEvent(r.Context(), eventId)
	if err != nil {
		http.Error(w, "event not found", http.StatusNotFound)
		return
	}

	var body sqlitestore.UpdateEventParams
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "error deconding workload", http.StatusBadRequest)
	}
	body.ID = eventId

	event, err := api.repo.UpdateEvent(r.Context(), body)
	if err != nil {
		fmt.Println("err: ", err)
		http.Error(w, "error updating event", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(event); err != nil {
		http.Error(w, "error encoding response", http.StatusInternalServerError)
		return
	}
}

// Soft Delete an event
// (DELETE /events/{eventId})
func (api API) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	stringId := chi.URLParam(r, "eventId")
	eventId := strings.TrimSpace(stringId)

	err := api.repo.SoftDeleteEvent(r.Context(), eventId)
	if err != nil {
		http.Error(w, "error deleting event", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}
