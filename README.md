# VideoStat - Content Intelligence & Processing System

VideoStat это backend-система для автоматизированного мониторинга, анализа и обработки видеоконтента.

Система реализует **event-driven pipeline**, который собирает данные о контенте с различных платформ, анализирует их и передаёт дальше в обработку (включая AI-генерацию и публикацию).

---

## 🚀 Назначение системы

Система решает задачи:

- сбор и хранение данных о блогерах и видео
- регулярный мониторинг контента (YouTube, Instagram, TikTok)
- анализ метрик и виральности
- асинхронная обработка данных
- подготовка данных для AI-пайплайна
- предоставление API для внешних клиентов и Telegram-бота

---

## 🧠 Архитектура

Проект построен по принципам:

- Clean Architecture
- Domain-Driven Design (DDD)
- CQRS
- Event-driven processing

### Слои системы

**Domain**
- Бизнес-логика (Blogger, Video)
- Стейт-машина обработки видео
- Доменные правила и инварианты

**Application**
- Use-cases (Commands / Queries)
- Оркестрация бизнес-процессов
- Независимость от инфраструктуры

**Infrastructure**
- PostgreSQL (pgx)
- RabbitMQ (очереди событий)
- Внешние API (Apify и др.)

**Transport**
- HTTP (Chi)
- gRPC (Protobuf)
- Telegram Bot API

---

## 🔁 Асинхронная модель

Система построена вокруг event-driven pipeline:

- API принимает входящие данные
- доменные события публикуются в RabbitMQ
- воркеры обрабатывают задачи независимо от API
- каждый этап обработки изолирован и масштабируется отдельно

---

## 🧩 Стейт-машина видео

Видео проходит строго определённые этапы обработки.

Переходы:

- контролируются доменной логикой
- проверяются на уровне бизнес-правил
- защищают систему от неконсистентных состояний

---

## 🛠 Технологический стек

- Go 1.25+
- PostgreSQL (pgx/v5)
- RabbitMQ
- gRPC (Buf, Protobuf)
- HTTP (Chi)
- Telegram Bot API
- slog (structured logging)
- Docker / docker-compose

---

## 📂 Структура проекта

```
├── api/              # Protobuf контракты
├── cmd/              # API и worker entrypoints
├── internal/
│   ├── app/          # DI и конфигурация
│   ├── application/  # Use-cases (CQRS)
│   ├── domain/       # Бизнес-логика
│   ├── infrastructure/ # БД, MQ, внешние сервисы
│   ├── platform/     # Клиенты и интеграции
│   └── worker/       # Фоновые процессы
└── migrations/       # SQL миграции
```

---

## 🚀 Быстрый старт

### 1. Установка зависимостей

- Go ≥ 1.25
- Docker + docker-compose
- Make
- Git
- protoc + buf

---

### 2. Клонирование

```bash
git clone <repo_url>
cd videostat
```

---

### 3. Настройка окружения

```bash
cp .env.example .env
```

Заполнить `.env` значениями.

---

### 4. Запуск инфраструктуры

```bash
make up
```

---

### 5. Миграции

```bash
make migrate-up
```

---

### 6. Генерация protobuf

```bash
make proto-gen
```

---

### 7. Запуск сервисов

```bash
make run                 # API
make run_worker_core     # core processing
make run_worker_cron     # scheduler
make run_worker_notify   # notifications
```

---

## 🔧 Makefile команды

### Основные

| Команда | Назначение |
|---|---|
| make up | инфраструктура |
| make down | остановка сервисов |
| make run | запуск API |

### Workers

| Команда | Назначение |
|---|---|
| make run_worker_core | обработка данных |
| make run_worker_cron | планировщик |
| make run_worker_notify | уведомления |

### Миграции

| Команда | Назначение |
|---|---|
| make migrate-up | применить миграции |
| make migrate-down | откатить |
| make migrate-create name=xxx | создать миграцию |

### Protobuf

| Команда | Назначение |
|---|---|
| make proto-gen | генерация кода |
| make proto-lint | проверка proto |
| make proto-breaking | проверка совместимости |

### Качество

| Команда | Назначение |
|---|---|
| make test-short | тесты |
| make lint | линтер |

---

## 🧪 Тестирование

- Table-driven tests
- Тестирование доменной логики
- Проверка state machine
- InMemory репозитории для use-case тестов

---

## 📈 Roadmap

- [x] Core event-driven architecture
- [x] RabbitMQ pipeline
- [x] Multi-transport API
- [ ] Deep video analysis (transcription, hooks)
- [ ] AI video generation pipeline
- [ ] Auto publishing system

---

## 📌 Примечание

Система спроектирована как масштабируемый backend с разделением ответственности и возможностью независимого масштабирования компонентов (API / workers / pipeline stages).