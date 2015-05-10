package main

type storeMock struct {
	ErrorToReturn     error
	AllConnectedUsers []*connectedUser
	ConnectedUser     *connectedUser
}

func (s *storeMock) GetAllUsers() ([]*connectedUser, error) {
	if s.ErrorToReturn != nil {
		return nil, s.ErrorToReturn
	}

	return s.AllConnectedUsers, nil
}

func (s *storeMock) GetUser(id ID) (*connectedUser, error) {
	if s.ErrorToReturn != nil {
		return nil, s.ErrorToReturn
	}
	return s.ConnectedUser, nil
}

func (s *storeMock) SaveUser(user *user) error {
	if s.ErrorToReturn != nil {
		return s.ErrorToReturn
	}
	return nil
}

func (s *storeMock) DeleteUser(id ID) error {
	if s.ErrorToReturn != nil {
		return s.ErrorToReturn
	}
	return nil
}

func (s *storeMock) CreateUserConnection(id1, id2 ID) error {
	if s.ErrorToReturn != nil {
		return s.ErrorToReturn
	}
	return nil
}

func (s *storeMock) DeleteUserConnection(id1, id2 ID) error {
	if s.ErrorToReturn != nil {
		return s.ErrorToReturn
	}
	return nil
}

func (s *storeMock) Reset() {
	s.ErrorToReturn = nil
	s.AllConnectedUsers = nil
	s.ConnectedUser = nil
}
