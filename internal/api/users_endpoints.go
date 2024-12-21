package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

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

	user, err := api.repo.CreateUser(r.Context(), body)
	if err != nil {
		http.Error(w, "error creating user", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(user)
}	

// Get an user
// (GET /users/{userId})
func (api API) GetUser(w http.ResponseWriter, r *http.Request) {
	stringId := chi.URLParam(r, "userId")
	fmt.Printf("Extracted userId: %s\n", stringId)
	
	stringId = strings.TrimSpace(stringId)
	userId, err := uuid.Parse(stringId)
	if err != nil {
		http.Error(w, "invalid userId format", http.StatusBadRequest)
		return
	}

	user, err := api.repo.GetUser(r.Context(), userId)
	if err != nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, "error encoding response", http.StatusInternalServerError)
	}
}
