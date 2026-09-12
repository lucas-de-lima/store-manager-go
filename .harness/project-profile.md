# Project Profile: Store Manager

## Identity

- **Product name:** Store Manager
- **Repository name:** store-manager
- **Short description:** Dropshipping sales management REST API built with Vanilla Go (net/http), MySQL via database/sql + go-sql-driver/mysql, following Handlers/Services/Repository layered architecture.
- **Problem statement:** Building a RESTful API to manage product and sales CRUD operations for a dropshipping business, requiring persistent storage, input validation, and layered architecture.
- **Legacy naming to ignore:** `sd-025-b-store-manager`, `tryber` references, Node.js/Express/`npm` conventions

## Domain

### Actors

- System administrator / API consumer (no authentication — open API)
- The API itself (consistent error handling, status codes, validation)

### Core concepts

- **Product** — sellable item with id and name
- **Sale** — a transaction containing one or more product/quantity line items
- **Sales_Product** — junction table linking sales to products with quantities (N:N)

### Entities

| Entity | Attributes | Persistence |
|---|---|---|
| Product | id (PK, auto), name | StoreManager.products |
| Sale | id (PK, auto), date (auto) | StoreManager.sales |
| SalesProduct | sale_id (FK), product_id (FK), quantity | StoreManager.sales_products |

### Business rules

- Product name must be at least 5 characters
- Product name is required
- Sale items require productId and quantity >= 1
- Sale productId must reference an existing product
- Cascading delete: product/sale deletion propagates to sales_products
- Product list ordered ascending by id
- Sales list ordered ascending by saleId, then productId

### Primary workflows

1. List all products / get product by id
2. Create product with name validation
3. Update product name by id
4. Delete product by id
5. Search products by name query param
6. Register sale with multiple line items
7. List all sales / get sale by id
8. Update sale line items
9. Delete sale by id
10. Unit test coverage for all layers (handlers, services, repository)

## Technical

### Application type

- Go (>=1.21) REST API
- Vanilla `net/http` — no framework (no Gin, Echo, Fiber, Chi)
- Handlers/Services/Repository layered architecture
- MySQL via `database/sql` + `go-sql-driver/mysql`

### Interfaces

- HTTP REST endpoints under `/products` and `/sales`
- JSON request/response bodies (encoding/json)
- Standard HTTP status codes (200, 201, 204, 400, 404, 422)

### Persistence

- MySQL database (`StoreManager`)
- Tables: `products`, `sales`, `sales_products`
- Connection via `database/sql` + `go-sql-driver/mysql` pool
- Raw SQL queries (parameterized — no ORM)

### External dependencies

| Dependency | Purpose |
|---|---|
| Go std library (`net/http`) | HTTP server, routing, handlers |
| Go std library (`database/sql`) | Database abstraction |
| Go std library (`encoding/json`) | JSON marshaling/unmarshaling |
| `github.com/go-sql-driver/mysql` | MySQL driver |
| Go std library (`testing`) | Unit testing |
| `github.com/stretchr/testify` | Test assertions (optional) |
| `github.com/DATA-DOG/go-sqlmock` | SQL mock for repository tests |

### Concurrency

- Go goroutines per HTTP request (net/http default)
- `database/sql` connection pool

### Operational requirements

- Docker Compose for local dev (golang + mysql containers)
- Environment variables via `.env` file
- Migration script: `go run cmd/migrate/main.go`
- Seed script: `go run cmd/seed/main.go`
- Build: `go build -o bin/store-manager ./cmd/server`

## Constraints

### Explicit

- Handlers/Repository/Services directory structure (`internal/handler/`, `internal/service/`, `internal/repository/`)
- Go formatting: `gofmt` / `go vet` compliance
- MySQL database prefix: `StoreManager.table_name`
- Vanilla Go only — no external HTTP frameworks

### Inferred

- RESTful URL conventions (`/resources` and `/resources/:id`)
- Consistent error body format: `{ "message": "<error>" }`
- Input validation before database access

### Non-goals

- Authentication/authorization
- User management
- Product categories or inventory tracking
- Pricing or payment processing
- Order fulfillment tracking

## Risks

| Risk | Mitigation |
|---|---|
| Raw SQL without ORM may lead to injection vulnerabilities | Use parameterized queries always (`database/sql` prepared statements) |
| Layers can become wrappers with no real abstraction | Design meaningful service logic, not pass-through |
| Test coverage targets increase gradually (5%->60%) | Incremental testing per requirement |
| Vanilla Go routing is path-based without wildcards | Manual URL path parsing with `strings.TrimPrefix` / manual parameter extraction |

## Open questions

None — project requirements are fully specified.