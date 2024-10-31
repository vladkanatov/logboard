# Build Logging
Cервер на Go для получения и отображения логов сборок через WebSocket. Сервер также поддерживает HTTP-запросы для добавления новых логов и автоматически архивирует их раз в неделю.

## API

### Push log

```http
  POST /logs
```

| Parameter | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `tab` | `string` | **Required**. Tab for displaying the log. Three options: packages-common, sdk, eap |
| `status` | `string` | **Required**. Three options: success, error, info |
| `data` | `string` | **Required**. String with data |

## Настройка WebSocket

В файле JS (frontend/js/app.js) убедитесь, что строка подключения к WebSocket указывает на правильный хост и порт вашего сервера:

```javascript
const ws = new WebSocket("ws://<ваш_хост>:8080/ws");
```

## Установка

### Шаг 1: Соберите Docker-образ

```bash
docker build . -t build-logging
```

### Шаг 2: Запустите docker-compose.yml

```bash
docker compose up -d
```
