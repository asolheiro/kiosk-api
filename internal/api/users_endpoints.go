package api

import (
	"encoding/json"
	"net/http"

	"github.com/asolheiro/kiosk-api/internal/pgstore"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

// Create a new user
// (POST /users)
func (api API) PostUser(w http.ResponseWriter, r *http.Request) {
	var body pgstore.CreateUserParams

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "error deconding JSON", http.StatusBadRequest)
		return
	}

	_, err := api.repo.CreateUser(r.Context(), body)
	if err != nil {
		http.Error(w, "error creating user", http.StatusBadRequest)
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode("user created successfully")
}	

// Get an user
// (GET /users/{userId})
func (api API) GetUser(w http.ResponseWriter, r *http.Request) {
	stringId := chi.URLParam(r, "userId")
	userId, err := uuid.Parse(stringId)
	if err != nil {
		http.Error(w, "error parsing parse", http.StatusBadRequest)
		return 
	}

	user, err := api.repo.GetUser(r.Context(), userId)
	if err != nil {
		http.Error(w, "error finding user", http.StatusBadRequest)
		return
	}

	jsonData, err := json.Marshal(user)
	if err != nil {
		http.Error(w, "error enconding json", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(jsonData)
}	