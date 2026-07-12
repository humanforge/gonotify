# Notification Service

Go microservice built with Echo v5, GORM/Postgres, Zap logging.

## Prerequisites

- Go 1.26+
- PostgreSQL 16+
- golang-migrate CLI (for migrations)

## Quick Start

```bash
# Copy and edit environment config (or use the existing .env)
cp .env.example .env

# Start Postgres and app via Docker Compose
docker compose -f deploy/docker-compose.yml up -d

# Or run locally:
go run ./cmd/server/main.go
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
├── cmd/server/main.go               # Entrypoint with explicit dependency wiring
├── internal/
│   ├── notification/                # Notification feature
│   │   ├── types.go                 # Domain struct, request/response DTOs
│   │   ├── store.go                 # Persistence layer (GORM)
│   │   ├── service.go               # Business logic
│   │   └── handler.go               # HTTP handlers
│   ├── platform/
│   │   ├── config/                  # Typed config, validation
│   │   ├── httpserver/              # Echo setup, middleware, health checks
│   │   ├── postgres/                # DB connection, pool, transactions
│   │   └── logging/                 # Zap logger setup
│   └── apperr/                      # Sentinels and typed error codes
├── migrations/                      # SQL migrations
├── deploy/                          # Dockerfile, docker-compose
└── .env
```

## Startup Order

1. Load and validate config
2. Initialize logger
3. Connect to database with connection pooling
4. Construct feature layer: `NewStore` → `NewService` → `NewHandler`
5. Register middleware chain
6. Register routes (feature handlers attached via `RegisterRoutes`)
7. Start HTTP server
8. On SIGINT/SIGTERM, gracefully drain connections, close DB, flush logs
