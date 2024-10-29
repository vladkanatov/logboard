# Указываем базовый образ с нужной версией Go
FROM golang:1.23 AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Инициализируем модуль Go и копируем исходные файлы
COPY main.go .
COPY database.go .

# Инициализация Go-модуля
RUN go mod init example.com/backend
RUN go mod tidy

# Собираем бинарный файл с именем backend
RUN go build -o backend main.go database.go

# Финальный образ для запуска бинарного файла
FROM debian:stable-slim

# Копируем бинарник из билдера
COPY --from=builder /app/backend /usr/local/bin/backend

# Устанавливаем рабочую директорию
WORKDIR /usr/local/bin

# Открываем порт, если он используется вашим приложением (например, 8080)
EXPOSE 8080

# Указываем команду для запуска приложения
CMD ["backend"]
