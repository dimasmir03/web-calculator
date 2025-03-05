# ВЕБ КАЛЬКУЛЯТОР С РАСПРЕДЕЛЕННЫМИ ВЫЧИСЛЕНИЯМИ

## Оглавление
- [Описание](#Описание)
- [Архитектура](#Архитектура)
- [API Документация](#api-документация)
- [Статусы выражений](#статусы-выражений)
- [Запуск](#запуск)
  - [Требования](#требования)
  - [Переменные окружения](#переменные-окружения)
  - [Сборка и запуск](#сборка-и-запуск)
- [Примеры запросов](#примеры-запросов)
- [Технологии](#технологии)

## Описание
Сервер-оркестратор для распределенного вычисления арифметических выражений. Использует RPN (Обратную Польскую Запись) и AST (Абстрактное Синтаксическое Дерево) для парсинка выражений.

## Архитектура
```mermaid
graph TD
    A[Пользователь] --> B[POST /calculate]
    B --> C{Валидация}
    C --> |OK| D[Создание AST]
    D --> E[Разбиение на задачи]
    E --> F[Очередь задач]
    F --> G[Агенты]
    G --> H[POST /internal/task]
    H --> I[Обновление статуса]
    I --> J[GET /expressions]
```

# API Документация
Доступна через Swagger UI: http://localhost:8080/swagger/index.html

### [Postman Collection](https://app.getpostman.com/join-team?invite_code=9cac2ae36844ef092a1cdc71606cb988f42f99edaf7d3ff684768b7782fee6eb)

# Статусы выражений


| Статус     | Описание              |
|------------|-----------------------| 
| pending    | Ожидает выполнения    |
| processing | В процессе вычисления |
| completed  | Успешно выполнено     |
| error      | Ошибка при вычислении |

# Запуск
## Требования
Необходимо иметь установленные программы:
* Go 1.24+
* Docker
* Make

## Переменные окружения
```bash
TIME_ADDITION_MS=1000       # Время сложения (мс)
TIME_SUBTRACTION_MS=1000    # Время вычитания (мс)
TIME_MULTIPLICATION_MS=1000 # Время умножения (мс)
TIME_DIVISION_MS=1000       # Время деления (мс)
PORT=8080                   # Порт сервера
```
## Сборка и запуск

* Через Makefile

```bash
make docker_network  # Создать сеть
make docker_build    # Собрать образ
make docker_run      # Запустить контейнер
```

* Docker Compose

```yaml
version: '3.8'
services:
  orchestrator:
    image: calc-server
    ports:
      - "8080:8080"
    environment:
      - TIME_ADDITION_MS=1000
      - TIME_SUBTRACTION_MS=1000
```

* Ручной запуск (Linux/Mac):

```bash
go build -o server ./cmd/server
TIME_ADDITION_MS=1000 ./server
```

# Примеры запросов

* Добавление выражения:

```bash
curl -X POST http://localhost:8080/api/v1/calculate \
  -H "Content-Type: application/json" \
  -d '{"expression": "2+2*2"}'
```

Ответ:

```json
{
    "id": "d78fe1c5-a656-46d9-825b-8f86b9fcecc4"
}
```

Ошибка:
`example: "2*(3"`
```bash
	expected 'RPar' type, got 'EOL';
```

# Технологии

- Echo Framework
- RPN/AST парсер
- Worker Pool
- Make
- Docker
- Docker-Compose
- Swagger
