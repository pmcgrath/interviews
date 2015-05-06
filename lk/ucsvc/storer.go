package main

type Storer interface {
	GetAllUsers() ([]*ConnectedUser, error)
	GetUser(id Id) (*ConnectedUser, error)
	SaveUser(user *User) error
	DeleteUser(id Id) error
	CreateUserConnection(id1, id2 Id) error
	DeleteUserConnection(id1, id2 Id) error
}
