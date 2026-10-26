
# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Общий Go Workspace
COPY go.work ./

# Модули проекта
COPY backend-service ./backend-service
COPY notification-service ./notification-service
COPY proto ./proto

# Собираем Notification Service
WORKDIR /app/notification-service

RUN CGO_ENABLED=0 GOOS=linux go build \
    -o /app/bin/notification ./cmd/notification

# Run stage
FROM alpine:3.20

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /app/bin/notification .

EXPOSE 50051

CMD ["./notification"]
