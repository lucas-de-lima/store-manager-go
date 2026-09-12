# Project Skill: Go HTTP Layer (Vanilla net/http)

## Context

The Store Manager uses Vanilla Go `net/http` for its REST API, not any third-party framework.

## Conventions

1. **Handler signature**: `func(w http.ResponseWriter, r *http.Request)`
2. **Route registration**: `http.HandleFunc(pattern, handler)` or `http.ServeMux`
3. **URL parameters**: Manual extraction via `strings.TrimPrefix(r.URL.Path, prefix)` and `strings.Split`
4. **Request body**: `json.NewDecoder(r.Body).Decode(&v)`
5. **Response encoding**: `json.NewEncoder(w).Encode(v)`
6. **Error responses**: `w.WriteHeader(status); json.NewEncoder(w).Encode(map[string]string{"message": msg})`
7. **Query parameters**: `r.URL.Query().Get("key")`

## Middleware pattern

```go
func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        log.Printf("%s %s", r.Method, r.URL.Path)
        next(w, r)
    }
}
```

## Layer boundaries

- Handlers must not contain business logic or SQL
- Handlers call service methods, receive results, write HTTP response
- Handlers map service errors to HTTP status codes