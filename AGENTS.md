# AGENTS.md - go-url-shortener

## Project Overview

Go URL shortener service using clean architecture. Go 1.25.1, standard library
`net/http` for HTTP, `database/sql` with `lib/pq` for PostgreSQL, and
`testcontainers-go` for integration tests.

Module path: `github.com/rcovery/go-url-shortener`

## Architecture

```
main.go                        # HTTP entrypoint, handler wiring
shorturl/                      # Domain layer
  shorturl.go                  # ShortURL domain struct
  id.go                        # ID type (UUIDv7 wrapper)
  idempotencykey.go            # IdempotencyKey type (UUIDv7 wrapper)
  repository.go                # Repository interface (Reader + Writer)
  service.go                   # Business logic
  service_test.go              # Integration tests
  postgres/
    repository.go              # PostgreSQL repository implementation
    repository_test.go         # Integration tests
internal/infra/postgres/
  connection.go                # DB connection helpers
  testutil.go                  # Testcontainers setup for tests
  migrations/*.sql             # SQL migration files (embedded via go:embed)
```

Domain types and interfaces live in `shorturl/`. Infrastructure adapters live in
`shorturl/postgres/` and `internal/infra/postgres/`. The `main.go` file is the
composition root.

## Build Commands

```bash
# Build
go build ./...

# Run the server (requires PostgreSQL via docker-compose)
docker compose up -d
go run main.go

# Format
gofmt -w .

# Vet
go vet ./...
```

## Test Commands

Integration tests require Docker to be running (testcontainers spins up
PostgreSQL containers automatically).

```bash
# Run all tests
go test ./...

# Run all tests (verbose)
go test -v ./...

# Run a single test function
go test ./shorturl/ -run TestCreate

# Run a single subtest
go test ./shorturl/ -run TestCreate/should_create_a_unique_shorturl

# Run tests in a specific package
go test ./shorturl/postgres/

# Run a single test in a specific package
go test ./shorturl/postgres/ -run TestSelect/Selecting_by_name

# Run all tests with coverage
go test ./... -coverprofile=coverage.out

# Run tests with race detector
go test -race ./...
```

## Database Migrations

Migrations are in `internal/infra/postgres/migrations/` and are embedded into
the test binary via `//go:embed`. For production, the Makefile uses `goose`:

```bash
# Requires .env with DBUSER, DBPASS, DBDATABASE
make migrateup
make migratedown
```

## Code Style Guidelines

### Formatting

- Use `gofmt` (standard Go formatting with tabs).
- No additional formatters or linters are configured; stick to `gofmt` and
  `go vet` defaults.

### Package Naming

- Lowercase, single-word package names: `shorturl`, `postgres`.
- No underscores or mixed case in package names.

### Import Organization

Group imports in this order, separated by blank lines:

1. Standard library
2. Third-party packages
3. Internal project packages

```go
import (
    "context"
    "database/sql"

    _ "github.com/lib/pq"

    "github.com/rcovery/go-url-shortener/shorturl"
)
```

- Use blank imports (`_`) for side-effect-only packages (drivers).
- Use snake_case aliases when needed to avoid conflicts:
  `infra_postgres "github.com/rcovery/go-url-shortener/internal/infra/postgres"`.

### Naming Conventions

- **Types**: PascalCase (`ShortURL`, `IdempotencyKey`, `Repository`, `Service`).
- **Exported functions**: PascalCase (`NewID`, `NewService`, `SetupContainer`).
- **Unexported functions**: camelCase (`getURLByName`, `createNewURL`).
- **Variables**: camelCase (`urlFound`, `connectionString`, `insertionErr`).
- **Constructor pattern**: `NewXxx` functions returning a pointer
  (`func NewService(repo Repository) *Service`).
- **Custom types**: String-based types with constructors
  (`type ID string` + `func NewID() (ID, error)`).

### Error Handling

- Return errors as the last return value.
- Check errors immediately with `if err != nil` and return early.
- Use `fmt.Errorf` for wrapping/creating error messages.
- Use `errors.Is()` for sentinel error comparisons (e.g., `sql.ErrNoRows`).
- In tests: use `t.Fatalf` for fatal setup errors, `t.Errorf` for assertion
  failures.
- Use `panic()` only in startup/initialization code (not in request paths).

### Testing Conventions

- Use external test packages: `package shorturl_test`, `package postgres_test`.
- Use table-driven subtests with `t.Run("descriptive name", ...)`.
- Subtest names are human-readable sentences:
  `"should create a unique shorturl"`, `"Selecting by name"`.
- Integration tests use testcontainers via
  `infra_postgres.SetupContainer(ctx, t)` which returns `(*sql.DB, *PostgresContainer)`.
- Always defer container termination:
  `defer infra_postgres.TerminateContainer(postgresContainer)`.
- Use the standard `testing` package only (no testify assertions).
- Assert with `t.Errorf("want %q, got %q", expected, actual)`.

### Interfaces

- Interfaces must abstract behaviors, not things. Name interfaces after what
  they do, not what they are (e.g., `Reader` and `Writer` instead of
  `Database`).
- Define interfaces in the domain package (`shorturl/repository.go`), not in
  the implementation package.
- Compose interfaces from smaller ones (`Repository` = `Reader` + `Writer`).
- All interface methods accept `context.Context` as the first parameter.

### SQL

- Use backtick raw strings for multi-line SQL.
- Use PostgreSQL-style positional placeholders (`$1`, `$2`).
- Indent SQL for readability.

### HTTP Handlers

- Use `http.HandleFunc` with Go 1.22+ path patterns (`"/{url_name}"`).
- Dispatch on method with `switch r.Method`.
- Use `log.Println` for logging (standard library `log` package).

### Dependencies

Key direct dependencies:
- `github.com/google/uuid` - UUIDv7 generation
- `github.com/lib/pq` - PostgreSQL driver
- `github.com/testcontainers/testcontainers-go` - Test containers
- `github.com/pressly/goose/v3` - Migration runner (used in test setup)
- `github.com/golang-migrate/migrate/v4` - Migration source driver

### Docker

- `docker-compose.yml` provides a local PostgreSQL instance (user: `dev`,
  password: `dev123`, database: `gourl`, port: `5432`).
- No Dockerfile for the Go application itself.
- Tests use testcontainers (no docker-compose needed for tests).

## Skills

This project has specialized skills defined in `.opencode/skills/`. When a task
matches one of the skills below, **always load and follow the skill instructions
instead of working from scratch**. Use the `skill` tool to load the appropriate
skill by name.

| Skill | Name | When to use |
| --- | --- | --- |
| Go Test | `go-test` | Writing, fixing, or improving integration and unit tests. Follows the project's red/green cycle, external test packages, table-driven subtests, and testcontainers conventions. |
| TODO Add | `todo-add` | Adding tasks, study topics, or notes to the project `TODO.md`. Handles categorization, duplicate checking, and proper formatting. |
| Git Commit | `git-commit` | Committing changes to git. Analyzes the working tree, decides whether to split or combine commits, uses conventional commit format, and asks for confirmation. |

**Rules:**
- If the user asks to write tests, load the `go-test` skill.
- If the user asks to add something to the TODO list or track a task/study
  topic, load the `todo-add` skill.
- If the user asks to commit changes, load the `git-commit` skill.
- When in doubt about whether a skill applies, load it -- the skill itself
  contains detailed instructions for deciding whether it is relevant.
