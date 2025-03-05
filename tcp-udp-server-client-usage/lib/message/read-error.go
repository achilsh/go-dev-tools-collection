package message

import "fmt"

// ReadError is the error returned in case of non-fatal parsing errors.
type ReadError struct {
	str string
}

func (e ReadError) Error() string {
	return e.str
}

func newError(format string, args ...interface{}) ReadError {
	return ReadError{
		str: fmt.Sprintf(format, args...),
	}
}

type EmptyPkgError struct {
	str string
}

func (e EmptyPkgError) Error() string {
	return e.str
}
func newEmptyPkgError(format string, args ...any) EmptyPkgError {
	return EmptyPkgError{
		str: fmt.Sprintf(format, args...),
	}
}

type IoTimeoutError struct {
	str string
}

func (e IoTimeoutError) Error() string {
	return e.str
}
func newIoTimeoutError(format string, args ...any) IoTimeoutError {
	return IoTimeoutError{
		str: fmt.Sprintf(format, args...),
	}
}
