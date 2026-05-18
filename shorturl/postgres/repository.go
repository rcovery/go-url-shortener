package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/rcovery/go-url-shortener/shorturl"
	"github.com/rcovery/go-url-shortener/shorturl/errs"
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

	var rawDBLink string
	var surl shorturl.SelectableShortURL

	scanErr := row.Scan(&surl.ID, &rawDBLink)
	if scanErr != nil {
		return surl, errs.NotFoundError.New(fmt.Sprintf("ByName: %v", scanErr))
	}

	link, linkErr := shorturl.NewLink(rawDBLink)
	if linkErr != nil {
		return surl, errs.NotFoundError.New(fmt.Sprintf("ByName: %v", scanErr))
	}

	surl.Link = link

	return surl, nil
}

func (r *Repository) SelectByIdempotencyKey(ctx context.Context, idempotencyKey shorturl.IdempotencyKey) (shorturl.SelectableShortURL, error) {
	row := r.DB.QueryRowContext(ctx, `
		SELECT id, link
		FROM shorturls
		WHERE idempotency_key = $1
			AND expires_at > NOW()
		LIMIT 1
	`, idempotencyKey)

	var rawDBLink string
	var surl shorturl.SelectableShortURL

	scanErr := row.Scan(&surl.ID, &rawDBLink)
	if scanErr != nil {
		return surl, errs.NotFoundError.New(fmt.Sprintf("ByIdempotencyKey: %v", scanErr))
	}

	link, linkErr := shorturl.NewLink(rawDBLink)
	if linkErr != nil {
		return surl, errs.NotFoundError.New(fmt.Sprintf("ByName: %v", scanErr))
	}

	surl.Link = link

	return surl, nil
}

func (r *Repository) Insert(ctx context.Context, id shorturl.ID, name string, link *shorturl.Link, idempotencyKey shorturl.IdempotencyKey) error {
	_, insertionErr := r.DB.ExecContext(ctx, `
		INSERT INTO shorturls
		(id, name, link, idempotency_key)
		VALUES
		($1, $2, $3, $4)
	`, id, name, link.String(), idempotencyKey,
	)

	if insertionErr != nil {
		return errs.NotCreatedErr.New(insertionErr.Error())
	}

	return nil
}
