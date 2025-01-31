# Указываем базовый образ для сборки
FROM golang:1.23 AS builder

# Создаем рабочую директорию
WORKDIR /app

# Копируем файлы проекта в контейнер
COPY . .

# Сборка CLI-приложения
RUN CGO_ENABLED=0 GOOS=linux go build -o logboard cmd/main.go

# Финальный этап: минимальный образ для запуска
FROM alpine:latest

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем скомпилированный бинарный файл из builder
COPY --from=builder /app/logboard /app/logboard

# Копируем папку static в стандартное место
COPY --from=builder /app/static /app/static

EXPOSE 8000

# Задание дефолтной команды при запуске
ENTRYPOINT ["/app/logboard"]