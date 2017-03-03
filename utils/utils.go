package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

// DecodeBody is used to JSON decode a body
func DecodeBody(resp *http.Response, out interface{}) error {
	dec := json.NewDecoder(resp.Body)
	return dec.Decode(out)
}

// EncodeBody is used to encode a request body
func EncodeBody(obj interface{}) (io.Reader, error) {
	buf := bytes.NewBuffer(nil)
	enc := json.NewEncoder(buf)
	if err := enc.Encode(obj); err != nil {
		return nil, err
	}
	return buf, nil
}

// GetStatus is a ...
func GetStatus(targetURL url.URL) (interface{}, error) {
	// Convert argument to REST call
	targetString := fmt.Sprintf("%s://%s/fru/api/about", targetURL.Scheme, targetURL.Host)

	// Send API call to validate that argument points to running server
	resp, err := http.Get(targetString)
	if err != nil {
		return nil, fmt.Errorf("Error sending API call: %s", err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Non-success status returned (%d): %s", resp.StatusCode, resp.Status)
	}

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading response body: %s", err)

	}

	return respBytes, nil
}
