# Project Skill: Go Testing Pattern

## Patterns

1. **Repository tests**: Use `go-sqlmock` to mock `database/sql` — no real database required
2. **Service tests**: Mock repository interfaces using testify mocks or hand-written stubs
3. **Handler tests**: Use `net/http/httptest.NewRecorder()` and `httptest.NewRequest()`
4. **Coverage measurement**: `go test -coverprofile=coverage.out ./internal/...` then `go tool cover -func=coverage.out`
5. **Test naming**: `TestFunctionName_Scenario` (e.g., `TestListProducts_Success`)
6. **Table-driven tests**: Use Go's slice-of-struct pattern for multiple cases

## Coverage targets per layer

| Layer | Minimum | Target |
|---|---|---|
| internal/handler/ | 5% → 60% incremental | One function per requirement |
| internal/service/ | 5% → 60% incremental | One function per requirement |
| internal/repository/ | 5% → 60% incremental | One function per requirement |

## File structure

```
internal/
  handler/
    product_handler_test.go
    sale_handler_test.go
  service/
    product_service_test.go
    sale_service_test.go
  repository/
    product_repository_test.go
    sale_repository_test.go
```