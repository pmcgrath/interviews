// Should really be using cucumber style tests, rather than this client, but would require a ruby dependency which I can't ask for ?
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type executionResult struct {
	StatusCode     int
	LocationHeader string
	Body           string
}

func (r executionResult) ExtractUserIdFromLocationHeader() string {
	// This is pretty lazy - does not check that location header macthed a pattern ... but for now
	return strings.TrimPrefix(r.LocationHeader, "/users/")
}

var (
	baseUrl           = flag.String("baseUrl", "http://localhost:8090/", "Service base url")
	userName          = flag.String("userName", "ted", "Service user name for basic authentication")
	password          = flag.String("password", "toe", "Service user password for basic authentication")
	pauseAfterApiCall = flag.Bool("pauseAfterApiCall", false, "Hit enter to continue after each api call, if not set will execute all commands with no pauses")
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
		userId := res.ExtractUserIdFromLocationHeader()
		getUser(userId)
		users[userName] = userId
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
	subUrl := "users"
	return wrappedExecute("GET", subUrl, "")
}

func createNewUser(name string) executionResult {
	fmt.Printf("\n\n*** Create user with name: %s\n", name)

	subUrl := "users"
	jsonPayload := `{"name":"` + name + `"}`
	result := wrappedExecute("POST", subUrl, jsonPayload)

	userId := result.ExtractUserIdFromLocationHeader()
	fmt.Printf("*** Created user with Id: %s\n", userId)

	return result
}

func getUser(userId string) executionResult {
	subUrl := "users/" + userId
	return wrappedExecute("GET", subUrl, "")
}

func getUserConnections(userId string) executionResult {
	subUrl := "users/" + userId + "/connections"
	return wrappedExecute("GET", subUrl, "")
}

func createConnection(id1, id2 string) executionResult {
	subUrl := "users/" + id1 + "/connections/" + id2
	return wrappedExecute("PUT", subUrl, "")
}

func deleteConnection(id1, id2 string) executionResult {
	subUrl := "users/" + id1 + "/connections/" + id2
	return wrappedExecute("DELETE", subUrl, "")
}

func wrappedExecute(method, path, payload string) executionResult {
	url := *baseUrl + path

	fmt.Printf("\t> %s on %s\n", method, url)

	result := execute(*userName, *password, method, url, payload)

	fmt.Printf("\t> STATUSCODE = %d\n", result.StatusCode)
	if result.Body != "" {
		// Assumes it is json - could have checked the Content-Type header to ensure
		fmt.Printf("\t> BODY\n")
		var ppBuf bytes.Buffer
		json.Indent(&ppBuf, []byte(result.Body), "", "  ")
		fmt.Println(string(ppBuf.Bytes()))
	}

	if *pauseAfterApiCall {
		fmt.Printf("\n\n? Hit enter to continue to next api call")
		var dummy string
		fmt.Scanf("%s", &dummy)
	}

	return result
}

func execute(userName, password, method, url, payload string) executionResult {
	result := executionResult{}

	payloadBytes := make([]byte, 0)
	if payload != "" {
		payloadBytes = []byte(payload)
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(payloadBytes))
	req.SetBasicAuth(userName, password)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	result.StatusCode = resp.StatusCode
	if locationHeaderValues, ok := resp.Header["Location"]; ok {
		result.LocationHeader = locationHeaderValues[0]
	}
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	result.Body = string(bodyBytes)

	return result
}
