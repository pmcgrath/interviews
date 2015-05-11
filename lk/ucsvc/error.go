package main

import (
	"errors"
	"log"
	"net/http"

	"github.com/labstack/echo"
)

var (
	errInvalidUserID                   = errors.New("Invalid user Id")
	errInvalidRequestData              = errors.New("Invalid request data")
	errNotAuthorized                   = errors.New("Not authorized")
	errNotAuthorizedInvalidCredentials = errors.New("Not authorized - invalid credentials")
	errNotFound                        = errors.New("Not found")
)

func globalErrorHandler(he *echo.HTTPError, c *echo.Context) {
	method := c.Request.Method
	path := c.Request.URL.Path
	log.Printf("Error detected for %s on %s Error details code: %d message: and error:%#v", method, path, he.Code, he.Message, he.Error)

	// Based on default one from echo.New function, altered to cater for specific errors
	// Default code
	if he.Code == 0 {
		he.Code = http.StatusInternalServerError
	}
	// Deal with specific errors
	if he.Error != nil {
		switch he.Error {
		case errInvalidUserID:
			he.Code = http.StatusBadRequest
			// Add specific error message ?
		case errInvalidRequestData:
			he.Code = http.StatusBadRequest
		case errNotAuthorized:
			he.Code = http.StatusUnauthorized
		case errNotAuthorizedInvalidCredentials:
			he.Code = http.StatusUnauthorized
			// Add specific error message ?
		case errNotFound:
			he.Code = http.StatusNotFound
		case echo.UnsupportedMediaType:
			// Echo has very limited media type support - See https://github.com/labstack/echo/blob/4068674a0b0fc16d6c33548445208d21cfbfa15b/echo.go#L135
			// It also only checks the for a prefix match
			he.Code = http.StatusUnsupportedMediaType
		}
	}
	// If no message, set based on code
	if he.Message == "" {
		he.Message = http.StatusText(he.Code)
	}

	http.Error(c.Response, he.Message, he.Code)
}
