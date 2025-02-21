#!make
-include .env

.PHONY: docker.up
docker.up:
	docker-compose --env-file .env -f docker-compose.yml -p audio-storage up

.PHONY: docker.down
docker.down:
	docker-compose -f docker-compose.yml -p audio-storage down

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
