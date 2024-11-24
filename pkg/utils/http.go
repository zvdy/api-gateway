package utils

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// SendJSONRequest sends an HTTP request with a JSON payload and returns the response.
func SendJSONRequest(method, url string, payload interface{}, headers map[string]string) (*http.Response, error) {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{}
	return client.Do(req)
}

// ReadResponseBody reads the body of an HTTP response and returns it as a string.
func ReadResponseBody(resp *http.Response) (string, error) {
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

// ParseJSONResponse parses the JSON body of an HTTP response into the given target.
func ParseJSONResponse(resp *http.Response, target interface{}) error {
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(target)
}
