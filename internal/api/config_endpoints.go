package api

import (
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi"
	"github.com/google/uuid"

	"github.com/asolheiro/kiosk-api/internal/sqlitestore"
)

type ConfigResponse struct {
	ID            string `json:"id"`
	TemplateImage string `json:"template_image"`
	Printer       string `json:"printer"`
	Orientation   string `json:"orientation"`
	PositionX     int    `json:"position_x"`
	PositionY     int    `json:"position_y"`
	FontSize      int    `json:"font_size"`
	WidthLimiter  int    `json:"width_limiter"`
}

func (ConfigResponse) fromConfig(config sqlitestore.Config) ConfigResponse {
	return ConfigResponse{
		ID:            config.ID,
		TemplateImage: config.TemplateImage.String,
		Printer:       config.Printer.String,
		Orientation:   config.Orientation.String,
		PositionX:     int(config.PositionX.Int64),
		PositionY:     int(config.PositionY.Int64),
		FontSize:      int(config.FontSize.Int64),
		WidthLimiter:  int(config.WidthLimiter.Int64),
	}
}

type CreateConfigParams struct {
	TemplateImage string `db:"template_image" json:"template_image"`
	Printer       string `db:"printer" json:"printer"`
	Orientation   string `db:"orientation" json:"orientation"`
	PositionX     int    `db:"position_x" json:"position_x"`
	PositionY     int    `db:"position_y" json:"position_y"`
	FontSize      int    `db:"font_size" json:"font_size"`
	WidthLimiter  int    `db:"width_limiter" json:"width_limiter"`
}

// Create a new config
// (POST /config)
func (api API) PostConfig(w http.ResponseWriter, r *http.Request) {
	var body CreateConfigParams
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		println(err)
		http.Error(w, "error deconding JSON", http.StatusBadRequest)
		return
	}

	config, err := api.repo.CreateConfig(r.Context(), sqlitestore.CreateConfigParams{
		ID:            uuid.New().String(),
		TemplateImage: body.TemplateImage,
		Printer:       body.Printer,
		Orientation:   body.Orientation,
		PositionX:     body.PositionX,
		PositionY:     body.PositionY,
		FontSize:      body.FontSize,
		WidthLimiter:  body.WidthLimiter,
	})
	if err != nil {
		http.Error(w, "error creating config", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(ConfigResponse{}.fromConfig(config)); err != nil {
		println(err.Error())
		http.Error(w, "error encoding response", http.StatusInternalServerError)
		return
	}
}

type UpdateConfigParams struct {
	TemplateImage string `db:"template_image" json:"template_image"`
	Printer       string `db:"printer" json:"printer"`
	Orientation   string `db:"orientation" json:"orientation"`
	PositionX     int    `db:"position_x" json:"position_x"`
	PositionY     int    `db:"position_y" json:"position_y"`
	FontSize      int    `db:"font_size" json:"font_size"`
	WidthLimiter  int    `db:"width_limiter" json:"width_limiter"`
}

// Update an config
// (PUT /config/{configId})
func (api API) PutConfig(w http.ResponseWriter, r *http.Request) {
	stringId := chi.URLParam(r, "configId")
	configId := strings.TrimSpace(stringId)

	_, err := api.repo.GetConfig(r.Context(), configId)
	if err != nil {
		println(err.Error())
		http.Error(w, "config not found", http.StatusBadRequest)
		return
	}

	var body UpdateConfigParams
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "error deconding workload", http.StatusBadRequest)
	}

	config, err := api.repo.UpdateConfig(r.Context(), sqlitestore.UpdateConfigParams{
		TemplateImage: sql.NullString{String: body.TemplateImage, Valid: body.TemplateImage != ""},
		Printer:       sql.NullString{String: body.Printer, Valid: body.Printer != ""},
		Orientation:   sql.NullString{String: body.Orientation, Valid: body.Orientation != ""},
		PositionX:     sql.NullInt64{Int64: int64(body.PositionX), Valid: body.PositionX != 0},
		PositionY:     sql.NullInt64{Int64: int64(body.PositionY), Valid: body.PositionY != 0},
		FontSize:      sql.NullInt64{Int64: int64(body.FontSize), Valid: body.FontSize != 0},
		WidthLimiter:  sql.NullInt64{Int64: int64(body.WidthLimiter), Valid: body.WidthLimiter != 0},
		ID:            configId,
	})
	if err != nil {
		http.Error(w, "error updating config", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(ConfigResponse{}.fromConfig(config)); err != nil {
		http.Error(w, "error encoding response", http.StatusInternalServerError)
		return
	}
}

type ImportDocumentParams struct {
	PathFile string `json:"path_file"`
}

// import guests to config
// (POST /import)
func (api API) ImportGuestsConfig(w http.ResponseWriter, r *http.Request) {

	var body ImportDocumentParams
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		println(err.Error())
		http.Error(w, "error decoding JSON", http.StatusBadRequest)
		return
	}

	file, err := os.Open(body.PathFile)
	if err != nil {
		println(err.Error())
		http.Error(w, "error opening file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		println(err.Error())
		http.Error(w, "error reading CSV file", http.StatusInternalServerError)
		return
	}

	// Skip the header
	if len(records) > 0 {
		records = records[1:]
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
			EventID:        uuid.New().String(),
		}

		_, err := api.repo.CreateGuest(r.Context(), guestParams)
		if err != nil {
			println(err.Error())
			http.Error(w, "error creating guest", http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message": "guests imported successfully"}`))
}

// Get an config
// (GET /config)
func (api API) GetConfig(w http.ResponseWriter, r *http.Request) {
	configs, err := api.repo.ListConfigs(r.Context())
	if err != nil {
		println(err.Error())
		http.Error(w, "error finding config", http.StatusInternalServerError)
		return
	}

	if len(configs) == 0 {
		http.Error(w, "no configs found", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(ConfigResponse{}.fromConfig(configs[0])); err != nil {
		println(err.Error())
		http.Error(w, "error encoding response", http.StatusInternalServerError)
		return
	}
}
