include .env

.PHONY: build up down ps bash test test-coverage

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

test:
	docker-compose exec go go test ./...

test-coverage:
	docker-compose exec go go test -v ./... -covermode=count -coverpkg=./... -coverprofile coverage/coverage.out
	docker-compose exec go go tool cover -html coverage/coverage.out -o coverage/coverage.html
	@echo 'See coverage/coverage.html'