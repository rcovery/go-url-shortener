package shorturl

type ShortURL struct {
	ID             ID
	Link           string
	Name           string
	IdempotencyKey IdempotencyKey
	ExpiresAt      string
}
