package shorturl

import "context"

type Reader interface {
	Select(ctx context.Context, id ID) (ID, error)
}

type Writer interface {
	Insert(ctx context.Context, shortURL ShortURL) (ID, error)
	Update(ctx context.Context, shortURL ShortURL) error
	Delete(ctx context.Context, id ID) error
}

type Repository interface {
	Reader
	Writer

	Close(ctx context.Context) error
}
