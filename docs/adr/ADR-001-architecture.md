# ADR-001: Flat Per-Module Architecture

## Status

Accepted

## Context

Need a Go backend architecture that is easy to navigate, testable, and scalable without over-engineering. The project requires clear separation between layers (data access, business logic, HTTP) while keeping related domain code co-located.

## Decision

Use a flat per-module structure where each domain (e.g., `user`, `product`) lives in `internal/<domain>/` and contains all its layers (entity, request, response, repository, service, handler) as separate files in the same package.

Manual DI wiring in `internal/server/server.go` — no Wire/Fx.

```
internal/
├── user/
│   ├── entity.go       # DB struct
│   ├── request.go      # Input DTOs
│   ├── response.go     # Output DTOs
│   ├── repository.go   # Repository interface + pgx implementation
│   ├── service.go      # Service interface + business logic
│   └── handler.go      # HTTP handlers + RegisterRoutes()
└── server/
    └── server.go       # Wires all domains, mounts routes
```

## Consequences

**Positive:**
- All code for a domain is in one place — easy to navigate
- No package import cycles
- Simple, explicit DI is readable and debuggable
- Interfaces on Repository and Service enable unit testing with mocks
- Adding a new domain follows a clear, repeatable pattern

**Negative:**
- Large domains may accumulate many files in one directory
- Manual DI wiring grows linearly with domains (acceptable for most projects)
- No automated enforcement of the pattern (rely on AGENTS.md guidelines)

## Related

- [RFC-001-go-boilerplate.md](../RFC-001-go-boilerplate.md)
