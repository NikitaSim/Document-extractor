# Document Extractor

Go-сервис для загрузки PDF/DOCX, извлечения текста и получения структурированных реквизитов через Ollama.

## Что извлекается

Сервис возвращает сведения только из раздела **«УЧАСТНИК ПРОЕКТА»** и игнорирует **«ОПЕРАТОР КЖНО»** и остальные разделы. Для отсутствующих полей возвращается пустая строка. Значения не должны додумываться моделью.

Поля: `organization`, `legal_address`, `ogrn`, `inn`, `kpp`, `bank_account`, `director_position`, `director_name`.

## Поток

`POST /api/v1/extract` → PDF/DOCX → извлечение текста → Ollama `/api/chat` → JSON по JSON Schema → HTTP response.

## Требования

- Go 1.23+
- Ollama с моделью, поддерживающей JSON Schema в `format` (например, `qwen3:8b`)

## Локальный запуск

```bash
ollama pull qwen3:8b
go mod tidy
go run ./cmd/server
```

Ollama по умолчанию ожидается на `http://localhost:11434`.

## Docker

```bash
docker compose up --build
```

Загрузка модели в контейнер Ollama:

```bash
docker compose exec ollama ollama pull qwen3:8b
```

## API

Проверка: `GET /health`

Извлечение PDF:

```bash
curl -X POST http://localhost:8080/api/v1/extract -F "file=@./document.pdf"
```

или DOCX:

```bash
curl -X POST http://localhost:8080/api/v1/extract -F "file=@./document.docx"
```

## Структура

```text
cmd/server             запуск HTTP-сервера
internal/config        конфигурация из переменных окружения
internal/handler       HTTP API
internal/parser        извлечение текста из PDF/DOCX
internal/ai/client.go  клиент Ollama Chat API
internal/ai/prompts.go системный промпт и JSON Schema
internal/ai/extractor.go преобразование ответа модели в DTO
internal/service       orchestration
internal/model         DTO документа и результата
```
