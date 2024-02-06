.PHONY: run-local run-docker-dev run-docker-prod down load-test

run-local:
	sqlc generate
	DB_HOST=localhost DB_POOL=10 air

run-docker-dev:
	docker compose up -d --build
	docker compose logs -f

run-docker-prod:
	docker compose up -d

down:
	docker compose down

load-test:
	./load-test.sh
