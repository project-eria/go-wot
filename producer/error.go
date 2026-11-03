package producer

import "errors"

type DataError struct {
	Message string
}

func (e *DataError) Error() string {
	return e.Message
}

type UnknownError struct {
	Message string
}

func (e *UnknownError) Error() string {
	return e.Message
}

type NotImplementedError struct {
	Message string
}

func (e *NotImplementedError) Error() string {
	return e.Message
}

// fromHandler keeps an error a handler meant for the client (DataError,
// NotImplementedError); any other is internal (SECURITY_AUDIT #11: its text
// goes to the log, never to the client).
func fromHandler(err error) error {
	var data *DataError
	var notImplemented *NotImplementedError
	if errors.As(err, &data) || errors.As(err, &notImplemented) {
		return err
	}
	return &UnknownError{Message: err.Error()}
}
