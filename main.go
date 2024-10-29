package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// Структура для логов
type LogEntry struct {
	Tab    string `json:"tab"`
	Status string `json:"status"`
	Data   string `json:"data"`
}

var clients = make(map[*websocket.Conn]bool) // Открытые соединения WebSocket
var broadcast = make(chan LogEntry)          // Канал для передачи логов

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Обработчик для приема логов от Logstash
func handleLogs(w http.ResponseWriter, r *http.Request) {
	var entry LogEntry
	if err := json.NewDecoder(r.Body).Decode(&entry); err != nil {
		http.Error(w, "Некорректные данные", http.StatusBadRequest)
		return
	}

	// Сохранение в базе данных
	if err := insertLog(entry.Tab, entry.Status, entry.Data); err != nil {
		log.Println("Ошибка при записи в БД:", err)
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}

	// Отправка в канал для передачи через WebSocket
	broadcast <- entry
	w.WriteHeader(http.StatusOK)
}

// Обработчик WebSocket
func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Ошибка при обновлении WebSocket:", err)
		return
	}
	defer ws.Close()
	clients[ws] = true

	// Чтение сообщений из WebSocket (если нужно)
	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			delete(clients, ws)
			break
		}
	}
}

// Отправка новых логов всем подключенным WebSocket клиентам
func handleMessages() {
	for {
		logEntry := <-broadcast

		message, err := json.Marshal(logEntry)
		if err != nil {
			log.Println("Ошибка при сериализации сообщения:", err)
			continue
		}

		// Отправка сообщения всем клиентам
		for client := range clients {
			err := client.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				log.Println("Ошибка при отправке сообщения клиенту:", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}

func main() {
	// Инициализация базы данных
	initDB()
	defer db.Close()

	// Запуск горутины для обработки WebSocket сообщений
	go handleMessages()

	http.HandleFunc("/logs", handleLogs)
	http.HandleFunc("/ws", handleWebSocket)

	log.Println("Сервер запущен на порту 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
