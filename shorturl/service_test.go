package shorturl_test

import (
	"context"
	"testing"

	infra_postgres "github.com/rcovery/go-url-shortener/internal/infra/postgres"
	"github.com/rcovery/go-url-shortener/shorturl"
	"github.com/rcovery/go-url-shortener/shorturl/postgres"
)

func TestCreate(t *testing.T) {
	t.Run("should create a unique shorturl", func(t *testing.T) {
		id, _ := shorturl.NewID()
		idempotencyKey, _ := shorturl.NewIdempotencyKey()
		name := "open-this-link-right-now"
		link := "https://google.com"

		ctx := context.Background()
		instance, postgresContainer := infra_postgres.SetupContainer(ctx, t)
		defer infra_postgres.TerminateContainer(postgresContainer)

		repo := postgres.NewRepository(instance)
		service := shorturl.NewService(repo)

		createdShorturl, creationErr := service.Create(ctx, id, idempotencyKey, name, link)
		if creationErr != nil {
			t.Errorf("cannot create a short URL %q", creationErr)
		}
		if createdShorturl == "" {
			t.Errorf("created URL is empty %q", createdShorturl)
		}

		id2, _ := shorturl.NewID()
		idempotencyKey2, _ := shorturl.NewIdempotencyKey()

		duplicatedURL, _ := service.Create(ctx, id2, idempotencyKey2, name, link)
		if createdShorturl == duplicatedURL {
			t.Errorf("created a duplicated URL %q", duplicatedURL)
		}
	})

	t.Run("should return error when name already exists", func(t *testing.T) {
		ctx := context.Background()
		instance, postgresContainer := infra_postgres.SetupContainer(ctx, t)
		defer infra_postgres.TerminateContainer(postgresContainer)

		repo := postgres.NewRepository(instance)
		service := shorturl.NewService(repo)

		id1, _ := shorturl.NewID()
		idempotencyKey1, _ := shorturl.NewIdempotencyKey()
		name := "taken-name"
		link := "https://example.com"

		_, firstErr := service.Create(ctx, id1, idempotencyKey1, name, link)
		if firstErr != nil {
			t.Fatalf("first Create failed unexpectedly: %v", firstErr)
		}

		id2, _ := shorturl.NewID()
		idempotencyKey2, _ := shorturl.NewIdempotencyKey()

		result, secondErr := service.Create(ctx, id2, idempotencyKey2, name, "https://other.com")
		if secondErr == nil {
			t.Errorf("expected an error when creating with duplicate name, got nil")
		}
		if result != "" {
			t.Errorf("want empty string, got %q", result)
		}
	})

	t.Run("should return existing link for same idempotency key", func(t *testing.T) {
		ctx := context.Background()
		instance, postgresContainer := infra_postgres.SetupContainer(ctx, t)
		defer infra_postgres.TerminateContainer(postgresContainer)

		repo := postgres.NewRepository(instance)
		service := shorturl.NewService(repo)

		id1, _ := shorturl.NewID()
		idempotencyKey, _ := shorturl.NewIdempotencyKey()
		name := "idempotent-link"
		link := "https://example.com/original"

		firstResult, firstErr := service.Create(ctx, id1, idempotencyKey, name, link)
		if firstErr != nil {
			t.Fatalf("first Create failed unexpectedly: %v", firstErr)
		}

		id2, _ := shorturl.NewID()

		secondResult, secondErr := service.Create(ctx, id2, idempotencyKey, name, link)
		if secondErr != nil {
			t.Errorf("expected no error for idempotent creation, got %v", secondErr)
		}
		if secondResult != firstResult {
			t.Errorf("want %q, got %q", firstResult, secondResult)
		}
	})

	t.Run("should get a link by his name", func(t *testing.T) {
		ctx := context.Background()
		instance, postgresContainer := infra_postgres.SetupContainer(ctx, t)
		defer infra_postgres.TerminateContainer(postgresContainer)

		repo := postgres.NewRepository(instance)
		service := shorturl.NewService(repo)

		id1, _ := shorturl.NewID()
		idempotencyKey, _ := shorturl.NewIdempotencyKey()
		name := "googlewebsitey2k"
		link := "https://google.com"

		createdLink, creationErr := service.Create(ctx, id1, idempotencyKey, name, link)
		if creationErr != nil {
			t.Fatalf("first Create failed unexpectedly: %v", creationErr)
		}

		selectedLink, selectErr := service.Select(ctx, name)
		if selectErr != nil {
			t.Errorf("expected no error for idempotent creation, got %v", selectErr)
		}
		if selectedLink != createdLink {
			t.Errorf("want %q, got %q", createdLink, selectedLink)
		}
	})
}
