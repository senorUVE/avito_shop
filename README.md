
## Reqire
подгружать через `.env-docker`_
| Env | Value | Описание |
|----------|----------|----------|
|GOOSE_DRIVER| postgres | `env` для миграций |
|GOOSE_DBSTRING| postgres dsn  | dsn postres sql для миграций |
| JWT_KEY_PRIVATE   | `rsa256`   | ключ шифрования   |
| JWT_KEY_PUBLIC    | `rsa256`   | ключ шифрования   |
|JWT_TOKEN_AC_TIME| время (5m) | время действия access токена |
|JWT_TOKEN_RE_TIME| время (24h) | время действия refresh токена |
|DB_DSN| postgres dsn | dsn для postgres sql |

Используется _jwt_, шифрования по схеме `RS256`
jwt : header, payload, signature
private key: header, payload --> signature
public key: signature --> header, payload


authdb создается при инициализации docker. Миграции накатываются при запуске контейнера.
+ Конфигурация для подключения к postgresql указывается в файле `.env.example`

__КАК ЗАПУСТИТЬ?__

ТРЕБУЮТСЯ УКАЗАННЫЕ env-ары


1) Генерим ключи `go run cmd/keygen/main.go` и добавляем их в `.env`

2) docker compose --build up -d

# Доп задания

* Провести нагрузочное тестирование полученного решения и приложить результаты тестирования 
* Реализовать интеграционное или E2E-тестирование для остальных сценариев  
* Описать конфигурацию линтера (.golangci.yaml в корне проекта для go, phpstan.neon для PHP или ориентируйтесь на свои, если используете другие ЯП для выполнения тестового)


СДЕЛАНЫ!

# Вопросы/Проблемы

При выполнении доп задания "Реализовать интеграционное или E2E-тестирование для остальных сценариев" сценарии явно не указаны, поэтому были выбраны 2 сценария:

1) Перевод монет сотруднику с недостающим балансом

2) Покупка мерча с недостающим балансом


# Нагрузочные тесты
```text
Running 30s test @ http://localhost:8080/api/auth
  10 threads and 1000 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    46.29ms   43.80ms 572.02ms   93.47%
    Req/Sec     2.53k     0.91k    5.53k    69.90%
  755423 requests in 30.10s, 91.49MB read
  Socket errors: connect 0, read 1793, write 3, timeout 0
  Non-2xx or 3xx responses: 755423
Requests/sec:  25100.78
Transfer/sec:      3.04MB
```

```text
Running 30s test @ http://localhost:8080/api/auth
  10 threads and 1000 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    48.21ms   47.66ms 474.80ms   91.33%
    Req/Sec     2.61k     0.99k    4.45k    72.96%
  770010 requests in 30.06s, 93.26MB read
  Socket errors: connect 0, read 2640, write 319, timeout 0
  Non-2xx or 3xx responses: 770010
Requests/sec:  25616.59
Transfer/sec:      3.10MB
```