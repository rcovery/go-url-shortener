package shorturl

import "context"

type Reader interface {
	SelectByName(ctx context.Context, name string) (string, error)
	SelectByIdempotencyKey(ctx context.Context, idempotencyKey IdempotencyKey) (string, error)
}

type Writer interface {
	Insert(ctx context.Context, id ID, name string, link string, idempotencyKey IdempotencyKey) error
}

type Repository interface {
	Reader
	Writer

	// Close(ctx context.Context) error
}
