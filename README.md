# QA Service

HTTP-сервис для вопросов и ответов. Позволяет создавать вопросы, оставлять ответы, получать списки и удалять записи. Хранение — PostgreSQL (GORM), валидация — go-playground/validator, миграции — goose.

## Возможности
- Создание/получение/удаление вопросов (`/questions`)
- Добавление/получение/удаление ответов (`/questions/{id}/answers`, `/answers/{id}`)
- Каскадное удаление ответов при удалении вопроса

## Технологии
- Go (HTTP + `net/http`)
- GORM + PostgreSQL
- goose для миграций
- go-playground/validator для валидации

## Подготовка окружения
`.env` (пример):
```
HTTP_HOST=0.0.0.0
HTTP_PORT=8081

DB_HOST=localhost
DB_PORT=5432
DB_USER=admin
DB_PASS=mypassword
DB_NAME=qaservice
```

## Локальный запуск
1) Создайте `.env`
2) Запустить базу данных:
``` docker compose up -d postgres ```
3) Запустить приложение: 
``` go run ./cmd/qaservice ```

## Docker Compose
В корне есть `docker-compose.yml`. Перед запуском подготовьте `.env`.
```
docker compose up --build
```
Порты по умолчанию: `5432` для БД и `8081` для приложения (см. файл compose, поменяйте маппинг при необходимости).

## Методы API

- `POST /questions` - создать вопрос `{ "text": "..." }`
- `GET /questions` - список вопросов
- `GET /questions/{id}` - вопрос с ответами
- `DELETE /questions/{id}` - удалить вопрос
- `POST /questions/{id}/answers` - создать ответ `{ "user_id": "...", "text": "..." }`
- `GET /answers/{id}` - получить ответ
- `DELETE /answers/{id}` - удалить ответ

Ошибки возвращаются в виде:
```
{
  "code": "validation_failed",
  "message": "Please check field values",
  "details": { "text": "required field" }
}
```

## Тесты
```
go test ./internal/repository/...
```
Тесты используют реальную базу `qaservice_test` (по умолчанию `postgres://admin:mypassword@localhost:5432`). При необходимости задайте свои переменные окружения перед запуском (см. `internal/repository/testdb`).

## Схема БД
Миграция `migrations/00001_init.sql` создаёт таблицы:
- `questions`: `id`, `text`, `created_at`
- `answers`: `id`, `question_id` (CASCADE delete), `user_id`, `text`, `created_at`
