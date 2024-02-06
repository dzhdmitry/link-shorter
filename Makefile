include .env

.PHONY: build up down ps bash

build:
	docker-compose build

up:
	docker-compose up -d

down:
	docker-compose down

ps:
	docker-compose ps

bash:
	docker-compose exec go bash
