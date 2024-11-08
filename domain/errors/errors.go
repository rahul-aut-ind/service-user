package errors

import (
	"fmt"
	"net/http"
	"strings"
)

const (
	// ErrCodeGeneric API Error for generic or non explicit errors
	ErrCodeGeneric = "GenericError"
	// ErrCodeBadRequest API Error code for bad request
	ErrCodeBadRequest = "BadRequest"
	// ErrCodeInvalidUserIDHeader API Error code for invalid header
	ErrCodeInvalidUserIDHeader = "InvalidHeader"
	// ErrCodeNoUser API Error code for no user exists
	ErrCodeNoUser = "NoUserFound"
	// The added to all error codes to prevent conflicting with other services
	errorMessageKeyPrefix = "service-user"
)

type Error struct {
	Code         string `json:"code"`
	Message      string `json:"message"`
	DebugMessage string `json:"debugMessage"`
}

func New(code string, err error) Error {
	return Error{
		Code:         code,
		Message:      fmt.Sprintf("%s.%s", errorMessageKeyPrefix, strings.ToLower(code)),
		DebugMessage: err.Error(),
	}
}

func (e Error) Error() string {
	return e.DebugMessage
}

func (e Error) HTTPCode() int {
	errCodeMap := map[string]int{
		ErrCodeBadRequest: http.StatusBadRequest,
		ErrCodeGeneric:    http.StatusInternalServerError,
		ErrCodeNoUser:     http.StatusNotFound,
	}
	if code, ok := errCodeMap[e.Code]; ok {
		return code
	}

	return http.StatusInternalServerError
}
