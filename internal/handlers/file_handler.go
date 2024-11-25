package handlers

import (
	"fmt"
	"log"
	"logboard/internal/models"
	"net/http"
	"os"
	"path/filepath"
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

	logEntry := data.Status + ": " + formatted + " " + data.Data + "\n"
	_, err = file.WriteString(logEntry)
	if err != nil {
		log.Printf("Failed to write log entry: %v", err)
	}
	return err
}

func ServeLogFile(tab string, w http.ResponseWriter, r *http.Request) error {
	// Путь к файлу лога
	filePath := "logs/" + tab + ".log"

	// Попытка открыть файл. Если файл не существует, он будет создан.
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		log.Printf("Error opening/creating log file: %v", err)
		return err
	}
	defer file.Close()

	// Если файл только что был создан (первоначальный доступ), пустой ответ (пока нет данных).
	if stat, _ := file.Stat(); stat.Size() == 0 {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte{}) // Пустой ответ
		return nil
	}

	// Если файл существует и доступен, отдаем его в ответ на HTTP-запрос
	http.ServeFile(w, r, filePath)
	return nil
}

func RenameLogFile(oldName, newName string) error {
	logDir := "logs" // Путь к папке с логами

	oldPath := filepath.Join(logDir, fmt.Sprintf("%s.log", oldName))
	newPath := filepath.Join(logDir, fmt.Sprintf("%s.log", newName))

	// Проверяем, существует ли старый файл
	if _, err := os.Stat(oldPath); os.IsNotExist(err) {
		return fmt.Errorf("file %s does not exist", oldPath)
	}

	// Переименовываем файл
	if err := os.Rename(oldPath, newPath); err != nil {
		return fmt.Errorf("failed to rename file: %v", err)
	}

	return nil
}
