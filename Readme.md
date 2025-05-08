# Распределенный вычислитель арифметических выражений

[![Go Report Card](https://goreportcard.com/badge/github.com/dimasmir03/web-calculator)](https://goreportcard.com/report/github.com/diamsmir03/web-calculator)
[![Docker Pulls](https://img.shields.io/docker/pulls/dimasmir/calc-server)](https://hub.docker.com/r/dimasmir/calc-server)
[![Docker Compose](https://img.shields.io/badge/Docker_Compose-2496ED?logo=docker&logoColor=white)](https://docs.docker.com/compose/)
[![Swagger](https://img.shields.io/badge/Swagger-85EA2D?logo=swagger&logoColor=black)](http://localhost:8080/swagger)

## Архитектура

### Компоненты

- Оркестратор (API, аутентификация, хранение данных)
- Агент (вычислитель выражений)
- Веб-интерфейс (пользовательский интерфейс)
- SQLite (хранение данных)

### Как работает

1. Оркестратор получает запросы от Веб-интерфейса
2. Оркестратор отправляет задачу Агенту
3. Агент вычисляет выражение и отправляет результат Оркестратору
4. Оркестратор сохраняет результат в SQLite
5. Оркестратор отправляет результат Веб-интерфейсу

## Запуск

### Требования

- Docker и Docker Compose
- Go 1.20+
- Make

### Запуск

1. Клонировать репозиторий
2. Перейти в папку `web-calculator`
3. Запустить `docker-compose up --build`

### Доступные endpoints

- Веб-интерфейс: http://localhost:8081
- Swagger UI: http://localhost:8080/swagger
- gRPC сервер: localhost:50051

## Аутентификация

### Регистрация

bash
curl -X POST http://localhost:8080/api/v1/register \
  -H "Content-Type: application/json" \
  -d '{"login":"user1", "password":"secret"}'

### Логин

bash
curl -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"login":"user1", "password":"secret"}'

### Использование токена

bash
curl -X POST http://localhost:8080/api/v1/calculate \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"expression": "2+2*2"}'

## API Документация

Полная документация доступна через Swagger UI:
http://localhost:8080/swagger

Основные endpoints:

POST /register - регистрация

POST /login - получение токена

POST /calculate - добавление выражения

GET /expressions - список выражений

## Переменные окружения

### Оркестратор

ini
DB_DSN=/data/calculator.db  # Путь к SQLite
JWT_SECRET=your-secret-key  # Ключ подписи JWT
GRPC_PORT=50051             # Порт gRPC сервера

### Агент

ini
SERVER_GRPC=orchestrator:50051  # gRPC адрес оркестратора
COMPUTING_POWER=4               # Количество горутин

## Технологии

### Сервер

- Go
- Echo Framework
- SQLite (GORM)
- JWT аутентификация
- gRPC для агентов

### Агент

- Go
- gRPC клиент
- Пул горутин

### Фронтенд

- jQuery
- HTML5
- CSS3
- JWT авторизация

### Инфраструктура

- Docker
- Docker Compose
- Автоматические миграции

## Спасибо за внимание!