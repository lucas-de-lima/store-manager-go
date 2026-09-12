## Vocabulary

| Term | Definition |
|---|---|
| Product | A sellable item in the catalog, with id and name fields |
| Sale | A transaction recording the sale of one or more products |
| SalesProduct | Junction record linking a sale to a product with quantity |
| Handler | Go HTTP handler function (func(w http.ResponseWriter, r *http.Request)) |
| Service | Business logic layer, encapsulates validation and orchestration |
| Repository | Data access layer, encapsulates SQL queries |
| Model | Shared Go struct type definitions |
| Vanilla Go | Go programming using only standard library packages |
| go-sql-driver/mysql | Pure Go MySQL driver implementing database/sql/driver |
| go-sqlmock | SQL mocking library for testing repository layer |
| httptest | Go standard library package for testing HTTP handlers |
| migration.sql | DDL script creating StoreManager database tables |
| seed.sql | DML script populating initial test data |

## Acronyms

| Acronym | Expansion |
|---|---|
| MSC | Model-Service-Controller (legacy Node.js pattern) |
| REST | Representational State Transfer |
| CRUD | Create, Read, Update, Delete |
| PK | Primary Key |
| FK | Foreign Key |
| N:N | Many-to-many relationship |

## SQL naming

- Database: `StoreManager`
- Table prefix in queries: `StoreManager.` (e.g., `StoreManager.products`)
- Tables: `products`, `sales`, `sales_products`