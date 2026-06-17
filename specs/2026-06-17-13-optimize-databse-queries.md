# Optimize Database Queries

**Issue:** [NemuCorp/demo-repo#13](https://github.com/NemuCorp/demo-repo/issues/13)

## Problem

The current database layer uses raw SQL with prepared statements but lacks benchmarking, query-plan analysis, and targeted optimizations. Several queries and the schema have identifiable inefficiencies that will degrade performance as the dataset grows:

- **Missing indexes** on high-traffic columns: `sessions.user_id`, `sessions.expires_at`, and `cart_items.user_id` are queried on every relevant request but have no index.
- **No pagination** on `ListProducts` — the query returns every row and will become slow with a large product catalog.
- **Prepared-statement churn** — every statement is prepared on startup regardless of query frequency, consuming server-side memory on the PostgreSQL connection.
- **No benchmarking infrastructure** — there is no mechanism to measure query latency, count scans, or compare before/after `EXPLAIN` plans.

## Goals

- Add missing database indexes to reduce sequential scans on frequently queried columns.
- Add pagination (limit/offset) to the product-listing query and API endpoint.
- Introduce a lightweight benchmarking harness (Go tests or script) that runs `EXPLAIN ANALYZE` on each query and reports plan shape and timing.
- Review existing queries for N+1 patterns, missing `WHERE` clauses, or redundant JOINs.
- Ensure all changes maintain backward compatibility with the existing API contracts.

## Approach

### 1. Add Missing Indexes

Create a migration `002_indexes.sql` that adds:

- `CREATE INDEX idx_sessions_user_id ON sessions(user_id);`
- `CREATE INDEX idx_sessions_expires_at ON sessions(expires_at);`
- `CREATE INDEX idx_cart_items_user_id ON cart_items(user_id);`

The `session_hash` index already exists. The `email` column is `UNIQUE` so PostgreSQL creates an implicit index. No index needed on `products` beyond the PK for the current query patterns.

### 2. Paginate Product Listing

Modify `db/product.go`:
- Replace the static `listProducts` prepared statement with a parameterized query accepting `LIMIT` and `OFFSET`.
- Update `ListProducts()` signature to `ListProducts(limit, offset int) ([]Product, error)`.
- Fall back to a sensible default (e.g., limit=50, offset=0) when no pagination params are provided.

Modify `handler/product.go`:
- Accept optional `limit` and `offset` query parameters in `List`.
- Pass them through to the DB layer.

### 3. Benchmarking Harness

Add a `db/bench_test.go` file with Go benchmark functions that:
- Connect to a test PostgreSQL instance (configurable via env).
- Insert seed data (100+ products, users, sessions).
- Run `EXPLAIN ANALYZE` on each query and capture the plan.
- Report sequential scans, index usage, and execution time.
- This harness is run via `go test -bench . ./server/db/` and is not part of the production binary.

### 4. Query Review

Audit each prepared statement:
- **auth.go**: `getSession` already filters by `session_hash` (indexed) and `expires_at` (new index). No changes needed.
- **cart.go**: `addItem` and `updateItem` CTEs are efficient single-round-trip patterns. No changes needed.
- **product.go**: `ListProducts` gains pagination (see above). Other product queries are simple PK lookups; no changes needed.

### 5. Connection Pool Tuning

Set `db.SetMaxOpenConns(25)` and `db.SetMaxIdleConns(5)` in `db.Init()` to prevent connection exhaustion under load.

## Risks

| Risk | Mitigation |
|------|------------|
| Index creation locks tables on large datasets | Since the project is early-stage with minimal data, migration impact is negligible; use `CONCURRENTLY` if data volume is a concern |
| Pagination breaks existing API consumers | Make `limit`/`offset` optional query params with defaults; existing callers see no change |
| Benchmark harness depends on a running PostgreSQL | Document the required `TEST_DB_DSN` env var; skip benchmarks if not set |
| Prepared statements for paginated queries can't be reused with different limits | Create a fresh prepared statement per unique limit/offset pair, or fall back to a non-prepared parameterized query for the listing endpoint |

## Validation

- [ ] Migration `002_indexes.sql` exists and creates the three new indexes.
- [ ] `EXPLAIN` on `SELECT ... FROM sessions WHERE user_id = $1` shows an index scan instead of a sequential scan.
- [ ] `EXPLAIN` on `SELECT ... FROM cart_items WHERE user_id = $1` shows an index scan.
- [ ] `GET /api/products?limit=20&offset=0` returns at most 20 products.
- [ ] `GET /api/products` (no params) returns at most 50 products (default limit).
- [ ] `go test -bench . ./server/db/ -run ^$` runs benchmarks and prints query plans.
- [ ] All existing handler tests pass without modification (aside from the new optional params).
- [ ] Connection pool settings are configured in `db.Init()`.
