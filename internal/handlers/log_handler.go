package handlers

import (
	"encoding/json"
	"log"
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

func HandleLogFile(w http.ResponseWriter, r *http.Request) {

	// Разрешаем CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	// Если это preflight-запрос, возвращаем пустой ответ
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Получаем значение параметра `tab` из URL
	tab := r.URL.Query().Get("tab")
	if tab == "" {
		http.Error(w, "Missing 'tab' parameter", http.StatusBadRequest)
		return
	}

	// Вызываем ServeLogFile и обрабатываем ошибки
	if err := ServeLogFile(tab, w, r); err != nil {
		http.Error(w, "Failed to serve log file", http.StatusInternalServerError)
		log.Printf("Error serving log file for tab '%s': %v", tab, err)
	}
}
