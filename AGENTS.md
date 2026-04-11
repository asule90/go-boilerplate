# AGENTS.md — AI Coding Guidelines

## Project Architecture Overview

This is a flat, module-based Go backend using Fiber. Each domain module lives in `internal/<domain>/` and contains:

- `entity.go` — database struct (maps to DB columns)
- `request.go` — input DTOs with validation tags
- `response.go` — output DTOs with JSON/example tags
- `repository.go` — `Repository` interface + `repository` implementation (data access)
- `service.go` — `Service` interface + `service` implementation (business logic)
- `handler.go` — `Handler` struct + HTTP methods + `RegisterRoutes()`

## How to Add a New Domain Module

### Step 1: Create the directory

```
mkdir internal/product
```

### Step 2: Create entity.go

```go
package product

import "time"

type Product struct {
    ID        string     `db:"id"         json:"id"`
    Name      string     `db:"name"       json:"name"`
    Price     float64    `db:"price"      json:"price"`
    CreatedAt time.Time  `db:"created_at" json:"created_at"`
    UpdatedAt time.Time  `db:"updated_at" json:"updated_at"`
    DeletedAt *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`
}
```

### Step 3: Create request.go

```go
package product

type CreateProductRequest struct {
    Name  string  `json:"name"  validate:"required,min=1,max=255"`
    Price float64 `json:"price" validate:"required,gt=0"`
}
```

### Step 4: Create response.go

```go
package product

type ProductResponse struct {
    ID    string  `json:"id"    example:"550e8400-e29b-41d4-a716-446655440000"`
    Name  string  `json:"name"  example:"Widget"`
    Price float64 `json:"price" example:"9.99"`
}

func ToResponse(p Product) ProductResponse {
    return ProductResponse{ID: p.ID, Name: p.Name, Price: p.Price}
}
```

### Step 5: Create repository.go

```go
package product

import (
    "context"
    "fmt"
    sq "github.com/Masterminds/squirrel"
    "github.com/sule/go-boilerplate/internal/db"
    "github.com/sule/go-boilerplate/pkg/errr"
)

type Repository interface {
    Create(ctx context.Context, p Product) (Product, error)
    GetByID(ctx context.Context, id string) (Product, error)
}

type repository struct{ qb *db.QueryBuilder }

func NewRepository(qb *db.QueryBuilder) Repository { return &repository{qb: qb} }

func (r *repository) GetByID(ctx context.Context, id string) (Product, error) {
    query := r.qb.BaseQuery("products", "id", "name", "price", "created_at", "updated_at", "deleted_at").
        Where(sq.Eq{"products.id": id})
    var p Product
    err := r.qb.ExecuteQueryRow(ctx, query).Scan(&p.ID, &p.Name, &p.Price, &p.CreatedAt, &p.UpdatedAt, &p.DeletedAt)
    if err != nil {
        return Product{}, fmt.Errorf("product.repo.GetByID: %w", errr.ParseDBError(err))
    }
    return p, nil
}
```

### Step 6: Create service.go, handler.go (follow user domain patterns)

### Step 7: Wire in server.go

In `internal/server/server.go`:

```go
productRepo := product.NewRepository(qb)
productSvc  := product.NewService(productRepo)
productHandler := product.NewHandler(productSvc, cfg, logger, fbAuth)
productHandler.RegisterRoutes(api)
```

## Naming Conventions

- **Packages**: lowercase, single word (e.g., `user`, `product`, `order`)
- **Interfaces**: `Repository`, `Service` (simple names, not `IRepository`)
- **Constructors**: `NewRepository`, `NewService`, `NewHandler`
- **Error wrapping**: `fmt.Errorf("domain.layer.Method: %w", err)`
- **DB columns**: snake_case in struct tags
- **JSON keys**: snake_case

## Error Handling Patterns

```go
// Repository: always wrap and parse DB errors
if err != nil {
    return User{}, fmt.Errorf("user.repo.Create: %w", errr.ParseDBError(err))
}

// Service: check for ErrNoRows explicitly
existing, err := s.repo.GetByFirebaseUID(ctx, uid)
if err != nil && !errors.Is(err, errr.ErrNoRows) {
    return UserResponse{}, fmt.Errorf("user.svc.Upsert: %w", err)
}

// Handler: use ParseErrHTTP for all service errors
resp, err := h.svc.GetByID(c.Context(), id)
if err != nil {
    return utils.ParseErrHTTP(c, err, nil, h.logger)
}

// Return custom status codes
return errr.New(http.StatusConflict, "email already in use")
return errr.NewF(http.StatusBadRequest, "invalid format: %s", field)
return errr.Wrap(http.StatusBadRequest, "invalid input", originalErr)
```

## Patterns to Use

✅ Interfaces for all repos and services  
✅ Manual DI wiring in server.go  
✅ `BaseQuery()` for soft-delete filtering  
✅ `ExecuteInsert/Update` with `RETURNING *`  
✅ `ParseDBError()` in every repository method  
✅ `ValidateStruct()` before service calls  
✅ Swagger annotations on every handler  
✅ `ParseErrHTTP()` in every handler  

## Patterns to Avoid

❌ Wire/Fx for dependency injection  
❌ Global variables for logger or DB connection  
❌ Returning raw `error` from handlers without `ParseErrHTTP`  
❌ Hard-coding SQL strings instead of using QueryBuilder  
❌ Skipping interface definitions for testability  
❌ Mutations of request structs inside service logic  

## Testing Guidelines

- Place tests in the same package as the file being tested (`_test.go` suffix)
- Use table-driven tests
- Mock `Repository` and `Service` interfaces for unit tests
- Integration tests should use a real DB (use `testcontainers-go` or a test database)
- Run: `go test ./...`
- Coverage: `go test -cover ./...`
