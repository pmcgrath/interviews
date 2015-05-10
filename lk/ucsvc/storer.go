package main

type storer interface {
	GetAllUsers() ([]*connectedUser, error)
	GetUser(id ID) (*connectedUser, error)
	SaveUser(user *user) error
	DeleteUser(id ID) error
	CreateUserConnection(id1, id2 ID) error
	DeleteUserConnection(id1, id2 ID) error
}
