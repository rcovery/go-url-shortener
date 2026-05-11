package postgres

import (
	"context"
	"database/sql"

	"github.com/rcovery/go-url-shortener/shorturl"
	"github.com/rcovery/go-url-shortener/shorturl/errors/notcreated"
	"github.com/rcovery/go-url-shortener/shorturl/errors/notfound"
)

type Repository struct {
	DB *sql.DB
}

func NewRepository(DB *sql.DB) *Repository {
	return &Repository{
		DB,
	}
}

func (r *Repository) SelectByName(ctx context.Context, name string) (shorturl.SelectableShortURL, error) {
	row := r.DB.QueryRowContext(ctx, `
		SELECT id, link
		FROM shorturls
		WHERE name = $1
			AND expires_at > NOW()
		LIMIT 1
	`, name)

	var link string
	var id shorturl.ID

	scanErr := row.Scan(&id, &link)
	if scanErr != nil {
		return shorturl.SelectableShortURL{}, notfound.New("ByName: Cannot scan link")
	}

	return shorturl.SelectableShortURL{
		ID:   id,
		Link: link,
	}, nil
}

func (r *Repository) SelectByIdempotencyKey(ctx context.Context, idempotencyKey shorturl.IdempotencyKey) (shorturl.SelectableShortURL, error) {
	row := r.DB.QueryRowContext(ctx, `
		SELECT id, link
		FROM shorturls
		WHERE idempotency_key = $1
			AND expires_at > NOW()
		LIMIT 1
	`, idempotencyKey)

	var id shorturl.ID
	var link string

	scanErr := row.Scan(&id, &link)
	if scanErr != nil {
		return shorturl.SelectableShortURL{}, notfound.New("ByIdempotencyKey: Cannot scan link")
	}

	return shorturl.SelectableShortURL{
		ID:   id,
		Link: link,
	}, nil
}

func (r *Repository) Insert(ctx context.Context, id shorturl.ID, name string, link string, idempotencyKey shorturl.IdempotencyKey) error {
	_, insertionErr := r.DB.ExecContext(ctx, `
		INSERT INTO shorturls
		(id, name, link, idempotency_key)
		VALUES
		($1, $2, $3, $4)
	`, id, name, link, idempotencyKey,
	)

	if insertionErr != nil {
		notcreated.New(insertionErr.Error())
	}

	return nil
}
