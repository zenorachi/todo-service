[![Golang](https://img.shields.io/badge/Go-v1.21-EEEEEE?logo=go&logoColor=white&labelColor=00ADD8)](https://go.dev/)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)

<div align="center">
    <h1>TO-DO service</h1>
    <h5>
        A microservice written in the Go programming language is designed to plan your daily tasks.
    </h5>
    <p>
        English | <a href="README.ru.md">Russian</a> 
    </p>
</div>

---

## Technologies used:
- [Golang](https://go.dev), [PostgreSQL](https://www.postgresql.org/)
- [Docker](https://www.docker.com/), [JWT Authentication](https://jwt.io/)
- [REST](https://ru.wikipedia.org/wiki/REST), [Swagger UI](https://swagger.io/tools/swagger-ui/)

---

## Navigation
* **[Installation](#installation)**
* **[Example of requests](#examples-of-requests)**
  * [Users](#users)
  * [Agenda](#agenda)
* **[Additional features](#additional-features)**

---

## Installation
```shell
git clone git@github.com:zenorachi/todo-service.git
```

---

## Getting started
1. **Setting up environment variables (create a .env file in the project root):**
```dotenv
# Database
export DB_HOST=
export DB_PORT=
export DB_USER=
export DB_NAME=
export DB_SSLMODE=
export DB_PASSWORD=

# Local database
export LOCAL_DB_PORT=

# Postgres service
export POSTGRES_PASSWORD=

# Password Hasher
export HASH_SALT=
export HASH_SECRET=

# GIN mode (optional, default - release)
export GIN_MODE=
```
> **Hint:**
if you are running the project using Docker, set `DB_HOST` to "**postgres**" (as the service name of Postgres in the docker-compose).

2. **Compile and run the project:**
```shell
make
```
3. **To test the service's functionality, you can navigate to the address
   http://localhost:8080/docs/index.html to access the Swagger documentation.**
> **Hint:** to complete the authorization in Swagger UI after receiving the JWT token, you need
to enter `Bearer <your_token>` (without "<" and ">" symbols) in the input field.

---

## Examples of requests

### Users
#### 1. Registration
* Request example:
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
* Response example:
```json
{
  "id": 1
}
```

#### 2. Authentication
* Request example:
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
* Response example:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTMzMDc4MjksInN1YiI6IjE0In0.N1QBZb1uVZQGJ7vROHhCdlaySu1o31yjTzFLnVk_XYw"
}
```

#### 3. Refresh token
* Request example:
```shell
curl -X 'GET' \
  'http://localhost:8080/api/v1/auth/refresh' \
  -H 'accept: application/json'
```
* Response example:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTMzMDc4NTIsInN1YiI6IjE0In0.cmXwc15TmNSI2mILSZjoqRhhtUN2AYZQu5had9OW07k"
}
```

### Agenda
#### 1. Create Task
* Request example:
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
* Response example:
```json
{
  "id": 1
}
```

#### 2. Get Task By ID
* Request example:
```shell
curl -X 'GET' \
  'http://localhost:8080/api/v1/agenda/1' \
  -H 'accept: application/json' \
  -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTU1Nzk1MzYsInN1YiI6IjIifQ.AhqXvtQPHMwp1Pv5Y9m6xXnMITUlRDnGo8oVm5DRvLU'
```
* Response example:
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

#### 3. Get All User Tasks
* Request example:
```shell
curl -X 'GET' \
  'http://localhost:8080/api/v1/agenda/get_all' \
  -H 'accept: application/json' \
  -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTU1Nzk1MzYsInN1YiI6IjIifQ.AhqXvtQPHMwp1Pv5Y9m6xXnMITUlRDnGo8oVm5DRvLU'
```
* Response example:
```json
{
  "tasks": [
      {
         "id": 1,
         "title": "Task 1",
         "description": "Description",
         "date": "2023-09-26T00:00:00Z",
         "status": "not done"
      },
      {
         "id": 2,
         "title": "Task 2",
         "description": "Description",
         "date": "2023-09-26T00:00:00Z",
         "status": "not done"
      }
  ]
}
```

#### 4. Get All Tasks With Pagination
* Request example:
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
* Response example:
```json
{
   "tasks": [
      {
         "id": 1,
         "title": "some task2",
         "description": "desc",
         "date": "2023-09-26T00:00:00Z",
         "status": "not done"
      },
      {
         "id": 2,
         "title": "some task3",
         "description": "desc",
         "date": "2023-09-26T00:00:00Z",
         "status": "not done"
      }
   ]
}
```

#### 5. Get All Tasks By Date
* Request example:
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
* Response example:
```json
{
   "tasks": [
      {
         "id": 1,
         "title": "some task1",
         "description": "desc",
         "date": "2023-09-26T00:00:00Z",
         "status": "not done"
      },
      {
         "id": 2,
         "title": "some task2",
         "description": "desc",
         "date": "2023-09-26T00:00:00Z",
         "status": "not done"
      },
      {
         "id": 3,
         "title": "some task3",
         "description": "desc",
         "date": "2023-09-26T00:00:00Z",
         "status": "not done"
      }
   ]
}
```

#### 6. Set Task Status
* Request example:
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
* Response example: *None*
> **Hint:** if the updating was successful, the server will return code 204 (NO CONTENT).

#### 7. Delete Task By ID
* Request example:
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
* Response example: *None*
> **Hint:** if the deleting was successful, the server will return code 204 (NO CONTENT).

#### 8. Delete All User Tasks
* Request example:
```shell
curl -X 'DELETE' \
  'http://localhost:8080/api/v1/agenda/delete_all' \
  -H 'accept: application/json' \
  -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTU1Nzk1MzYsInN1YiI6IjIifQ.AhqXvtQPHMwp1Pv5Y9m6xXnMITUlRDnGo8oVm5DRvLU'
```
* Response example: *None*
> **Hint:** if the deleting was successful, the server will return code 204 (NO CONTENT).

---

## Additional features
1. **Run tests**
```shell
make test
```
2. **Run the linter**
```shell
make lint
```
3. **Create migration files**
```shell
make migrate-create
```
4. **Migrations up / down**
```shell
make migrate-up
```
```shell
make migrate-down
```
5. **Stop all running containers**
```shell
make stop
```
