# Сервис Аутентификации

## Описание проекта

Сервис Аутентификации - это микросервис на языке Go, предназначенный для выдачи и обновления пар токенов доступа (Access Token) и обновления (Refresh Token) для аутентификации пользователей. Сервис предоставляет REST API для взаимодействия с другими системами.

## Основные функции

- **Выдача токенов**: Генерация пары Access и Refresh токенов для указанного идентификатора пользователя (UUID).
- **Обновление токенов**: Обновление пары токенов на основе предоставленных токенов с проверкой IP-адреса клиента и отправкой email-уведомления в случае изменения IP.

## Технологии

- **Язык программирования**: Go 1.24.0
- **Основные зависимости**:
  - `github.com/golang-jwt/jwt/v5` - для работы с JWT токенами
  - `golang.org/x/crypto` - для криптографических операций
  - `github.com/jackc/pgx/v5` - для работы с PostgreSQL

## Структура проекта

- **`cmd/`**: Точка входа в приложение (`main.go`)
- **`internal/`**: Внутренние пакеты сервиса
  - **`internal/auth/`**: Логика генерации и валидации токенов, отправки email
  - **`internal/http_handlers/`**: Обработчики HTTP-запросов
  - **`internal/storage/`**: Работа с базой данных (PostgreSQL), миграции

## API

Эндпоинты:

- **GET `/token/{userId}`**: Выдача пары токенов для указанного userId (UUID)
- **POST `/refresh`**: Обновление пары токенов

Сервер по умолчанию запускается на порте `8080`.

## Установка и запуск

1. **Клонируйте репозиторий**:
   ```bash
   git clone <URL репозитория>
   cd auth-service
   ```

2. **Установите зависимости**:
   ```bash
   go mod tidy
   ```

3. **Настройте окружение**:
   - Убедитесь, что у вас настроена база данных PostgreSQL
   - Настройте необходимые переменные окружения (см. документацию или код)

4. **Запустите сервис**:
   ```bash
   go run ./cmd/main.go
   ```

## Docker

Проект содержит `Dockerfile` для контейнеризации сервиса.

## Тестирование

Проект включает тесты для ключевых компонентов:
- Тесты для генерации токенов (`internal/auth/token_generator_test.go`)
- Тесты для работы с базой данных (`internal/storage/pg_storage_test.go`)
- Тесты для HTTP-обработчиков (`internal/http_handlers/http_handlers_test.go`)

Для запуска тестов используйте:
```bash
go test ./...
```
