#!make
-include .env

.PHONY: docker.up
docker.up:
	docker-compose --env-file .env -f docker-compose.yml -p audio-storage up -d --build

.PHONY: docker.down
docker.down:
	docker-compose -f docker-compose.yml -p audio-storage down

.PHONY: postgres.login
postgres.login:
	docker-compose exec postgres psql "postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@localhost:${POSTGRES_PORT}/${POSTGRES_DB_NAME}"

.PHONY: k6-test
k6-test:
	docker run --rm -i \
		-v $(shell pwd)/e2e/audio-storage-test.js:/audio-storage-test.js \
		-v $(shell pwd)/e2e/testdata:/testdata \
		-e BASE_URL="http://host.docker.internal:8080" \
		grafana/k6 run /audio-storage-test.js

.PHONY: migrate.dep
migrate.dep:
ifeq ($(shell uname), Linux)
	@curl -sSfL https://github.com/golang-migrate/migrate/releases/download/v4.17.1/migrate.linux-amd64.tar.gz | tar zxf - --directory /tmp \
	&& cp /tmp/migrate .
else ifeq ($(shell uname), Darwin)
	@curl -sSfL https://github.com/golang-migrate/migrate/releases/download/v4.17.1/migrate.darwin-amd64.tar.gz | tar zxf - --directory /tmp
	&& cp /tmp/migrate .
else
	@echo "Your OS is not supported."
endif
