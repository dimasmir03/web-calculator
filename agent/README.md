# Web Calculator Agent

## Оглавление
- [Описание](#описание)
- [Принцип работы](#принцип-работы)
- [Запуск](#запуск)
  - [Требования](#требования)
  - [Переменные окружения](#переменные-окружения)
  - [Сборка и запуск](#сборка-и-запуск)
- [Масштабирование](#масштабирование)
- [Технологии](#технологии)

## Описание
Агент для выполнения арифметических операций. Получает задачи от оркестратора, выполняет с заданной задержкой и возвращает результат.

## Принцип работы
1. Запрос задачи: `GET /internal/task`
2. Вычисление операции
3. Отправка результата: `POST /internal/task`
4. Повтор каждые `TIME_WAIT_MS` мс при ошибках

## Запуск

### Требования
- Go 1.24+
- Доступ к оркестратору
- Make
- Docker

### Переменные окружения

```bash
SERVER_URL=http://host.docker.internal:8080  # URL оркестратора
COMPUTING_POWER=4                           # Количество горутин
TIME_WAIT_MS=1000                           # Задержка при ошибках (мс)
```

# Сборка и запуск
- Через Makefile:

```bash
make docker_build
make docker_run
```

- Docker Compose:

```yaml
agent:
  image: calc-agent
  environment:
    - SERVER_URL=http://orchestrator:8080
    - COMPUTING_POWER=4
```

Запуск нескольких агентов (Bash):

```bash
for i in {1..3}; do
  docker run -d --name agent_$i calc-agent
done
```

Масштабирование
Увеличивайте COMPUTING_POWER или количество контейнеров:

```bash
# Linux/Mac
COMPUTING_POWER=8 ./agent

# Windows
$env:COMPUTING_POWER=8; .\agent.exe
```

# Технологии
- Go Routines
- HTTP Client
- Docker
- Worker Pool

Пример команды для тестов:
## Тестирование
```bash
# Запуск всех тестов
make test

# Покрытие кода
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```