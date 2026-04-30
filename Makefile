include .env
export

env-up:
	docker compose up -d concert-postgres concert-redis concert-minio concert-minio-setup

env-down:
	docker compose down

env-cleanup:
	@echo "Delete ALL data (Postgres & Minio)? [y/N]"; \
	read ans; \
	if [ "$$ans" = "y" ]; then \
		docker compose down && \
		rm -rf ./out/pgdata ./out/miniodata && \
		echo "Environment cleaned"; \
	else \
		echo "Cleanup cancelled"; \
	fi

env-port-forward:
	@docker compose up -d port-forwarder-postgres port-forwarder-redis port-forwarder-minio port-forwarder-minio-console

env-port-close:
	@docker compose down port-forwarder-postgres port-forwarder-redis port-forwarder-minio port-forwarder-minio-console

migrate-create:
	@docker compose run --rm concert-postgres-migrate \
		create \
		-ext sql \
		-dir /migrations \
		-seq "$(seq)"

migrate-up:
	@make migrate-action action=up

migrate-down:
	@make migrate-action action=down

migrate-action:
	@docker compose run --rm concert-postgres-migrate \
		-path /migrations \
		-database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@concert-postgres:5432/${POSTGRES_DB}?sslmode=disable \
		"${action}"


test-db:
	@echo "--- Testing PostgreSQL Connection ---"
	@docker exec -it concert-env-postgres psql -U $(POSTGRES_USER) -d $(POSTGRES_DB) -c "SELECT 1 as connection_ok;"
	@echo "--- Checking Tables ---"
	@docker exec -it concert-env-postgres psql -U $(POSTGRES_USER) -d $(POSTGRES_DB) -c "\dt"

test-redis:
	@echo "--- Testing Redis ---"
	@docker exec -it concert-env-redis redis-cli ping

test-s3:
	@echo "--- Testing Minio API ---"
	@docker run --rm curlimages/curl:7.85.0 -I -s http://host.docker.internal:9000/minio/health/live | findstr "HTTP/"

test-all: test-db test-redis test-s3

migrate-force:
	@docker compose run --rm concert-postgres-migrate \
		-path /migrations \
		-database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@concert-postgres:5432/${POSTGRES_DB}?sslmode=disable \
		force $(v)

app-run: export LOGGER_FOLDER := $(CURDIR)/out/logs
app-run:
	@go mod tidy && go run cmd/server/main.go