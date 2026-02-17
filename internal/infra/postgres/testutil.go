package postgres

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"

	"github.com/testcontainers/testcontainers-go"

	"github.com/pressly/goose/v3"

	tc_postgres "github.com/testcontainers/testcontainers-go/modules/postgres"
)

//go:embed migrations/*.sql
var EmbedMigrations embed.FS

func SetupContainer(ctx context.Context, t *testing.T) (*sql.DB, *tc_postgres.PostgresContainer) {
	dbName := os.Getenv("DBDATABASE")
	dbUser := os.Getenv("DBUSER")
	dbPassword := os.Getenv("DBPASS")
	dbSslmode := os.Getenv("DBSSLMODE")

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

	connectionString, connectionStringErr := postgresContainer.ConnectionString(ctx, fmt.Sprintf("sslmode=%s", dbSslmode))
	if connectionStringErr != nil {
		t.Fatalf("Cannot get connection string: %v", connectionStringErr.Error())
	}

	instance := SetupDatabase(ctx, t, connectionString)
	SetupMigrations(ctx, t, instance)

	return instance, postgresContainer
}

func SetupDatabase(ctx context.Context, t *testing.T, connectionString string) *sql.DB {
	instance, instanceErr := NewDatabaseConnection(connectionString)
	if instanceErr != nil {
		t.Fatalf("Cannot get database instance: %v", instanceErr.Error())
	}

	return instance
}

func SetupMigrations(ctx context.Context, t *testing.T, instance *sql.DB) {
	goose.SetBaseFS(EmbedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		panic(err)
	}

	if err := goose.Up(instance, "migrations"); err != nil {
		panic(err)
	}
}

func TerminateContainer(postgresContainer *tc_postgres.PostgresContainer) {
	if err := testcontainers.TerminateContainer(postgresContainer); err != nil {
		log.Printf("failed to terminate container: %s", err)
	}
}
