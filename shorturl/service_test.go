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

		duplicatedUrl, _ := service.Create(ctx, id2, idempotencyKey2, name, link)
		// if duplicatedErr == nil {
		// 	t.Errorf("created a duplicated URL %q", creationErr)
		// }
		if createdShorturl == duplicatedUrl {
			t.Errorf("created a duplicated URL %q", duplicatedUrl)
		}
	})
}
