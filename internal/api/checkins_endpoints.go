package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/asolheiro/kiosk-api/internal/pgstore"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

// Create a new check in
// (POST /checkin)
func (api API) PostCheckIn(w http.ResponseWriter, r *http.Request) {
	var body pgstore.CreateCheckInParams

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "error deconding JSON", http.StatusBadRequest)
		return
	}

	checkIn, err := api.repo.CreateCheckIn(r.Context(), body)
	if err != nil {
		http.Error(w, "error creating check in", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(checkIn)
}	

// Get a checkIn
// (GET /checkin/{checkIn})
func (api API) GetCheckIn(w http.ResponseWriter, r *http.Request) {
	stringId := chi.URLParam(r, "checkInId")
	stringId = strings.TrimSpace(stringId)
	checkInId, err := uuid.Parse(stringId)
	if err != nil {
		http.Error(w, "invalid checkInId ", http.StatusBadRequest)
		return
	}

	checkIn, err := api.repo.GetCheckIn(r.Context(), checkInId)
	if err != nil {
		http.Error(w, "checkIn not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(checkIn); err != nil {
		http.Error(w, "error encoding response", http.StatusInternalServerError)
		return
	}
}

// List users
// (GET /users
func (api API) ListCheckIns(w http.ResponseWriter, r *http.Request) {
	checkIn, err := api.repo.ListCheckIns(r.Context())
	if err != nil {
		http.Error(w, "error finding check in's", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(checkIn); err != nil {
		http.Error(w, "error encoding response", http.StatusInternalServerError)
		return 
	}
}


