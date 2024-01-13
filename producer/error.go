package producer

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
