# Initial Repo Structure

**Issue:** [NemuCorp/demo-repo#2](https://github.com/NemuCorp/demo-repo/issues/2)

## Problem

The `demo-repo` project needs a well-defined initial folder structure for both the client (React with TypeScript) and server (Go with Gin). The server requires organized packages for errors, handlers, database access, and logging, with clear conventions for how each layer should be built and how they should grow over time.

## Goals

- Establish a clean monorepo layout with `client/` (React + TypeScript) and `server/` (Go + Gin + PostgreSQL).
- Define server-side package boundaries: `myerrors`, `handler`, `db`, `logger`, plus entrypoints `main.go` and `cmd.go`.
- Enforce a pattern of per-domain handler and database structs (e.g., `AuthHandler`, `CartHandler`) with methods bound to them.
- Use raw PostgreSQL queries (no ORM abstractions) with prepared statements initialized at startup.
- Hash passwords and session tokens; support multiple sessions per user.
- Keep the structure flat but scalable: when a package grows too large, add subdirectories without changing the top-level conventions.

## Approach

### Top-Level Layout

```
demo-repo/
  client/          # React + TypeScript frontend
  server/          # Go + Gin backend
    main.go        # init(), DB & logger init, router configuration
    cmd.go         # CLI commands: DB UP, DOWN, CLEAN, IMPORT, EXPORT
    myerrors/      # Custom sentinel errors (errors.New("product not found"))
    handler/       # HTTP handlers (auth.go, cart.go, product.go, helpers.go)
    db/            # Database access (migrations/, auth.go, cart.go, product.go)
    logger/        # Logger initialization and logging modes
```

### Server Packages

#### `myerrors/`
- Define all sentinel errors as package-level variables using `errors.New(...)`.
- Examples: `ErrProductNotFound`, `ErrUnauthorized`, `ErrCartEmpty`.

#### `handler/`
- Each domain gets its own file: `auth.go`, `cart.go`, `product.go`.
- Shared utilities go in `helpers.go`.
- Each domain defines a handler struct (e.g., `AuthHandler`, `CartHandler`) holding any needed dependencies (e.g., a DB connection or domain-specific DB struct).
- All handler methods are bound to these structs.
- When a domain file grows too large, convert it into a subdirectory (e.g., `handler/auth/`) and split logic across files within that directory.

#### `db/`
- Each domain gets its own file: `auth.go`, `cart.go`, `product.go`.
- Each domain defines a DB struct (e.g., `AuthDB`, `CartDB`) holding its prepared statements.
- On init, all prepared statements are created and stored in the struct.
- Queries use raw SQL via `database/sql` — no ORM, no query builder.
- Migrations live in `db/migrations/` as versioned SQL files.
- When a domain file grows too large, convert it into a subdirectory (e.g., `db/auth/`) and split logic across files within that directory.

#### `logger/`
- Expose an `Init()` function to configure the logger.
- Support multiple logging modes (e.g., development, production).

#### `main.go`
- `init()` function: initialize DB connection, create prepared statements, initialize logger.
- Configure Gin router, register handler routes.

#### `cmd.go`
- CLI commands for database lifecycle:
  - `UP` — run pending migrations
  - `DOWN` — rollback the last migration
  - `CLEAN` — drop all tables
  - `IMPORT` — import data from file
  - `EXPORT` — export data to file

### Security Conventions

- Passwords are hashed before storage (e.g., bcrypt).
- Session tokens are hashed before storage.
- A single user can hold multiple active sessions (no global uniqueness constraint on session entries).

### Prepared Statement Strategy

- During `init()` in `main.go` (or explicitly during DB init), all prepared statements are created and cached on their respective domain DB structs.
- Handlers receive the DB structs and call the prepared statements, never constructing ad-hoc query strings at runtime.

## Risks

| Risk | Mitigation |
|------|------------|
| Prepared statements tied to a single connection pool lifecycle | Re-prepare on reconnect if using a connection pool wrapper |
| Flat file structure becoming unwieldy as domains grow | Convention to split into subdirectories when a file exceeds a reasonable length |
| Raw SQL leading to duplication across domains | Keep queries scoped to their domain DB struct; share helpers only via `helpers.go` patterns |
| Session-hash collisions for multiple sessions per user | Include a unique session ID (e.g., UUID) alongside the hashed token |

## Validation

- [ ] `server/` directory exists with all specified packages and files.
- [ ] `main.go` compiles and starts a Gin server.
- [ ] `cmd.go` accepts subcommands and prints usage for each DB operation.
- [ ] Prepared statements are registered at startup and usable by domain DB structs.
- [ ] Auth flow hashes passwords on registration and verifies on login.
- [ ] Multiple session entries can be created for the same user.
- [ ] Migrations run via the `UP` command produce correct table schemas.
- [ ] No ORM or query-builder dependency is present in `go.mod`.
