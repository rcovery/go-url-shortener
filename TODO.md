# TODO

> `main.go` is legacy and will be rewritten last. Issues specific to it are not
> tracked here — they will be resolved when it is replaced.

## Critical Bugs

- [x] Fix `GetConnectionFromEnv()` reading `DBNAME` env var — `.env` defines `DBDATABASE`
- [x] Fix `db.Query` used for INSERT in `shorturl/postgres/repository.go` — use `db.Exec`
- [x] Fix `*sql.Rows` never closed in `SelectByName` and `SelectByIdempotencyKey`

## Architecture

- [ ] Make `Service.Repo` unexported (`repo` instead of `Repo`)
- [ ] Remove `database/sql` import from `shorturl/service.go` — define a domain-level "not found" error
- [ ] Have repository return `*ShortURL` instead of raw `string`
- [ ] Use `db.QueryContext(ctx, ...)` instead of `db.Query(...)` to respect context cancellation
- [ ] Use the `ShortURL` domain struct (currently defined but never instantiated)

## Database

- [ ] Add UNIQUE constraint on `shorturls.name`
- [ ] Add UNIQUE constraint on `shorturls.idempotency_key`
- [ ] Add index on `shorturls.name`
- [ ] Add index on `shorturls.idempotency_key`
- [ ] Add index on `shorturls.expires_at` (or composite index on `name, expires_at`)
- [ ] Add NOT NULL constraint on `idempotency_key` column
- [ ] Change `link` column from VARCHAR(255) to TEXT
- [ ] Add trigger or application logic to update `updated_at` on modification
- [ ] Make URL expiration configurable (currently hardcoded to 1 day)
- [ ] Add cleanup mechanism for expired URLs

## Error Handling

- [ ] Add `db.Ping()` in `NewDatabaseConnection` to verify connectivity
- [ ] Configure connection pool (`SetMaxOpenConns`, `SetMaxIdleConns`, `SetConnMaxLifetime`)

## Security

- [ ] Add input validation on URL creation (valid URL format, name length/characters)
- [ ] Make `sslmode` configurable (currently hardcoded to `disable`)
- [ ] Replace real credentials in `.env.example` with placeholders

## Production Readiness (for new main.go)

- [ ] Add graceful shutdown (signal handling + `http.Server.Shutdown()`)
- [ ] Add health check endpoint (`/health` or `/readyz`)
- [ ] Add structured logging (use `log/slog`)
- [ ] Add request/access logging middleware
- [ ] Make server port and bind address configurable via env vars
- [ ] Return JSON response bodies on success and error
- [ ] Validate `Content-Type: application/json` on POST requests
- [ ] Limit request body size (`http.MaxBytesReader`)
- [ ] Add HTTP server timeouts (`ReadTimeout`, `WriteTimeout`, `IdleTimeout`)
- [ ] Add rate limiting
- [ ] Add CORS headers

## Tests

- [ ] Add service tests for error paths, idempotent creation, and invalid inputs
- [ ] Add repository tests for edge cases (non-existent name, expired URL, duplicates)
- [ ] Add uniqueness test for `IdempotencyKey` (like `ID` test already has)
- [ ] Add tests for `internal/infra/postgres/connection.go`
- [ ] Remove commented-out test code in `service_test.go`
- [ ] Check errors from `NewID()` / `NewIdempotencyKey()` in tests instead of discarding

## Build & Deployment

- [ ] Fix Makefile migration targets — references `database/migrations` but path is `internal/infra/postgres/migrations`
- [ ] Fix Makefile `testcoverage` target (missing output file argument)
- [ ] Add a Dockerfile for the Go application
- [ ] Pin PostgreSQL image version in `docker-compose.yml` (currently `postgres:latest`)
- [ ] Pin PostgreSQL image version in `testutil.go` (currently `postgres:latest`)
- [ ] Fix Docker healthcheck to check correct database (`gourl` not `postgres`)
- [ ] Remove leftover `golang-migrate/migrate/v4` dependency from `go.mod`
- [ ] Move `goose` from indirect to direct dependency in `go.mod`
- [ ] Add CI/CD pipeline configuration

## Research / Design Questions

- [ ] Evaluate creating a `Name` type with `NewName()` constructor for parsing/validation
- [ ] Evaluate creating a `Link` type with `NewLink()` constructor for parsing/validation

## Feature Ideas

- [ ] Implement a search/dashboard screen for looking up links and exporting reports

## Completed

- [x] Implement repository integration tests
- [x] Implement repository functions
