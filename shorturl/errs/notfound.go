package errs

import "errors"

type errNotFound struct {
	Message string
}

func (err errNotFound) Error() string {
	return err.Message
}

func (err errNotFound) New(msg string) errNotFound {
	err.Message = msg
	return err
}

func (err errNotFound) Is(target error) bool {
	return errors.As(target, &errNotFound{})
}

var NotFoundError = errNotFound{}
