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

// Creates a response from status code and the content/body (as bytes)
func CreateResponse(status Status, content []byte) *Response {
	res := &Response{}
	res.SetStatusLine(status)
	res.SetDefaultHeaders(len(content))
	res.Body = content 
	return res
}

// Returns an empty 200 OK response
func Ok() *Response {
	return &Response {
		HttpVersion: "HTTP/1.1",
		StatusCode: StatusOk,
		Headers: Headers{"Connection": "close"},
		Body: []byte{},
	}
}

// Converts the response struct into the complete byte representation
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

func (r *Response) String() string {
	return string(r.Marshall())
}

func (r *Response) Print() {
	fmt.Print("[[RESPONSE]]\n" + r.String() + "\n")
}

const VERSION = "HTTP/1.1"

func (r *Response) SetStatusLine(statusCode Status) error {
	r.HttpVersion = VERSION
	r.StatusCode  = statusCode
	return nil
}

func (r *Response) SetHeader(header string, value string) error {
	return r.Headers.Add(header, value)
} // FIX: Sometimes you may want to override as opposed to just adding

func (r *Response) SetDefaultHeaders(contentLength int) {
	headers := Headers{}

	length := strconv.Itoa(contentLength)
	headers.Add("Content-Length", length)
	headers.Add("Connection", "close")
	headers.Add("Content-Type", "text/plain")
	r.Headers = headers
}
//TODO: Refactor this - maybe to use SetHeader
// Maybe add more options to CreateResponse instead (i.e. content-type, etc)

// Writes everything at once to the writer (i.e. net.Connection)
func (r *Response) Write(writer io.Writer) {
	marshalledResponse := r.Marshall()
	writer.Write(marshalledResponse)
}
		
// TODO: 
// - Additional Headers: i.e. cache-control, datetime, etc
// - Support for streamed response (i.e. to allow for chunked encode)
