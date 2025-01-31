package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/asolheiro/kiosk-api/internal/imageService"
)

type PostPrintParams struct {
	Title         string `json:"title"`
	Description   string `json:"description"`
	TemplateImage string `json:"template_image"`
}

// Create a new config
// (POST /printer)
func (api API) PostPrint(w http.ResponseWriter, r *http.Request) {
	var body PostPrintParams

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		fmt.Println(err)
		http.Error(w, "error deconding JSON", http.StatusBadRequest)
		return
	}

	mainConfig, err := api.repo.GetFirstConfig(r.Context())
	if err != nil {
		http.Error(w, "Config not found", http.StatusBadRequest)
		return
	}
	err = imageService.DefaultPrint(mainConfig.TemplateImage.String, body.Title, body.Description, int(mainConfig.PositionX.Int64), int(mainConfig.PositionY.Int64), int(mainConfig.FontSize.Int64))
	if err != nil {
		http.Error(w, "error getting image", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]string{"message": "Image saved"})
}
