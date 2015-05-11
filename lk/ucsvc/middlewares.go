package main

import "github.com/labstack/echo"

func authenticate(userName, password string) bool {
	// Should be going to a store, but this is good enough for now
	return userName == "ted" && password == "toe"
}

func basicAuth(h echo.HandlerFunc) echo.HandlerFunc {
	// See https://tools.ietf.org/html/rfc7235#page-7
	// Pending - Set WWW-authenticate response header in the event of a failure
	return echo.HandlerFunc(func(c *echo.Context) *echo.HTTPError {
		userName, password, ok := c.Request.BasicAuth()
		if !ok {
			return &echo.HTTPError{Error: errNotAuthorized}
		}

		if !authenticate(userName, password) {
			return &echo.HTTPError{Error: errNotAuthorizedInvalidCredentials}
		}

		return h(c)
	})
}
