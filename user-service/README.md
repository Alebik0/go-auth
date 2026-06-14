# User Service API

A simple REST API for user management built with Go and Gin.

## Prerequisites

* Go 1.25.5
* Git

## Build Setup

### Generate Swagger Files

```bash
go install github.com/swaggo/swag/cmd/swag@latest
swag init
```

### Build Executable

```bash
go build -o user-service .
```

### Run Application

```bash
./user-service
```

Or run directly:

```bash
go run .
```

## Base Path

```text
/api/v1
```

## Endpoints

### Create User

```http
POST /api/v1/user
```

Example request:

```json
{
  "name": "John Doe",
  "description": "Just a random description"
}
```

---

### Get User

```http
GET /api/v1/user/{id}
```

---

### Update User

```http
PUT /api/v1/user/{id}
```

Example request:

```json
{
  "name": "Jane Doe",
  "description": "New description"
}
```

---

### Delete User

```http
DELETE /api/v1/user/{id}
```

## Route Summary

| Method | Endpoint            |
| ------ | ------------------- |
| POST   | `/api/v1/user`      |
| GET    | `/api/v1/user/{id}` |
| PUT    | `/api/v1/user/{id}` |
| DELETE | `/api/v1/user/{id}` |

## Environment variables

* `POSTGRES_HOST` - database host
* `POSTGRES_PORT` - database port
* `POSTGRES_USER` - database user
* `POSTGRES_PASSWORD` - database user password
* `POSTGRES_DB` - database name
* `USER_SERVICE_PORT` - service running port
