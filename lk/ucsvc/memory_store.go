package main

import "sync"

// Should really be using a data structure that supports graph like access - bi-directional, but since for now this will do,
// the interface is the important thing at this stage
type memoryStore struct {
	userIndex            map[Id]*User
	userConnectionsIndex map[Id]map[Id]struct{}
	sync.RWMutex
}

func (s *memoryStore) GetAllUsers() ([]*ConnectedUser, error) {
	s.RLock()
	defer s.RUnlock()

	connectedUsers := make([]*ConnectedUser, 0, len(s.userIndex))
	for _, user := range s.userIndex {
		connectedUser, err := s.GetUser(user.Id)
		if err != nil {
			return nil, err
		}

		connectedUsers = append(connectedUsers, connectedUser)
	}

	return connectedUsers, nil
}

func (s *memoryStore) GetUser(id Id) (*ConnectedUser, error) {
	s.RLock()
	defer s.RUnlock()

	user, ok := s.userIndex[id]
	if !ok {
		return nil, ErrNotFound
	}

	var connections []*User
	if cids, ok := s.userConnectionsIndex[user.Id]; ok {
		connections = make([]*User, 0, len(cids))
		for cid, _ := range cids {
			cuser := s.userIndex[cid]
			connections = append(connections, cuser)
		}
	}
	connectedUser := &ConnectedUser{User: user, Connections: connections}

	return connectedUser, nil
}

func (s *memoryStore) SaveUser(user *User) error {
	s.Lock()
	defer s.Unlock()

	s.userIndex[user.Id] = user // Upsert

	return nil
}

func (s *memoryStore) DeleteUser(id Id) error {
	s.Lock()
	defer s.Unlock()

	if ids, ok := s.userConnectionsIndex[id]; ok {
		for cid, _ := range ids {
			delete(s.userConnectionsIndex[cid], id)
		}
		delete(s.userConnectionsIndex, id)
	}

	delete(s.userIndex, id)

	return nil
}

func (s *memoryStore) CreateUserConnection(id1, id2 Id) error {
	s.Lock()
	defer s.Unlock()

	_, ok := s.userConnectionsIndex[id1]
	if !ok {
		s.userConnectionsIndex[id1] = make(map[Id]struct{})
	}
	s.userConnectionsIndex[id1][id2] = struct{}{}

	_, ok = s.userConnectionsIndex[id2]
	if !ok {
		s.userConnectionsIndex[id2] = make(map[Id]struct{})
	}
	s.userConnectionsIndex[id2][id1] = struct{}{}

	return nil
}

func (s *memoryStore) DeleteUserConnection(id1, id2 Id) error {
	s.Lock()
	defer s.Unlock()

	_, ok := s.userConnectionsIndex[id1]
	if ok {
		delete(s.userConnectionsIndex[id1], id2)
	}

	_, ok = s.userConnectionsIndex[id2]
	if ok {
		delete(s.userConnectionsIndex[id2], id1)
	}

	return nil
}

func NewMemoryStore() Storer {
	return &memoryStore{
		userIndex:            make(map[Id]*User),
		userConnectionsIndex: make(map[Id]map[Id]struct{}),
	}
}
