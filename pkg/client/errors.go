package client

import "fmt"

var (
	ErrBadStatus = fmt.Errorf("bad status code")
)

type APIError struct {
	StatusCode int
	Err        error
}

func (e *APIError) Error() string {
	return fmt.Sprintf("APIError: %s", e.Err.Error())
}

func (e *APIError) Unwrap() error {
	return e.Err
}

func (e *APIError) Is(target error) bool {
	t, ok := target.(*APIError)
	if !ok {
		return false
	}
	return t.StatusCode == e.StatusCode || e.Err == t.Err
}

func newError(statusCode int, err error) error {
	return &APIError{
		StatusCode: statusCode,
		Err:        err,
	}
}
