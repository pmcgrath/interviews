// +build integration

// See http://peter.bourgon.org/go-in-production/#testing-and-validation for seperating integration tests
// Test with : 	go test -tags integration -v
//		go test -tags "test integration" -v
package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

type httpExecutionResult struct {
	StatusCode     int
	LocationHeader string
	Body           []byte
	Error          error
}

type testUser map[string]interface{}

type apiResult struct {
	httpExecutionResult
	NewUserId string
	Users     []testUser
	User      testUser
}

var (
	baseUrl         = "http://localhost:8090/"
	userName        = "ted"
	password        = "toe"
	jsonContentType = "application/json"
	emptyPayload    = make([]byte, 0)
	validUserId     = string(newID())
)

func getAllUsersViaApi() apiResult {
	subUrl := "users"
	result := apiResult{httpExecutionResult: executeSimpleHttp("GET", subUrl)}
	if result.Error == nil {
		result.Error = json.Unmarshal(result.Body, &result.Users)
	}
	return result
}

func createNewUserViaApi(name string) apiResult {
	subUrl := "users"
	payload := []byte(`{"name":"` + name + `"}`)
	result := apiResult{httpExecutionResult: executeHttp(userName, password, "POST", subUrl, jsonContentType, payload)}
	result.NewUserId = strings.TrimPrefix(result.LocationHeader, "/users/")
	return result
}

func getUserViaApi(userId string) apiResult {
	subUrl := "users/" + userId
	result := apiResult{httpExecutionResult: executeSimpleHttp("GET", subUrl)}
	if result.Error == nil {
		result.Error = json.Unmarshal(result.Body, &result.User)
	}
	return result
}

func getUserConnectionsViaApi(userId string) apiResult {
	subUrl := "users/" + userId + "/connections"
	result := apiResult{httpExecutionResult: executeSimpleHttp("GET", subUrl)}
	if result.Error == nil {
		result.Error = json.Unmarshal(result.Body, &result.Users)
	}
	return result
}

func createUserConnectionsViaApi(id1, id2 string) apiResult {
	subUrl := "users/" + id1 + "/connections/" + id2
	result := apiResult{httpExecutionResult: executeSimpleHttp("PUT", subUrl)}
	return result
}

func deleteUserConnectionsViaApi(id1, id2 string) apiResult {
	subUrl := "users/" + id1 + "/connections/" + id2
	result := apiResult{httpExecutionResult: executeSimpleHttp("DELETE", subUrl)}
	return result
}

func executeSimpleHttp(method, subUrl string) httpExecutionResult {
	return executeHttp(userName, password, method, subUrl, jsonContentType, emptyPayload)
}

func executeHttp(userName, password, method, subUrl, contentType string, payload []byte) httpExecutionResult {
	result := httpExecutionResult{}

	credential := base64.StdEncoding.EncodeToString([]byte(userName + ":" + password))
	authorization := "Basic " + credential

	url := baseUrl + subUrl

	req, err := http.NewRequest(method, url, bytes.NewBuffer(payload))
	req.Header.Set("Accept", contentType)
	req.Header.Set("Authorization", authorization)
	req.Header.Set("Content-Type", contentType)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		result.Error = err
		return result
	}
	defer resp.Body.Close()

	result.StatusCode = resp.StatusCode
	if locationHeaderValues, ok := resp.Header["Location"]; ok {
		result.LocationHeader = locationHeaderValues[0]
	}
	result.Body, result.Error = ioutil.ReadAll(resp.Body)

	return result
}
