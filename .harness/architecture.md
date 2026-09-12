# Architecture Decision: Store Manager — Go Layered REST API

## Status

Proposed

## Context

The system is a dropshipping sales management HTTP API with two bounded resources (products, sales), one MySQL persistence store, and per-layer test coverage requirements. The previous Node.js/Express implementation followed MSC (Model-Service-Controller). The Go rewrite must maintain equivalent layer separation using Vanilla Go std library.

Hard constraints:
- Vanilla Go only — no Gin, Echo, Fiber, Chi
- MySQL persistence via database/sql + go-sql-driver/mysql
- Prescribed directory structure for handlers, services, repository
- Per-layer test coverage targets (5%–60%)

## Decision

Select **Go Layered REST API (A2 specialization)** with three internal layers:

```text
HTTP Transport (handlers — net/http)
         ↓
Business Logic (services)
         ↓
Data Access (repository — database/sql)
```

### Layer responsibilities

| Layer | Directory | Responsibilities | Tools |
|---|---|---|---|
| Handlers | `internal/handler/` | HTTP request parsing, response writing, status codes, error mapping | net/http, encoding/json |
| Services | `internal/service/` | Business rules, validation, orchestration | Pure Go |
| Repository | `internal/repository/` | SQL queries, data access, parameterized statements | database/sql, go-sql-driver/mysql |
| Models | `internal/model/` | Domain structs, shared types | Go types |

### Why this architecture fits

| Requirement | Alignment |
|---|---|
| Vanilla Go constraint | Each layer uses only std library + go-sql-driver/mysql |
| Clean separation | Handlers own HTTP, services own rules, repository owns SQL |
| Testability | Repository mocked with go-sqlmock, services mocked at repository, handlers call services |
| Incremental coverage | Layers provide natural per-boundary coverage targets |
| Prescribed structure | `internal/` enforces visibility; handler/service/repository directories map to layers |
| MySQL persistence | repository layer encapsulates all SQL behind Go interfaces |

## Simplest Viable Architecture

A single `main.go` with inline handler functions, a service struct, and direct SQL calls would be simpler but would violate the per-layer coverage requirement and make isolation testing impossible.

## Alternatives Considered

### Simple Application (A1)

Would place HTTP, business logic, and persistence in the same package. Rejected because per-layer coverage targets cannot be measured independently and the project's architectural requirements explicitly mandate separation.

### Hexagonal (A4)

Ports/adapters pattern would require interfaces for the repository and handler layers. While the repository interface is a valid pattern, full hexagonal with inbound/outbound port abstractions adds ceremony without current benefit — only one transport (HTTP) and one database (MySQL) exist.

### Modular Monolith (A5)

Two modules (products, sales) could be justified by domain separation, but the simple CRUD nature and shared junction table (sales_products) make strong module boundaries artificial.

## Rejected Alternatives

| Alternative | Rejection reason |
|---|---|
| Simple Application (A1) | No layer isolation for coverage measurement |
| Hexagonal (A4) | Single transport + single database, no adapter interchangeability needed |
| Modular Monolith (A5) | Domain too small to justify bounded module boundaries |

## Consequences

### Positive

- Clear separation of concerns with Go idiomatic patterns
- Each layer testable independently via interface mocking
- Repository interface enables testability without real database
- Services remain pure Go with no HTTP or SQL awareness
- Handlers remain thin: parse request -> call service -> write response

### Negative

- Vanilla Go routing requires manual path parsing (no wildcards)
- More boilerplate than a framework — URL parameter extraction, JSON encoding/decoding, error-to-status-code mapping
- Middleware pattern requires manual wrapping (no express-like middleware chain)

### Operational

- Server startup via `net/http.Server` with route registration
- Environment configuration through os.Getenv

### Testing

- Repository tests: go-sqlmock for database mocking
- Service tests: mock repository interfaces, pure Go assertions
- Handler tests: net/http/httptest for request/response recording

## Architecture-driving requirements

1. Vanilla Go constraint (no Gin, Echo, Fiber, Chi)
2. MySQL raw SQL persistence
3. Per-layer test coverage targets (5%–60%)
4. Independent testability of each layer via mocking

## Evolution triggers

- Adding authentication -> consider middleware pattern with net/http Handler wrapper
- Multiple transport adapters -> evaluate Hexagonal (A4)
- Multiple bounded domains -> evaluate Modular Monolith (A5)
- Complex routing needs -> evaluate Go 1.22+ http.ServeMux pattern matching

## Risks

- Manual URL routing is error-prone — ensure rigorous handler tests
- Service layer must not become pass-through — enforce meaningful business logic
- SQL injection via raw queries — enforce parameterized statements