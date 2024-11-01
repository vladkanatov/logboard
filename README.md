# Build Logging
Cервер на Go для получения и отображения логов сборок через WebSocket. Сервер также поддерживает HTTP-запросы для добавления новых логов и автоматически архивирует их раз в неделю.

## API

### Push log

```http
  POST /log
```

| Parameter | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `tab` | `string` | **Required**. Tab for displaying the log. Three options: packages-common, sdk, eap |
| `status` | `string` | **Required**. Three options: success, error, info |
| `data` | `string` | **Required**. String with data |

