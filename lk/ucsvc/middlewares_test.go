package main

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/labstack/echo"
)

func TestBasicAuthViaMux(t *testing.T) {
	isAuthorized := false
	mux := echo.New()
	mux.Use(basicAuth)
	mux.Get("/users/1", func(c *echo.Context) *echo.HTTPError {
		isAuthorized = true
		return nil
	})

	credential := base64.StdEncoding.EncodeToString([]byte("ted:toe"))
	header := http.Header(map[string][]string{"Authorization": []string{"Basic " + credential}})
	res := httptest.NewRecorder()
	req := &http.Request{Method: "GET", URL: &url.URL{Path: "/users/1"}, Header: header}
	mux.ServeHTTP(res, req)

	if !isAuthorized {
		t.Error("!!!")
	}
}

func TestBasicAuth(t *testing.T) {
	testCases := []struct {
		Scheme         string
		Credentials    string
		ExpectedResult bool
		ExpectedError  error
	}{
		{"Basic", "ted:toe", true, nil},
		{"Unknown", "ted:toe", false, errNotAuthorized},
		{"Basic", "", false, errNotAuthorized},
		{"Basic", "ted", false, errNotAuthorized},
		{"Basic", "ted+toe", false, errNotAuthorized},
		{"Basic", "ted:", false, errNotAuthorizedInvalidCredentials},
	}
	for i, testCase := range testCases {
		isAuthorized, err := executeAuthentication(testCase.Scheme, testCase.Credentials)

		if isAuthorized != testCase.ExpectedResult {
			t.Errorf("Unexpected result for index %d", i)
		}
		if err != testCase.ExpectedError {
			t.Errorf("Unexpected error [%s] for index %d", err, i)
		}
	}
}

func executeAuthentication(scheme, credentials string) (bool, error) {
	isAuthorized := false
	basicAuthWrappedHandlerFunc := basicAuth(func(c *echo.Context) *echo.HTTPError {
		isAuthorized = true
		return nil
	})

	credential := base64.StdEncoding.EncodeToString([]byte(credentials))
	header := http.Header(map[string][]string{"Authorization": []string{scheme + " " + credential}})
	req := &http.Request{Method: "GET", URL: &url.URL{Path: "/users/1"}, Header: header}
	c := &echo.Context{Request: req}

	err := basicAuthWrappedHandlerFunc(c)
	if err != nil {
		return isAuthorized, err.Error
	}

	return isAuthorized, nil
}
