package shorturl

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Create(ctx context.Context, id ID, idempotencyKey IdempotencyKey, name string, link string) (string, error) {
	urlFound, urlError := s.repo.SelectByIdempotencyKey(ctx, idempotencyKey)
	if urlError != nil && !errors.Is(urlError, sql.ErrNoRows) {
		return "", urlError
	}
	if urlFound != "" {
		return urlFound, nil
	}

	urlFound, urlError = s.repo.SelectByName(ctx, name)
	if urlError != nil && !errors.Is(urlError, sql.ErrNoRows) {
		return "", urlError
	}
	if urlFound != "" {
		return "", fmt.Errorf("cannot create a new URL with %q", name)
	}

	insertedErr := s.repo.Insert(ctx, id, name, link, idempotencyKey)
	if insertedErr != nil {
		return "", insertedErr
	}

	return link, nil
}

func (s *Service) Select(ctx context.Context, name string) (string, error) {
	urlFound, urlError := s.repo.SelectByName(ctx, name)
	if urlError != nil && !errors.Is(urlError, sql.ErrNoRows) {
		return "", urlError
	}
	if urlFound == "" {
		return "", fmt.Errorf("cannot retrieve this URL with %q", name)
	}

	return urlFound, nil
}
