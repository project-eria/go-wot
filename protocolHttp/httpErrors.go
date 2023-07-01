package protocolHttp

import "net/http"

type errorReturn struct {
	ErrorType  string
	HttpStatus int
}

// https://heycam.github.io/webidl/#idl-DOMException-error-names
var (
	NotSupportedError = errorReturn{
		ErrorType:  "NotSupportedError",
		HttpStatus: http.StatusBadRequest,
	}
	NotFoundError = errorReturn{
		ErrorType:  "NotFoundError",
		HttpStatus: http.StatusNotFound,
	}
	EncodingError = errorReturn{
		ErrorType:  "EncodingError",
		HttpStatus: http.StatusBadRequest,
	}
	UnknownError = errorReturn{
		ErrorType:  "UnknownError",
		HttpStatus: http.StatusInternalServerError,
	}
	NotAllowedError = errorReturn{
		ErrorType:  "NotAllowedError",
		HttpStatus: http.StatusUnauthorized,
	}
	DataError = errorReturn{
		ErrorType:  "DataError",
		HttpStatus: http.StatusBadRequest,
	}
)
