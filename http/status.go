package http

import (
	"fmt"
)

type Status int
const (
	StatusOk            Status = 200
	StatusCreated              = 201
	StatusBadRequest           = 400
	StatusUnauthorised         = 401
	StatusNotFound             = 404
	StatusInternalError        = 500
)

var statusMessage = map[Status]string{
	200: "Ok",
	201: "Created",
	400: "Bad Request",
	401: "Not authenticated",
	404: "Not Found",
	500: "Internal Server Error",
}

func (s Status) String() string {
	n := int(s)
	msg, _ := statusMessage[s]
	return fmt.Sprintf("%d", n) + SPACE + msg
}
