# DevOps Interview – HTTP App

A simple HTTP application with a `/health` endpoint that uses **PostgreSQL** as the primary database and **Redis** as a caching layer.

## Requirements

- Go 1.21+
- PostgreSQL (e.g. running on `localhost:5432`)
- Redis (e.g. running on `localhost:6379`)

## Configuration

| Variable       | Default                                      | Description                |
|----------------|----------------------------------------------|----------------------------|
| `HTTP_ADDR`    | `:8080`                                     | HTTP server listen address |
| `DATABASE_URL` | `postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable` | PostgreSQL connection URL |
| `REDIS_ADDR`   | `localhost:6379`                            | Redis server address       |

## Run

```bash
go run .
```

Or with custom config:

```bash
export DATABASE_URL="postgres://user:pass@localhost:5432/mydb?sslmode=disable"
export REDIS_ADDR="localhost:6379"
go run .
```

## Health endpoint

**GET /health**

Returns JSON with overall status and the status of PostgreSQL and Redis:

- **200** – All dependencies healthy (`"status": "ok"`)
- **503** – One or more dependencies unhealthy (`"status": "degraded"`)

Example response (healthy):

```json
{
  "status": "ok",
  "postgres": "healthy",
  "redis": "healthy"
}
```

Example response (degraded):

```json
{
  "status": "degraded",
  "postgres": "healthy",
  "redis": "unhealthy",
  "details": {
    "redis": "connection refused"
  }
}
```

## Quick start with Docker

```bash
# Start PostgreSQL and Redis
docker run -d --name postgres -p 5432:5432 -e POSTGRES_PASSWORD=postgres postgres:16-alpine
docker run -d --name redis -p 6379:6379 redis:7-alpine

# Run the app
go run .
```

Then: `curl http://localhost:8080/health`
