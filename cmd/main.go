package main

import (
	"log"
	"logboard/internal/handlers"
	"net/http"
	"os"
)

func init() {
	logDir := "logs"

	// Проверяем, существует ли директория logs
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		// Создаем директорию logs
		err := os.Mkdir(logDir, 0755)
		if err != nil {
			log.Fatalf("Failed to create logs directory: %v", err)
		}
		log.Printf("Created directory: %s", logDir)
	}
}

func main() {
	http.HandleFunc("/log", handlers.HandleRequest)         // Ожидает POST-запросы с `tab` в теле
	http.HandleFunc("/logs", handlers.HandleWebSocketLogs)  // Ожидает GET-запросы для чтения логов
	http.Handle("/", http.FileServer(http.Dir("./static"))) // Отдает статические файлы для фронтенда
	http.HandleFunc("/all-logs", handlers.HandleLogFile)
	http.HandleFunc("/rename-tab", handlers.HandleRenameTab)
	log.Println("Server started on :8000")
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}
}
