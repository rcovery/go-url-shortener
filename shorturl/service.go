package shorturl

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type Service struct {
	Repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		Repo: repo,
	}
}

func (s *Service) Create(ctx context.Context, id ID, idempotencyKey IdempotencyKey, name string, link string) (string, error) {
	urlFound, urlError := s.Repo.SelectByIdempotencyKey(ctx, idempotencyKey)
	if urlError != nil {
		return "", urlError
	}
	if urlFound != "" {
		return urlFound, nil
	}

	urlFound, urlError = s.Repo.SelectByName(ctx, name)
	if urlError != nil && !errors.Is(urlError, sql.ErrNoRows) {
		return "", urlError
	}
	if urlFound != "" {
		return "", fmt.Errorf("cannot create a new URL with %q", name)
	}

	insertedErr := s.Repo.Insert(ctx, id, name, link, idempotencyKey)
	if insertedErr != nil {
		return "", insertedErr
	}

	return link, nil
}
