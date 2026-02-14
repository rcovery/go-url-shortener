package shorturl_test

import (
	"testing"

	"github.com/rcovery/go-url-shortener/shorturl"
)

func TestNewIdempotencyKey(t *testing.T) {
	t.Run("should create a new IdempotencyKey", func(t *testing.T) {
		IdempotencyKey, err := shorturl.NewIdempotencyKey()
		if err != nil {
			t.Fatalf("NewID() %v", err)
		}

		if IdempotencyKey == "" {
			t.Errorf("Expected a IdempotencyKey, received nothing")
		}
	})
}
