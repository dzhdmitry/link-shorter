# Сокращатель ссылок 

## Требования

- go 1.21

или

- docker-compose
- make

## Сборка и запуск

```bash
go build -o tmp/app
./tmp/app -host=<HOST> -port=<PORT> -key_max_length=<MAX_LENGTH_OF_LINK_KEY>
```

или

```bash
make build
make up
```

## Тесты

```bash
make test
```
