# Auth-service

Сервис авторизации по email + OTP.

Что реализовано:
- `POST /auth/start` - отправка OTP на почту
- `POST /auth/verify` - проверка OTP, создание локального пользователя при первом входе, отправка Kafka-события на создание пользователя, возврат токена
- `POST /auth/check` - проверка Bearer-токена
- `POST /auth/logout` - удаление токена по email
- `GET /health` - healthcheck

## Пример запросов

### Отправить OTP

```bash
curl -X POST http://localhost:8080/auth/start \
  -H 'Content-Type: application/json' \
  -d '{"email":"user@example.com"}'
```

### Подтвердить OTP

```bash
curl -X POST http://localhost:8080/auth/verify \
  -H 'Content-Type: application/json' \
  -d '{"email":"user@example.com","code":"ABC123"}'
```

### Проверить токен

```bash
curl -X POST http://localhost:8080/auth/check \
  -H 'Authorization: Bearer <token>'
```

### Выйти

```bash
curl -X POST http://localhost:8080/auth/logout \
  -H 'Content-Type: application/json' \
  -d '{"email":"user@example.com"}'
```

## Важные env

Смотри `.env.example`.

Минимум для запуска:
- Postgres
- SMTP сервер или Mailpit
- Kafka broker
- миграции из `migrations/`


## Persistence

The service uses GORM with PostgreSQL for `users`, `tokens`, and `mails` repositories. After pulling the changes, run `go mod tidy` to fetch the new GORM dependencies before `go build ./...`.
