package notfound

type NotFound struct {
	Message string
}

func (err *NotFound) Error() string {
	return err.Message
}

func New(msg string) *NotFound {
	return &NotFound{
		Message: msg,
	}
}
