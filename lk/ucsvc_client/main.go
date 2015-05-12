// Should really be using cucumber style tests, rather than this client, but would require a ruby dependency which I don't want here
package main

import (
	"flag"
	"fmt"
	"strings"
)

var (
	baseURL           = flag.String("baseURL", "https://localhost:8090/", "Service base URL")
	userName          = flag.String("userName", "ted", "Service user name for basic authentication")
	password          = flag.String("password", "toe", "Service user password for basic authentication")
	pauseAfterAPICall = flag.Bool("pauseAfterAPICall", false, "Hit enter to continue after each api call, if not set will execute all commands with no pauses")
)

func init() {
	flag.Parse()
}

func main() {
	getAllUsers()

	// Create users
	userNames := [...]string{"Tom Toe", "JIM BOB", "Terry Smith", "Joe Dolan", "Jerry", "Billy", "Todd"}
	users := make(map[string]string, len(userNames))
	for _, userName := range userNames {
		res := createNewUser(userName)
		userID := extractUserIDFromLocationHeader(res.LocationHeader)
		getUser(userID)
		users[userName] = userID
	}

	// Create connections
	createConnection(users[userNames[0]], users[userNames[2]])
	createConnection(users[userNames[0]], users[userNames[5]])

	// Get user connections
	getUserConnections(users[userNames[0]])

	// Get all users
	getAllUsers()

	// Delete connection
	// Show user first
	getUser(users[userNames[0]])
	// Delete connection
	deleteConnection(users[userNames[0]], users[userNames[5]])
	// Show user - connection gone
	getUser(users[userNames[0]])

	// Bad api calls
	// Unexpected route - 404
	wrappedExecute("GET", "unknownroute", "")
	// Create user no data - should be 400
	wrappedExecute("POST", "users", "{}")
	// Create user with incomplete data, no name field in json data - should be 400
	wrappedExecute("POST", "users", `{"age": 222}`)
	// Create user using wrong method\verb - should be 405, but is 404 due to echo package ISSUE - should be returning a 405 but due to the way it does a route find it gives a 404
	wrappedExecute("PUT", "users", `{"name": "ted"}`)
	// Get unknown user where Id is not valid - Should be 400
	wrappedExecute("GET", "users/u1", "")
	// Get unknown user where Id is valid but user does not exist - Should be 404
	wrappedExecute("GET", "users/f47ac10b-58cc-0372-8567-0e02b2c3d479", "")
	// Get unknown user connections where Id is valid but user does not exist - Should be 404
	wrappedExecute("GET", "users/f47ac10b-58cc-0372-8567-0e02b2c3d479/connections", "")
	// Create connections for invalid user Id - Should be 400
	wrappedExecute("PUT", "users/f47a/connections/f47ac10b-58cc-0372-8567-0e02b2c3d479", "")
	// Create connections for unknow user Id - Should be 404
	wrappedExecute("PUT", "users/f47ac10b-58cc-0372-8567-0e02b2c3d479/connections/a47ac10b-58cc-0372-8567-0e02b2c3d470", "")
	// Unknown route - used to cause echo package panic - is now fixed
	wrappedExecute("PUT", "users/f47ac10b-58cc-0372-8567-0e02b2c3d479/connections/a47ac10b-58cc-0372-8567-0e02b2c3d470/aaaaaa", "")
}

func getAllUsers() executionResult {
	subURL := "users"
	return wrappedExecute("GET", subURL, "")
}

func createNewUser(name string) executionResult {
	fmt.Printf("\n\n*** Create user with name: %s\n", name)

	subURL := "users"
	jsonPayload := `{"name":"` + name + `"}`
	result := wrappedExecute("POST", subURL, jsonPayload)

	userID := extractUserIDFromLocationHeader(result.LocationHeader)
	fmt.Printf("*** Created user with Id: %s\n", userID)

	return result
}

func getUser(userID string) executionResult {
	subURL := "users/" + userID
	return wrappedExecute("GET", subURL, "")
}

func getUserConnections(userID string) executionResult {
	subURL := "users/" + userID + "/connections"
	return wrappedExecute("GET", subURL, "")
}

func createConnection(id1, id2 string) executionResult {
	subURL := "users/" + id1 + "/connections/" + id2
	return wrappedExecute("PUT", subURL, "")
}

func deleteConnection(id1, id2 string) executionResult {
	subURL := "users/" + id1 + "/connections/" + id2
	return wrappedExecute("DELETE", subURL, "")
}

func extractUserIDFromLocationHeader(header string) string {
	return strings.TrimPrefix(header, "/users/")
}
