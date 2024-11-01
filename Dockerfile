# Стадия 1: сборка
FROM golang:1.23-alpine AS builder

# Устанавливаем зависимости для сборки
RUN apk add --no-cache git gcc musl-dev

# Создаем рабочую директорию
WORKDIR /app

# Копируем файлы go.mod и go.sum для установки зависимостей
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем остальные файлы проекта
COPY . .

# Собираем бинарный файл с флагами оптимизации
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o /app-server .

# Стадия 2: финальный контейнер
FROM alpine:latest

# Создаем пользователя для безопасного запуска
RUN adduser -D -u 1000 appuser

# Устанавливаем директорию для приложения
WORKDIR /home/appuser

# Копируем бинарник из стадии сборки
COPY --from=builder /app-server .

# Копируем статические файлы (frontend)
COPY --from=builder /app/static /home/appuser/static

# Меняем права на исполняемый файл и статические файлы
RUN chown -R appuser:appuser /home/appuser

# Переходим на пользователя с ограниченными правами
USER appuser

# Экспортируем порт для доступа к приложению
EXPOSE 8080

# Запуск приложения
CMD ["./app-server"]
