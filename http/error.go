package http

import (
	"fmt"
)

type HttpError struct {
	responseStatus Status
	errorMessage   string
}

// Implement the Go error interface
func (err HttpError) Error() string {
	return err.errorMessage
}

func errorAsHttpError(err error) HttpError {
	value, isHttpError := err.(HttpError)
	if isHttpError {
		return value
	}

	fmt.Print(err)

	return HttpError{
		responseStatus: StatusInternalError,
		errorMessage:   "An unknown internal error occurred",
	}
	//Note: Probably don't want to expose internal errors to end-users
}
