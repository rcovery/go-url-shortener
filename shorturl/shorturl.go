package shorturl

type ShortURL struct {
	ID             ID
	Link           string
	Name           string
	IdempotencyKey IdempotencyKey
	ExpiresAt      string
}

type SelectableShortURL struct {
	ID   ID
	Link string
}

type NotFound struct {
	Code    int
	Message int
}
