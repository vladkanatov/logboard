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
