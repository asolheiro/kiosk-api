package api

import (
	"encoding/json"
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
	stringId = strings.TrimSpace(stringId)
	userId, err := uuid.Parse(stringId)
	if err != nil {
		http.Error(w, "invalid userId", http.StatusBadRequest)
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
		return
	}
}

// List users
// (GET /users
func (api API) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := api.repo.ListUsers(r.Context())
	if err != nil {
		http.Error(w, "error finding users", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(users); err != nil {
		http.Error(w, "error encoding response", http.StatusInternalServerError)
		return 
	}
}


// Update an user
// (PUT /users/{userId})
func (api API) PutUser( w http.ResponseWriter, r *http.Request) {
	stringId := chi.URLParam(r, "userId")
	stringId = strings.TrimSpace(stringId)
	userId, err := uuid.Parse(stringId)
	
	if err != nil {
		http.Error(w, "invalid userId", http.StatusBadRequest)
		return
	}

	_, err = api.repo.GetUser(r.Context(), userId)
	if err != nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	var body pgstore.UpdateUserParams
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "error deconding workload", http.StatusBadRequest)
	}
	body.ID = userId

	user, err := api.repo.UpdateUser(r.Context(), body)
	if err != nil {
		http.Error(w, "error updating user", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, "error encoding response", http.StatusInternalServerError)
		return
	}
}

// Soft Delete an user
// (DELETE /users/{userId})
func (api API) DeleteUser(w http.ResponseWriter, r *http.Request) {
	stringId := chi.URLParam(r, "userId")
	stringId = strings.TrimSpace(stringId)
	
	userId, err := uuid.Parse(stringId)
	if err != nil {
		http.Error(w, "invalid userId", http.StatusBadRequest)
		return
	}

	err = api.repo.SoftDeleteUser(r.Context(), userId)
	if err != nil {
		http.Error(w, "error deleting user", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}
