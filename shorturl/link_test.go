package shorturl_test

import (
	"testing"

	"github.com/rcovery/go-url-shortener/shorturl"
	"github.com/stretchr/testify/assert"
)

func TestNewLink(t *testing.T) {
	t.Run("should create a new valid ShortURL Link", func(t *testing.T) {
		rawURL := "https://google.com"
		link, err := shorturl.NewLink(rawURL)

		assert.NoError(t, err)
		// if err != nil {
		// 	t.Fatalf("NewID() %v", err)
		// }

		assert.Equal(t, rawURL, link.String())
	})
	t.Run("should not instantiate HTTP link", func(t *testing.T) {
		rawURL := "http://google.com"
		link, err := shorturl.NewLink(rawURL)

		assert.Error(t, err, "We cannot use HTTP links")
		assert.Nil(t, link)
	})
}
