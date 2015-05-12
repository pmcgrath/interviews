package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/afex/hystrix-go/hystrix"
)

type executionResult struct {
	StatusCode     int
	LocationHeader string
	Body           string
}

func init() {
	// This config is arbitary - would need to be much more aware of service SLAs etc set these
	hystrix.ConfigureCommand("makeHTTPRequest", hystrix.CommandConfig{
		Timeout:               100,
		MaxConcurrentRequests: 10,
		ErrorPercentThreshold: 5,
	})
}

func wrappedExecute(method, path, payload string) executionResult {
	url := *baseURL + path

	fmt.Printf("\t> %s on %s\n", method, url)

	result, err := execute(*userName, *password, method, url, payload)
	if err != nil {
		fmt.Printf("\t> ERROR DETECTED: %s\n", err)
		return result
	}

	fmt.Printf("\t> STATUSCODE = %d\n", result.StatusCode)
	if result.Body != "" {
		// Assumes it is json - could have checked the Content-Type header to ensure
		fmt.Printf("\t> BODY\n")
		var ppBuf bytes.Buffer
		json.Indent(&ppBuf, []byte(result.Body), "", "  ")
		fmt.Println(string(ppBuf.Bytes()))
	}

	if *pauseAfterAPICall {
		fmt.Printf("\n\n? Hit enter to continue to next api call")
		var dummy string
		fmt.Scanf("%s", &dummy)
	}

	return result
}

func execute(userName, password, method, url, payload string) (result executionResult, err error) {
	result = executionResult{}

	var payloadBytes []byte
	if payload != "" {
		payloadBytes = []byte(payload)
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return
	}

	req.SetBasicAuth(userName, password)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	// Didn't bother with cert file here, see integration tests for cert file usage
	tlsConfig := &tls.Config{InsecureSkipVerify: true}
	transport := &http.Transport{TLSClientConfig: tlsConfig}
	client := &http.Client{Transport: transport, Timeout: 100 * time.Millisecond}

	resp, err := makeHTTPRequest(client, req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	result.StatusCode = resp.StatusCode
	if locationHeaderValues, ok := resp.Header["Location"]; ok {
		result.LocationHeader = locationHeaderValues[0]
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	result.Body = string(bodyBytes)

	return
}

func makeHTTPRequest(client *http.Client, req *http.Request) (*http.Response, error) {
	resultChan := make(chan *http.Response, 1)
	errors := hystrix.Go(
		"makeHTTPRequest",
		func() error {
			resp, err := client.Do(req)
			if err == nil {
				resultChan <- resp
				return nil
			}
			return err
		},
		func(err error) error {
			// Should wrap with specific error
			return err
		})

	select {
	case result := <-resultChan:
		return result, nil
	case err := <-errors:
		return nil, err
	}
}
