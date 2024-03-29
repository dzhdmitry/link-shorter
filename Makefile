include .env

.PHONY: build up down ps logs bash bash-root migration migrate test test-coverage lint swagger

build:
	docker-compose build

up:
	docker-compose up -d

down:
	docker-compose down

ps:
	docker-compose ps

logs:
	docker-compose logs go

bash:
	docker-compose exec go bash

bash-root:
	docker-compose exec -u root go bash

migration:
	docker-compose exec go migrate create -seq -ext=.sql -dir=./migrations ${name}

migrate:
	docker-compose exec go migrate -path=./migrations -database=${DB_DSN} up

test:
	docker-compose exec go go test -p 1 ./...

test-coverage:
	docker-compose exec go go test -p 1 ./... -covermode=set -coverprofile tmp/coverage.out
	docker-compose exec go go tool cover -html tmp/coverage.out -o tmp/coverage.html
	docker-compose exec go rm tmp/coverage.out
	@echo '   *** See tmp/coverage.html ***'

lint:
	docker-compose run --rm golangci golangci-lint run -v

swagger:
	docker-compose exec go swag init -q
