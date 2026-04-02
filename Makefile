ENV_FILE ?= .env
BUF_VERSION := 1.36.0

define run_with_env
	@set -a; source $(ENV_FILE); set +a; $(1)
endef

up:
	docker compose up -d db rabbitmq

down:
	docker compose down

run:
	$(call run_with_env,go run ./cmd/api)

run_worker_notify:
	$(call run_with_env,go run ./cmd/worker_notify)

run_worker_cron:
	$(call run_with_env,go run ./cmd/worker_cron)

test-short:
	$(call run_with_env,go test -v -race -count=1 ./...)

lint:
	$(call run_with_env,golangci-lint run)

migrate-up:
	docker compose run --rm migrate up

migrate-down:
	docker compose run --rm migrate down 1

migrate-version:
	docker compose run --rm migrate version

migrate-create:
	docker compose run --rm migrate create -ext sql -dir /migrations $(name)


proto-gen:
	docker compose run --rm buf generate

proto-lint:
	docker compose run --rm buf lint

proto-breaking:
	docker compose run --rm buf breaking --against '.git#branch=main'
