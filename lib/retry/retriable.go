package retry

type RetriableError struct {
	Err error
}

func Retriable(err error) error {
	return &RetriableError{
		Err: err,
	}
}

func (err *RetriableError) Error() string {
	return err.Err.Error()
}

func (err *RetriableError) Unwrap() error {
	return err.Err
}
