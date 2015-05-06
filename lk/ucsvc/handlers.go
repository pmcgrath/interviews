package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
)

type UserHandler struct {
	uc *UserConnections
}

func (uh *UserHandler) getUsers(c *echo.Context) error {
	users, err := uh.uc.GetAll()
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, users)
}

func (uh *UserHandler) createUser(c *echo.Context) error {
	newUser := &NewUser{}
	if err := c.Bind(newUser); err != nil {
		return err
	}

	id, err := uh.uc.Create(newUser)
	if err != nil {
		return err
	}

	location := fmt.Sprintf("/users/%s", id)
	c.Response.Header().Set("Location", location)
	c.Response.WriteHeader(http.StatusCreated)

	return nil
}

func (uh *UserHandler) getUser(c *echo.Context) error {
	id := Id(c.P(0))

	user, err := uh.uc.Get(id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, user)
}

func (uh *UserHandler) deleteUser(c *echo.Context) error {
	id := Id(c.P(0))

	return uh.uc.Delete(id)
}

func (uh *UserHandler) getUserConnections(c *echo.Context) error {
	id := Id(c.P(0))

	user, err := uh.uc.Get(id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, user.Connections)
}

func (uh *UserHandler) createUserConnection(c *echo.Context) error {
	id1 := Id(c.P(0))
	id2 := Id(c.P(1))

	return uh.uc.CreateConnection(id1, id2)
}

func (uh *UserHandler) deleteUserConnection(c *echo.Context) error {
	id1 := Id(c.P(0))
	id2 := Id(c.P(1))

	return uh.uc.DeleteConnection(id1, id2)
}
