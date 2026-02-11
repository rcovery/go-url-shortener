package postgres_test

import (
	"context"
	"database/sql"
	"log"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	internal_postgres "github.com/rcovery/go-url-shortener/internal/infra/postgres"
	"github.com/rcovery/go-url-shortener/shorturl"
	"github.com/rcovery/go-url-shortener/shorturl/postgres"

	"github.com/testcontainers/testcontainers-go"
	tc_postgres "github.com/testcontainers/testcontainers-go/modules/postgres"

	migrate_postgres "github.com/golang-migrate/migrate/v4/database/postgres"
)

func TestSelect(t *testing.T) {
	ctx := context.Background()

	instance, postgresContainer := SetupContainer(ctx, t)
	defer func() {
		if err := testcontainers.TerminateContainer(postgresContainer); err != nil {
			log.Printf("failed to terminate container: %s", err)
		}
	}()

	repo := postgres.Repository{
		instance,
	}

	id, _ := shorturl.NewID()
	idempotencyKey, _ := shorturl.NewIdempotencyKey()
	name, link := "RCovery", "https://neocities.org"

	insertErr := repo.Insert(ctx, id, name, link, idempotencyKey)
	if insertErr != nil {
		t.Fatalf("There was an Insert Error %q", insertErr.Error())
	}

	foundShorturl, err := repo.SelectByName(ctx, name)
	if err != nil {
		t.Errorf("Cannot get URL by name, instead got %q", err)
	}
	if foundShorturl == "" {
		t.Errorf("Empty short url, maybe try to insert before")
	}

	if foundShorturl != link {
		t.Errorf("want %q, got %q", link, foundShorturl)
	}
}

func SetupContainer(ctx context.Context, t *testing.T) (*sql.DB, *tc_postgres.PostgresContainer) {
	dbName := "gourl"
	dbUser := "dev"
	dbPassword := "dev123"

	postgresContainer, containerErr := tc_postgres.Run(
		ctx, "postgres:latest",

		tc_postgres.BasicWaitStrategies(),
		tc_postgres.WithPassword(dbPassword),
		tc_postgres.WithUsername(dbUser),
		tc_postgres.WithDatabase(dbName),
	)

	if containerErr != nil {
		t.Fatalf("An error occurred when starting container: %v", containerErr.Error())
	}

	testcontainers.CleanupContainer(t, postgresContainer)

	connectionString, connectionStringErr := postgresContainer.ConnectionString(ctx, "sslmode=disable")
	if connectionStringErr != nil {
		t.Fatalf("Cannot get connection string: %v", connectionStringErr.Error())
	}

	instance := SetupDatabase(ctx, t, connectionString)
	SetupMigrations(ctx, t, instance)

	return instance, postgresContainer
}

func SetupDatabase(ctx context.Context, t *testing.T, connectionString string) *sql.DB {
	instance, instanceErr := internal_postgres.NewDatabaseConnection(connectionString)
	if instanceErr != nil {
		t.Fatalf("Cannot get database instance: %v", instanceErr.Error())
	}

	return instance
}

func SetupMigrations(ctx context.Context, t *testing.T, instance *sql.DB) {
	driver, migrateErr := migrate_postgres.WithInstance(instance, &migrate_postgres.Config{})
	if migrateErr != nil {
		t.Fatalf("Cannot get migration instance: %v", migrateErr.Error())
	}

	m, migrateErr := migrate.NewWithDatabaseInstance(
		"file://../../database/migrations",
		"postgres", driver)
	if migrateErr != nil {
		t.Fatalf("Cannot get migration database instance: %v", migrateErr.Error())
	}

	if err := m.Up(); err != nil {
		t.Fatalf("Cannot run migrations: %v", err.Error())
	}
}
