# VideoStat

Система автоматической генерации вертикальных видео (9:16) для TikTok/Reels на основе донорских роликов.

📖 **[Руководство пользователя](GUIDE.md)** — как работать с Telegram ботом, запускать генерацию и разбираться с проблемами. Анализирует донорское видео, создаёт аватар с той же речью через HeyGen, генерирует динамичный B-roll через Kling и собирает финальное вертикальное видео (аватар внизу + B-roll фон) через Shotstack.

---

## Быстрый старт

### 1. Зависимости

| Зависимость | Версия | Нужна для |
|-------------|--------|-----------|
| Go | ≥ 1.25.6 | сборка и запуск сервисов |
| Docker + Docker Compose v2 | последняя | инфраструктура, миграции, protobuf |
| Make | любая | команды Makefile |
| golangci-lint | любая | `make lint` |
| ngrok | любая | локальная разработка (публичный URL для Shotstack) |

> Миграции (`make migrate-up`) и генерация protobuf (`make proto-gen`) запускаются **внутри Docker** — отдельно ставить `golang-migrate` и `buf` не нужно.

#### Go

**macOS:**
```bash
brew install go
```

**Linux:**
```bash
wget https://go.dev/dl/go1.25.6.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.25.6.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc
```

Проверить: `go version`

#### Docker + Docker Compose v2

**macOS:** установить [Docker Desktop](https://www.docker.com/products/docker-desktop/) — Compose v2 включён.

**Linux:**
```bash
# Docker
curl -fsSL https://get.docker.com | sh
sudo usermod -aG docker $USER
newgrp docker

# Docker Compose v2 (плагин)
sudo apt-get install docker-compose-plugin   # Debian/Ubuntu
sudo yum install docker-compose-plugin       # RHEL/CentOS
```

Проверить: `docker compose version` (именно `docker compose`, не `docker-compose`)

#### Make

**macOS:**
```bash
xcode-select --install
```

**Linux:**
```bash
sudo apt-get install make    # Debian/Ubuntu
sudo yum install make        # RHEL/CentOS
```

#### golangci-lint

**macOS:**
```bash
brew install golangci-lint
```

**Linux:**
```bash
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin
```

Проверить: `golangci-lint --version`

> Нужен только для `make lint`. Для запуска проекта не требуется.

#### ngrok

**macOS:**
```bash
brew install ngrok/ngrok/ngrok
```

**Linux:**
```bash
curl -sSL https://ngrok-agent.s3.amazonaws.com/ngrok.asc | sudo tee /etc/apt/trusted.gpg.d/ngrok.asc >/dev/null
echo "deb https://ngrok-agent.s3.amazonaws.com buster main" | sudo tee /etc/apt/sources.list.d/ngrok.list
sudo apt update && sudo apt install ngrok
```

Требует регистрации на [ngrok.com](https://ngrok.com) и авторизации: `ngrok config add-authtoken <токен>`

> Нужен только для локальной разработки. В production не требуется.

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

Shotstack скачивает видео по публичному URL. MinIO на `localhost` ему недоступен — нужен публичный туннель.

```bash
ngrok http 9000
```

Скопировать URL вида `https://abc123.ngrok.io` и добавить в `.env`:

```
S3_PUBLIC_URL=https://abc123.ngrok.io
```

> В production `S3_PUBLIC_URL` совпадает с адресом публичного S3/R2.
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

> `S3_PUBLIC_URL` используется для формирования ссылок, которые передаются в Shotstack и Kling. Если не указан — используется `S3_ENDPOINT`.

### Telegram Bot

| Переменная | Описание | Пример |
|------------|---------|--------|
| `TELEGRAM_BOT_TOKEN` | токен бота от @BotFather | |
| `TELEGRAM_ALLOWED_USERNAMES` | whitelist username'ов через запятую, без `@`. Если не задано — бот открыт для всех | `gregoryKarman,otherUser` |

### LLM (OpenAI)

| Переменная | Описание | Пример |
|------------|---------|--------|
| `LLM_PROVIDER` | провайдер | `openai` |
| `OPENAI_TOKEN` | API ключ OpenAI | `sk-...` |
| `OPENAI_MODEL` | модель | `gpt-4.1-mini` |

Рекомендуемые модели (по соотношению качество/цена): `gpt-4.1-mini`, `gpt-4.1`, `gpt-4o`.

Используется **только** для написания B-roll промптов — сегментация транскрипта выполняется в Go (каждые ~3 секунды), LLM получает готовые сегменты и пишет промпт для каждого.

### AssemblyAI

| Переменная | Описание |
|------------|---------|
| `ASSEMBLYAI_TOKEN` | API ключ для транскрибации видео с пословными таймингами |

### HeyGen

| Переменная | Описание |
|------------|---------|
| `HEYGEN_API_KEY` | API ключ |
| `HEYGEN_AVATAR_ID` | ID аватара (из личного кабинета HeyGen) |
| `HEYGEN_VOICE_ID` | ID голоса (из личного кабинета HeyGen) |

Генерирует горизонтальное видео с аватаром **1280×720** (HD, 16:9). Аватар размещается в нижней части финального вертикального видео (~40% высоты).

### Kling

| Переменная | Описание | Пример |
|------------|---------|--------|
| `KLING_ACCESS_KEY_ID` | Access Key ID из дашборда Kling | |
| `KLING_SECRET_KEY` | Secret Key | |
| `KLING_MODEL` | модель генерации | `kling-v1-6` |

Доступные модели: `kling-v1`, `kling-v1-5`, `kling-v1-6`, `kling-v2-1`, `kling-v2-5-turbo`.

Генерирует B-roll клипы в формате **9:16** (вертикальный). Длительность клипа: **5 секунд** (сегмент < 8 сек) или **10 секунд** (сегмент ≥ 8 сек). Аутентификация через JWT (HMAC-SHA256), генерируется автоматически.

> При ошибке rate limit (код 1303) или 429 — сабмит прерывается, оставшиеся pending сегменты досылаются на следующем тике крона.

### Shotstack

| Переменная | Описание | Пример |
|------------|---------|--------|
| `SHOTSTACK_API_KEY` | API ключ | из личного кабинета shotstack.io |
| `SHOTSTACK_BASE_URL` | базовый URL **без** `/render` | `https://api.shotstack.io/stage` |

Собирает финальное вертикальное видео **720×1280** (9:16):
- Нижний слой: B-roll клипы (фон, без звука)
- Верхний слой: аватар 16:9 внизу экрана (~40% высоты)

В `stage` режиме на видео есть watermark — бесплатно для тестирования. Для production использовать `https://api.shotstack.io/v1`.

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
  [2a] Аватар (HeyGen 1280×720)       [2b] B-roll (OpenAI → Kling 9:16)
        │  SSML из words                   │  Go делит на ~3-сек сегменты
        │  те же паузы и интонации          │  LLM пишет промпт для каждого
        ▼                                  │  Kling генерирует клип (5 или 10 сек)
  [3] Загрузка в S3                        ▼
        │                           [3] Kling клипы готовы
        └──────────────┬────────────────────┘
                       ▼
              [4] Сборка (Shotstack 720×1280)
                   B-roll клипы — фон (весь экран)
                   Аватар — нижние 40% экрана
                       │
                       ▼
              [5] Telegram уведомление → пользователю
```

### Детальный поток событий

| Шаг | Событие / Крон | Обработчик | Результат |
|-----|---------------|-----------|-----------|
| 1 | Команда `/start_process_video` | Telegram → `StartProcessVideo` | Сохраняет `VideoWatcher(video_id, chat_id)`, запускает пайплайн |
| 2 | `VideoProcessingStarted` | `worker_core` → `FetchVideoSources` | Находит прямую ссылку на видео |
| 3 | `VideoSourceFound` | `worker_core` → `AnalyzeVideo` | AssemblyAI транскрибирует, возвращает `words` с ms-таймингами |
| 4a | `VideoAnalyzeDone` | `worker_core` → `SubmitVideoGeneration` | Строит SSML из `words`, отправляет в HeyGen |
| 4b | `VideoAnalyzeDone` | `worker_core` → `GenerateBrollSegments` | Go делит транскрипт на ~3-сек сегменты, LLM пишет промпты, Kling получает задачи |
| 5 | cron `*/1` | `PollVideoGenerations` | HeyGen готов → скачивает и загружает в S3 |
| 6 | cron `*/1` | `PollBrollGenerations` | Kling клипы готовы → если аватар тоже готов, запускает сборку |
| 7 | cron `*/1` | `RetryPendingBrollSubmissions` | Досылает pending B-roll сегменты в Kling (после rate limit) |
| 8 | cron `*/1` | `TriggerPendingCompositions` | Страховка: если оба готовы, но сборка не запустилась — запускает |
| 9 | cron `*/1` | `PollCompositions` | Shotstack готов → сохраняет `result_url` |
| 10 | `VideoGenerationDone/Error` | `worker_notify` | Telegram уведомление пользователю |

### Сегментация B-roll (двухшаговый подход)

1. **Go** механически делит массив `words` (из AssemblyAI) на сегменты по ~3000 мс, сохраняя точные `start_ms`/`end_ms`
2. **LLM** получает готовые сегменты с текстом и временными метками → пишет только `broll_prompt` для каждого

Такой подход гарантирует точность таймингов и правильное количество сегментов (~7-8 для 25-секундного видео).

Промпты требуют: динамичная съёмка, американский сеттинг, реальные люди (не CGI), без слоумо, без азиатских локаций и персонажей.

### SSML и паузы аватара

Из массива `words` строится SSML для HeyGen:
- пауза ≥ 400 мс → `<break time="Xs"/>`
- пауза < 400 мс → пробел

Результат: аватар воспроизводит речь с теми же паузами что и оригинал.

### Повторная генерация

Если видео уже было обработано, при повторном `/start_process_video` Telegram показывает диалог подтверждения:
- **Да** → сбрасывает все данные (broll, compositions, generations, watchers) и запускает заново
- **Нет** → отмена

---

## Статусы видео

```
new
 └── processing
       └── generation_processing
             ├── ready            (HeyGen завершил, видео в S3)
             └── generation_failed
```

---

## Сервисы

| Сервис | Запуск | Назначение |
|--------|--------|-----------|
| `worker_core` | `make run_worker_core` | Обрабатывает события из RabbitMQ: анализ, отправка в HeyGen/Kling |
| `worker_cron` | `make run_worker_cron` | Polling каждую минуту: HeyGen, Kling, Shotstack; retry pending сегментов; страховочный триггер сборки; обновление блогеров в 3:00 |
| `worker_notify` | `make run_worker_notify` | Слушает события готовности/ошибки, шлёт Telegram уведомления |
| API | `make run` | HTTP + gRPC API |

### Cron расписание

| Задача | Расписание | Назначение |
|--------|-----------|-----------|
| `PollVideoGenerations` | каждую минуту | статус HeyGen генераций |
| `PollBrollGenerations` | каждую минуту | статус Kling клипов |
| `PollCompositions` | каждую минуту | статус Shotstack сборок |
| `RetryPendingBrollSubmissions` | каждую минуту | досылает pending B-roll после rate limit |
| `TriggerPendingCompositions` | каждую минуту | запускает сборку если оба потока завершились |
| `RefreshAllBloggers` | 03:00 ежедневно | обновление списка видео блогеров |

---

## Telegram бот

| Команда / Кнопка | Описание |
|-----------------|---------|
| `/start_process_video <url>` | Запустить генерацию по ссылке на видео-донор |
| Создать блогера | Добавить блогера по ссылке на канал |
| Список блогеров | Показать всех блогеров |
| Список видео | Показать все видео |

Если видео уже обрабатывалось — бот покажет кнопки подтверждения перед перезапуском.
После запуска бот пришлёт уведомление когда видео будет готово или если произошла ошибка.

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
│   │   ├── kling/              # клиент Kling API (JWT auth, 9:16, 5/10 сек)
│   │   ├── llm/openai/         # генерация B-roll промптов
│   │   ├── notifier/           # Telegram уведомления
│   │   ├── repository/blogger/ # PostgreSQL + InMemory
│   │   ├── shotstack/          # клиент Shotstack API (720×1280)
│   │   ├── storage/s3/         # клиент S3/MinIO
│   │   └── videogenerator/heygen/ # клиент HeyGen API (1280×720)
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
| `video_generations` | статус генерации HeyGen, S3 URL аватара |
| `video_broll_segments` | B-roll сегменты: тайминги, промпты, Kling статус, URL клипов |
| `video_compositions` | статус сборки Shotstack, итоговый URL финального видео |
| `video_watchers` | связь video_id → chat_id для уведомлений |
| `video_prompts` | промпты для HeyGen (текст из транскрипта) |

---

## Известные ограничения

- **Kling rate limit**: не более ~5 параллельных задач. При превышении (код 1303) сабмит прерывается, pending сегменты досылаются на следующем тике (каждую минуту).
- **Kling длительность**: принимает только **5** или **10** секунд. Сегменты < 8 сек → 5-секундный клип, ≥ 8 сек → 10-секундный.
- **Shotstack рендер**: занимает 5–15 минут в зависимости от длины видео.
- **ngrok**: при рестарте меняется URL — нужно обновить `S3_PUBLIC_URL` в `.env` и перезапустить воркеры.
- **SHOTSTACK_BASE_URL**: указывать **без** `/render` в конце (например `https://api.shotstack.io/stage`, не `https://api.shotstack.io/stage/render`).
