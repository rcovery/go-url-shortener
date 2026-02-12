package testutil

import (
	"context"
	"database/sql"
	"testing"

	"github.com/testcontainers/testcontainers-go"

	"github.com/golang-migrate/migrate/v4"
	migrate_postgres "github.com/golang-migrate/migrate/v4/database/postgres"
	internal_postgres "github.com/rcovery/go-url-shortener/internal/infra/postgres"
	tc_postgres "github.com/testcontainers/testcontainers-go/modules/postgres"
)

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
