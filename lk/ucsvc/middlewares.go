package main

import "github.com/labstack/echo"

func authenticate(userName, password string) bool {
	// Should be going to a store, but this is good enough for now
	return userName == "ted" && password == "toe"
}

func basicAuth(h echo.HandlerFunc) echo.HandlerFunc {
	return echo.HandlerFunc(func(c *echo.Context) error {
		userName, password, ok := c.Request.BasicAuth()
		if !ok {
			return errNotAuthorized
		}

		if !authenticate(userName, password) {
			return errNotAuthorizedInvalidCredentials
		}

		return h(c)
	})
}
