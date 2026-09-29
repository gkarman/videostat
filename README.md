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

Вместе с инфраструктурой поднимается мониторинг: Prometheus, Grafana, Loki, Tempo, Alloy — см. [Наблюдаемость](#наблюдаемость-метрики-логи-трейсы).

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

### 8. Открыть мониторинг

Grafana: http://localhost:3000 (`admin` / `admin`) → дашборд **Videostat Overview**.
Логи всех сервисов, метрики и трейсы — там же, подробнее в разделе [Наблюдаемость](#наблюдаемость-метрики-логи-трейсы).

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

### Мониторинг, логи, трейсы

| Переменная | Описание | Пример |
|------------|---------|--------|
| `LOG_LEVEL` | уровень логов | `debug` |
| `LOG_DIR` | папка для файлов логов (их забирает Alloy → Loki); пусто — логи только в терминал | `logs` |
| `METRICS_WORKER_CORE_ADDR` | адрес `/metrics` worker_core | `localhost:9101` |
| `METRICS_WORKER_NOTIFY_ADDR` | адрес `/metrics` worker_notify | `localhost:9102` |
| `METRICS_WORKER_CRON_ADDR` | адрес `/metrics` worker_cron | `localhost:9103` |
| `TRACING_OTLP_ENDPOINT` | куда слать трейсы (OTLP/gRPC); пусто — трейсинг выключен. Указывать IP, не `localhost` | `127.0.0.1:4317` |
| `TRACING_SAMPLE_RATIO` | доля сохраняемых трейсов (1 — все) | `1` |
| `GRAFANA_USER` / `GRAFANA_PASS` | логин и пароль Grafana | `admin` / `admin` |

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

## Наблюдаемость: метрики, логи, трейсы

Три источника информации о работе системы, все смотрятся в одной Grafana и связаны между собой:

| | Что отвечает | Хранилище | Как попадает |
|---|---|---|---|
| **Метрики** | сколько и как быстро (RPS, ошибки, длина очередей, видео по статусам) | Prometheus | Prometheus сам забирает `/metrics` у сервисов каждые 15 с |
| **Логи** | что именно происходило | Loki | сервисы пишут в `logs/*.log`, Alloy читает файлы и отправляет в Loki |
| **Трейсы** | из чего состоял один запрос и где он потратил время | Tempo | сервисы шлют спаны по OTLP в Alloy, Alloy пересылает в Tempo |

```
 Go-сервисы (на хосте, make run*)                 Docker (make up)

 /metrics   ◄────── каждые 15 с ──────  Prometheus ◄── RabbitMQ, postgres-exporter ─┐
 logs/*.log ──────►  Alloy  ──────────► Loki ───────────────────────────────────────┼──► Grafana :3000
 OTLP :4317 ──────►  Alloy  ──────────► Tempo ──────────────────────────────────────┘
```

Весь мониторинг поднимается вместе с остальной инфраструктурой: `make up`. Конфиги лежат в `monitoring/`.

### Где смотреть

| Что | Адрес | Зачем |
|---|---|---|
| **Grafana** | http://localhost:3000 (`admin` / `admin`) | основное место: дашборды, поиск по логам и трейсам |
| — дашборд **Videostat Overview** | http://localhost:3000/d/videostat-overview | всё ли в порядке: алерты, сервисы, видео, очереди, воркеры, cron, внешние API, Postgres, ошибки в логах |
| — дашборд **Videostat API** | http://localhost:3000/d/videostat-api | подробно про HTTP API: RPS, ошибки, время ответа, Go runtime |
| — **Explore** | Grafana → Explore | произвольные запросы: Prometheus (PromQL), Loki (LogQL), Tempo (TraceQL) |
| — **Alerting** | Grafana → Alerting → Alert rules | все правила алертов (и Prometheus, и Loki) и их состояние |
| Prometheus | http://localhost:9090 | `Status → Target health` — кто из сервисов доступен; `Alerts` — правила алертов |
| Alloy | http://localhost:12345 | схема конвейера логов и трейсов, ошибки компонентов |
| RabbitMQ | http://localhost:15672 | очереди и сообщения (метрики очередей есть и в Grafana) |

### Метрики

Каждый сервис отдаёт метрики по HTTP, Prometheus забирает их сам (список целей — `monitoring/prometheus/prometheus.yml`):

| Сервис | Адрес `/metrics` | Настройка |
|---|---|---|
| `api` | http://localhost:8081/metrics | тот же порт, что HTTP API |
| `worker_core` | http://localhost:9101/metrics | `METRICS_WORKER_CORE_ADDR` |
| `worker_notify` | http://localhost:9102/metrics | `METRICS_WORKER_NOTIFY_ADDR` |
| `worker_cron` | http://localhost:9103/metrics | `METRICS_WORKER_CRON_ADDR` |
| RabbitMQ | `rabbitmq:15692` (внутри docker) | встроенный плагин `rabbitmq_prometheus` |
| Postgres | `postgres-exporter:9187` (внутри docker) | контейнер `postgres-exporter` |

**Свои метрики приложения:**

| Метрика | Тип | Метки | Где пишется |
|---|---|---|---|
| `http_requests_total` | Counter | `method`, `route`, `status` | `transport/http/middleware/metrics.go` |
| `http_request_duration_seconds` | Histogram | `method`, `route` | там же |
| `http_requests_in_flight` | Gauge | — | там же |
| `worker_messages_total` | Counter | `event_type`, `result` (`ok`/`error`/`skipped`) | `worker/router.go` — каждое событие из RabbitMQ |
| `worker_message_duration_seconds` | Histogram | `event_type` | там же |
| `cron_job_runs_total` | Counter | `task`, `result` | `worker/cron/worker.go` (`runJob`) |
| `cron_job_duration_seconds` | Histogram | `task` | там же |
| `cron_job_last_success_timestamp_seconds` | Gauge | `task` | там же — время последнего успешного запуска |
| `external_api_requests_total` | Counter | `provider`, `status` (HTTP-код или `error`) | `infrastructure/metrics/external_api.go` — все вызовы HeyGen, Kling, Shotstack, AssemblyAI, Apify, OpenAI, OpenRouter, Anthropic |
| `external_api_request_duration_seconds` | Histogram | `provider` | там же |
| `videostat_videos` | Gauge | `status` | `infrastructure/metrics/video_collector.go` — считается запросом в БД при каждом сборе, только в `worker_cron` |
| `videostat_videos_failed` | Gauge | `stage` | там же |

Плюс стандартные `go_*` и `process_*` (горутины, память, CPU) у каждого сервиса.

**Полезные запросы (PromQL):**

```promql
sum by (job) (rate(http_requests_total[5m]))                                   # RPS
sum by (event_type) (rate(worker_messages_total{result="error"}[5m]))          # ошибки обработки событий
histogram_quantile(0.95, sum by (le, provider) (rate(external_api_request_duration_seconds_bucket[5m])))  # p95 внешних API
time() - cron_job_last_success_timestamp_seconds                               # сколько секунд назад задача прошла успешно
rabbitmq_detailed_queue_messages_ready                                         # сколько сообщений ждёт в очередях
videostat_videos                                                               # видео по статусам
```

### Алерты

Правила — файлами в git, отправка уведомлений пока **не настроена** (Alertmanager не подключён): сработавшие алерты видны в Grafana → Alerting, на http://localhost:9090/alerts и на дашборде Overview.

| Файл | Алерты |
|---|---|
| `monitoring/prometheus/rules/alerts.yml` | **critical:** `ServiceDown`, `PostgresDown`, `QueueNoConsumers`, `CronTaskStale`, `HTTP5xxHigh`; **warning:** `QueueBacklog`, `WorkerMessageErrors`, `VideosFailed`, `RefreshBloggersStale`, `HTTPLatencyHigh`, `ExternalAPIErrors`, `PostgresConnectionsHigh` |
| `monitoring/loki/rules/alerts.yml` | `PanicInLogs` (critical) — паника, перехваченная Recovery; `ErrorLogsSpike` (warning) |

Проверить синтаксис правил Prometheus перед коммитом:

```bash
docker run --rm --entrypoint promtool -v "$PWD/monitoring/prometheus:/p:ro" prom/prometheus:v3.5.0 check rules /p/rules/alerts.yml
```

### Логи

Сервисы пишут JSON-логи (`slog`) одновременно в терминал и в файл `logs/<сервис>.log` (если задан `LOG_DIR`). Файлы ротируются: 50 МБ, 3 архива, 7 дней. Alloy читает файлы и отправляет в Loki, хранение — 7 дней.

**Метки Loki** (по ним фильтруется быстро): `service` (`api`, `worker_core`, `worker_notify`, `worker_cron` — из имени файла), `level`, `env`.
Всё остальное — поля внутри JSON, фильтруются через `| json`.

**Соглашения по полям логов:**

- ключи только в **snake_case** (`video_id`, не `videoID`) — проверяет линтер `sloglint` (`make lint`);
- ошибка — всегда ключ `error`;
- `video_id`, `blogger_id`, `event_type`, `event_id`, `trace_id` добавляются **автоматически** в точках входа (обработка события, HTTP, cron, бот) через `logger.WithField` / `logger.WithTraceID` — в командах их передавать не нужно, достаточно брать логгер через `logger.FromContext(ctx)`;
- имена общих полей — константы в `internal/infrastructure/logger/fields.go`.

**Полезные запросы (LogQL, Grafana → Explore → Loki):**

```logql
{env="local", level="ERROR"}                                    # все ошибки всех сервисов
{service="worker_core"} |= "timeout"                            # поиск по тексту
{env="local"} | json | video_id="<id>"                          # весь путь видео через все сервисы
{env="local"} | json | trace_id="<id>"                          # все логи одного трейса
{env="local", level="ERROR"} | json | line_format "{{.service}} {{.video_id}} {{.msg}}: {{.error}}"
sum by (service) (count_over_time({env="local", level="ERROR"}[5m]))   # ошибки по сервисам (график)
```

### Трейсинг

Трейсинг — через **OpenTelemetry** (код не зависит от хранилища), хранение — **Tempo**. Включается переменной `TRACING_OTLP_ENDPOINT`; если она пустая, трейсинг выключен и код работает вхолостую.

**Что создаёт спаны:**

| Что | Имя спана | Где |
|---|---|---|
| Сообщение / команда в Telegram-боте | `telegram /<команда>`, `telegram callback` | `transport/telegram/command/router.go` |
| Входящий HTTP-запрос | `GET /test/videos-apify/{id}` | `transport/http/middleware/tracing.go` |
| Входящий gRPC-вызов | имя метода | `platform/grps_server.go` |
| Публикация в RabbitMQ | `publish <routing key>` | `infrastructure/mq/rabbit_publisher.go` |
| Обработка события из RabbitMQ | `handle <event_type>` (атрибут `video_id`) | `worker/router.go` |
| Запуск cron-задачи | `cron <task>` | `worker/cron/worker.go` |
| SQL-запрос | `SELECT`, `UPDATE`… (текст — в `db.query.text`) | `infrastructure/db/postgres.go` (`otelpgx`) |
| Вызов внешнего API | `heygen POST`, `anthropic POST`… | `infrastructure/metrics/external_api.go` |

**Трейс проходит через RabbitMQ:** publisher кладёт заголовок `traceparent` (W3C) в заголовки AMQP-сообщения, consumer его читает — команда в боте, публикация и обработка в `worker_core` видны одним деревом.
Цепочка обрывается там, где работа продолжается по расписанию: после отправки в HeyGen/Kling/Shotstack статус опрашивает cron, и это новый трейс. Весь жизненный цикл видео ищется по `video_id`.

**Полезные запросы (TraceQL, Grafana → Explore → Tempo):**

```
{ resource.service.name = "worker_core" }
{ span.video_id = "<id>" }                       # все трейсы обработки видео
{ status = error }                               # трейсы с ошибками
{ duration > 5s }                                # медленные
{ name =~ "heygen.*" && duration > 5s }          # медленные вызовы HeyGen
```

### Как всё связано

Метрики, логи и трейсы связаны через `trace_id` и одинаковые имена сервисов (`job` в Prometheus = `service` в Loki = `service.name` в Tempo):

- **метрика → трейс:** на графиках времени (p95) есть точки-exemplars — клик открывает трейс конкретного медленного запроса;
- **трейс → логи:** в трейсе у спана кнопка «Logs for this span» — логи именно этого трейса;
- **лог → трейс:** у строки лога с `trace_id` кнопка «Открыть трейс»;
- **метрика → логи:** на дашборде Overview клик по сервису / по линии ошибок открывает его логи за тот же период.

Типичный разбор: на Overview вырос p95 обработки → клик по точке на пике → в трейсе видно, какой вызов занял время → «Logs for this span» → причина в логах.

### Как добавить своё

- **Метрику** — `promauto.NewCounterVec/NewHistogramVec/...` рядом с кодом, который её обновляет (см. `worker/metrics.go`). В метки — только значения из небольшого фиксированного набора (статус, тип события, провайдер); **никаких id** — каждое значение метки создаёт отдельный временной ряд. Для длительностей — `metrics.ObserveWithTrace(ctx, …)`, чтобы появились exemplars.
- **Поле в логах** — ключ в snake_case. Если поле должно быть во всех логах операции — `ctx = logger.WithField(ctx, key, value)` в начале операции.
- **Спан** — `ctx, span := tracer.Start(ctx, "имя"); defer span.End()` и передавать этот `ctx` дальше. Если где-то посередине взять `context.Background()`, трейс разорвётся.
- **Алерт** — правило в `monitoring/prometheus/rules/alerts.yml` (PromQL) или `monitoring/loki/rules/alerts.yml` (LogQL). Prometheus перечитывает конфиг без рестарта: `curl -X POST localhost:9090/-/reload`.
- **Дашборд** — собрать в Grafana, затем `Edit → Settings → JSON Model` и сохранить в `monitoring/grafana/dashboards/*.json`. Правки, сделанные только в UI, в git не попадают.

### Если что-то не работает

| Симптом | Что проверить |
|---|---|
| Сервис `DOWN` в Prometheus | сервис запущен? `curl localhost:<порт>/metrics`; после правки кода — перезапустить сервис |
| Правка `prometheus.yml` не применилась | `curl -X POST localhost:9090/-/reload`; проверить Status → Configuration (Docker Desktop синхронизирует файлы с задержкой — иногда нужен повторный reload) |
| Нет логов в Loki | в `.env` задан `LOG_DIR=logs`? появились файлы в `logs/`? ошибки в UI Alloy (http://localhost:12345) |
| Нет трейсов | в `.env` задан `TRACING_OTLP_ENDPOINT=127.0.0.1:4317` — **именно IP**: с `localhost` gRPC-экспорт на macOS зависал. Для диагностики — `otel.SetErrorHandler(...)`: по умолчанию SDK теряет спаны молча |
| Трейс «рвётся» на части | где-то в цепочке вызов идёт с `context.Background()` вместо переданного `ctx` |
| Нет точек-exemplars на графиках | сервис перезапущен с новым кодом? `/metrics` должен отдаваться в формате OpenMetrics (`metrics.Handler()`) |
| Новый контейнер мониторинга не поднялся после `make down` / `make up` | он должен быть в списке сервисов `make up` в `Makefile` |

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
| `make up` | запустить docker-compose: postgres, rabbitmq, minio и мониторинг (prometheus, grafana, loki, tempo, alloy, postgres-exporter) |
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
│   │   ├── logger/             # slog: файл логов, поля контекста (video_id, trace_id)
│   │   ├── metrics/            # метрики внешних API, VideoCollector, exemplars, /metrics
│   │   ├── mq/                 # RabbitMQ publisher/consumer (+ traceparent в заголовках)
│   │   ├── llm/openai/         # генерация B-roll промптов
│   │   ├── notifier/           # Telegram уведомления
│   │   ├── repository/blogger/ # PostgreSQL + InMemory
│   │   ├── shotstack/          # клиент Shotstack API (720×1280)
│   │   ├── storage/s3/         # клиент S3/MinIO
│   │   └── videogenerator/heygen/ # клиент HeyGen API (1280×720)
│   ├── platform/               # фабрики клиентов, InitTracing, сервер /metrics воркеров
│   └── worker/
│       ├── core/               # обработчики событий RabbitMQ
│       ├── cron/               # cron задачи
│       └── notify/             # обработчики уведомлений
├── migrations/                 # SQL миграции
├── logs/                       # файлы логов сервисов (в git не попадают)
└── monitoring/
    ├── prometheus/             # prometheus.yml (что собирать) + rules/alerts.yml
    ├── grafana/                # источники данных и дашборды (provisioning)
    ├── loki/                   # loki.yml + rules/alerts.yml (алерты по логам)
    ├── tempo/                  # tempo.yml
    └── alloy/                  # config.alloy: сбор логов из logs/ и приём трейсов
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
