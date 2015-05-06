package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

var (
	baseUrl = "http://localhost:8090/"
)

func main() {
	fmt.Printf("\n\n*** Get all users\n")
	getAllUsers()

	// Create users
	userNames := [...]string{"Tom Toe", "JIM BOB", "Terry Smith", "Joe Dolan", "Jerry", "Billy", "Todd"}
	users := make(map[string]string, len(userNames))
	for _, userName := range userNames {
		fmt.Printf("\n\n*** Create user with name: %s\n", userName)
		userId := createNewUser(userName)
		fmt.Printf("*** Created user with id: %s\n", userId)

		fmt.Printf("\n\n*** Get user with id: %s \n", userId)
		getUser(userId)

		users[userName] = userId
	}

	// Create connections
	createConnection(users[userNames[0]], users[userNames[2]])
	createConnection(users[userNames[0]], users[userNames[5]])

	// Get all users
	fmt.Printf("\n\n*** Get all users\n")
	getAllUsers()

	// Delete connection
	// Show user first
	getUser(users[userNames[0]])
	// Delete connection
	deleteConnection(users[userNames[0]], users[userNames[5]])
	// Show user - connection gone
	getUser(users[userNames[0]])
}

func getAllUsers() {
	subUrl := "users"
	body, _ := execute("GET", subUrl, "")
	printJsonBody(body)
}

func createNewUser(name string) string {
	subUrl := "users"
	jsonPayload := `{"name":"` + name + `"}`
	_, locationHeader := execute("POST", subUrl, jsonPayload)

	return strings.TrimPrefix(locationHeader, "/users/")
}

func getUser(userId string) {
	subUrl := "users/" + userId
	body, _ := execute("GET", subUrl, "")
	printJsonBody(body)
}

func createConnection(id1, id2 string) {
	subUrl := "users/" + id1 + "/connections/" + id2
	execute("PUT", subUrl, "")
}

func deleteConnection(id1, id2 string) {
	subUrl := "users/" + id1 + "/connections/" + id2
	execute("DELETE", subUrl, "")
}

func execute(method, path, payload string) (body, locationHeader string) {
	url := baseUrl + path
	fmt.Printf("\t> %s on %s\n", method, url)

	payloadBytes := make([]byte, 0)
	if payload != "" {
		payloadBytes = []byte(payload)
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(payloadBytes))
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if locationHeaderValues, ok := resp.Header["Location"]; ok {
		locationHeader = locationHeaderValues[0]
	}

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	body = string(bodyBytes)

	fmt.Printf("\t> response Status: %s\n", resp.Status)
	fmt.Printf("\t> response Body: %v\n", body)

	return
}

func printJsonBody(jsonString string) {
	var ppBuf bytes.Buffer
	json.Indent(&ppBuf, []byte(jsonString), "", "  ")
	fmt.Println(string(ppBuf.Bytes()))
}
