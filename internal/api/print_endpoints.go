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
	PositionX     int    `json:"position_x"`
	PositionY     int    `json:"position_y"`
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

	err := imageService.CustomPrint(body.TemplateImage, body.Title, body.Description, body.PositionX, body.PositionY)
	if err != nil {
		http.Error(w, "error getting image", http.StatusBadRequest)
		return
	}

	err = imageService.DefaultPrint(body.TemplateImage, body.Title, body.Description, body.PositionX, body.PositionY)
	if err != nil {
		http.Error(w, "error getting image", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]string{"message": "Image saved"})
}
