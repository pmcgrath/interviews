package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
)

var userConns *userConnections

func getUsers(c *echo.Context) *echo.HTTPError {
	users, err := userConns.GetAll()
	if err != nil {
		return &echo.HTTPError{Error: err}
	}

	return c.JSON(http.StatusOK, users)
}

func createUser(c *echo.Context) *echo.HTTPError {
	newUser := &newUser{}
	if err := c.Bind(newUser); err != nil {
		return err
	}

	id, err := userConns.Create(newUser)
	if err != nil {
		return &echo.HTTPError{Error: err}
	}

	location := fmt.Sprintf("/users/%s", id)
	c.Response.Header().Set("Location", location)
	c.Response.WriteHeader(http.StatusCreated)

	return nil
}

func getUser(c *echo.Context) *echo.HTTPError {
	id := ID(c.P(0))

	user, err := userConns.Get(id)
	if err != nil {
		return &echo.HTTPError{Error: err}
	}

	return c.JSON(http.StatusOK, user)
}

func deleteUser(c *echo.Context) *echo.HTTPError {
	id := ID(c.P(0))

	if err := userConns.Delete(id); err != nil {
		return &echo.HTTPError{Error: err}
	}

	return nil
}

func getUserConnections(c *echo.Context) *echo.HTTPError {
	id := ID(c.P(0))

	user, err := userConns.Get(id)
	if err != nil {
		return &echo.HTTPError{Error: err}
	}

	return c.JSON(http.StatusOK, user.Connections)
}

func createUserConnection(c *echo.Context) *echo.HTTPError {
	id1 := ID(c.P(0))
	id2 := ID(c.P(1))

	err := userConns.CreateConnection(id1, id2)
	if err != nil {
		return &echo.HTTPError{Error: err}
	}

	c.Response.WriteHeader(http.StatusCreated)

	return nil
}

func deleteUserConnection(c *echo.Context) *echo.HTTPError {
	id1 := ID(c.P(0))
	id2 := ID(c.P(1))

	err := userConns.DeleteConnection(id1, id2)
	if err != nil {
		return &echo.HTTPError{Error: err}
	}

	return nil
}
