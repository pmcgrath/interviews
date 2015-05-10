package main

import (
	"strings"

	"code.google.com/p/go-uuid/uuid"
)

var (
	emptyID = ID("")
)

type ID string

func newID() ID {
	return ID(uuid.New())
}

func (i ID) IsValid() bool {
	idAsString := string(i)
	parsedID := uuid.Parse(idAsString)

	if parsedID == nil {
		return false
	}

	return idAsString == parsedID.String()
}

type newUser struct {
	Name string `json:"name"`
}

func (u *newUser) IsValid() bool {
	if strings.TrimSpace(u.Name) == "" {
		return false
	}

	return true
}

type user struct {
	ID   ID     `json:"id,ID"`
	Name string `json:"name"`
}

func (u *user) IsValid() bool {
	if !u.ID.IsValid() {
		return false
	}
	if strings.TrimSpace(u.Name) == "" {
		return false
	}

	return true
}

type connectedUser struct {
	*user
	Connections []*user `json:"connections"`
}
