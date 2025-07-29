package http

import (
	"fmt"
	"io"
	"slices"
	"strings"
)

type Request struct {
	HttpMethod    string
	HttpVersion   string
	Path          string
	Headers       Headers
	Body          []byte
	ContentLength int
}

var validVersions = []string{
	"HTTP/1.1",
}

var acceptedMethods = []string{
	"GET",
	"PUT",
	"POST",
	"OPTIONS",
	"PATCH",
	"HEAD",
	"DELETE",
	"TRACE",
	"CONNECT",
}

const BUF_SIZE = 1024

// Takes a reader (stream of bytes) and parses it into a HTTP Request
func requestFromStream(stream io.Reader) (*Request, error) {
	buffer := make([]byte, BUF_SIZE)
	p := newParser()
	for p.state != COMPLETE {
		n, err := stream.Read(buffer)
		if err == io.EOF {
			p.state = COMPLETE
			break
			//TODO: Should I return an early exit error?
		}

		err = p.parse(buffer[:n])
		if err != nil {
			return nil, err
		}
	}

	request := p.request
	return request, nil
}

func (r *Request) Print() {
	var headers strings.Builder
	for key, value := range r.Headers {
		headers.WriteString(key + ": " + value + "\n")
	}

	content := "[[Request]]\n" +
		"Method: " + r.HttpMethod + "\n" +
		"Path: " + r.Path + "\n" +
		"Version: " + r.HttpVersion + "\n" +
		"[[Headers]]\n" +
		headers.String() +
		"[[BODY]]\n" +
		string(r.Body) + "\n"

	fmt.Printf(content)
}

func (r *Request) validate() error {
	if !slices.Contains(validVersions, r.HttpVersion) {
		return ERROR_INCORRECT_VERSION
	}

	if !slices.Contains(acceptedMethods, r.HttpMethod) {
		return ERROR_METHOD_INVALID
	}

	//TODO: Extend validation

	return nil
}

var (
	ERROR_INCORRECT_VERSION = HttpError{
		responseStatus: StatusBadRequest,
		errorMsg: "Http version is incorrect",
	}

	ERROR_METHOD_INVALID = HttpError{
		responseStatus: StatusBadRequest,
		errorMsg: "Http method is not one of the accepted options",
	}
)

// TODO:
// - Extend validation (i.e. validate path inline w/ https://www.rfc-editor.org/rfc/rfc9112.html#section-3.2-6)
// - Add support for parsing chunked encoding
// - Should params be handled at parsing or at routing?
