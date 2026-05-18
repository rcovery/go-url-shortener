package errs

type errNotCreated struct {
	Message string
}

func (err *errNotCreated) Error() string {
	return err.Message
}

func (err *errNotCreated) New(msg string) *errNotCreated {
	err.Message = msg
	return err
}

var NotCreatedErr = errNotCreated{}
