# VideoStat

Система автоматической генерации видео на основе донорских роликов. Анализирует донорское видео, создаёт аватар с той же речью через HeyGen, генерирует B-roll через Kling и собирает финальное видео (PiP) через Shotstack.

---

## Быстрый старт

### 1. Зависимости

- Go ≥ 1.25
- Docker + docker-compose
- Make
- ngrok (для локальной разработки с Shotstack)

### 2. Клонировать репозиторий

```bash
git clone <repo_url>
cd videostat
```

### 3. Настроить окружение

```bash
cp .env.example .env
```

Заполнить `.env` — описание всех переменных в разделе [Переменные окружения](#переменные-окружения).

### 4. Поднять инфраструктуру

```bash
make up
```

Автоматически создаётся MinIO bucket `videostat` с публичным доступом на чтение.

### 5. Настроить ngrok для S3 (локальная разработка)

Shotstack скачивает видео по URL. MinIO на `localhost` ему недоступен — нужен публичный туннель.

```bash
ngrok http 9000
```

Скопировать URL вида `https://abc123.ngrok.io` и добавить в `.env`:

```
S3_PUBLIC_URL=https://abc123.ngrok.io
```

> В production `S3_PUBLIC_URL` совпадает с `S3_ENDPOINT` (публичный S3/R2).
> При каждом рестарте ngrok URL меняется — нужно обновить `.env` и перезапустить воркеры.

### 6. Применить миграции

```bash
make migrate-up
```

### 7. Запустить сервисы

```bash
make run                 # API
make run_worker_core     # обработка событий
make run_worker_cron     # планировщик
make run_worker_notify   # Telegram уведомления
```

---

## Переменные окружения

### База данных

| Переменная | Описание | Пример |
|------------|---------|--------|
| `DB_HOST` | хост PostgreSQL | `localhost` |
| `DB_PORT` | порт | `5432` |
| `DB_USER` | пользователь | `postgres` |
| `DB_PASS` | пароль | `postgres` |
| `DB_NAME` | имя базы | `videostat` |
| `DB_SSLMODE` | режим SSL | `disable` |

### RabbitMQ

| Переменная | Описание | Пример |
|------------|---------|--------|
| `RABBITMQ_USER` | пользователь | `guest` |
| `RABBITMQ_PASS` | пароль | `guest` |
| `RABBITMQ_HOST` | хост | `localhost` |
| `RABBITMQ_PORT` | порт | `5672` |
| `RABBITMQ_EXCHANGE` | exchange | `videostat` |

### S3 / MinIO

| Переменная | Описание | Пример |
|------------|---------|--------|
| `S3_ENDPOINT` | адрес S3 API (внутренний) | `http://localhost:9000` |
| `S3_PUBLIC_URL` | публичный URL для внешних сервисов | `https://abc123.ngrok.io` (локально) |
| `S3_ACCESS_KEY` | access key | `minioadmin` |
| `S3_SECRET_KEY` | secret key | `minioadmin` |
| `S3_BUCKET` | имя bucket | `videostat` |
| `S3_REGION` | регион | `us-east-1` |

> `S3_PUBLIC_URL` используется для формирования ссылок которые передаются в Shotstack. Если не указан — используется `S3_ENDPOINT`.

### Telegram Bot

| Переменная | Описание |
|------------|---------|
| `TELEGRAM_BOT_TOKEN` | токен бота от @BotFather |

### LLM (OpenAI)

| Переменная | Описание | Пример |
|------------|---------|--------|
| `LLM_PROVIDER` | провайдер | `openai` |
| `OPENAI_TOKEN` | API ключ OpenAI | `sk-...` |
| `OPENAI_MODEL` | модель | `gpt-4o-mini` |

Используется для генерации B-roll промптов: LLM делит транскрипт на сегменты 5-8 сек и придумывает промпт для видеоклипа каждого сегмента.

### AssemblyAI

| Переменная | Описание |
|------------|---------|
| `ASSEMBLYAI_TOKEN` | API ключ для транскрибации видео |

### HeyGen

| Переменная | Описание |
|------------|---------|
| `HEYGEN_API_KEY` | API ключ |
| `HEYGEN_AVATAR_ID` | ID аватара |
| `HEYGEN_VOICE_ID` | ID голоса |

Генерирует видео с аватаром. Получить `AVATAR_ID` и `VOICE_ID` можно в личном кабинете HeyGen.

### Kling

| Переменная | Описание |
|------------|---------|
| `KLING_ACCESS_KEY_ID` | Access Key ID из дашборда Kling |
| `KLING_SECRET_KEY` | Secret Key |

Генерирует B-roll видеоклипы по текстовым промптам (5 сек, 16:9). Аутентификация через JWT (HMAC-SHA256), генерируется автоматически.

### Shotstack

| Переменная | Описание | Пример |
|------------|---------|--------|
| `SHOTSTACK_API_KEY` | API ключ | из личного кабинета shotstack.io |
| `SHOTSTACK_ENV` | окружение | `stage` (sandbox) или `v1` (production) |

Собирает финальное видео: B-roll клипы как фон + аватар PiP (30%, правый нижний угол). В `stage` режиме на видео есть watermark — бесплатно для тестирования.

### Apify

| Переменная | Описание |
|------------|---------|
| `APIFY_TOKEN` | токен для парсинга YouTube/TikTok/Instagram |
| `APIFY_HOST` | хост Apify |

---

## Как работает пайплайн

### Общая схема

```
Пользователь (Telegram)
  └── /start_process_video <url>
        │
        ▼
  [1] Анализ (AssemblyAI)
        │  transcript + words с ms-таймингами
        ▼
  [2a] Аватар (HeyGen)               [2b] B-roll сегменты (OpenAI → Kling)
        │  SSML из words                   │  LLM делит на блоки 5-8 сек
        │  тот же текст, те же паузы        │  для каждого блока — промпт
        ▼                                  │  Kling генерирует видеоклип
  [3] Загрузка в S3                        ▼
        │                           [3] Kling клипы готовы
        └──────────────┬────────────────────┘
                       ▼
              [4] Сборка (Shotstack)
                   B-roll как фон
                   аватар PiP в правом нижнем углу (30%)
                       │
                       ▼
              [5] Telegram уведомление → пользователю
```

### Детальный поток событий

| Шаг | Событие | Обработчик | Результат |
|-----|---------|-----------|-----------|
| 1 | Команда `/start_process_video` | Telegram → `StartProcessVideo` | Сохраняет `VideoWatcher(video_id, chat_id)`, запускает пайплайн |
| 2 | `VideoProcessingStarted` | `worker_core` → `FetchVideoSources` | Находит прямую ссылку на видео |
| 3 | `VideoSourceFound` | `worker_core` → `AnalyzeVideo` | AssemblyAI транскрибирует, возвращает `words` с ms-таймингами |
| 4a | `VideoAnalyzeDone` | `worker_core` → `SubmitVideoGeneration` | Строит SSML из `words`, отправляет в HeyGen |
| 4b | `VideoAnalyzeDone` | `worker_core` → `GenerateBrollSegments` → `SubmitBrollGenerations` | OpenAI делит транскрипт на сегменты, Kling генерирует клипы |
| 5 | cron `*/3` | `worker_cron` → `PollVideoGenerations` | Когда HeyGen готов — скачивает и загружает в S3 |
| 6 | cron `*/3` | `worker_cron` → `PollBrollGenerations` | Когда все Kling клипы готовы — запускает `ComposeFinalVideo` |
| 7 | `ComposeFinalVideo` | Shotstack | Собирает финальное видео (B-roll фон + аватар PiP) |
| 8 | cron `*/3` | `worker_cron` → `PollCompositions` | Когда Shotstack готов — сохраняет `result_url` |
| 9 | `VideoGenerationDone/Error` | `worker_notify` | Telegram уведомление пользователю |

### SSML и паузы

Из массива `words` (каждое слово с `start`/`end` в мс) строится SSML для HeyGen:
- пауза ≥ 400 мс → `<break time="Xs"/>` (в секундах)
- пауза < 400 мс → пробел

Результат: аватар воспроизводит речь с теми же интонациями и паузами что и донор.

### Повторная генерация

Если анализ уже есть (`video_analysis`) — при повторном `/start_process_video` пайплайн пропускает анализ и сразу запускает генерацию. Работает из любого статуса видео.

---

## Статусы видео (state machine)

```
new
 └── processing
       └── generation_processing
             ├── ready            (HeyGen завершил, видео в S3)
             └── generation_failed
                   └── generation_processing  (повторная попытка)
```

Принудительный сброс в `generation_processing` через `MarkForRegeneration()` — доступен из любого статуса если есть `video_analysis`.

---

## Сервисы

| Сервис | Запуск | Назначение |
|--------|--------|-----------|
| `worker_core` | `make run_worker_core` | Обрабатывает события из RabbitMQ: анализ, HeyGen, Kling |
| `worker_cron` | `make run_worker_cron` | Polling: HeyGen каждые 3 мин, Kling каждые 3 мин, Shotstack каждые 3 мин, обновление блогеров в 3:00 |
| `worker_notify` | `make run_worker_notify` | Слушает события готовности/ошибки, шлёт Telegram уведомления |
| API | `make run` | HTTP + gRPC API |

---

## Telegram бот

| Команда / Кнопка | Описание |
|-----------------|---------|
| `/start_process_video <url>` | Запустить генерацию по ссылке на видео-донор |
| Создать блогера | Добавить блогера по ссылке на канал |
| Список блогеров | Показать всех блогеров |
| Список видео | Показать все видео |

После запуска генерации бот пришлёт уведомление когда видео будет готово (или если произошла ошибка).

---

## Makefile

### Инфраструктура

| Команда | Назначение |
|---------|-----------|
| `make up` | запустить docker-compose (postgres, rabbitmq, minio) |
| `make down` | остановить |

### Запуск сервисов

| Команда | Назначение |
|---------|-----------|
| `make run` | API сервер |
| `make run_worker_core` | воркер обработки событий |
| `make run_worker_cron` | воркер планировщика |
| `make run_worker_notify` | воркер уведомлений |

### Миграции

| Команда | Назначение |
|---------|-----------|
| `make migrate-up` | применить все миграции |
| `make migrate-down` | откатить последнюю |
| `make migrate-create name=xxx` | создать новую миграцию |

### Качество

| Команда | Назначение |
|---------|-----------|
| `make test-short` | запустить тесты |
| `make lint` | линтер |
| `make proto-gen` | генерация protobuf |

---

## Структура проекта

```
├── cmd/
│   ├── api/                    # HTTP + gRPC сервер
│   ├── worker_core/            # воркер событий
│   ├── worker_cron/            # воркер планировщика
│   └── worker_notify/          # воркер уведомлений
├── internal/
│   ├── app/                    # инициализация сервисов (DI)
│   ├── application/
│   │   └── blogger/command/    # use-cases
│   ├── config/                 # конфигурация
│   ├── domain/blogger/         # доменные сущности и репозиторий
│   ├── infrastructure/
│   │   ├── kling/              # клиент Kling API
│   │   ├── llm/openai/         # клиент OpenAI
│   │   ├── notifier/           # Telegram уведомления
│   │   ├── repository/blogger/ # PostgreSQL + InMemory
│   │   ├── shotstack/          # клиент Shotstack API
│   │   └── storage/s3/         # клиент S3/MinIO
│   ├── platform/               # фабрики клиентов
│   └── worker/
│       ├── core/               # обработчики событий RabbitMQ
│       ├── cron/               # cron задачи
│       └── notify/             # обработчики уведомлений
└── migrations/                 # SQL миграции
```

---

## Таблицы БД

| Таблица | Назначение |
|---------|-----------|
| `bloggers` | блогеры |
| `videos` | видео со статусом |
| `video_analysis` | результат анализа AssemblyAI (words, transcript) |
| `video_generations` | статус генерации HeyGen, S3 URL |
| `video_broll_segments` | B-roll сегменты: промпты, Kling статус, URL клипов |
| `video_compositions` | статус сборки Shotstack, итоговый URL |
| `video_watchers` | связь video_id → chat_id для уведомлений |
| `video_prompts` | (не используется активно) промпты от LLM |
