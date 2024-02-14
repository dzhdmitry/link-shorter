# Сокращатель ссылок 

## Установка и запуск (docker)

### Установка

* Скопировать <code>.env.example</code> в <code>.env</code>
* `make build up migrate down`

### Запуск

`make up`

## Установка и запуск (нативно)

<details>
  <summary>Установка и запуск (нативно)</summary>

### Требования

* go 1.21
* postgreSQL
* [golang-migrate](https://github.com/golang-migrate/migrate)

### Установка

`go migrate -path=./migrations -database=${DB_DSN} up`

### Запуск

`go run .` (`go run . -help` для просмотра опций)
</details>

## Тесты

```bash
make test
make test-coverage
```
