#!make
-include .env

.PHONY: docker.up
docker.up:
	docker-compose -f docker-compose.yml -p audio-storage up -d

.PHONY: docker.down
docker.down:
	docker-compose -f docker-compose.yml -p audio-storage down