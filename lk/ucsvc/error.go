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

func globalErrorHandler(code int, err error, c *echo.Context) {
	if err == nil {
		return
	}

	method := c.Request.Method
	path := c.Request.URL.Path
	log.Printf("Error detected for %s on %s : %#v", method, path, err)

	switch err {
	case errInvalidUserID:
		http.Error(c.Response, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	case errInvalidRequestData:
		http.Error(c.Response, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	case errNotAuthorized:
		http.Error(c.Response, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
	case errNotAuthorizedInvalidCredentials:
		http.Error(c.Response, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
	case errNotFound:
		http.Error(c.Response, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	case echo.UnsupportedMediaType:
		// Echo has very limited media type support - See https://github.com/labstack/echo/blob/4068674a0b0fc16d6c33548445208d21cfbfa15b/echo.go#L135
		// It also only checks the for a prefix match
		http.Error(c.Response, http.StatusText(http.StatusUnsupportedMediaType), http.StatusUnsupportedMediaType)
	default:
		http.Error(c.Response, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
