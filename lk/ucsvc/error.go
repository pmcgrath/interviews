package main

import (
	"errors"
	"log"
	"net/http"

	"github.com/labstack/echo"
)

var (
	ErrInvalidUserId      = errors.New("Invalid user Id")
	ErrInvalidRequestData = errors.New("Invalid request data")
	ErrNotFound           = errors.New("Not found")
)

func globalErrorHandler(code int, err error, c *echo.Context) {
	if err == nil {
		return
	}

	method := c.Request.Method
	path := c.Request.URL.Path
	log.Printf("Error detected for %s on %s : %#v", method, path, err)

	switch err {
	case ErrInvalidUserId:
		http.Error(c.Response, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	case ErrInvalidRequestData:
		http.Error(c.Response, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	case ErrNotFound:
		http.Error(c.Response, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	case echo.UnsupportedMediaType:
		// Echo has very limited media type support - See https://github.com/labstack/echo/blob/4068674a0b0fc16d6c33548445208d21cfbfa15b/echo.go#L135
		// It also only checks the for a prefix match
		http.Error(c.Response, http.StatusText(http.StatusUnsupportedMediaType), http.StatusUnsupportedMediaType)
	default:
		http.Error(c.Response, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
