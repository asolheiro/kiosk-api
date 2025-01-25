package api

import (
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi"

	"github.com/asolheiro/kiosk-api/internal/sqlitestore"
)

// Create a new config
// (POST /config)
func (api API) PostConfig(w http.ResponseWriter, r *http.Request) {
	var body sqlitestore.CreateConfigParams

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "error deconding JSON", http.StatusBadRequest)
		return
	}

	config, err := api.repo.CreateConfig(r.Context(), body)
	if err != nil {
		http.Error(w, "error creating config", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(config)
}

// Update an config
// (PUT /config/{configId})
func (api API) PutConfig(w http.ResponseWriter, r *http.Request) {
	stringId := chi.URLParam(r, "configId")
	configId := strings.TrimSpace(stringId)

	_, err := api.repo.GetGuest(r.Context(), configId)
	if err != nil {
		http.Error(w, "config not found", http.StatusNotFound)
		return
	}

	var body sqlitestore.UpdateGuestParams
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "error deconding workload", http.StatusBadRequest)
	}
	body.ID = configId

	config, err := api.repo.UpdateGuest(r.Context(), body)
	if err != nil {
		http.Error(w, "error updating config", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(config); err != nil {
		http.Error(w, "error encoding response", http.StatusInternalServerError)
		return
	}
}

type ImportDocumentParams struct {
	path  string
	event string
}

// import guests to config
// (POST /config/{configId}/import)
func (api API) ImportGuestsConfig(w http.ResponseWriter, r *http.Request) {
	stringId := chi.URLParam(r, "configId")
	configId := strings.TrimSpace(stringId)

	_, err := api.repo.GetConfig(r.Context(), configId)
	if err != nil {
		http.Error(w, "config not found", http.StatusNotFound)
		return
	}

	var body ImportDocumentParams
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "error decoding JSON", http.StatusBadRequest)
		return
	}

	file, err := os.Open(body.path)
	if err != nil {
		http.Error(w, "error opening file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		http.Error(w, "error reading CSV file", http.StatusInternalServerError)
		return
	}

	for _, record := range records {
		if len(record) < 2 {
			http.Error(w, "invalid CSV format", http.StatusBadRequest)
			return
		}

		guestParams := sqlitestore.CreateGuestParams{
			FullName:       record[0],
			Email:          sql.NullString{String: record[1], Valid: record[1] != ""},
			DocumentNumber: record[2],
			Occupation:     sql.NullString{String: record[3], Valid: record[1] != ""},
			ProfilePicture: sql.NullString{String: record[4], Valid: record[1] != ""},
			EventID:        body.event,
		}

		_, err := api.repo.CreateGuest(r.Context(), guestParams)
		if err != nil {
			http.Error(w, "error creating guest", http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "guests imported successfully"}`))
}

// Get an config
// (GET /config)
func (api API) GetConfig(w http.ResponseWriter, r *http.Request) {
	configs, err := api.repo.ListConfigs(r.Context())
	if err != nil {
		http.Error(w, "error finding config", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if len(configs) == 0 {
		http.Error(w, "no configs found", http.StatusBadRequest)
		return
	}

	if err := json.NewEncoder(w).Encode(configs[0]); err != nil {
		http.Error(w, "error encoding response", http.StatusInternalServerError)
		return
	}
}
