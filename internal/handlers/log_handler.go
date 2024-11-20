package handlers

import (
	"encoding/json"
	"logboard/internal/models"
	"net/http"
)

// HandleRequest обрабатывает входящие POST-запросы для записи логов.
func HandleRequest(w http.ResponseWriter, r *http.Request) {
	var requestData models.RequestData
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if requestData.Tab == "" {
		http.Error(w, "Tab parameter is required in the request body", http.StatusBadRequest)
		return
	}

	if err := LogRequest(requestData); err != nil {
		http.Error(w, "Error logging request", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
