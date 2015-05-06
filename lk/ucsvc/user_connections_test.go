package main

type StoreMock struct {
}

func (s *StoreMock) GetAllUsers() ([]*ConnectedUser, error) {
	return nil, nil
}

func (s *StoreMock) GetUser(id Id) (*ConnectedUser, error) {
	return nil, nil
}

func (s *StoreMock) CreateUser(newUser *NewUser) (*ConnectedUser, error) {
	return nil, nil
}

func (s *StoreMock) DeleteUser(id Id) error {
	return nil
}

func (s *StoreMock) CreateUserConnection(id1, id2 Id) error {
	return nil
}

func (s *StoreMock) DeleteUserConnection(id1, id2 Id) error {
	return nil
}
