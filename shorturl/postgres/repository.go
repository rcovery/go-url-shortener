package postgres

import (
	"context"
	"database/sql"
	"log"

	"github.com/rcovery/go-url-shortener/shorturl"
)

type Repository struct {
	DB *sql.DB
}

func NewRepository(DB *sql.DB) *Repository {
	return &Repository{
		DB,
	}
}

func (r *Repository) SelectByName(ctx context.Context, name string) (string, error) {
	row := r.DB.QueryRowContext(ctx, `
		SELECT link
		FROM shorturls
		WHERE name = $1
			AND expires_at > NOW()
		LIMIT 1
	`, name)

	var link string

	scanErr := row.Scan(&link)
	if scanErr != nil {
		log.Println("ByName: Cannot scan link")
		return "", scanErr
	}

	return link, nil
}

func (r *Repository) SelectByIdempotencyKey(ctx context.Context, idempotencyKey shorturl.IdempotencyKey) (string, error) {
	row := r.DB.QueryRowContext(ctx, `
		SELECT link
		FROM shorturls
		WHERE idempotency_key = $1
			AND expires_at > NOW()
		LIMIT 1
	`, idempotencyKey)

	var link string
	scanErr := row.Scan(&link)
	if scanErr != nil {
		log.Println("ByIdempotencyKey: Cannot scan link")
		return "", scanErr
	}

	return link, nil
}

func (r *Repository) Insert(ctx context.Context, id shorturl.ID, name string, link string, idempotencyKey shorturl.IdempotencyKey) error {
	_, insertionErr := r.DB.ExecContext(ctx, `
		INSERT INTO shorturls
		(id, name, link, idempotency_key)
		VALUES
		($1, $2, $3, $4)
	`, id, name, link, idempotencyKey,
	)

	return insertionErr
}
