package main

import "testing"

var (
	theStoreMock          = &storeMock{}
	userConnectionsToTest = &userConnections{Store: theStoreMock}
)

func TestUserConnections_GetAll(t *testing.T) {
	theStoreMock.Reset()

	// Not much to test here as this is a pass through function
	allUsers, err := userConnectionsToTest.GetAll()
	if allUsers != nil {
		t.Error("Expected allUsers to be nil")
	}
	if err != nil {
		t.Error("Expected error to be nil")
	}
}

func TestUserConnections_Get_Where_Invalid_Id(t *testing.T) {
	theStoreMock.Reset()

	user, err := userConnectionsToTest.Get(emptyID)

	if user != nil {
		t.Error("Expected user to be nil")
	}
	if err != errInvalidUserID {
		t.Error("Expected invalid user Id error")
	}
}

func TestUserConnections_Get_Where_Store_Error(t *testing.T) {
	theStoreMock.Reset()
	theStoreMock.ErrorToReturn = errNotFound

	id := newID()
	user, err := userConnectionsToTest.Get(id)

	if user != nil {
		t.Error("Expected user to be nil")
	}
	if err != errNotFound {
		t.Error("Expected not found error")
	}
}

func TestUserConnections_Get_Where_User_Not_Found(t *testing.T) {
	id := newID()

	theStoreMock.Reset()
	theStoreMock.ConnectedUser = &connectedUser{user: &user{ID: id, Name: "Ted Toe"}}

	user, err := userConnectionsToTest.Get(id)

	if user != theStoreMock.ConnectedUser {
		t.Error("Expected user to match")
	}
	if err != nil {
		t.Error("Expected no error")
	}
}

func TestUserConnections_Create_Where_Invalid_Data(t *testing.T) {
	theStoreMock.Reset()

	newUser := &newUser{}

	userID, err := userConnectionsToTest.Create(newUser)

	if userID != "" {
		t.Error("Expected no userId")
	}
	if err != errInvalidRequestData {
		t.Error("Expected invalid request data")
	}
}
