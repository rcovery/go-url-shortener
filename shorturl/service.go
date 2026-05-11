package shorturl

import (
	"context"
	"errors"
	"fmt"

	"github.com/rcovery/go-url-shortener/shorturl/errors/notfound"
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
	_, ok := errors.AsType[*notfound.NotFound](urlError)

	if urlError != nil && ok {
		return "", urlError
	}
	if urlFound.ID != "" {
		return urlFound.Link, nil
	}

	urlFound, urlError = s.repo.SelectByName(ctx, name)
	_, ok = errors.AsType[*notfound.NotFound](urlError)

	if urlError != nil && ok {
		return "", urlError
	}
	if urlFound.ID != "" {
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
	_, ok := errors.AsType[*notfound.NotFound](urlError)

	if urlError != nil && ok {
		return "", urlError
	}
	if urlFound.ID == "" {
		return "", fmt.Errorf("cannot retrieve this URL with %q", name)
	}

	return urlFound.Link, nil
}
