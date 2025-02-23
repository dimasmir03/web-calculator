# Используем многоступенчатую сборку для уменьшения размера образа
FROM golang:1.24-alpine as builder

WORKDIR /app

# Копируем зависимости и собираем их
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код оркестратора
COPY . .

# Компилируем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./cmd/server

# Финальный образ
FROM alpine:latest

WORKDIR /app

# Копируем бинарник из предыдущей стадии
COPY --from=builder /server .
COPY --from=builder /app/.env ./.env

# Порт для API оркестратора
EXPOSE 8080

# Запуск оркестратора
CMD ["./server"]