package http

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
)

type Response struct {
	HttpVersion string
	StatusCode  Status
	Headers     Headers
	Body        []byte
}

// Converts the response struct into a valid byte representation for a HTTP/1.1 response
func (r *Response) Marshall() []byte {
	buf := bytes.Buffer{}
	statusline := r.HttpVersion + SPACE + r.StatusCode.String()

	buf.WriteString(statusline)
	buf.Write(CRLF)

	for key, value := range r.Headers {
		buf.WriteString(key + ": " + value)
		buf.Write(CRLF)
	}
	buf.Write(CRLF)

	buf.Write(r.Body)

	return buf.Bytes()
}

// Writes the responce to the writer (i.e. net.Connection)
func (r *Response) Write(writer io.Writer) {
	marshalledResponse := r.Marshall()
	writer.Write(marshalledResponse)
}

// Converts the response struct into a string representation. This string representation should be a valid HTTP/1.1 response.
func (r *Response) String() string {
	return string(r.Marshall())
}

func (r *Response) Print() {
	fmt.Print("[[RESPONSE]]\n" + r.String() + "\n")
}

const VERSION = "HTTP/1.1"

// Creates a response from status code and the content/body (as bytes).
func CreateResponse(status Status, content []byte) *Response {
	res := &Response{
		HttpVersion: VERSION,
		StatusCode: status,
		Headers: Headers{"Connection": "close"},
		Body: content,
	}

	length := strconv.Itoa(len(content))
	res.SetHeader("Content-Length", length)

	res.SetHeader("Content-Type", "text/plain") //Can be overwritten later

	res.Body = content 
	return res
}

// Returns an empty 200 OK response
func Ok() *Response {
	return &Response {
		HttpVersion: VERSION,
		StatusCode: StatusOk,
		Headers: Headers{"Connection": "close"},
		Body: []byte{},
	}
}

// Adds header, if it exists will concatenate onto existing value
func (r *Response) AddHeader(header string, value string) error {
	return r.Headers.Add(header, value)
}

// Adds header, if it exists will override existing value
func (r *Response) SetHeader(header string, value string) error {
	return r.Headers.Set(header, value)
}
		
// TODO:
// - Additional Headers: i.e. cache-control, datetime, etc
// - Support for streamed response (i.e. to allow for chunked encode)
// - More structured way to handle content-type?
