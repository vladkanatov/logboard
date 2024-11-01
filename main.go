package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

// Структура для хранения данных запроса
type RequestData struct {
	Tab    string `json:"tab"`
	Status string `json:"status"`
	Data   string `json:"data"`
}

var logMutex sync.Mutex

// Функция для инициализации директории логов и создания файла при необходимости
func ensureLogFile(tab string) (*os.File, error) {
	// Проверка существования директории logs
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		if err := os.Mkdir("logs", 0755); err != nil {
			log.Fatalf("Failed to create logs directory: %v", err)
		}
	}

	filename := "logs/" + tab + ".log"
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	return file, nil
}

// Функция для записи лога, добавляющая разделитель для нового дня
func logRequest(data RequestData) error {
	logMutex.Lock()
	defer logMutex.Unlock()

	file, err := ensureLogFile(data.Tab)
	if err != nil {
		log.Printf("Failed to open log file: %v", err)
		return err
	}
	defer file.Close()

	// Добавляем разделитель для нового дня в 00:00
	now := time.Now()
	fileInfo, err := file.Stat()
	if err == nil && fileInfo.Size() > 0 {
		// Проверяем, нужно ли добавить новую дату
		if lastModified := fileInfo.ModTime(); lastModified.Day() != now.Day() {
			_, _ = file.WriteString("-----------------------------------\n" + now.Format("2006-01-02") + "\n")
		}
	}

	logEntry := data.Status + ": " + data.Data + "\n"
	_, err = file.WriteString(logEntry)
	if err != nil {
		log.Printf("Failed to write log entry: %v", err)
	}
	return err
}

// Обработчик для входящих данных
func handleRequest(w http.ResponseWriter, r *http.Request) {
	var requestData RequestData
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if requestData.Tab == "" {
		http.Error(w, "Tab parameter is required in the request body", http.StatusBadRequest)
		return
	}

	if err := logRequest(requestData); err != nil {
		http.Error(w, "Error logging request", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Обработчик для отправки логов в ответе на запросы фронтенда
func handleLogRead(w http.ResponseWriter, r *http.Request) {
	tab := r.URL.Query().Get("tab")
	if tab == "" {
		http.Error(w, "Tab parameter is required", http.StatusBadRequest)
		return
	}

	file, err := os.Open("logs/" + tab + ".log")
	if os.IsNotExist(err) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte{}) // Пустой ответ, если файл не существует
		return
	} else if err != nil {
		http.Error(w, "Error opening log file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	http.ServeFile(w, r, "logs/"+tab+".log")
}

func main() {
	http.HandleFunc("/log", handleRequest)                  // Ожидает POST-запросы с `tab` в теле
	http.HandleFunc("/logs", handleLogRead)                 // Ожидает GET-запросы для чтения логов
	http.Handle("/", http.FileServer(http.Dir("./static"))) // Отдает статические файлы для фронтенда
	log.Println("Server started on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
