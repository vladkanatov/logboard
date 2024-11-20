package main

import (
	"log"
	"logboard/internal/handlers"
	"net/http"
)

func main() {
	http.HandleFunc("/log", handlers.HandleRequest)         // Ожидает POST-запросы с `tab` в теле
	http.HandleFunc("/logs", handlers.HandleWebSocketLogs)  // Ожидает GET-запросы для чтения логов
	http.Handle("/", http.FileServer(http.Dir("./static"))) // Отдает статические файлы для фронтенда
	log.Println("Server started on :8000")
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}
}
