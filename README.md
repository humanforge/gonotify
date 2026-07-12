# Notification Service

Go microservice built with Echo v5, GORM/Postgres, Zap logging.

## Prerequisites

- Go 1.26+
- PostgreSQL 16+
- golang-migrate CLI (for migrations)

## Quick Start

```bash
# Copy and edit environment config
cp .env.example .env

# Start Postgres and app via Docker Compose
docker compose -f deploy/docker-compose.yml up -d

# Or run locally:
go run ./cmd/notification-service/main.go
```

## Environment Variables

| Variable | Required | Default | Description |
|---|---|---|---|
| `ENV` | No | `local` | Runtime environment (`local`, `development`, `staging`, `production`) |
| `PORT` | No | `:8080` | HTTP listen address |
| `DATABASE_CONN_URL` | **Yes** | — | PostgreSQL connection string |
| `DB_MAX_IDLE_CONNS` | No | `10` | Max idle DB connections |
| `DB_MAX_OPEN_CONNS` | No | `100` | Max open DB connections |
| `DB_CONN_MAX_LIFETIME` | No | `1h` | Max connection lifetime |
| `DB_CONN_MAX_IDLE_TIME` | No | `30m` | Max idle connection time |
| `LOG_ROOT_PATH` | No | `./logs` | Log file directory |
| `LOG_LEVEL` | No | `info` | Log level (debug, info, warn, error) |
| `LOG_MAX_SIZE_MB` | No | `100` | Log rotation max size (MB) |
| `LOG_MAX_AGE_DAYS` | No | `30` | Log retention days |
| `LOG_MAX_BACKUPS` | No | `5` | Max rotated log files |
| `SHUTDOWN_TIMEOUT` | No | `15s` | Graceful shutdown timeout |
| `CORS_ALLOWED_ORIGINS` | No | `*` | Comma-separated CORS origins |

## API Endpoints

- `GET /healthz` — Liveness probe (always 200)
- `GET /readyz` — Readiness probe (pings DB, 503 if unhealthy)
- `GET /api/v1/healthz` — Versioned health check
- `GET /api/v1/readyz` — Versioned readiness check

## Makefile Commands

```bash
make build       # Build binary
make run         # Build and run
make test        # Run tests with race detection
make lint        # Run golangci-lint
make vet         # Run go vet
make clean       # Remove build artifacts
make docker-build # Build Docker image
```

## Project Structure

```
.
├── cmd/notification-service/main.go     # Entrypoint
├── internal/
│   ├── config/          # Typed config, validation
│   ├── logger/          # Zap logger setup
│   ├── database/        # GORM connection, pool, transactions
│   ├── server/          # Echo setup, routes, graceful shutdown
│   ├── middlewares/      # Request logging, recover, CORS, etc.
│   ├── handlers/        # HTTP handlers, error handler
│   ├── domain/          # Business logic interfaces
│   └── repository/      # DB implementations
├── migrations/          # SQL migrations
├── deploy/              # Dockerfile, docker-compose
└── .env.example
```

## Startup Order

1. Load and validate config
2. Initialize logger
3. Connect to database with connection pooling
4. Register middleware chain
5. Register routes
6. Start HTTP server
7. On SIGINT/SIGTERM, gracefully drain connections, close DB, flush logs
