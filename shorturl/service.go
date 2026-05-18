package shorturl

import (
	"context"
	"errors"
	"fmt"

	"github.com/rcovery/go-url-shortener/shorturl/errs"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Create(ctx context.Context, id ID, idempotencyKey IdempotencyKey, name string, link *Link) (*Link, error) {
	urlFound, urlError := s.repo.SelectByIdempotencyKey(ctx, idempotencyKey)
	if urlError != nil && !errors.Is(urlError, errs.NotFoundError) {
		return nil, urlError
	}
	if urlFound.ID != "" {
		return urlFound.Link, nil
	}

	urlFound, urlError = s.repo.SelectByName(ctx, name)
	if urlError != nil && !errors.Is(urlError, errs.NotFoundError) {
		return nil, urlError
	}
	if urlFound.ID != "" {
		return nil, fmt.Errorf("cannot create a new URL with %q", name)
	}

	insertedErr := s.repo.Insert(ctx, id, name, link, idempotencyKey)
	if insertedErr != nil {
		return nil, insertedErr
	}

	return link, nil
}

func (s *Service) Select(ctx context.Context, name string) (*Link, error) {
	urlFound, urlError := s.repo.SelectByName(ctx, name)
	if urlError != nil && errors.Is(urlError, errs.NotFoundError) {
		return nil, urlError
	}
	if urlFound.ID == "" {
		return nil, fmt.Errorf("cannot retrieve this URL with %q", name)
	}

	return urlFound.Link, nil
}
