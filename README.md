# Мессенджер (Backend)

Бэкенд-часть мессенджера, реализованная на Go.

## Что сделано и планеы на будущее

Подробный список задач и их статус можно найти в [документации по задачам](docs/tasks.md).

## Архитектура проекта

Проект следует принципам чистой архитектуры и имеет следующую структуру:

```
├── api/           # API спецификации и документация
├── cmd/           # Точки входа приложения
│   ├── migrator/  # Утилита для миграций БД
│   └── server/    # Основной сервер приложения
├── internal/      # Внутренняя логика приложения
│   ├── config/    # Конфигурация приложения
│   └── database/  # Работа с базой данных
├── migrations/    # SQL миграции
└── pkg/           # Переиспользуемые пакеты
    ├── errors/    # Обработка ошибок
    └── logger/    # Логирование
```

## Требования

- Docker и Docker Compose
- Go 1.20 или выше (для локальной разработки)
- PostgreSQL (запускается в контейнере)

## Запуск проекта

### Через Docker Compose

1. Клонируйте репозиторий:
```bash
git clone <repository-url>
cd messenger
```

2. Создайте файл с переменными окружения:
```bash
cp .yaml.example .yaml
```

3. Запустите приложение:
```bash
docker-compose up -d
```

Приложение будет доступно по адресу `http://localhost:8080`

### Локальный запуск

1. Установите зависимости:
```bash
go mod download
```

2. Настройте переменные окружения:
```bash
cp .env.example .env
```

3. Запустите PostgreSQL:
```bash
docker-compose up -d postgres
```

4. Запустите сервер:
```bash
go run cmd/server/main.go
```

## Миграции

### Через Docker

Миграции автоматически применяются при запуске через Docker Compose. Для ручного запуска:

```bash
docker-compose run --rm migrator
```

### Локально

Для применения миграций локально:

```bash
go run cmd/migrator/main.go
```

## Конфигурация

Основные настройки приложения находятся в файле `.yaml`. Пример конфигурации:

```yaml
# Server Configuration
app:
  host: localhost
  port: 8080

# Database Configuration
db:
  host: localhost
  port: 5432
  user: postgres
  password: postgres
  name: messenger
  sslmode: disable
```

## Разработка

Проект использует:
- Go для бэкенда
- PostgreSQL для хранения данных
- JWT для аутентификации
- WebSocket для real-time коммуникации
- Docker для контейнеризации
- Swagger для документации API