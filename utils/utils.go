package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/dellemc-symphony/workflow-cli/transport"
	log "github.com/sirupsen/logrus"
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

// GetURL is a ...
func GetURL(targetURL url.URL) (interface{}, error) {
	client, err := transport.NewClient(targetURL.String())
	if err != nil {
		log.Fatalf(err.Error())
	}

	// Convert argument to REST call
	targetString := fmt.Sprintf("%s://%s/fru/api/%s", targetURL.Scheme, targetURL.Host, targetURL.Path)

	// Send API call to validate that argument points to running server
	req, err := http.NewRequest(http.MethodGet, targetString, nil)
	if err != nil {
		log.Warnf("%s", err)
	}

	resp, err := client.Do(req)
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
