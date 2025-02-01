package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/asolheiro/kiosk-api/internal/imageService"
	"github.com/go-chi/chi"
	"github.com/godoes/printers"
)

type PostPrintParams struct {
	Title       string `json:"title"`
	Description string `json:"description"`
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
	result, err := imageService.DefaultPrint(mainConfig.TemplateImage.String, body.Title, body.Description, int(mainConfig.PositionX.Int64), int(mainConfig.PositionY.Int64), int(mainConfig.FontSize.Int64))
	if err != nil {
		http.Error(w, "error getting image", http.StatusBadRequest)
		return
	}

	err = imageService.Print(result, mainConfig.Printer.String)
	if err != nil {
		http.Error(w, "error getting image", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]string{"message": "Image saved"})
}

// Create a new config
// (GET /printers)
func (api API) GetPrinters(w http.ResponseWriter, r *http.Request) {
	names, err := printers.ReadNames()
	if err != nil {
		http.Error(w, "error getting printers", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(names)
}

// Create a new config
// (GET /printers/{printerName})
func (api API) GetPrinterInfo(w http.ResponseWriter, r *http.Request) {
	printerName := chi.URLParam(r, "printerName")
	println("Printer name: ", printerName)
	p, err := printers.Open(printerName)
	if err != nil {
		println(err)
		http.Error(w, "error getting printers", http.StatusBadRequest)
		return
	}

	defer func() {
		_ = p.Close()
	}()

	info, err := p.DriverInfo()
	if err != nil {
		println(err.Error())
		http.Error(w, "error getting printers", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(info)
}
