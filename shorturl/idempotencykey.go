package shorturl

import "github.com/google/uuid"

// IdempotencyKey This package defines what an IdempotencyKey will be
type IdempotencyKey string

func NewIdempotencyKey() (IdempotencyKey, error) {
	newuuid, err := uuid.NewV7()
	return IdempotencyKey(newuuid.String()), err
}
