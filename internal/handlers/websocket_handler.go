package handlers

import (
	"bufio"
	"log"
	"logboard/internal/models"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			// Разрешаем соединения с любого источника
			return true
		},
	}

	// Используем карту, чтобы хранить клиентов для каждой вкладки
	clients = make(map[string]map[*websocket.Conn]bool) // вкладка -> клиенты
	mu      sync.Mutex                                  // мьютекс для безопасного доступа к клиентам
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

	// Добавляем нового клиента в список для данной вкладки
	mu.Lock()
	if clients[tab] == nil {
		clients[tab] = make(map[*websocket.Conn]bool)
	}
	clients[tab][conn] = true
	mu.Unlock()

	log.Printf("New WebSocket client connected for tab: %s", tab)

	// Отправляем содержимое файла логов при подключении
	if err := sendLogFile(tab, conn); err != nil {
		log.Printf("Error sending log file: %v", err)
	}

	// Постоянно отправляем новые строки из файла (tailing)
	if err := tailLogFile(tab, conn); err != nil {
		log.Printf("Error tailing log file: %v", err)
	}

	// Ожидаем сообщений от клиента (например, для обработки ошибок)
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			log.Printf("WebSocket read error: %v", err)
			break
		}
	}

	// Удаляем клиента при разрыве соединения
	mu.Lock()
	delete(clients[tab], conn)
	mu.Unlock()
	log.Printf("WebSocket client disconnected from tab: %s", tab)
}

// tailLogFile следит за файлом в реальном времени и отправляет новые строки через WebSocket.
func tailLogFile(tab string, conn *websocket.Conn) error {
	filePath := "logs/" + tab + ".log"
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Позиция курсора на начало файла
	stat, err := file.Stat()
	if err != nil {
		return err
	}

	offset := stat.Size()

	for {
		// Считываем новые строки, если файл увеличился
		file.Seek(offset, 0)

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

		// Обновляем offset для отслеживания новых данных в файле
		stat, err := file.Stat()
		if err != nil {
			return err
		}
		offset = stat.Size()

		// Пауза перед повторной проверкой (чтобы не нагружать процессор)
		time.Sleep(1 * time.Second)
	}
}

// broadcastToClients отправляет данные всем подключенным WebSocket клиентам для вкладки.
func broadcastToClients(tab string, data models.RequestData) {
	mu.Lock()
	defer mu.Unlock()

	for client := range clients[tab] {
		err := client.WriteJSON(data)
		if err != nil {
			log.Printf("Failed to send WebSocket message to client: %v", err)
			client.Close()
			delete(clients[tab], client)
		}
	}
}

// broadcastToAllClients отправляет данные всем подключенным WebSocket клиентам.
func broadcastToAllClients(data models.RequestData) {
	mu.Lock()
	defer mu.Unlock()

	for tab, clientsForTab := range clients {
		for client := range clientsForTab {
			err := client.WriteJSON(data)
			if err != nil {
				log.Printf("Failed to send WebSocket message to client on tab %s: %v", tab, err)
				client.Close()
				delete(clientsForTab, client)
			}
		}
	}
}
