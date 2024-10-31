package main

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/websocket"
)

// Data - структура данных
type Data struct {
	Tab    string `json:"tab"`
	Status string `json:"status"`
	Data   string `json:"data"`
}

var client *websocket.Conn // единственный клиент
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// WebSocket-коннектор
func handleConnections(w http.ResponseWriter, r *http.Request) {
	var err error
	client, err = upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Ошибка подключения к WebSocket:", err)
		return
	}
	defer client.Close()

	for {
		// Просто ждём, чтобы не блокировать горутину
		_, _, err := client.ReadMessage()
		if err != nil {
			log.Println("Ошибка чтения сообщения:", err)
			break
		}
	}
}

// Обработчик для POST-запросов на добавление логов
func handleLogs(w http.ResponseWriter, r *http.Request) {
	var data Data

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Не удалось прочитать тело запроса", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	err = json.Unmarshal(body, &data)
	if err != nil {
		http.Error(w, "Неверный формат JSON", http.StatusBadRequest)
		return
	}

	insertAndBroadcastData(data.Tab, data.Status, data.Data)

	w.WriteHeader(http.StatusCreated)
}

// Функция для добавления данных в лог-файл и отправки через WebSocket
func insertAndBroadcastData(tab, status, data string) {
	logDir := "./logs/"
	err := os.MkdirAll(logDir, os.ModePerm)
	if err != nil {
		log.Printf("Ошибка при создании директории логов: %v\n", err)
		return
	}

	logFile := fmt.Sprintf("%sbuild_log_%s.log", logDir, time.Now().Format("2006-01-02"))
	logEntry := fmt.Sprintf("%s [%s] %s: %s\n", time.Now().Format("15:04:05"), tab, status, data)

	// Запись в лог-файл
	f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Ошибка записи лога: %v\n", err)
		return
	}
	defer f.Close()

	if _, err := f.WriteString(logEntry); err != nil {
		log.Printf("Ошибка записи в лог-файл: %v\n", err)
	}

	// Отправка данных через WebSocket
	if client != nil {
		message, _ := json.Marshal(Data{Tab: tab, Status: status, Data: data})
		log.Println("Отправка сообщения через WebSocket:", string(message))
		err := client.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Println("Ошибка отправки сообщения:", err)
			// client.Close() // временно закомментируем, чтобы увидеть, не вызывает ли это проблемы
			client = nil
		}
	} else {
		log.Println("Нет активного WebSocket-соединения для отправки сообщения")
	}
}

// Ежедневное обновление
func dailyUpdate() {
	for {
		now := time.Now().Format("2006-01-02")
		insertAndBroadcastData("packages-common", "info", now)
		time.Sleep(24 * time.Hour)
	}
}

// Добавляем функцию для архивирования логов
func archiveLogs(files []string, archiveName string) error {
	archive, err := os.Create(archiveName)
	if err != nil {
		return fmt.Errorf("не удалось создать архив: %v", err)
	}
	defer archive.Close()

	zipWriter := zip.NewWriter(archive)
	defer zipWriter.Close()

	for _, file := range files {
		if err := addFileToZip(zipWriter, file); err != nil {
			return fmt.Errorf("ошибка при добавлении файла в архив: %v", err)
		}
		if err := os.Remove(file); err != nil { // Удаляем файл после добавления в архив
			log.Printf("не удалось удалить файл %s: %v", file, err)
		}
	}
	return nil
}

// Функция для добавления файла в архив
func addFileToZip(zipWriter *zip.Writer, filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	wr, err := zipWriter.Create(filepath.Base(filename))
	if err != nil {
		return err
	}

	_, err = io.Copy(wr, file)
	return err
}

// Функция еженедельной архивации
func weeklyArchive() {
	for {
		// Проверяем каждое воскресенье в полночь
		now := time.Now()
		nextWeek := now.AddDate(0, 0, 7-int(now.Weekday()))
		time.Sleep(time.Until(nextWeek))

		logDir := "./logs/"
		files, err := filepath.Glob(logDir + "build_log_*.log")
		if err != nil {
			log.Println("ошибка при получении списка логов:", err)
			continue
		}

		if len(files) > 0 {
			archiveName := fmt.Sprintf("%s/archive_%s.zip", logDir, now.Format("2006-01-02"))
			if err := archiveLogs(files, archiveName); err != nil {
				log.Println("ошибка при архивировании:", err)
			} else {
				log.Printf("архив %s создан", archiveName)
			}
		}
	}
}

func main() {
	http.HandleFunc("/ws", handleConnections)
	http.HandleFunc("/logs", handleLogs)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./frontend/index.html")
	})

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./frontend"))))

	go dailyUpdate()
	go weeklyArchive() // запускаем еженедельное архивирование

	insertAndBroadcastData("packages-common", "success", "2024-10-18 23:13:11 trs.auth 2.14.6 by antipov_sv")

	log.Println("Сервер запущен на порту 8080")
	http.ListenAndServe(":8080", nil)
}
