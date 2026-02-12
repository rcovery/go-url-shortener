package postgres_test

import (
	"context"
	"log"
	"testing"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/rcovery/go-url-shortener/internal/infra/testutil"
	"github.com/rcovery/go-url-shortener/shorturl"
	"github.com/rcovery/go-url-shortener/shorturl/postgres"

	"github.com/testcontainers/testcontainers-go"
)

func TestSelect(t *testing.T) {
	ctx := context.Background()

	instance, postgresContainer := testutil.SetupContainer(ctx, t)
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
