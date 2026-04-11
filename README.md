# go-boilerplate

A production-ready Go backend boilerplate using Fiber, PostgreSQL, and clean architecture patterns.

## Features

- 🚀 **Fiber v2** — Fast, Express-inspired web framework
- 🐘 **PostgreSQL** with `pgx/v5` driver and connection pooling
- 🔨 **Squirrel** query builder with soft-delete support
- 🔐 **Dual auth modes** — Firebase Auth or JWT (configurable via env)
- 📝 **Structured logging** with Uber Zap + fiberzap middleware
- ✅ **Validation** with go-playground/validator
- 🗄️ **Database migrations** with golang-migrate
- 📖 **Swagger** documentation auto-generated with swaggo
- 🐳 **Docker** multi-stage build + docker-compose
- 🔄 **Live reload** with Air
- 🧰 **Cobra CLI** for server start, migrations, version info

## Tech Stack

| Concern | Library |
|---------|---------|
| Web framework | github.com/gofiber/fiber/v2 |
| DB driver | github.com/jackc/pgx/v5 |
| Query builder | github.com/Masterminds/squirrel |
| CLI | github.com/spf13/cobra |
| Logger | go.uber.org/zap |
| Migrations | github.com/golang-migrate/migrate/v4 |
| Validation | github.com/go-playground/validator/v10 |
| Auth (Firebase) | firebase.google.com/go/v4 |
| Auth (JWT) | github.com/golang-jwt/jwt/v5 |
| Swagger | github.com/gofiber/swagger + github.com/swaggo/swag |

## Prerequisites

- Go 1.23+
- Docker & Docker Compose
- PostgreSQL 16+ (or use docker-compose)
- `make` (optional but recommended)

## Quick Start

```bash
# 1. Clone and rename the module
git clone https://github.com/yourorg/go-boilerplate.git
cd go-boilerplate

# Optional: rename the module
find . -type f -name "*.go" | xargs sed -i 's|github.com/sule/go-boilerplate|github.com/yourorg/yourproject|g'
go mod edit -module github.com/yourorg/yourproject

# 2. Copy environment config
cp .env.example .env
# Edit .env with your values

# 3. Start dependencies
docker-compose up -d postgres

# 4. Run migrations
make migrate-up

# 5. Start the server
make run
# or with live reload:
make dev
```

## Project Structure

```
.
├── cmd/
│   └── app/
│       ├── db/           # DB migration commands
│       ├── root.go       # Cobra root command
│       ├── serve.go      # Server start command
│       └── version.go    # Version command
├── config/
│   └── config.go         # Typed config loaded from env
├── db/
│   └── migrations/       # SQL migration files
├── docs/                 # Auto-generated Swagger docs
├── internal/
│   ├── db/               # DB pool, transactions, query builder
│   ├── middleware/        # Auth, CORS, logging, common stack
│   ├── server/           # App wiring and route mounting
│   ├── user/             # User domain module
│   └── utils/            # Response helpers, validation, pagination
├── pkg/
│   ├── errr/             # Error types and DB error parsing
│   ├── logger/           # Zap logger initialization
│   ├── null/             # Nil pointer helpers
│   └── types/            # Shared types (Pagination, SortDirection)
├── version/              # Version variables (injected at build time)
├── .air.toml             # Live reload config
├── .env.example          # Environment variable template
├── Dockerfile            # Multi-stage Docker build
├── docker-compose.yaml   # Local dev services
├── main.go               # Entry point
└── Makefile              # Development commands
```

## How to Add a New Module

See [AGENTS.md](AGENTS.md) for detailed step-by-step instructions.

Quick summary:
1. Create `internal/<domain>/` directory
2. Add `entity.go`, `request.go`, `response.go`
3. Add `repository.go` with interface + implementation
4. Add `service.go` with interface + implementation
5. Add `handler.go` with HTTP handlers and `RegisterRoutes()`
6. Wire it in `internal/server/server.go`

## Configuration Reference

| Variable | Default | Description |
|----------|---------|-------------|
| `APP_ENV` | `development` | Environment (development/production) |
| `PORT` | `8080` | HTTP server port |
| `APP_NAME` | `go-boilerplate` | Application name |
| `TIMEZONE` | `UTC` | Log timestamp timezone |
| `DATABASE_URL` | — | PostgreSQL connection URL |
| `AUTH_MODE` | `firebase` | Auth mode: `firebase` or `jwt` |
| `FIREBASE_CREDENTIALS` | — | Path to Firebase service account JSON |
| `JWT_SECRET` | — | HMAC secret for JWT signing |
| `JWT_LIFETIME_MINUTES` | `60` | JWT token lifetime |

## Auth Modes

### Firebase Auth
Set `AUTH_MODE=firebase` and `FIREBASE_CREDENTIALS=/path/to/service-account.json`.

Protected routes expect `Authorization: Bearer <firebase-id-token>`.

Use `middleware.UID(c)` to get the Firebase UID in handlers.

### JWT Auth
Set `AUTH_MODE=jwt` and `JWT_SECRET=your-secret`.

Protected routes expect `Authorization: Bearer <jwt-token>`.

Use `middleware.UserID(c)` to get the subject claim in handlers.

## Development Commands

```bash
make build          # Build binary to bin/
make build-alpine   # Build static binary for Linux
make run            # Run server directly
make dev            # Run with live reload (air)
make swag           # Generate Swagger docs
make test           # Run tests
make lint           # Run golangci-lint
make migrate-up     # Apply all pending migrations
make migrate-down   # Roll back last migration
make tools          # Install dev tools (air, swag, golangci-lint)
```

## API Documentation

Swagger UI is available at [http://localhost:8080/swagger/](http://localhost:8080/swagger/) when the server is running.

Generate/update docs:
```bash
make swag
```

## Deployment

### Docker

```bash
docker build -t go-boilerplate .
docker run -p 8080:8080 --env-file .env go-boilerplate
```

### Docker Compose (full stack)

```bash
docker-compose up -d
```

Includes PostgreSQL, the app, and pgAdmin (port 5050).

### Binary

```bash
make build-alpine
./bin/boilerplate serve
```

### Environment Variables

All configuration is via environment variables (12-factor). See `.env.example` for the full list.
