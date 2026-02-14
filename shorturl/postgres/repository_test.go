package postgres_test

import (
	"context"
	"testing"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	infra_postgres "github.com/rcovery/go-url-shortener/internal/infra/postgres"
	"github.com/rcovery/go-url-shortener/shorturl"
	"github.com/rcovery/go-url-shortener/shorturl/postgres"
)

func TestSelect(t *testing.T) {
	t.Run("Selecting by name", func(t *testing.T) {
		ctx := context.Background()

		instance, postgresContainer := infra_postgres.SetupContainer(ctx, t)
		defer infra_postgres.TerminateContainer(postgresContainer)

		repo := postgres.NewRepository(instance)

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
	})

	t.Run("Selecting by idempotency_key", func(t *testing.T) {
		ctx := context.Background()

		instance, postgresContainer := infra_postgres.SetupContainer(ctx, t)
		defer infra_postgres.TerminateContainer(postgresContainer)

		repo := postgres.NewRepository(instance)

		id, _ := shorturl.NewID()
		idempotencyKey, _ := shorturl.NewIdempotencyKey()
		name, link := "RCovery", "https://neocities.org"

		insertErr := repo.Insert(ctx, id, name, link, idempotencyKey)
		if insertErr != nil {
			t.Fatalf("There was an Insert Error %q", insertErr.Error())
		}

		foundShorturl, err := repo.SelectByIdempotencyKey(ctx, idempotencyKey)
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
