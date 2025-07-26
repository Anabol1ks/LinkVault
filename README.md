# LinkVault

Многофункциональный сервис для сокращения ссылок с поддержкой JWT-аутентификации, статистики переходов, cron-очистки и удобным REST API.

## Возможности
- Регистрация и аутентификация пользователей (JWT)
- Создание коротких ссылок (анонимно и с привязкой к пользователю)
- Статистика по ссылкам: общее количество переходов, уникальные IP, география, график по дням
- Очистка истёкших и неактуальных ссылок по CRON
- Swagger-документация (автоматически генерируется)
- Логирование (zap)
- Docker/Docker Compose для быстрого запуска


---

## Быстрый старт

### 1. Клонирование репозитория
```bash
git clone https://github.com/Anabol1ks/LinkVault.git
cd LinkVault
```

### 2. Настройка переменных окружения

- Для локального запуска скопируйте `.env.example` в `.env` и при необходимости измените значения:

```bash
cp .env.example .env
```

- Для Docker Compose обязательно проверьте:
  - DB_HOST=linkv-db
  - DB_PORT=5432
  - DB_USER=postgres
  - DB_PASSWORD=linkv12341
  - DB_NAME=linkv-db
  - DB_SSLMODE=disable

- Для локального запуска без Docker:
  - DB_HOST=localhost
  - DB_PORT=5432
  - ... (остальные значения аналогично)

---

## Запуск локально (без Docker)

1. Установите Go 1.24+
2. Установите PostgreSQL 17+ и создайте базу данных с параметрами из `.env`
3. Запустите миграции и приложение:

```bash
go mod download
go run ./cmd/main.go
```

---

## Запуск через Docker Compose

```bash
docker-compose up --build
```

- Приложение будет доступно на http://localhost:8080
- Swagger: http://localhost:8080/swagger/index.html
- База данных: порт 5432

---


---

## Запуск через Makefile

Для удобства доступны команды:

```bash
make build   # Сборка бинарника
make run     # Запуск приложения
```

---

## Основные технологии
- Go (Gin, GORM, JWT, zap, cron, swag)
- PostgreSQL
- Docker, Docker Compose

---

## Структура проекта
```
cmd/            # Точка входа (main.go)
docs/           # Swagger-документация
internal/
  config/       # Загрузка и парсинг конфигов
  croncleaner/  # Очистка старых ссылок по CRON
  handler/      # HTTP-хендлеры
  jwt/          # JWT-логика
  logger/       # Логирование
  middleware/   # Middleware для Gin
  models/       # GORM-модели
  repository/   # Работа с БД
  response/     # Структуры ответов API
  router/       # Настройка роутера
  service/      # Бизнес-логика
  storage/      # Подключение и миграции БД
```

---

## API
- Swagger: http://localhost:8080/swagger/index.html
- Примеры запросов и ответы описаны в Swagger и в комментариях к коду.

---

## Важно
- Для production-режима обязательно измените секреты и параметры БД!
- Не храните реальные секреты в публичных репозиториях.
- Для деплоя на сервер используйте собранный Docker-образ из ghcr.io или свой registry.

---

## Автор
- [Anabol1ks](https://github.com/Anabol1ks)

---

## Лицензия
MIT
