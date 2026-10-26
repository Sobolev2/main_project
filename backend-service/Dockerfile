
# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Копируем общий Go Workspace
COPY go.work ./

# Копируем модули проекта
COPY backend-service ./backend-service
COPY notification-service ./notification-service
COPY proto ./proto

# Собираем основной Backend
WORKDIR /app/backend-service

RUN CGO_ENABLED=0 GOOS=linux go build \
    -o /app/bin/project ./cmd/project

RUN CGO_ENABLED=0 GOOS=linux go build \
    -o /app/bin/migrator ./cmd/migrator

# Run stage
FROM alpine:3.20

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /app/bin/project .
COPY --from=builder /app/bin/migrator .
COPY --from=builder /app/backend-service/migrations ./migrations

EXPOSE 8080

CMD ["./project"]
