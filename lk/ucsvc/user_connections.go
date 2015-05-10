package main

func newUserConnections(store storer) *userConnections {
	return &userConnections{
		Store: store,
	}
}

type userConnections struct {
	Store storer
}

func (u *userConnections) GetAll() ([]*connectedUser, error) {
	return u.Store.GetAllUsers()
}

func (u *userConnections) Get(id ID) (*connectedUser, error) {
	if !id.IsValid() {
		return nil, errInvalidUserID
	}

	return u.Store.GetUser(id)
}

func (u *userConnections) Create(newUser *newUser) (ID, error) {
	if !newUser.IsValid() {
		return emptyID, errInvalidRequestData
	}

	id := newID()
	user := &user{
		ID:   id,
		Name: newUser.Name,
	}

	if err := u.Store.SaveUser(user); err != nil {
		return emptyID, err
	}

	return id, nil
}

func (u *userConnections) Delete(id ID) error {
	if !id.IsValid() {
		return errInvalidUserID
	}

	if _, err := u.Store.GetUser(id); err != nil {
		return err
	}

	return u.Store.DeleteUser(id)
}

func (u *userConnections) CreateConnection(id1, id2 ID) error {
	if !id1.IsValid() {
		return errInvalidUserID
	}

	if !id2.IsValid() {
		return errInvalidUserID
	}

	if _, err := u.Store.GetUser(id1); err != nil {
		return err
	}

	if _, err := u.Store.GetUser(id2); err != nil {
		return err
	}

	return u.Store.CreateUserConnection(id1, id2)
}

func (u *userConnections) DeleteConnection(id1, id2 ID) error {
	if !id1.IsValid() {
		return errInvalidUserID
	}

	if !id2.IsValid() {
		return errInvalidUserID
	}

	if _, err := u.Store.GetUser(id1); err != nil {
		return err
	}

	if _, err := u.Store.GetUser(id2); err != nil {
		return err
	}

	return u.Store.DeleteUserConnection(id1, id2)
}
