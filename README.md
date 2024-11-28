# YouTube Thumbnail gRPC Service

## Описание
gRPC-сервис для загрузки превью видеороликов YouTube. При повторном запросе для одного и того же видео возвращается закэшированное изображение из SQLite.

## Особенности
- Кэширование запросов с использованием SQLite.
- Асинхронная загрузка изображений через CLI с флагом `--async`.
- Реализация с использованием `gRPC` и `Protobuf`.

---

## Установка и запуск

### Предварительные требования
- Установленный Go (версии 1.20 и выше)
- `protoc` (Protocol Buffers Compiler)
- SQLite3

### Шаги установки
1. Клонируйте репозиторий:
   ```bash
   git clone <repository_url>
   cd youtube-thumbnail-grpc
Установите зависимости:

```bash
go mod tidy```
Сгенерируйте gRPC-код из .proto:

```bash
protoc --go_out=. --go-grpc_out=. proto/thumbnail.proto```

Запуск сервера
Запустите сервер:

```bash
go run cmd/server/main.go```
Сервер начнет слушать на localhost:50051.


