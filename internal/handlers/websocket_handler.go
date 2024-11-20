package handlers

import (
	"bufio"
	"log"
	"logboard/internal/models"
	"net/http"
	"os"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			// Разрешить соединения с любого источника
			return true
		},
	}

	clients = make(map[*websocket.Conn]bool) // Активные WebSocket клиенты
	mu      sync.Mutex                       // Мьютекс для безопасного доступа к клиентам
)

// HandleWebSocketLogs обрабатывает WebSocket соединения.
func HandleWebSocketLogs(w http.ResponseWriter, r *http.Request) {
	tab := r.URL.Query().Get("tab")
	if tab == "" {
		http.Error(w, "Tab parameter is required", http.StatusBadRequest)
		return
	}

	// Обновление до WebSocket соединения
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}
	defer conn.Close()

	// Добавляем нового клиента
	mu.Lock()
	clients[conn] = true
	mu.Unlock()

	log.Printf("New WebSocket client connected for tab: %s", tab)

	// Отправляем содержимое файла логов при подключении
	if err := sendLogFile(tab, conn); err != nil {
		log.Printf("Error sending log file: %v", err)
	}

	// Удаляем клиента при разрыве соединения
	mu.Lock()
	delete(clients, conn)
	mu.Unlock()
}

// sendLogFile отправляет содержимое файла логов клиенту.
func sendLogFile(tab string, conn *websocket.Conn) error {
	filePath := "logs/" + tab + ".log"
	file, err := os.Open(filePath)
	if os.IsNotExist(err) {
		// Если файл не существует, просто отправляем пустой ответ
		return conn.WriteMessage(websocket.TextMessage, []byte("No logs available"))
	} else if err != nil {
		return err
	}
	defer file.Close()

	// Считываем файл построчно и отправляем клиенту
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if err := conn.WriteMessage(websocket.TextMessage, []byte(line)); err != nil {
			log.Printf("Error sending log line: %v", err)
			return err
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Error reading log file: %v", err)
		return err
	}

	return nil
}

// broadcastToClients отправляет данные всем подключенным WebSocket клиентам.
func broadcastToClients(data models.RequestData) {
	mu.Lock()
	defer mu.Unlock()

	for client := range clients {
		err := client.WriteJSON(data)
		if err != nil {
			log.Printf("Failed to send WebSocket message: %v", err)
			client.Close()
			delete(clients, client)
		}
	}
}
