# AvitoCoinStore

## Инструкция по запуску

Для запуска сервиса нужен установленный docker, тогда достаточно выполнить следующую команду:

    docker-compose up --build

После запуска контейнера, сервер будет доступен в `localhost:8080`

## Использованный стек

Язык сервиса: Go

База данных: PostgreSQL

## Выполненно

- Сервис соответствует [API](./static/schema.json)

- После первой и последующих авторизаций генерируется и выдаётся личный JWT токен для пользователя. Для генерации используется скрытый от внешнего наблюдателя `user_id` и секретное слово объявленное в файле .env

- Написаны юнит тесты для хендлеров и функций взаимодействующих с базой данных (подробнее в отдельной папке [tests](./tests/))

- Написано интеграционное тестирование для разных сценариев работы (см. [tests](./tests/))

- Проведено нагрузочное тестирование (см. [tests](./tests/))

- Конфигурацию линтера можно посмотреть в файле [.golangci.yaml](./.golangci.yaml)

- Большинство функций имеют краткое описание того что они делают и для чего нужны

## Примечания

- Вообще по-хорошему файл [.env](./.env) не должен быть запушен и более того, .env файл отдельно перечислен в [.gitignore](./.gitignore), однако я решил, что лучше будет если он будет доступен сразу после клонирования в репозиторий. Тем более что там не лежит ничего ценного, да и проект это просто выполнение поставленного задания.

- Для тестирования отдельных частей кода приходилось делать специальные интерфейсы, что среди прочего привело к появлению обертки над pgxpool.Pool, однако позволило добиться большего покрытия тестами.
