package main

type UserConnections struct {
	store Storer
}

func NewUserConnections(store Storer) *UserConnections {
	return &UserConnections{
		store: store,
	}
}

func (u *UserConnections) GetAll() ([]*ConnectedUser, error) {
	return u.store.GetAllUsers()
}

func (u *UserConnections) Get(id Id) (*ConnectedUser, error) {
	if !id.IsValid() {
		return nil, ErrInvalidUserId
	}

	return u.store.GetUser(id)
}

func (u *UserConnections) Create(newUser *NewUser) (Id, error) {
	if !newUser.IsValid() {
		return Id(""), ErrInvalidRequestData
	}

	lid := NewId()
	user := &User{
		Id:   lid,
		Name: newUser.Name,
	}

	if err := u.store.SaveUser(user); err != nil {
		return Id(""), err
	}

	return lid, nil
}

func (u *UserConnections) Delete(id Id) error {
	if !id.IsValid() {
		return ErrInvalidUserId
	}

	if _, err := u.store.GetUser(id); err != nil {
		return err
	}

	return u.store.DeleteUser(id)
}

func (u *UserConnections) CreateConnection(id1, id2 Id) error {
	if !id1.IsValid() {
		return ErrInvalidUserId
	}

	if !id2.IsValid() {
		return ErrInvalidUserId
	}

	if _, err := u.store.GetUser(id1); err != nil {
		return err
	}

	if _, err := u.store.GetUser(id2); err != nil {
		return err
	}

	return u.store.CreateUserConnection(id1, id2)
}

func (u *UserConnections) DeleteConnection(id1, id2 Id) error {
	if !id1.IsValid() {
		return ErrInvalidUserId
	}

	if !id2.IsValid() {
		return ErrInvalidUserId
	}

	if _, err := u.store.GetUser(id1); err != nil {
		return err
	}

	if _, err := u.store.GetUser(id2); err != nil {
		return err
	}

	return u.store.DeleteUserConnection(id1, id2)
}
