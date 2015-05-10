package main

import "sync"

func newMemoryStore() storer {
	return &memoryStore{
		UserIndex:            make(map[ID]*user),
		UserConnectionsIndex: make(map[ID]map[ID]struct{}),
	}
}

// Should really be using a data structure that supports graph like access - bi-directional, but since we are going to use a seperate
// persistence service this will do for now, the interface is the important thing at this stage
type memoryStore struct {
	UserIndex            map[ID]*user
	UserConnectionsIndex map[ID]map[ID]struct{}
	sync.RWMutex
}

func (s *memoryStore) GetAllUsers() ([]*connectedUser, error) {
	s.RLock()
	defer s.RUnlock()

	connectedUsers := make([]*connectedUser, 0, len(s.UserIndex))
	for _, user := range s.UserIndex {
		connectedUser, err := s.GetUser(user.ID)
		if err != nil {
			return nil, err
		}

		connectedUsers = append(connectedUsers, connectedUser)
	}

	return connectedUsers, nil
}

func (s *memoryStore) GetUser(id ID) (*connectedUser, error) {
	s.RLock()
	defer s.RUnlock()

	theUser, ok := s.UserIndex[id]
	if !ok {
		return nil, errNotFound
	}

	var connections []*user
	if cids, ok := s.UserConnectionsIndex[theUser.ID]; ok {
		connections = make([]*user, 0, len(cids))
		for cid := range cids {
			cuser := s.UserIndex[cid]
			connections = append(connections, cuser)
		}
	}
	connectedUser := &connectedUser{user: theUser, Connections: connections}

	return connectedUser, nil
}

func (s *memoryStore) SaveUser(user *user) error {
	s.Lock()
	defer s.Unlock()

	s.UserIndex[user.ID] = user // Upsert

	return nil
}

func (s *memoryStore) DeleteUser(id ID) error {
	s.Lock()
	defer s.Unlock()

	if ids, ok := s.UserConnectionsIndex[id]; ok {
		for cid := range ids {
			delete(s.UserConnectionsIndex[cid], id)
		}
		delete(s.UserConnectionsIndex, id)
	}

	delete(s.UserIndex, id)

	return nil
}

func (s *memoryStore) CreateUserConnection(id1, id2 ID) error {
	s.Lock()
	defer s.Unlock()

	_, ok := s.UserConnectionsIndex[id1]
	if !ok {
		s.UserConnectionsIndex[id1] = make(map[ID]struct{})
	}
	s.UserConnectionsIndex[id1][id2] = struct{}{}

	_, ok = s.UserConnectionsIndex[id2]
	if !ok {
		s.UserConnectionsIndex[id2] = make(map[ID]struct{})
	}
	s.UserConnectionsIndex[id2][id1] = struct{}{}

	return nil
}

func (s *memoryStore) DeleteUserConnection(id1, id2 ID) error {
	s.Lock()
	defer s.Unlock()

	_, ok := s.UserConnectionsIndex[id1]
	if ok {
		delete(s.UserConnectionsIndex[id1], id2)
	}

	_, ok = s.UserConnectionsIndex[id2]
	if ok {
		delete(s.UserConnectionsIndex[id2], id1)
	}

	return nil
}
