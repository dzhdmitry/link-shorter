# Сокращатель ссылок 

## Установка и запуск

<table>
    <thead>
        <tr>
            <th>Docker-compose</th>
            <th>Native</th>
        </tr>
    </thead>
    <tbody>
        <tr>
            <td colspan="2" align="center">Требования</td>
        </tr>
        <tr>
            <td>
                docker-compose<br>make
            </td>
            <td>
                go 1.21<br>
                postgreSQL<br>
                <a href="https://github.com/golang-migrate/migrate">golang-migrate</a>
            </td>
        </tr>
        <tr>
            <td colspan="2" align="center">Установка</td>
        </tr>
        <tr>
            <td>
                <pre>
make build up migrate down</pre>
            </td>
            <td>
                <pre>
go migrate -path=./migrations -database=${DB_DSN} up</pre>
            </td>
        </tr>
        <tr>
            <td colspan="2" align="center">Запуск</td>
        </tr>
        <tr>
            <td>
                <pre>
make up</pre>
            </td>
            <td>
                <pre>
go run .</pre>
                (<code>go run . -help</code> для просмотра опций)
            </td>
        </tr>
      </tbody>
</table>

## Тесты

```bash
make test
make test-coverage
```
