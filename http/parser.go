package http

import (
	"bytes"
	"strconv"
	"strings"
)

var CRLF = []byte("\r\n")
const SPACE = " "

type State int

const (
	INIT State = iota
	HEADERS
	MESSAGE
	COMPLETE
)

type requestParser struct {
	state         State
	request       *Request
	line          []byte
}

func newParser() *requestParser {
	return &requestParser{
		state: INIT,
		request: &Request{
			Headers: Headers{},
			Body:    []byte{},
			ContentLength: 0,
		},
		line:          []byte{},
	}
}

// Works towards parsing a request.
// Does not require data to contain all the request bytes (see request.requestFromStream) - will work with whats its got and remember its spot.
func (p *requestParser) parse(data []byte) error {
	for _, char := range data {
		p.line = append(p.line, char)
		if bytes.HasSuffix(p.line, CRLF) {
			p.line = bytes.TrimSuffix(p.line, CRLF)
			err := p.parseLine()
			if err != nil {
				return err
			}
			p.line = []byte{}
		}
		if p.state == MESSAGE {
			p.parseBody()
		}
	}
	return nil
}

func (p *requestParser) parseLine() error {
	switch p.state {
	case INIT:
		return p.parseRequestLine()
	case HEADERS:
		if len(p.line) == 0 {
			p.state = MESSAGE
			return nil
		}
		return p.parseHeader()
	}
	return nil
}

func (p *requestParser) parseRequestLine() error {
	line := string(p.line)
	components := strings.Split(line, SPACE)
	if len(components) != 3 {
		return ERROR_MALFORMED_RL
	}

	p.request.HttpMethod  = components[0]
	p.request.Path  	  = components[1]
	p.request.HttpVersion = components[2]

	p.state = HEADERS
	return p.request.validate()
}

func (p *requestParser) parseHeader() error {
	line := string(p.line)
	header, value, success := strings.Cut(line, ":")
	if !success {
		return ERROR_MALFORMED_HEADERS
	}

	header = strings.TrimPrefix(header, SPACE)
	value = strings.TrimSpace(value)

	err := p.request.Headers.Add(header, value)
	if err != nil {
		return err
	}
	return nil
}

func (p *requestParser) parseBody() error {
	if p.request.ContentLength == 0 {
		cl, exists := p.request.Headers.Get("content-length")
		if !exists {
			p.state = COMPLETE
			return nil
		}

		l, err := strconv.Atoi(cl)
		if err != nil {
			return ERROR_LENGTH_NOT_NUM
		}
		p.request.ContentLength = l
	}
	//FIX: Edge case "content-length" exists but is zero	

	body := p.line
	length := p.request.ContentLength

	if len(body) > length {
		return ERROR_BODY_LENGTH
	}

	if len(body) == length {
		p.request.Body = body
		p.state = COMPLETE
	}

	return nil
}

var (
	ERROR_BODY_LENGTH = HttpError{
		responseStatus: StatusBadRequest,	
		errorMsg: "Error with body length or content-length header",
	}

	ERROR_LENGTH_NOT_NUM = HttpError{
		responseStatus: StatusBadRequest,
		errorMsg: "Header 'content-length' was not a valid number",
	}

	ERROR_MALFORMED_RL = HttpError{
		responseStatus: StatusBadRequest,
		errorMsg: "Request line does not have three components",
	}

	ERROR_MALFORMED_HEADERS = HttpError{
		responseStatus: StatusBadRequest,
		errorMsg: "Header line is missing semicolon",
	}
)
