package main

import (
	"errors"
	"net/http"

	"github.com/labstack/echo"
)

var (
	errInvalidUserID      = errors.New("Invalid user Id")
	errInvalidRequestData = errors.New("Invalid request data")
	errNotFound           = errors.New("Not found")
)

func mapToEchoError(err error) *echo.HTTPError {
	code := 0
	if err != nil {
		switch err {
		case errInvalidUserID:
			code = http.StatusBadRequest
		case errInvalidRequestData:
			code = http.StatusBadRequest
		case errNotFound:
			code = http.StatusNotFound
		case echo.UnsupportedMediaType:
			// Echo has very limited support for content negotiation
			code = http.StatusUnsupportedMediaType
		}
	}
	return &echo.HTTPError{Code: code}
}
