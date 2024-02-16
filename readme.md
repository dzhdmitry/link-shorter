# Сокращатель ссылок

## Установка и запуск (docker)

* Скопировать <code>.env.example</code> в <code>.env</code>
* `make build up migrate down`
* `make up`

## Установка и запуск (нативно)

<details>
  <summary>Установка и запуск (нативно)</summary>

### Требования

* go 1.21
* postgreSQL
* redis
* [golang-migrate](https://github.com/golang-migrate/migrate)

### Установка и запуск

`go migrate -path=./migrations -database=postgres://user:password@host:5432/short_links?sslmode=disable up` (заменить `user`, `password`, `host` требуемыми значениями)

`go run .` (`go run . -help` для просмотра опций)
</details>

## Тесты и линтеры

```bash
make test
make test-coverage
make lint
```

## Описание

Простой сокращатель ссылок, преобразует ссылку, например https://github.com/dzhdmitry?tab=repositories&language=go,
в короткую ссылку вида `http://host/d3s`.

Для каждой ссылки генерируется кототкий уникальный токен, состоящий из набора [0-9a-z] и получающийся путём применения [биективной функции](http://en.wikipedia.org/wiki/Bijection) к порядковому номеру ссылки.
То есть, для 1-й ссылки будет токен `1`, для 10-й - `a`, для 10000-й - `7ps` и т.д.
При поиске полной ссылки её номер получается обратным преобразованием: `7ps` -> №10000, `d3s4c` -> №22011420 и т.д.

Особенности:

1. С сервисом можно работать JSON-запросами, получая JSON в ответ, есть batch-запросы, можно гененировать/получать множесто ссылок.
2. Все запросы логируются в stdout, невалидные запросы обрабатывабтся, отдаётся корректный ответ.
3. Поддерживается GZIP-сжатие данных http-запросов.
4. Может хранить данные в двух режимах:
   * в памяти с синхронным и асинхронным сохранением в файл, при запуске может восстанавливаться из файла, при остановке "дожидается" асинхронных задач
   * в postgreSQL
5. Может кэшировать данные 
   * в памяти, реализована статегия вытеснения [LFU](https://en.wikipedia.org/wiki/Least_frequently_used) при заполнении кэша.
   * в Redis.
6. Может ограничивать кол-во запросов к сервису от одного IP, при превышении предела отдаёт HTTP-код 429.
7. Параметры (тип хранилища, параметры соединения с бд, объём кеша, кол-во запросов) задаются переменными окружения и Args командной строки.

### Endpoint-ы:

1. `POST http://localhost/generate`
   * запрос: `{"URL": "http://a.com/path"}`
   * ответ: `{"link": "http://localhost/go/7pv"}`

2. `GET http://localhost/go/:key`
   * ответ: `{"link": "http://localhost/go/7pv"}`

3. `POST http://localhost/batch/generate`
   * запрос: `["http://a.com/path", "http://b.com/path"]`
   * ответ: `{"links": {"http://a.com/path": "http://localhost/go/7pv", "http://b.com/path": "http://localhost/go/5v6"}}`

4. `GET http://localhost/batch/go`
   * запрос: `["7pv", "5v6"]`
   * ответ: `{"links": {"7pv": "http://a.com/path", "5v6":"http://b.com/path"}}`

Все входящие ссылки должны быть валидными URL-ами.
Ссылки в запросе `/batch/generate` не должны повторяться.

## Лицензия

[MIT](https://github.com/dzhdmitry/link-shorter?tab=MIT-1-ov-file)
