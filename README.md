# Social Platform Backend

Backend социальной сети в стиле Instagram на Go: регистрация и JWT-авторизация, посты, лайки, комментарии, подписки и друзья, личные чаты, лента. Учебный/личный проект.

## Возможности

| Область | Что реализовано |
|---|---|
| Авторизация | регистрация, логин, access JWT + refresh-токен (в БД хранится только хэш), выход с устройства, выход со всех устройств кроме текущего |
| Пользователи | профиль (`/me`), обновление, удаление, смена пароля (bcrypt), поиск, получение по id и по username |
| Посты | создание, чтение, обновление, удаление, посты пользователя, лента |
| Лента | сначала посты друзей, затем посты остальных пользователей |
| Лайки и комментарии | лайк/снятие лайка, счётчики, список лайкнувших, статус лайка, CRUD комментариев |
| Подписки и друзья | подписка/отписка, списки и счётчики, статус подписки. При взаимной подписке пара превращается в «друзей» |
| Чаты | создание чата 1-на-1, сообщения (отправка, правка, удаление), отметка «прочитано», счётчик непрочитанных |
| Кэш | Redis-кэш пользователя по `id` и по `username`, инвалидация при обновлении и удалении |
| События | при лайке в Kafka (топик `post-events`) отправляется событие `PostLiked`; отдельный consumer читает его |

## Стек

Go 1.25, Gin, PostgreSQL (pgx/v5, пул соединений), goose (миграции), Redis (go-redis), Kafka (segmentio/kafka-go), JWT (golang-jwt), bcrypt, slog, testify, Docker / docker-compose, манифесты Kubernetes.

## Структура

```
cmd/
  project/               # точка входа HTTP API
  migrator/              # запуск миграций (goose)
  post_events_consumer/  # consumer Kafka-событий о лайках
internal/
  app/          # сборка приложения: БД, Redis, Kafka, роуты
  config/       # конфигурация из переменных окружения
  controllers/  # HTTP-хэндлеры (зависят от интерфейсов репозиториев) + тесты
  repository/   # доступ к PostgreSQL/Redis, интерфейсы *Repo, общий Store
  routes/       # маршруты, группы public/private
  middleware/   # проверка Bearer JWT
  models/, dto/ # модели и структуры запросов
  kafka/, events/, storage/, mocks/, errs/
migrations/     # 00001..00004: users, posts/comments/likes/chats/messages, followers, refresh_tokens
k8s/            # deployment, configMap, postgres
```

## Запуск локально

1. Поднять инфраструктуру:
   ```bash
   docker-compose up -d postgres redis kafka
   ```
2. Конфигурация:
   ```bash
   cp .env.example .env
   ```
   Переменные: `APP_NAME`, `PUBLIC_API_HOST`, `PUBLIC_API_PORT`, `PG_HOST`, `PG_PORT`, `PG_USER`, `PG_PASS`, `PG_DB`, `PG_MAX_CONN`, `POSTGRES_DSN` (нужен мигратору), `JWT_SECRET`.
3. Миграции:
   ```bash
   go run ./cmd/migrator ./migrations up
   ```
4. API:
   ```bash
   go run ./cmd/project
   ```
   Сервер слушает порт из `PUBLIC_API_PORT` (по умолчанию 8080).
5. Consumer событий (по желанию):
   ```bash
   go run ./cmd/post_events_consumer
   ```

> Адреса Redis и Kafka сейчас заданы в коде как `localhost:6379` и `localhost:9092`, поэтому API запускается локально, а не внутри контейнера `app` из docker-compose.

## Тесты

```bash
go test ./...
```

Тесты хэндлеров авторизации и пользователей используют моки репозиториев (`internal/mocks`): успешные сценарии, ошибки валидации, внутренние ошибки.

## API

Все маршруты под префиксом `/api`. Публичные: `register`, `login`, `refresh`. Остальные требуют заголовок `Authorization: Bearer <access_token>`.

**Auth**
- `POST /register`, `POST /login`, `POST /refresh` (публичные)
- `POST /logout`, `POST /logout/all-devices`

**Users**
- `GET /me`, `PUT /user`, `DELETE /user`, `POST /password`
- `GET /users`, `GET /users/search`, `GET /user/:id`, `GET /user/username/:username`

**Posts**
- `POST /posts`, `GET /posts`, `GET /posts/:id`, `PATCH /posts/:id`, `DELETE /posts/:id`
- `GET /users/:id/posts`, `GET /feed`

**Likes / Comments**
- `POST|DELETE /posts/:id/like`, `GET /posts/:id/likes`, `GET /posts/:id/likes/count`, `GET /posts/:id/like-status`, `GET /users/:id/likes`
- `POST|GET /posts/:id/comments`, `GET /posts/:id/comments/count`, `PATCH|DELETE /comments/:id`

**Follow / Friends**
- `POST|DELETE /:id/follow`, `GET /:id/followers`, `GET /:id/following`, `GET /:id/followers/count`, `GET /:id/following/count`, `GET /:id/follow-status`
- `GET /friends`, `GET /:id/friends`, `GET /:id/friends/count`, `GET /:id/friend-status`

**Chats / Messages**
- `POST|GET /users/:id/chat`, `GET /chats`, `DELETE /chats/:id`
- `POST|GET /chats/:id/messages`, `PATCH|DELETE /messages/:id`, `PATCH /messages/:id/status/read`, `GET /chats/:id/messages/unread/count`

## Известные ограничения / планы

- Лента без пагинации и сортируется по id; два запроса к БД на ленту.
- Адреса Redis и Kafka не вынесены в конфиг.
- Нет graceful shutdown, пути в `Makefile` (`run`, `build`) нужно привести к `cmd/project`.
- Consumer пока только выводит события в консоль.
- Нет загрузки изображений: пост содержит только текст.