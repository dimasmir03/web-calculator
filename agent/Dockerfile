# Используем многоступенчатую сборку
FROM golang:1.24-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код агента
COPY . .

# Компилируем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -o /agent ./cmd/agent

# Финальный образ
FROM alpine:latest

WORKDIR /app

COPY --from=builder /agent .
COPY --from=builder /app/.env ./.env

# Запуск агента
CMD ["./agent"]