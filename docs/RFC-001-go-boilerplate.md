# RFC-001: Go Backend Boilerplate

## Problem Statement

Need a production-ready Go backend template with clean architecture patterns, dual auth support (Firebase/JWT), PostgreSQL with pgx/v5, and a complete developer toolchain.

## Proposed Approach

Flat per-module structure under `internal/` with manual DI wiring, interfaces for testability, and a shared QueryBuilder for type-safe SQL construction.

## Scope

**In:** Fiber web framework, pgx/v5, squirrel query builder, cobra CLI, zap logger, golang-migrate, validator, Firebase + JWT auth, Swagger, Docker.  
**Out:** Wire/Fx DI, gRPC, message queues, caching layer.

## Key Decisions

1. **Flat module structure** — each domain has entity/request/response/repository/service/handler in one directory
2. **Manual DI** — explicit wiring in server.go, no magic
3. **Dual auth** — AUTH_MODE env switches between Firebase and JWT at runtime
4. **QueryBuilder** — wraps squirrel with soft-delete, pagination, and transaction-aware execution
5. **StatusCodeError** — carries HTTP status codes through the error chain

## Todo

- [x] Core infrastructure (db, middleware, utils, pkg)
- [x] User domain module (full CRUD)
- [x] CLI commands (serve, version, db migrate)
- [x] Docker + docker-compose
- [x] Makefile + air config
- [x] AGENTS.md coding guidelines
