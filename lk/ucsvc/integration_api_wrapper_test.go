// +build integration

// See http://peter.bourgon.org/go-in-production/#testing-and-validation for seperating integration tests
// Test with : 	go test -tags integration -v
//		go test -tags "test integration" -v
package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
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
	baseUrl         = "https://localhost:8090/"
	userName        = "ted"
	password        = "toe"
	jsonContentType = "application/json"
	emptyPayload    = make([]byte, 0)
	validUserId     = string(newID())
	certFilePath    = "../certs/server.crt"
	timeout         = time.Duration(50 * time.Millisecond)
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

	url := baseUrl + subUrl

	req, err := http.NewRequest(method, url, bytes.NewBuffer(payload))
	req.SetBasicAuth(userName, password)
	req.Header.Set("Accept", contentType)
	req.Header.Set("Content-Type", contentType)

	// Deal with https certificate, could have been lazy and just set TLS config's InsecureSkipVerify
	roots := x509.NewCertPool()
	if strings.HasPrefix(url, "https") {
		certFileData, err := ioutil.ReadFile(certFilePath)
		if err != nil {
			panic(err)
		}
		if !roots.AppendCertsFromPEM(certFileData) {
			panic("Couldn't load certificate file")
		}
	}

	tlsConfig := &tls.Config{RootCAs: roots}
	transport := &http.Transport{TLSClientConfig: tlsConfig}
	client := &http.Client{Transport: transport, Timeout: timeout}

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
