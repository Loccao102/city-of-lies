# ADR 0005: pgx + sqlc Over Heavy ORM

### Status
Accepted

### Context
Heavy ORMs often obscure query performance, generate suboptimal queries for complex relational joins (such as belief graphs and world event streams), and lead to N+1 performance pitfalls.

### Decision
Use native PostgreSQL 18 with `pgx/v5` connection pooling and `sqlc` for compile-time generated type-safe Go code from hand-crafted SQL queries.

### Consequences
- Explicit, auditable SQL schemas and queries.
- Zero runtime reflection overhead.
- Strict compile-time type checking between SQL schema and Go types.
