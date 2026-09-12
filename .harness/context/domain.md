## Domain: Store Management — Dropshipping Sales

### Purpose

This project provides a RESTful API for managing products and sales in a dropshipping business context. The API supports CRUD operations for products and sales, with validation rules specific to dropshipping workflow requirements.

### Context boundaries

- The system operates without authentication — it is an open API
- No user/customer management — only product catalog and sales transactions
- No pricing, payment processing, or inventory tracking
- Sales record product-quantity pairs referencing validated products

### Key domain events

- Product created / updated / deleted
- Sale registered with line items
- Sale updated (line items replaced)
- Sale deleted

### Technical domain

- Go 1.21+ project binary: compiled, single-process
- MySQL persistence via database/sql connection pool
- Unit tests run against mocked SQL layer (no real database required)
- Integration tests may use real MySQL via Docker

### Architecture layers

| Layer | Package | Dependencies |
|---|---|---|
| Handler | internal/handler | net/http, encoding/json |
| Service | internal/service | Handler (called by), Repository (calls) |
| Repository | internal/repository | database/sql, go-sql-driver/mysql |
| Model | internal/model | None (shared types) |

### Error handling convention

All error responses follow `{"message": "<error text>"}` format. HTTP status codes:

| Condition | Status |
|---|---|
| Success (list/get) | 200 |
| Created | 201 |
| Deleted (no body) | 204 |
| Validation error | 400 |
| Resource not found | 404 |
| Business rule violation | 422 |