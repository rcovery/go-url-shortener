package notcreated

type NotCreatedError struct {
	Message string
}

func (err *NotCreatedError) Error() string {
	return err.Message
}

func New(msg string) *NotCreatedError {
	return &NotCreatedError{
		Message: msg,
	}
}
