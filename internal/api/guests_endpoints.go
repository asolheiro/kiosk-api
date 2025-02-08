package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/asolheiro/kiosk-api/internal/sqlitestore"
	"github.com/go-chi/chi"
)

// Create a new guest
// (POST /guest)
func (api API) PostGuest(w http.ResponseWriter, r *http.Request) {
	var body sqlitestore.CreateGuestParams

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "error deconding JSON", http.StatusBadRequest)
		return
	}

	guest, err := api.repo.CreateGuest(r.Context(), body)
	if err != nil {
		http.Error(w, "error creating guest", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(guest)
}

// Get an guest
// (GET /guest/{guestId})
func (api API) GetGuest(w http.ResponseWriter, r *http.Request) {
	stringId := chi.URLParam(r, "guestId")
	guestId := strings.TrimSpace(stringId)

	guest, err := api.repo.GetGuest(r.Context(), guestId)
	if err != nil {
		http.Error(w, "guest not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(guest); err != nil {
		http.Error(w, "error encoding response", http.StatusInternalServerError)
		return
	}
}

// Get an guest by document number
// (GET /guest/{documentnumber}/document)
func (api API) GetGuestByDocument(w http.ResponseWriter, r *http.Request) {
	documentNumber := chi.URLParam(r, "documentNumber")

	guest, err := api.repo.GetGuestByDocumentNumber(r.Context(), documentNumber)
	if err != nil {
		http.Error(w, "guest not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(guest); err != nil {
		http.Error(w, "error encoding response", http.StatusInternalServerError)
		return
	}
}

// List guests
// (GET /guests)
func (api API) ListGuests(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	println("quering", query)
	if query == "" {
		http.Error(w, "error finding guests", http.StatusBadRequest)
		return

	}
	guests, err := api.repo.GetGuestByDocumentNumber(r.Context(), query)
	if err != nil {
		http.Error(w, "No result found", http.StatusBadRequest)
		return
	}
	fmt.Printf("guests: %+v\n", guests)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(guests); err != nil {
		println(err.Error())
		http.Error(w, "error encoding response", http.StatusInternalServerError)
		return
	}
}

// Update an guest
// (PUT /guest/{guestId})
func (api API) PutGuest(w http.ResponseWriter, r *http.Request) {
	stringId := chi.URLParam(r, "guestId")
	guestId := strings.TrimSpace(stringId)

	_, err := api.repo.GetGuest(r.Context(), guestId)
	if err != nil {
		http.Error(w, "guest not found", http.StatusNotFound)
		return
	}

	var body sqlitestore.UpdateGuestParams
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "error deconding workload", http.StatusBadRequest)
	}
	body.ID = guestId

	guest, err := api.repo.UpdateGuest(r.Context(), body)
	if err != nil {
		http.Error(w, "error updating guest", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(guest); err != nil {
		http.Error(w, "error encoding response", http.StatusInternalServerError)
		return
	}
}

// Soft Delete an guest
// (DELETE /guest/{guestId})
func (api API) DeleteGuest(w http.ResponseWriter, r *http.Request) {
	stringId := chi.URLParam(r, "guestId")
	guestId := strings.TrimSpace(stringId)

	err := api.repo.SoftDeleteGuest(r.Context(), guestId)
	if err != nil {
		http.Error(w, "error deleting guest", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}
