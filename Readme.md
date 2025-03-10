# НЕРУГАЙТЕ Я ВСЁ ИСПРАВЛЮ, напиши issue что готовы дать время на исправления всех недочетов

## СПАСИБО

## @dimajustgo

# Распределённый вычислитель арифметических выражений

[![Go Report Card](https://goreportcard.com/badge/github.com/dimasmir03/web-calculator)](https://goreportcard.com/report/github.com/diamsmir03/web-calculator)
[![Docker Pulls](https://img.shields.io/docker/pulls/dimasmir/calc-server)](https://hub.docker.com/r/dimasmir/calc-server)
[![Docker Compose](https://img.shields.io/badge/Docker_Compose-2496ED?logo=docker&logoColor=white)](https://docs.docker.com/compose/)
[![Swagger](https://img.shields.io/badge/Swagger-85EA2D?logo=swagger&logoColor=black)](http://localhost:8080/swagger)

## Содержание

- [Архитектура системы](#архитектура-системы)
- [Компоненты системы](#компоненты-системы)
- [Запуск системы](#запуск-системы)
  - [Требования](#требования)
  - [Docker Compose](#docker-compose)
  - [Ручная сборка](#ручная-сборка)
- [API Документация](#api-документация)
- [Статусы выражений](#статусы-выражений)
- [Примеры запросов](#примеры-запросов)
- [Переменные окружения](#переменные-окружения)
- [Особенности работы](#особенности-работы)
- [Технологии](#технологии)

## Архитектура системы

```mermaid
sequenceDiagram
    participant UI as Веб-интерфейс
    participant Orchestrator
    participant Agent1
    participant Agent2

    UI->>Orchestrator: POST /calculate
    Orchestrator->>Agent1: Задача 1
    Orchestrator->>Agent2: Задача 2
    Agent1->>Orchestrator: Результат 1
    Agent2->>Orchestrator: Результат 2
    Orchestrator->>UI: Обновление статуса
```

## Компоненты системы

| Компонент     | Описание                           | Порт |
| ------------- | ---------------------------------- | ---- |
| Оркестратор   | Управление задачами и выражениями  | 8080 |
| Агент         | Выполнение арифметических операций | -    |
| Веб-интерфейс | Пользовательский интерфейс         | 8081 |

### 🖥 Оркестратор (Сервер)

1. Принимает выражения от пользователей
2. Преобразует выражения в задачи
3. Управляет очередью задач
4. Возвращает статусы вычислений

### ⚡ Агент (Вычислитель)

1. Получает задачи Оркестратор (Сервера)
2. Получает задачи из очереди
3. Выполняет арифметические операции
4. Работает параллельно в нескольких горутинах
5. Возвращает результаты на сервер

### Веб-интерфейс

1. Отправляет выражения серверу
2. Выводит статусы вычислений

## Запуск системы

### Требования

Docker и Docker Compose  
Go 1.24+
Make

### Docker Compose

```bash
# 1. Скачайте репозиторий
  git clone https://github.com/diamsmir03/web-calculator.git
  cd web-calculator

# 2. Запустите систему
  docker-compose up --build
# 3. Открыть в браузере:
#    - Интерфейс: http://localhost:8081
#    - Документация: http://localhost:8080/swagger
```

### Ручная сборка

```bash
# 1. Скачайте репозиторий
  git clone https://github.com/dimasmir03/web-calculator.git
  cd web-calculator
# Сервер
  cd server
  go build -o server.exe ./cmd/server
  ./server.exe

# Агент (в отдельном терминале)
  cd agent
  go build -o agent.exe ./cmd/agent
  ./agent.exe
```

## API Документация

Доступна через Swagger UI:
<http://localhost:8080/swagger/index.html>

### Коллекция запросов для тестов [Postman Collection](https://app.getpostman.com/join-team?invite_code=9cac2ae36844ef092a1cdc71606cb988f42f99edaf7d3ff684768b7782fee6eb)

## Статусы выражений

| Статус     | Описание              | HTTP Код |
| ---------- | --------------------- | -------- |
| pending    | Ожидает выполнения    | 201      |
| processing | В процессе вычисления | 200      |
| completed  | Успешно выполнено     | 200      |
| error      | Ошибка при вычислении | 422      |

## Примеры запросов

- Добавить выражение

```bash
  curl --location 'localhost/api/v1/calculate' \
  --header 'Content-Type: application/json' \
  --data '{
    "expression": "(2+3)*4"
  }'
```

Ответ:

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000"
}
```

- Получить статус:

```bash
  curl http://localhost:8080/api/v1/expressions/550e8400-e29b-41d4-a716-446655440000
```

Ответ:

```json
{
  "expressions": [
    {
      "id": "2225fdb3-2e84-4308-8bb1-cc5e20ba77e4",
      "expression": "2*3",
      "status": "pending",
      "result": 0
    },
    {
      "id": "7631907e-9ec0-49e6-94b6-4b165dfe4404",
      "expression": "(2+5)*4",
      "status": "pending",
      "result": 0
    }
  ]
}
```

## Переменные окружения

1.Сервер

```ini
  PORT=8080
  TIME_ADDITION_MS=1000       # Время сложения (мс)
  TIME_SUBTRACTION_MS=1000    # Время вычитания
  TIME_MULTIPLICATION_MS=1000 # Время умножения
  TIME_DIVISION_MS=1000       # Время деления
```

2.Агент

```ini
  SERVER_URL=http://host.docker.internal:8080  # URL оркестратора
  COMPUTING_POWER=4                           # Количество горутин
  TIME_WAIT_MS=1000                           # Задержка при ошибках
```

## Особенности работы

Парсинг выражений: Используется RPN (Обратная Польская Запись) и AST

Параллельные вычисления: Агенты работают независимо друг от друга

Обработка ошибок: Автоматический повтор запросов при сбоях

Масштабирование:

```bash
  # Увеличить количество агентов
  docker-compose up --scale agent=3
```

## Технологии

- Сервер:
- Go,
- Echo Framework,
- Swagger
- Агент:
- Go Routines,
- HTTP Client
- Парсинг:
- AST,
- RPN
- Деплой:
- Make
- Docker,
- Docker Compose
