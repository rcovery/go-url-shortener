---
name: go-test
description: Implement integration and unit tests in Go following the project's established patterns, conventions, and clean architecture
---

## What I do

Write integration and unit tests for Go code in this project, following the
established conventions: external test packages, table-driven subtests,
standard `testing` package only, testcontainers for database tests, and clean
architecture boundaries.

## When to use me

Use this skill when the user asks to:

- Write tests for new or existing code
- Add test coverage to a package
- Create integration tests that need a database
- Create unit tests for domain logic or value types
- Fix or improve existing tests
- Any variation of "test this", "add tests for...", "write tests for..."

## Architecture context

Before writing any test, understand the project's clean architecture:

```
shorturl/                      # Domain layer (interfaces + logic)
  shorturl.go                  # ShortURL domain struct
  id.go                        # ID type (UUIDv7 wrapper)
  idempotencykey.go            # IdempotencyKey type (UUIDv7 wrapper)
  repository.go                # Repository interface (Reader + Writer)
  service.go                   # Business logic (Service struct)
shorturl/postgres/             # Infrastructure adapter
  repository.go                # PostgreSQL Repository implementation
internal/infra/postgres/       # Test infrastructure
  testutil.go                  # SetupContainer, SetupDatabase, SetupMigrations, TerminateContainer
  connection.go                # DB connection helpers
  migrations/*.sql             # Embedded SQL migrations (goose)
main.go                        # HTTP composition root
```

Key interfaces defined in the domain layer (`shorturl/repository.go`):

```go
type Reader interface {
    SelectByName(ctx context.Context, name string) (string, error)
    SelectByIdempotencyKey(ctx context.Context, idempotencyKey IdempotencyKey) (string, error)
}

type Writer interface {
    Insert(ctx context.Context, id ID, name string, link string, idempotencyKey IdempotencyKey) error
}

type Repository interface {
    Reader
    Writer
}
```

Module path: `github.com/rcovery/go-url-shortener`

## Strategy: Red/Green

This project enforces the **red/green** testing strategy. Every test must be
seen failing before it is seen passing. This ensures tests actually validate
behavior and are not false positives.

The cycle for **each test** (or subtest) is:

1. **RED** -- Write the test first. The test must fail (either because the
   production code does not exist yet, returns wrong values, or the expected
   behavior is not implemented). Run the test and **confirm it fails**. If the
   test passes immediately, the test is not validating anything new -- revise
   the assertion or reconsider whether the test is needed.

2. **GREEN** -- Write or modify the minimal production code needed to make the
   failing test pass. Run the test again and **confirm it passes**.

Repeat the red/green cycle for every test or subtest before moving to the next
one. Do NOT batch-write multiple tests and then make them all pass at once --
go one at a time.

## Steps

1. **Read the code under test** -- Read the source file(s) being tested to
   understand the types, functions, methods, inputs, outputs, and error
   conditions. Identify every code path and edge case.

2. **Read existing tests** -- Read any existing `_test.go` files in the same
   package to understand what is already covered and avoid duplication. Match
   the style and naming conventions of existing tests.

3. **Classify the test type** -- Determine whether each test should be:
   - **Unit test**: Tests pure logic, value types, or domain functions that do
     not require external dependencies (database, network, filesystem).
   - **Integration test**: Tests code that interacts with PostgreSQL through
     the repository layer. Requires testcontainers.

4. **Write one test (RED)** -- Write a single test or subtest following all
   conventions described below. Create the test file if it does not exist, or
   add to the existing test file. Then run it:
   ```bash
   go test ./path/to/package/ -run TestFunctionName/subtest_name -v
   ```
   **The test MUST fail.** If it passes immediately, either:
   - The behavior is already implemented and the test adds no value -- skip it
     or write a more specific assertion that does fail.
   - The assertion is wrong (e.g., testing the wrong condition) -- fix it.

   Confirm the failure output before proceeding.

5. **Make it pass (GREEN)** -- Write or modify the minimal production code to
   make the failing test pass. Then run the test again:
   ```bash
   go test ./path/to/package/ -run TestFunctionName/subtest_name -v
   ```
   **The test MUST pass.** If it still fails, fix the production code (not the
   test) until it passes. Do NOT weaken assertions to force a pass.

6. **Repeat** -- Go back to step 4 for the next test or subtest. Continue the
   red/green cycle until all planned tests are written and passing.

7. **Run the full suite** -- After all individual red/green cycles are done,
   run the entire test suite to ensure nothing is broken:
   ```bash
   go test ./path/to/package/ -v
   ```

8. **Run vet** -- Ensure the code passes `go vet`:
   ```bash
   go vet ./path/to/package/
   ```

### When the code already exists

When writing tests for code that is already implemented (adding coverage to
existing functions), the red/green cycle still applies. To get a legitimate
RED phase:

- **Test an uncovered edge case** that the current implementation does not
  handle correctly (e.g., empty input, duplicate name, expired URL).
- **Test error conditions** that may not be wired up yet.
- **Temporarily break the production code** (e.g., comment out a return
  statement or change a value) to confirm the test actually catches the
  failure, then restore the code. Document this in the test run output.
- If after honest effort the test passes immediately because the code is
  already correct, that is acceptable -- note it and move on. The goal is to
  avoid blindly writing tests that could never fail.

## Test conventions

### File placement and package naming

- Test files live next to the source file they test: `foo.go` -> `foo_test.go`
- **Always use external test packages** to test only the exported API:
  - Source in `package shorturl` -> Test in `package shorturl_test`
  - Source in `package postgres` -> Test in `package postgres_test`
- This enforces that tests interact with the code the same way consumers do

### Import organization

Group imports in this order, separated by blank lines:

1. Standard library
2. Third-party packages
3. Internal project packages

```go
import (
    "context"
    "testing"

    _ "github.com/golang-migrate/migrate/v4/source/file"
    _ "github.com/lib/pq"
    infra_postgres "github.com/rcovery/go-url-shortener/internal/infra/postgres"
    "github.com/rcovery/go-url-shortener/shorturl"
    "github.com/rcovery/go-url-shortener/shorturl/postgres"
)
```

- Use blank imports (`_`) for side-effect-only packages (drivers)
- Use snake_case aliases for long import paths:
  `infra_postgres "github.com/rcovery/go-url-shortener/internal/infra/postgres"`

### Test function naming

- Top-level test functions: `TestXxx` where `Xxx` matches the function or
  method being tested (e.g., `TestCreate`, `TestSelect`, `TestNewID`)
- Subtest names are human-readable sentences describing the behavior:
  - `"should create a unique shorturl"`
  - `"should create a new ShortURL ID"`
  - `"Selecting by name"`
  - `"Selecting by idempotency_key"`

### Table-driven subtests

Use `t.Run("descriptive name", func(t *testing.T) { ... })` for subtests:

```go
func TestNewID(t *testing.T) {
    t.Run("should create a new ShortURL ID", func(t *testing.T) {
        ID, err := shorturl.NewID()
        if err != nil {
            t.Fatalf("NewID() %v", err)
        }

        if ID == "" {
            t.Errorf("Expected a ShortURL ID, received nothing")
        }
    })

    t.Run("should create an unique ID", func(t *testing.T) {
        ID1, err1 := shorturl.NewID()
        if err1 != nil {
            t.Fatalf("NewID() %v", err1)
        }
        // ...
    })
}
```

### Assertions

- Use the standard `testing` package only. **No testify or other assertion
  libraries.**
- Use `t.Fatalf(...)` for fatal setup errors that prevent the test from
  continuing (e.g., container startup failure, insert failure in setup)
- Use `t.Errorf(...)` for assertion failures where the test can continue
- Follow the pattern `"want %q, got %q"` for comparison failures:

```go
if foundShorturl != link {
    t.Errorf("want %q, got %q", link, foundShorturl)
}
```

### Error checking

- Always check errors immediately after the call
- Use `t.Fatalf` for errors in test setup/preconditions:

```go
insertErr := repo.Insert(ctx, id, name, link, idempotencyKey)
if insertErr != nil {
    t.Fatalf("There was an Insert Error %q", insertErr.Error())
}
```

- Use `t.Errorf` for errors in the actual assertion:

```go
createdShorturl, creationErr := service.Create(ctx, id, idempotencyKey, name, link)
if creationErr != nil {
    t.Errorf("cannot create a short URL %q", creationErr)
}
```

## Integration test patterns

### Testcontainers setup

Integration tests use testcontainers to spin up a real PostgreSQL instance.
Docker must be running. Use the helper from
`internal/infra/postgres/testutil.go`:

```go
func TestSomething(t *testing.T) {
    t.Run("descriptive behavior", func(t *testing.T) {
        ctx := context.Background()

        instance, postgresContainer := infra_postgres.SetupContainer(ctx, t)
        defer infra_postgres.TerminateContainer(postgresContainer)

        repo := postgres.NewRepository(instance)
        // ... test code ...
    })
}
```

**Key points:**
- `SetupContainer` creates a PostgreSQL container, connects to it, and runs
  all migrations automatically
- Always `defer infra_postgres.TerminateContainer(postgresContainer)` to clean
  up
- Each subtest that needs a database should create its own container for
  isolation
- The import alias is `infra_postgres` for
  `github.com/rcovery/go-url-shortener/internal/infra/postgres`

### Required blank imports for database tests

Repository tests in `shorturl/postgres/` need these blank imports for drivers:

```go
import (
    _ "github.com/golang-migrate/migrate/v4/source/file"
    _ "github.com/lib/pq"
)
```

Service-level integration tests in `shorturl/` only need the infra_postgres
import (drivers are loaded transitively).

### Repository test pattern

When testing repository methods, first insert test data, then query it:

```go
func TestSelect(t *testing.T) {
    t.Run("Selecting by name", func(t *testing.T) {
        ctx := context.Background()

        instance, postgresContainer := infra_postgres.SetupContainer(ctx, t)
        defer infra_postgres.TerminateContainer(postgresContainer)

        repo := postgres.NewRepository(instance)

        // Arrange: create test data
        id, _ := shorturl.NewID()
        idempotencyKey, _ := shorturl.NewIdempotencyKey()
        name, link := "RCovery", "https://neocities.org"

        // Arrange: insert precondition data
        insertErr := repo.Insert(ctx, id, name, link, idempotencyKey)
        if insertErr != nil {
            t.Fatalf("There was an Insert Error %q", insertErr.Error())
        }

        // Act
        foundShorturl, err := repo.SelectByName(ctx, name)

        // Assert
        if err != nil {
            t.Errorf("Cannot get URL by name, instead got %q", err)
        }
        if foundShorturl == "" {
            t.Errorf("Empty short url, maybe try to insert before")
        }
        if foundShorturl != link {
            t.Errorf("want %q, got %q", link, foundShorturl)
        }
    })
}
```

### Service integration test pattern

Service tests wire up the real repository and test business logic end-to-end:

```go
func TestCreate(t *testing.T) {
    t.Run("should create a unique shorturl", func(t *testing.T) {
        // Arrange
        id, _ := shorturl.NewID()
        idempotencyKey, _ := shorturl.NewIdempotencyKey()
        name := "open-this-link-right-now"
        link := "https://google.com"

        ctx := context.Background()
        instance, postgresContainer := infra_postgres.SetupContainer(ctx, t)
        defer infra_postgres.TerminateContainer(postgresContainer)

        repo := postgres.NewRepository(instance)
        service := shorturl.NewService(repo)

        // Act
        createdShorturl, creationErr := service.Create(ctx, id, idempotencyKey, name, link)

        // Assert
        if creationErr != nil {
            t.Errorf("cannot create a short URL %q", creationErr)
        }
        if createdShorturl == "" {
            t.Errorf("created URL is empty %q", createdShorturl)
        }
    })
}
```

## Unit test patterns

### Value type tests (no database needed)

For domain value types like `ID` and `IdempotencyKey`, test creation and
uniqueness without any infrastructure:

```go
package shorturl_test

import (
    "testing"

    "github.com/rcovery/go-url-shortener/shorturl"
)

func TestNewID(t *testing.T) {
    t.Run("should create a new ShortURL ID", func(t *testing.T) {
        ID, err := shorturl.NewID()
        if err != nil {
            t.Fatalf("NewID() %v", err)
        }

        if ID == "" {
            t.Errorf("Expected a ShortURL ID, received nothing")
        }
    })

    t.Run("should create an unique ID", func(t *testing.T) {
        ID1, err1 := shorturl.NewID()
        if err1 != nil {
            t.Fatalf("NewID() %v", err1)
        }

        ID2, err2 := shorturl.NewID()
        if err2 != nil {
            t.Fatalf("NewID() %v", err2)
        }

        if ID1 == ID2 {
            t.Errorf("The IDs are equal! %s / %s", ID1, ID2)
        }
    })
}
```

## Running tests

```bash
# Run all tests (requires Docker for integration tests)
go test ./...

# Run all tests verbose
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

## Important rules

- ALWAYS follow the red/green cycle: see the test fail (RED) before making it
  pass (GREEN). Never skip the RED phase.
- ALWAYS work one test at a time through the red/green cycle. Do NOT
  batch-write multiple tests and then make them all pass at once.
- ALWAYS use external test packages (`package foo_test`, not `package foo`)
- ALWAYS use the standard `testing` package -- no testify, no gomock, no other
  assertion libraries
- ALWAYS use `t.Run(...)` subtests with descriptive human-readable names
- ALWAYS use `t.Fatalf` for setup failures, `t.Errorf` for assertion failures
- ALWAYS defer container termination in integration tests
- ALWAYS create a fresh container per subtest for isolation
- ALWAYS check errors immediately after the call that returns them
- ALWAYS verify tests compile and pass before finishing
- NEVER skip the RED phase -- if a test passes immediately without any code
  change, investigate whether the test is actually asserting the right thing
- NEVER weaken an assertion to force a GREEN phase -- fix the production code
  instead
- NEVER use `t.Parallel()` unless the user explicitly requests it
- NEVER import test utilities from application packages -- use `internal/infra/postgres/testutil.go`
- NEVER use mocks unless the user explicitly requests them -- prefer real
  database integration tests via testcontainers
- NEVER add third-party test dependencies (testify, gomock, etc.)
- Format with `gofmt` conventions (tabs, standard Go formatting)
