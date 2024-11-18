package test

import (
	"bytes"
	"logboard/internal/handlers"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestHandleRequest_StatusCode(t *testing.T) {
	// Подготовка запроса
	requestBody := []byte(`{"tab": "test-tab", "status": "success", "data": "Test log entry"}`)
	req := httptest.NewRequest(http.MethodPost, "/log", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	// Создание ResponseRecorder для записи ответа
	rr := httptest.NewRecorder()

	// Вызов обработчика
	handlers.HandleRequest(rr, req)

	// Проверка статуса ответа
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %v, got %v", http.StatusOK, rr.Code)
	}

	// Удаляем файл лога после теста
	defer os.RemoveAll("logs")
}

func TestHandleLogRead_StatusCode(t *testing.T) {
	// Создание тестового лог-файла
	os.Mkdir("logs", 0755)
	defer os.RemoveAll("logs") // Удаляем директорию после теста

	logFilePath := "logs/test-tab.log"
	err := os.WriteFile(logFilePath, []byte("Test log entry\n"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test log file: %v", err)
	}

	// Подготовка запроса
	req := httptest.NewRequest(http.MethodGet, "/logs?tab=test-tab", nil)
	rr := httptest.NewRecorder()

	// Вызов обработчика
	handlers.HandleLogRead(rr, req)

	// Проверка статуса ответа
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %v, got %v", http.StatusOK, rr.Code)
	}
}

func TestHandleLogRead_StatusCode_NoFile(t *testing.T) {
	// Подготовка запроса для несуществующего файла
	req := httptest.NewRequest(http.MethodGet, "/logs?tab=nonexistent", nil)
	rr := httptest.NewRecorder()

	// Вызов обработчика
	handlers.HandleLogRead(rr, req)

	// Проверка статуса ответа
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %v, got %v", http.StatusOK, rr.Code)
	}
}
