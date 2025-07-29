package http

import (
	"bytes"
	"strings"
)

type Headers map[string]string

func (H Headers) Add(header string, value string) (error) {
	if !isValidFieldName(header) {
		return ERROR_HEADER_FIELD_NAME
	}

	if len(value) == 0 {
		return ERROR_HEADER_NO_VALUE
	}

	fieldName := strings.ToLower(header)
	oldValue, exists := H[fieldName]
	if exists {
		H[fieldName] = oldValue + ", " + value
	} else{
		H[fieldName] = value
	}

	return nil
}

func (H Headers) Get(header string) (string, bool) {
	fieldName := strings.ToLower(header)
	value, exists := H[fieldName]
	return value, exists
}

func isValidFieldName(header string) bool {
	if strings.HasSuffix(header, SPACE){
		return false
	}

	for _, char := range header {
		if !isValidChar(char){
			return false
		}
	}

	return len(header) >= 1
}

func isValidChar(char rune) bool {
	return bytes.ContainsRune(allowedSpecialChars, rune(char)) ||
   		'0' <= char && char <= '9' ||
		'a' <= char && char <= 'z' ||
		'A' <= char && char <= 'Z'
}

var allowedSpecialChars = []byte{
'!', '#', '$', '%', '&', '\'', '*', '+', '-', '.', '^', '_', '`', '|', '~',
}

var (
	ERROR_HEADER_NO_VALUE = HttpError{
		responseStatus: StatusBadRequest,
		errorMsg: "A header entry contains no value",
	}

	ERROR_HEADER_FIELD_NAME = HttpError{
		responseStatus: StatusBadRequest,
		errorMsg: "Header field name is not valid",
	}
)
