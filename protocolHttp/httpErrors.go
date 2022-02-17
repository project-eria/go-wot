package protocolHttp

import "net/http"

type errorReturn struct {
	errorType  string
	httpStatus int
}

// https://heycam.github.io/webidl/#idl-DOMException-error-names
var (
	NotSupportedError = errorReturn{
		errorType:  "NotSupportedError",
		httpStatus: http.StatusBadRequest,
	}
	NotFoundError = errorReturn{
		errorType:  "NotFoundError",
		httpStatus: http.StatusNotFound,
	}
	EncodingError = errorReturn{
		errorType:  "EncodingError",
		httpStatus: http.StatusBadRequest,
	}
	UnknownError = errorReturn{
		errorType:  "UnknownError",
		httpStatus: http.StatusInternalServerError,
	}
	NotAllowedError = errorReturn{
		errorType:  "NotAllowedError",
		httpStatus: http.StatusUnauthorized,
	}
	DataError = errorReturn{
		errorType:  "DataError",
		httpStatus: http.StatusBadRequest,
	}
)
