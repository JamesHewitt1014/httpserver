package http

import (
	"fmt"
)

type HttpError struct {
	responseStatus Status
	errorMsg       string
}

// Implement the Go error interface
func (err HttpError) Error() string {
	return err.errorMsg
}

func ResponseFromError(err error) *Response {
	httpError := errorAsHttpError(err)
	message := []byte(httpError.errorMsg)
	return CreateResponse(httpError.responseStatus, message)
}

func errorAsHttpError(err error) HttpError {
	value, isHttpError := err.(HttpError)
	if isHttpError {
		return value
	}

	fmt.Print(err)

	return HttpError{
		responseStatus: StatusInternalError,
		errorMsg: "An unknown internal error occurred",
	}
	//Note: Probably don't want to expose internal errors to end-users
}
