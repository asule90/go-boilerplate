package errr

import "fmt"

// StatusCodeError is an error that carries an HTTP status code.
type StatusCodeError struct {
	Code    int
	Message string
	Err     error
}

func (e *StatusCodeError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s", e.Message, e.Err.Error())
	}
	return e.Message
}

func (e *StatusCodeError) Unwrap() error {
	return e.Err
}

// New creates a StatusCodeError with the given HTTP status code and message.
func New(code int, message string) *StatusCodeError {
	return &StatusCodeError{Code: code, Message: message}
}

// NewF creates a StatusCodeError with a formatted message.
func NewF(code int, format string, args ...interface{}) *StatusCodeError {
	return &StatusCodeError{Code: code, Message: fmt.Sprintf(format, args...)}
}

// Wrap creates a StatusCodeError wrapping an existing error.
func Wrap(code int, message string, err error) *StatusCodeError {
	return &StatusCodeError{Code: code, Message: message, Err: err}
}

// WrapF creates a StatusCodeError wrapping an existing error with a formatted message.
func WrapF(code int, err error, format string, args ...interface{}) *StatusCodeError {
	return &StatusCodeError{Code: code, Message: fmt.Sprintf(format, args...), Err: err}
}
