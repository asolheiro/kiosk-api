package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/asolheiro/kiosk-api/internal/pgstore"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

// Create a new guest
// (POST /guest)
func (api API) PostGuest(w http.ResponseWriter, r *http.Request) {
	var body pgstore.CreateGuestParams

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
	stringId = strings.TrimSpace(stringId)
	guestId, err := uuid.Parse(stringId)
	if err != nil {
		http.Error(w, "invalid guestId", http.StatusBadRequest)
		return
	}

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
	guests, err := api.repo.ListGuests(r.Context())
	if err != nil {
		http.Error(w, "error finding guests", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(guests); err != nil {
		http.Error(w, "error encoding response", http.StatusInternalServerError)
		return 
	}
}


// Update an guest
// (PUT /guest/{guestId})
func (api API) PutGuest( w http.ResponseWriter, r *http.Request) {
	stringId := chi.URLParam(r, "guestId")
	stringId = strings.TrimSpace(stringId)
	guestId, err := uuid.Parse(stringId)
	
	if err != nil {
		http.Error(w, "invalid guestId", http.StatusBadRequest)
		return
	}

	_, err = api.repo.GetGuest(r.Context(), guestId)
	if err != nil {
		http.Error(w, "guest not found", http.StatusNotFound)
		return
	}

	var body pgstore.UpdateGuestParams
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
	stringId = strings.TrimSpace(stringId)
	
	guestId, err := uuid.Parse(stringId)
	if err != nil {
		http.Error(w, "invalid guestId", http.StatusBadRequest)
		return
	}

	err = api.repo.SoftDeleteGuest(r.Context(), guestId)
	if err != nil {
		http.Error(w, "error deleting guest", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}
