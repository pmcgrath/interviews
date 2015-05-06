package main

import (
	"strings"

	"code.google.com/p/go-uuid/uuid"
)

type Id string

func NewId() Id {
	return Id(uuid.New())
}

func (i Id) IsValid() bool {
	idAsString := string(i)
	parsedId := uuid.Parse(idAsString)

	if parsedId == nil {
		return false
	}

	return idAsString == parsedId.String()
}

type NewUser struct {
	Name string `json:"name"`
}

func (u *NewUser) IsValid() bool {
	if strings.TrimSpace(u.Name) == "" {
		return false
	}

	return true
}

type User struct {
	Id   Id     `json:"id,Id"`
	Name string `json:"name"`
}

func (u *User) IsValid() bool {
	if !u.Id.IsValid() {
		return false
	}
	if strings.TrimSpace(u.Name) == "" {
		return false
	}

	return true
}

type ConnectedUser struct {
	*User
	Connections []*User `json:"connections"`
}
