package handlers

import (
	"log"
	"logboard/internal/models"
	"net/http"
	"os"
	"sync"
	"time"
)

var logMutex sync.Mutex

// ensureLogFile создает файл лога, если он не существует.
func ensureLogFile(tab string) (*os.File, error) {
	// Проверка существования директории logs
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		if err := os.Mkdir("logs", 0755); err != nil {
			log.Printf("Failed to create logs directory: %v", err)
			return nil, err
		}
	}

	filename := "logs/" + tab + ".log"
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	return file, nil
}

// LogRequest записывает данные запроса в файл.
func LogRequest(data models.RequestData) error {
	logMutex.Lock()
	defer logMutex.Unlock()

	file, err := ensureLogFile(data.Tab)
	if err != nil {
		log.Printf("Failed to open log file: %v", err)
		return err
	}
	defer file.Close()

	// Добавляем разделитель для нового дня
	t := time.Now()
	formatted := t.Format("2006-01-02 15:04:05")

	fileInfo, err := file.Stat()
	if err == nil && fileInfo.Size() > 0 {
		if lastModified := fileInfo.ModTime(); lastModified.Day() != t.Day() {
			_, _ = file.WriteString("-----------------------------------\n" + t.Format("2006-01-02") + "\n")
		}
	}

	logEntry := data.Status + ": " + formatted + " - " + data.Data + "\n"
	_, err = file.WriteString(logEntry)
	if err != nil {
		log.Printf("Failed to write log entry: %v", err)
	}
	return err
}

// ServeLogFile отдает файл лога в ответе HTTP-запроса.
func ServeLogFile(tab string, w http.ResponseWriter, r *http.Request) error {
	filePath := "logs/" + tab + ".log"
	file, err := os.Open(filePath)
	if os.IsNotExist(err) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte{}) // Пустой ответ, если файл не существует
		return nil
	} else if err != nil {
		log.Printf("Error opening log file: %v", err)
		return err
	}
	defer file.Close()

	http.ServeFile(w, r, filePath)
	return nil
}
