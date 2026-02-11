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

func (r *Repository) SelectByName(ctx context.Context, name string) (string, error) {
	row, err := r.DB.Query(`
		SELECT link
		FROM shorturls
		WHERE name = $1
			AND expires_at > NOW()
		LIMIT 1
	`, name)
	if err != nil {
		log.Println("Error in the first if")
		return "", err
	}

	var link string

	// defer row.Close()
	row.Next()
	err = row.Scan(&link)
	if err != nil {
		log.Println("Error in the second if")
		return "", err
	}

	return link, nil
}

func (r *Repository) Insert(ctx context.Context, id shorturl.ID, name string, link string, idempotencyKey shorturl.IdempotencyKey) error {
	_, insertionErr := r.DB.Query(`
		INSERT INTO shorturls
		(id, name, link, idempotency_key)
		VALUES
		($1, $2, $3, $4)
	`, id, name, link, idempotencyKey,
	)

	return insertionErr
}
