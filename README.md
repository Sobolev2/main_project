# semen_project — Backend на Go

Backend социальной сети на Go: REST API поверх PostgreSQL, слоистая архитектура (controllers → repository → routes, JWT-мидлварь на приватных маршрутах).

## Стек

| Категория | Технологии |
|---|---|
| Язык | Go 1.25 |
| HTTP | Gin |
| База данных | PostgreSQL, pgx/v5 |
| Миграции | Goose |
| Аутентификация | JWT (golang-jwt/jwt/v5) |
| Конфигурация | godotenv, caarlos0/env |
| Контейнеризация | Docker, docker-compose |

## Быстрый старт

### Docker (рекомендуется)

```bash
docker-compose up --build
```

Сервис — http://localhost:8080, PostgreSQL — localhost:5432.

### Локально, без Docker

```bash
cp .env.example .env          # заполнить значения
go run ./cmd/migrator migrations up
go run ./cmd/project
```

## Переменные окружения

Полный список — в `.env.example`: `APP_NAME`, `PUBLIC_API_HOST`, `PUBLIC_API_PORT`, `PG_HOST`, `PG_PORT`, `PG_USER`, `PG_PASS`, `PG_DB`, `PG_MAX_CONN`, `POSTGRES_DSN`, `JWT_SECRET`.

## Структура проекта

```
.
├── cmd
│   ├── project            # точка входа основного приложения
│   └── migrator           # отдельный бинарник для прогона миграций (goose)
├── internal
│   ├── app                # инициализация приложения
│   ├── config              # чтение конфигурации
│   ├── controllers          # HTTP-хендлеры (auth, user, post, comment, like, friend, follow, message, chat)
│   ├── dto                  # структуры запросов по доменам
│   ├── errs                 # общие ошибки
│   ├── middleware            # JWT auth middleware
│   ├── models                # доменные модели
│   ├── repository             # доступ к данным (pgx)
│   ├── routes                 # регистрация маршрутов по доменам
│   └── storage                # подключение к PostgreSQL
├── migrations               # SQL-миграции (Goose)
├── Dockerfile
└── docker-compose.yml
```

## Функциональность

Все приватные маршруты защищены JWT-мидлварью (`internal/middleware/auth.go`).

- **Аутентификация** — регистрация, логин, refresh-токен, logout, logout со всех устройств кроме текущего
- **Пользователи** — просмотр/обновление профиля, смена пароля, удаление аккаунта
- **Посты** — создание, обновление, удаление
- **Комментарии** — создание, обновление
- **Лайки** — на посты
- **Друзья** — система друзей (взаимная связь)
- **Подписки** — система подписчиков (одностороння связь)
- **Личные сообщения** — отправка, обновление
- **Чаты** — создание чата между двумя пользователями. Реализовано через обычный REST (нет зависимости под WebSocket), не real-time — стоит явно так и писать в резюме/README, чтобы не создавать впечатление real-time чата

## Разработка

```bash
go vet ./...
go test ./...
```
