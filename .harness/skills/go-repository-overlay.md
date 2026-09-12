# Project Skill: Go Repository Layer (database/sql)

## Context

The Store Manager uses `database/sql` with `go-sql-driver/mysql` for MySQL access.

## Conventions

1. **Connection**: `sql.Open("mysql", dsn)` with connection pooling
2. **Parameterized queries**: Always use `?` placeholders, never string formatting
3. **Repository interface**: Each repository defines an interface for testability
4. **Scanning**: `rows.Scan(&dest)` mapping columns to struct fields
5. **Transactions**: `db.Begin()` for multi-table operations (sales + sales_products)
6. **LastInsertId**: `result.LastInsertId()` for auto-generated PKs

## Test pattern

```go
// go-sqlmock example
db, mock, err := sqlmock.New()
repo := NewProductRepository(db)
mock.ExpectQuery("SELECT \\* FROM StoreManager\\.products").
    WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
        AddRow(1, "Martelo de Thor"))
```

## Query prefix convention

All table references in SQL must use `StoreManager.` prefix:
- `SELECT * FROM StoreManager.products`
- `INSERT INTO StoreManager.products (name) VALUES (?)`