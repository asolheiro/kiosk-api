package utils

import (
	"encoding/json"
	"net/http"
)

func HealthCheck (w http.ResponseWriter, r *http.Request) {
	healthy := map[string]string{
		"status": "healthy",
	}
	
	jsonData, err := json.Marshal(healthy)
	if err != nil {
		http.Error(w, "error encoding json", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}