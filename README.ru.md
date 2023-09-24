[![Golang](https://img.shields.io/badge/Go-v1.21-EEEEEE?logo=go&logoColor=white&labelColor=00ADD8)](https://go.dev/)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)

<div align="center">
    <h1>TO-DO сервис</h1>
    <h5>
        Микросервис, написанный на Golang, предназначенный для планирования ваших ежедневных задач.
    </h5>
    <p>
        <a href="README.md">English</a> | Russian 
    </p>
</div>

---

## Используемые технологии:
- [Golang](https://go.dev), [PostgreSQL](https://www.postgresql.org/)
- [Docker](https://www.docker.com/), [JWT Аутентификация](https://jwt.io/)
- [REST](https://ru.wikipedia.org/wiki/REST), [Swagger UI](https://swagger.io/tools/swagger-ui/)

---

## Навигация
* **[Загрузка](#Загрузка)**
* **[Примеры запросов](#примеры-запросов)**
    * [Пользователи](#пользователи)
    * [TODO-лист](#todo-лист)
* **[Дополнительные возможности](#дополнительные-возможности)**

---

## Загрузка
```shell
git clone git@github.com:zenorachi/todo-service.git
```

---

## Начало работы
1. **Настройка переменных окружения (создайте файл .env в корне проекта):**
```dotenv
# База данных
export DB_HOST=
export DB_PORT=
export DB_USER=
export DB_NAME=
export DB_SSLMODE=
export DB_PASSWORD=

# Локальная БД
export LOCAL_DB_PORT=

# Сервис postgres
export POSTGRES_PASSWORD=

# Секреты для пароля
export HASH_SALT=
export HASH_SECRET=

# GIN мод (необязательно, по умолчанию - release)
export GIN_MODE=
```
> **Подсказка:** если вы запускаете проект с помощью Docker, установите `DB_HOST`=postgres (как имя сервиса Postgres в docker-compose).

2. **Запуск сервиса:**
```shell
make
```
3. **Чтобы протестировать работу сервиса, можно перейти по адресу
   http://localhost:8080/docs/index.html для получения Swagger документации.**
> **Подсказка:** чтобы пройти авторизацию в Swagger UI после получения JWT токена, в поле для ввода необходимо
ввести `Bearer <полученный_токен>` (без символов "<" и ">").

---

## Примеры запросов

### Пользователи
#### 1. Регистрация
* Пример запроса:
```shell
curl -X 'POST' \
  'http://localhost:8080/api/v1/auth/sign-up' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "login": "maksim-go",
  "email": "maksim-go@gmail.com",
  "password": "qwerty123"
}'
```
* Пример ответа:
```json
{
  "id": 1
}
```

#### 2. Аутентификация
* Пример запроса:
```shell
curl -X 'POST' \
  'http://localhost:8080/api/v1/auth/sign-in' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "login": "maksim-go",
  "password": "qwerty123"
}'
```
* Пример ответа:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTMzMDc4MjksInN1YiI6IjE0In0.N1QBZb1uVZQGJ7vROHhCdlaySu1o31yjTzFLnVk_XYw"
}
```

#### 3. Обновление токена
* Пример запроса:
```shell
curl -X 'GET' \
  'http://localhost:8080/api/v1/auth/refresh' \
  -H 'accept: application/json'
```
* Пример ответа:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTMzMDc4NTIsInN1YiI6IjE0In0.cmXwc15TmNSI2mILSZjoqRhhtUN2AYZQu5had9OW07k"
}
```

### TODO-лист
#### 1. Создание задачи
* Пример запроса:
```shell
curl -X 'POST' \
  'http://localhost:8080/api/v1/agenda/create' \
  -H 'accept: application/json' \
  -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTU1Nzk1MzYsInN1YiI6IjIifQ.AhqXvtQPHMwp1Pv5Y9m6xXnMITUlRDnGo8oVm5DRvLU' \
  -H 'Content-Type: application/json' \
  -d '{
  "date": "2023-Sep-26",
  "description": "Description",
  "status": "not done",
  "title": "Task 1"
}'
```
* Пример ответа:
```json
{
  "id": 1
}
```

#### 2. Получение задачи по ID
* Пример запроса:
```shell
curl -X 'GET' \
  'http://localhost:8080/api/v1/agenda/1' \
  -H 'accept: application/json' \
  -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTU1Nzk1MzYsInN1YiI6IjIifQ.AhqXvtQPHMwp1Pv5Y9m6xXnMITUlRDnGo8oVm5DRvLU'
```
* Пример ответа:
```json
{
   "task": {
      "title": "Task 1",
      "description": "Description",
      "date": "2023-09-26T00:00:00Z",
      "status": "not done"
   }
}
```

#### 3. Получение всех задач пользователя
* Пример запроса:
```shell
curl -X 'GET' \
  'http://localhost:8080/api/v1/agenda/get_all' \
  -H 'accept: application/json' \
  -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTU1Nzk1MzYsInN1YiI6IjIifQ.AhqXvtQPHMwp1Pv5Y9m6xXnMITUlRDnGo8oVm5DRvLU'
```
* Пример ответа:
```json
{
  "tasks": [
      {
         "title": "Task 1",
         "description": "Description",
         "date": "2023-09-26T00:00:00Z",
         "status": "not done"
      },
      {
         "title": "Task 2",
         "description": "Description",
         "date": "2023-09-26T00:00:00Z",
         "status": "not done"
      }
  ]
}
```

#### 4. Получение всех задач пользователя с пагинацией и определенным статусом
* Пример запроса:
```shell
curl -X 'GET' \
  'http://localhost:8080/api/v1/agenda/get_by_date?page=5&status=not_done' \
  -H 'accept: application/json' \
  -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTU1Nzk1MzYsInN1YiI6IjIifQ.AhqXvtQPHMwp1Pv5Y9m6xXnMITUlRDnGo8oVm5DRvLU' \
  -H 'Content-Type: application/json' \
  -d '{
  "limit": 2,
  "offset": 2
}'
```
* Пример ответа:
```json
{
   "tasks": [
      {
         "title": "some task2",
         "description": "desc",
         "date": "2023-09-26T00:00:00Z",
         "status": "not done"
      },
      {
         "title": "some task3",
         "description": "desc",
         "date": "2023-09-26T00:00:00Z",
         "status": "not done"
      }
   ]
}
```

#### 5. Получение всех задач по дате
* Пример запроса:
```shell
curl -X 'GET' \
  'http://localhost:8080/api/v1/agenda/get_by_date' \
  -H 'accept: application/json' \
  -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTU1Nzk1MzYsInN1YiI6IjIifQ.AhqXvtQPHMwp1Pv5Y9m6xXnMITUlRDnGo8oVm5DRvLU' \
  -H 'Content-Type: application/json' \
  -d '{
  "date": "2023-Sep-26",
  "limit": 5,
  "offset": 1
}'
```
* Пример ответа:
```json
{
   "tasks": [
      {
         "title": "some task1",
         "description": "desc",
         "date": "2023-09-26T00:00:00Z",
         "status": "not done"
      },
      {
         "title": "some task2",
         "description": "desc",
         "date": "2023-09-26T00:00:00Z",
         "status": "not done"
      },
      {
         "title": "some task3",
         "description": "desc",
         "date": "2023-09-26T00:00:00Z",
         "status": "not done"
      }
   ]
}
```

#### 6. Установить новый статус задачи
* Пример запроса:
```shell
curl -X 'PUT' \
  'http://localhost:8080/api/v1/agenda/set_status' \
  -H 'accept: application/json' \
  -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTU1Nzk1MzYsInN1YiI6IjIifQ.AhqXvtQPHMwp1Pv5Y9m6xXnMITUlRDnGo8oVm5DRvLU' \
  -H 'Content-Type: application/json' \
  -d '{
  "status": "done",
  "task_id": 10
}'
```
* Пример ответа: *None*
> **Подсказка:** если обновление прошло успешно, сервер вернет 204 код (NO CONTENT).

#### 7. Удалить задачу по ID
* Пример запроса:
```shell
curl -X 'DELETE' \
  'http://localhost:8080/api/v1/agenda/delete_by_id' \
  -H 'accept: application/json' \
  -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTU1Nzk1MzYsInN1YiI6IjIifQ.AhqXvtQPHMwp1Pv5Y9m6xXnMITUlRDnGo8oVm5DRvLU' \
  -H 'Content-Type: application/json' \
  -d '{
  "task_id": 10
}'
```
* Пример ответа: *None*
> **Подсказка:** если удаление прошло успешно, сервер вернет 204 код (NO CONTENT).

#### 8. Удалить все задачи пользователя
* Пример запроса:
```shell
curl -X 'DELETE' \
  'http://localhost:8080/api/v1/agenda/delete_all' \
  -H 'accept: application/json' \
  -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTU1Nzk1MzYsInN1YiI6IjIifQ.AhqXvtQPHMwp1Pv5Y9m6xXnMITUlRDnGo8oVm5DRvLU'
```
* Пример ответа: *None*
> **Подсказка:** если удаление прошло успешно, сервер вернет 204 код (NO CONTENT).

---

## Дополнительные возможности
1. **Запуск тестов**
```shell
make test
```
2. **Запуск линтера**
```shell
make lint
```
3. **Создание файлов миграций**
```shell
make migrate-create
```
4. **Миграции вверх / вниз**
```shell
make migrate-up
```
```shell
make migrate-down
```
5. **Остановка всех запущенных контейнеров**
```shell
make stop
```
