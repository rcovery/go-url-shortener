package shorturl

import "github.com/google/uuid"

type ID string

func NewID() (ID, error) {
	newuuid, err := uuid.NewV7()
	return ID(newuuid.String()), err
}
