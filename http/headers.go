package http

import (
	"bytes"
	"strings"
)

type Headers map[string]string

// Adds a header, if header already exists - will concatenate the value
func (H Headers) Add(header string, value string) (error) {
	err := checkHeaderForErrors(header, value)
	if err != nil {
		return err
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

// Adds a header, if header already exists - will override the value
func (H Headers) Set(header string, value string) (error) {
	err := checkHeaderForErrors(header, value)
	if err != nil {
		return err
	}

	fieldName := strings.ToLower(header)
	H[fieldName] = value

	return nil
}

func (H Headers) Get(header string) (string, bool) {
	fieldName := strings.ToLower(header)
	value, exists := H[fieldName]
	return value, exists
}

func checkHeaderForErrors(header string, value string) (error) {
	if !isValidFieldName(header) {
		return ERROR_HEADER_FIELD_NAME
	}

	if len(value) == 0 {
		return ERROR_HEADER_NO_VALUE
	}

	return nil
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

//NOTE:
//Std library uses map[string][]string rather than map[string]string - this allows for multi values to be done without concatenating the string (easier management)
//I was thinking that a different implementation might also help with response creation... i.e. pass headers as individual options
