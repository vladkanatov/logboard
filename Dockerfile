# Используем официальный образ Go для сборки
FROM golang:1.20 as builder

# Устанавливаем рабочую директорию в контейнере
WORKDIR /app

# Копируем Go-файлы в контейнер
COPY . .

# Сборка Go-приложения
RUN go mod download
RUN go build -o backend main.go database.go

# Создаем минимальный образ для запуска приложения
FROM alpine:latest

# Устанавливаем зависимости SQLite для работы с базой данных
RUN apk --no-cache add sqlite

WORKDIR /app

# Копируем скомпилированное приложение из builder-а
COPY --from=builder /app/backend .

# Запуск приложения
CMD ["./backend"]
