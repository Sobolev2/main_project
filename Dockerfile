# ---- Build stage ----
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Отдельный слой под зависимости — пересобираются только при изменении go.mod/go.sum
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Статическая сборка: бинарник без зависимости от libc из финального образа
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/project ./cmd/project
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/migrator ./cmd/migrator

# ---- Run stage ----
FROM alpine:3.20

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /app/bin/project .
COPY --from=builder /app/bin/migrator .
COPY --from=builder /app/migrations ./migrations

EXPOSE 8080

CMD ["./project"]
