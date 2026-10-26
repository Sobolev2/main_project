# Social Platform

Бэкенд социальной сети в стиле Instagram на Go, построенный как набор взаимодействующих сервисов: REST API, отдельный сервис уведомлений, обмен событиями через Kafka и синхронные вызовы по gRPC. К бэкенду прилагается SPA-клиент на React + TypeScript.

![Go](https://img.shields.io/badge/Go-1.25-00ADD8?logo=go&logoColor=white)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15-4169E1?logo=postgresql&logoColor=white)
![Redis](https://img.shields.io/badge/Redis-7-DC382D?logo=redis&logoColor=white)
![Kafka](https://img.shields.io/badge/Kafka-4.1-231F20?logo=apachekafka&logoColor=white)
![gRPC](https://img.shields.io/badge/gRPC-Protobuf-244c5a)
![Docker](https://img.shields.io/badge/Docker-compose-2496ED?logo=docker&logoColor=white)
![React](https://img.shields.io/badge/React-19-61DAFB?logo=react&logoColor=black)

## Архитектура

```mermaid
flowchart LR
    FE[Frontend<br/>React + Vite] -->|REST /api| BE

    subgraph Backend
        BE[backend-service<br/>Gin]
        NS[notification-service]
    end

    BE --> PG1[(PostgreSQL<br/>platform)]
    BE --> R[(Redis<br/>кэш пользователей)]
    BE -->|событие message.sent| K{{Kafka<br/>topic: message-events}}
    K --> NS
    NS --> PG2[(PostgreSQL<br/>notifications)]
    BE -->|gRPC: список, счётчик, прочитано| NS
```

**Как работают уведомления.** Когда пользователь отправляет сообщение, `backend-service` сохраняет его в БД и публикует событие `message.sent` в Kafka. `notification-service` читает событие, создаёт уведомление в собственной базе и подтверждает offset только после успешной записи. Клиент получает уведомления через REST `backend-service`, который обращается к сервису уведомлений по gRPC.

### Ключевые инженерные решения

| Решение | Зачем |
|---|---|
| Отдельный сервис уведомлений со своей БД | Независимое развитие и масштабирование, слабая связность через события |
| Идемпотентность консьюмера: `UNIQUE (type, message_id)` + `ON CONFLICT DO NOTHING`, offset коммитится после записи | At-least-once доставка без дублей уведомлений при повторной обработке |
| Ключ сообщения в Kafka = `receiver_id` | Порядок событий для одного получателя сохраняется в рамках партиции |
| Отправка сообщения не зависит от Kafka: ошибка публикации логируется, но сообщение уже сохранено | Недоступность брокера не ломает основной сценарий чата |
| gRPC-вызовы с таймаутом 5 с, при недоступности сервиса клиент получает `502` | Сбой одного сервиса не подвешивает API |
| Контракт сервисов в отдельном модуле `proto/` | Единый источник правды, сгенерированный код общий для клиента и сервера |
| Go workspace (`go.work`): `backend-service`, `notification-service`, `proto` | Несколько модулей в одном репозитории без костылей с `replace` |
| Access JWT + refresh-токены, в БД хранится только хэш refresh-токена | Утечка БД не даёт готовых токенов; поддержка выхода с устройства и со всех устройств кроме текущего |
| Пароли хранятся как bcrypt-хэши | Стандартная защита учётных данных |
| Redis-кэш пользователя по `id` и `username` с инвалидацией при обновлении и удалении | Разгрузка БД на самых частых чтениях |
| Контроллеры зависят от интерфейсов репозиториев, в тестах подставляются моки | Хэндлеры тестируются без базы данных |
| Миграции через goose, мигратор запускается как одноразовый сервис в compose (`service_completed_successfully`) | Схема применяется автоматически и до старта API |
| Многоступенчатые Docker-сборки (`golang:1.25-alpine` → `alpine:3.20`) | Небольшие итоговые образы |

## Возможности

| Область | Что реализовано |
|---|---|
| Авторизация | Регистрация, логин, обновление токена, выход с устройства, выход со всех устройств кроме текущего |
| Пользователи | Профиль (`/me`), обновление, удаление, смена пароля, поиск, получение по id и по username |
| Посты | Создание, чтение, обновление, удаление, посты пользователя, лента |
| Лента | Сначала посты друзей, затем остальные |
| Лайки и комментарии | Лайк и снятие лайка, счётчики, список лайкнувших, статус лайка, CRUD комментариев |
| Подписки и друзья | Подписка и отписка, списки и счётчики. При взаимной подписке пользователи становятся друзьями |
| Чаты | Чат 1-на-1, отправка, правка и удаление сообщений, отметка «прочитано», счётчик непрочитанных |
| Уведомления | О новых сообщениях: список, счётчик непрочитанных, отметка «прочитано» |
| Клиент | Регистрация, вход, дашборд, создание поста |

## Стек

**Backend:** Go 1.25, Gin, pgx/v5 (пул соединений), goose, go-redis, kafka-go, gRPC + Protobuf, golang-jwt, bcrypt, slog, testify

**Frontend:** React 19, TypeScript, Vite, React Router

**Инфраструктура:** PostgreSQL 15, Redis 7, Apache Kafka 4.1, Docker / Docker Compose, манифесты Kubernetes

## Структура репозитория

```
Backend/
├── backend-service/            # основной REST API
│   ├── cmd/
│   │   ├── project/            # точка входа API
│   │   └── migrator/           # запуск миграций (goose)
│   ├── internal/
│   │   ├── app/                # сборка приложения: PostgreSQL, Redis, Kafka, gRPC-клиент
│   │   ├── config/             # конфигурация из переменных окружения
│   │   ├── controllers/        # HTTP-обработчики + тесты
│   │   ├── repository/         # доступ к данным, интерфейсы *Repo
│   │   ├── routes/             # маршруты (public / private)
│   │   ├── middleware/         # проверка Bearer JWT
│   │   ├── kafka/, events/     # продюсер и описание событий
│   │   └── models/, dto/, mocks/, storage/, errs/
│   └── migrations/             # SQL-миграции (goose)
├── notification-service/       # сервис уведомлений
│   ├── cmd/notification/       # точка входа: gRPC-сервер + Kafka consumer
│   ├── internal/               # config, controllers, grpc, kafka, repository, storage
│   └── migrations/
├── proto/                      # notification.proto и сгенерированный код
├── k8s/                        # манифесты Kubernetes
├── docker-compose.yml
├── Makefile
└── go.work
Frontend/                       # React + TypeScript (структура по слоям: app, pages, shared)
```

## Быстрый старт

Требования: Docker с Compose v2. Для локальной разработки дополнительно Go 1.25+ и Node.js.

### Вариант 1. Всё в Docker

```bash
cd Backend
docker compose up -d --build
```

Поднимаются PostgreSQL, Redis, Kafka, мигратор (применяет миграции и завершается), `backend-service` на `http://localhost:8080` и `notification-service` (gRPC на `127.0.0.1:50051`).

Сервису уведомлений нужна своя база. Пока она не создаётся автоматически, подготовьте её один раз:

```bash
docker compose exec -T postgres psql -U postgres -c "CREATE DATABASE notifications;"
docker compose exec -T postgres psql -U postgres -d notifications \
  < notification-service/migrations/00001_create_notifications.sql
docker compose restart notification-service
```

Проверка:

```bash
curl -X POST http://localhost:8080/api/register \
  -H 'Content-Type: application/json' \
  -d '{"username":"demo","email":"demo@example.com","password":"Passw0rd!"}'
```

Точные поля запроса смотрите в `backend-service/internal/dto/authentication_requests/`.

### Вариант 2. Локальная разработка

```bash
cd Backend
docker compose up -d postgres redis kafka   # только инфраструктура
cp .env.example .env                        # затем задайте свой JWT_SECRET
make migrate-up                             # миграции backend-service
make run                                    # API на :8080
```

Сервис уведомлений запускается отдельно, ему нужен `POSTGRES_DSN` с базой `notifications`:

```bash
cd Backend/notification-service
POSTGRES_DSN='postgres://postgres:1234@localhost:5432/notifications?sslmode=disable' \
  go run ./cmd/notification
```

### Frontend

```bash
cd Frontend
cp .env.example .env     # VITE_API_URL=http://localhost:8080
npm install
npm run dev              # http://localhost:3000
```

### Команды Makefile (из каталога `Backend/`)

| Команда | Действие |
|---|---|
| `make run` / `make build` | Запуск и сборка backend-service |
| `make test` | Тесты |
| `make migrate-up` / `migrate-down` | Применить и откатить миграции |
| `make migrate-create NAME=...` | Создать новую миграцию |
| `make docker-up` / `docker-down` / `docker-logs` | Управление контейнерами |
| `make lint` / `make fmt` | Линтер (golangci-lint) и форматирование |

## Конфигурация

| Переменная | Назначение | По умолчанию |
|---|---|---|
| `APP_NAME` | Имя приложения | `application` |
| `PUBLIC_API_HOST`, `PUBLIC_API_PORT` | Адрес API | `localhost`, `8080` |
| `PG_HOST`, `PG_PORT`, `PG_USER`, `PG_PASS`, `PG_DB` | Подключение к PostgreSQL | обязательные |
| `PG_MAX_CONN` | Размер пула соединений | `20` |
| `JWT_SECRET` | Секрет подписи JWT | обязательная |
| `REDIS_ADDR` | Адрес Redis | `localhost:6379` |
| `KAFKA_BROKERS` | Адрес Kafka | `localhost:9092` |
| `NOTIFICATION_GRPC_ADDR` | Адрес сервиса уведомлений | `localhost:50051` |
| `POSTGRES_DSN` | DSN для мигратора и сервиса уведомлений | обязательная для них |

> Значения в `docker-compose.yml` (пароль БД, `JWT_SECRET`) предназначены только для локальной разработки.

## API

Все маршруты находятся под префиксом `/api`. Публичные: `register`, `login`, `refresh`. Остальные требуют заголовок `Authorization: Bearer <access_token>`.

| Группа | Маршруты |
|---|---|
| Auth | `POST /register`, `POST /login`, `POST /refresh`, `POST /logout`, `POST /logout/all-devices` |
| Users | `GET /me`, `PUT /user`, `DELETE /user`, `POST /password`, `GET /users`, `GET /users/search`, `GET /user/:id`, `GET /user/username/:username` |
| Posts | `POST /posts`, `GET /posts`, `GET /posts/:id`, `PATCH /posts/:id`, `DELETE /posts/:id`, `GET /users/:id/posts`, `GET /feed` |
| Likes | `POST /posts/:id/like`, `DELETE /posts/:id/like`, `GET /posts/:id/likes`, `GET /posts/:id/likes/count`, `GET /posts/:id/like-status`, `GET /users/:id/likes` |
| Comments | `POST /posts/:id/comments`, `GET /posts/:id/comments`, `GET /posts/:id/comments/count`, `PATCH /comments/:id`, `DELETE /comments/:id` |
| Follow | `POST /:id/follow`, `DELETE /:id/follow`, `GET /:id/followers`, `GET /:id/following`, `GET /:id/followers/count`, `GET /:id/following/count`, `GET /:id/follow-status` |
| Friends | `GET /friends`, `GET /:id/friends`, `GET /:id/friends/count`, `GET /:id/friend-status` |
| Chats | `POST /users/:id/chat`, `GET /users/:id/chat`, `GET /chats`, `DELETE /chats/:id` |
| Messages | `POST /chats/:id/messages`, `GET /chats/:id/messages`, `GET /chats/:id/messages/unread/count`, `PATCH /messages/:id`, `DELETE /messages/:id`, `PATCH /messages/:id/status/read` |
| Notifications | `GET /notifications`, `GET /notifications/unread/count`, `PATCH /notifications/:id/read` |

### Контракт gRPC

`notification.v1.NotificationService` (`proto/notification.proto`): `GetNotifications`, `GetUnreadCount`, `MarkAsRead`.

## Тесты

```bash
cd Backend && make test
```

Хэндлеры авторизации и пользователей покрыты тестами с моками репозиториев (`internal/mocks`): успешные сценарии, ошибки валидации и внутренние ошибки.

## Планы развития

- Пагинация ленты и списков, меньше запросов к БД на ленту
- Graceful shutdown для `backend-service` (в сервисе уведомлений он уже реализован)
- Автоматическое создание БД `notifications` и применение её миграций в `docker-compose`
- Dead-letter очередь и повторные попытки в консьюмере вместо остановки на повреждённом событии
- CORS-middleware или прокси в Vite для работы клиента с API из браузера
- Больше типов уведомлений (лайки, комментарии, подписки), пагинация уведомлений по параметрам запроса
- TLS для gRPC-соединения между сервисами
- Расширение тестов: остальные контроллеры и `notification-service`
- Манифесты k8s для Redis, Kafka и сервиса уведомлений, CI (линтер и тесты)
- Загрузка изображений к постам, интерфейс чатов и уведомлений во фронтенде
