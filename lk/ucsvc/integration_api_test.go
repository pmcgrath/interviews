// +build integration

// See http://peter.bourgon.org/go-in-production/#testing-and-validation for seperating integration tests
// Test with : 	go test -tags integration -v
//		go test -tags "test integration" -v
package main

import "testing"

func TestGetAllUsersWithNoAuthorization(t *testing.T) {
	result := executeHttp("", password, "GET", "users", jsonContentType, nil)

	if result.StatusCode != 401 {
		t.Errorf("Status code [%d] is not 401", result.StatusCode)
	}
}

func TestGetAllUsers(t *testing.T) {
	result := getAllUsersViaApi()

	if result.StatusCode != 200 {
		t.Errorf("Status code [%d] is not 200", result.StatusCode)
	}
	if result.Error != nil {
		t.Errorf("Unexpected error: %#v", result.Error)
	}
}

func TestCreateNewUser(t *testing.T) {
	name := "Tom"

	result := createNewUserViaApi(name)

	if result.StatusCode != 201 {
		t.Errorf("Status code [%d] is not 201", result.StatusCode)
	}
	if result.Error != nil {
		t.Errorf("Unexpected error: %#v", result.Error)
	}
	if result.NewUserId == "" {
		t.Error("The new user Id from location header should not be empty")
	}
}

func TestGetUserForCreatedUser(t *testing.T) {
	name := "Tom"
	userIds := createTestUsers("Tom")

	result := getUserViaApi(userIds[0])

	if result.StatusCode != 200 {
		t.Errorf("Status code [%d] is not 200", result.StatusCode)
	}
	if result.Error != nil {
		t.Errorf("Unexpected error: %#v", result.Error)
	}
	if result.User["id"] != userIds[0] {
		t.Errorf("User Id [%s] is not the same as creation user Id [%s]", result.User["id"], userIds[0])
	}
	if result.User["name"] != name {
		t.Errorf("User name [%s] is not the same as name used for creation [%s]", result.User["name"], name)
	}
	if result.User["connections"] != nil {
		t.Error("Should have no connections")
	}
}

func TestCreateUserConnections(t *testing.T) {
	userIds := createTestUsers("Ted", "Toe")

	result := createUserConnectionsViaApi(userIds[0], userIds[1])

	if result.StatusCode != 201 {
		t.Errorf("Status code [%d] is not 201", result.StatusCode)
	}
	if result.Error != nil {
		t.Errorf("Unexpected error: %#v", result.Error)
	}
}

func TestGetUserConnectionsAfterConnectionCreation(t *testing.T) {
	userIds := createTestUsers("Ted", "Toe", "Joe")
	createUserConnectionsViaApi(userIds[0], userIds[2])

	getResult1 := getUserConnectionsViaApi(userIds[0])
	getResult2 := getUserConnectionsViaApi(userIds[2])

	if getResult1.StatusCode != 200 {
		t.Errorf("Status code[1] [%d] is not 200", getResult1.StatusCode)
	}
	if getResult2.StatusCode != 200 {
		t.Errorf("Status code [%d] is not 200", getResult2.StatusCode)
	}
	if getResult1.Error != nil {
		t.Errorf("Unexpected error: %#v", getResult1.Error)
	}
	if getResult2.Error != nil {
		t.Errorf("Unexpected error: %#v", getResult2.Error)
	}
	if len(getResult1.Users) != 1 {
		t.Errorf("Unexpected connections count: %d", len(getResult1.Users))
	}
	if len(getResult2.Users) != 1 {
		t.Errorf("Unexpected connections count: %d", len(getResult2.Users))
	}
	if getResult1.Users[0]["id"] != userIds[2] {
		t.Errorf("Missing connection for: %s", userIds[2])
	}
	if getResult2.Users[0]["id"] != userIds[0] {
		t.Errorf("Missing connection for: %s", userIds[0])
	}
}

func TestDeleteUserConnection(t *testing.T) {
	userIds := createTestUsers("Ted", "Toe", "Joe")

	result := deleteUserConnectionsViaApi(userIds[0], userIds[1])

	if result.StatusCode != 200 {
		t.Errorf("Status code [%d] is not 200", result.StatusCode)
	}
	if result.Error != nil {
		t.Errorf("Unexpected error: %#v", result.Error)
	}
}

func TestGetUserConnectionsAfterConnectionDeletion(t *testing.T) {
	userIds := createTestUsers("Ted", "Toe", "Joe")
	createUserConnectionsViaApi(userIds[0], userIds[2])
	deleteUserConnectionsViaApi(userIds[0], userIds[2])

	getResult1 := getUserConnectionsViaApi(userIds[0])
	getResult2 := getUserConnectionsViaApi(userIds[2])

	if getResult1.StatusCode != 200 {
		t.Errorf("Status code[1] [%d] is not 200", getResult1.StatusCode)
	}
	if getResult2.StatusCode != 200 {
		t.Errorf("Status code [%d] is not 200", getResult2.StatusCode)
	}
	if len(getResult1.Users) != 0 {
		t.Errorf("Unexpected connections count: %d", len(getResult1.Users))
	}
	if len(getResult2.Users) != 0 {
		t.Errorf("Unexpected connections count: %d", len(getResult2.Users))
	}
}

func TestCreateNewUserWithNoData(t *testing.T) {
	result := executeHttp(userName, password, "POST", "users", jsonContentType, []byte(`{}`))

	if result.StatusCode != 400 {
		t.Errorf("Status code [%d] is not 400", result.StatusCode)
	}
}

func TestCreateNewUserWithEmptyName(t *testing.T) {
	result := executeHttp(userName, password, "POST", "users", jsonContentType, []byte(`{"name": ""}`))

	if result.StatusCode != 400 {
		t.Errorf("Status code [%d] is not 400", result.StatusCode)
	}
}

func TestGetUserWhereUserDoesNotExist(t *testing.T) {
	result := getUserViaApi(validUserId)

	if result.StatusCode != 404 {
		t.Errorf("Status code [%d] is not 404", result.StatusCode)
	}
}

func TestUrlThatUsedToCauseEchoToPanicButIsNowFixed(t *testing.T) {
	result := executeSimpleHttp("PUT", "users/f47ac10b-58cc-0372-8567-0e02b2c3d479/connections/a47ac10b-58cc-0372-8567-0e02b2c3d470/connections")

	if result.StatusCode != 404 {
		t.Errorf("Status code [%d] is not 404", result.StatusCode)
	}
}

func TestCreateUserWithUnsupportedMediaType(t *testing.T) {
	result := executeHttp(userName, password, "POST", "users", "xml", nil)

	if result.StatusCode != 415 {
		t.Errorf("Status code [%d] is not 415", result.StatusCode)
	}
}

func createTestUsers(names ...string) []string {
	userIds := make([]string, 0)
	for _, name := range names {
		result := createNewUserViaApi(name)
		userIds = append(userIds, result.NewUserId)
	}

	return userIds
}
