package shorturl

type ShortURL struct {
	ID             ID
	Link           *Link
	Name           string
	IdempotencyKey IdempotencyKey
	ExpiresAt      string
}

type SelectableShortURL struct {
	ID   ID
	Link *Link
}
