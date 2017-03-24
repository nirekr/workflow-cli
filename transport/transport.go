package transport

import (
	"crypto/tls"
	"net/http"
	"net/url"

	cleanhttp "github.com/hashicorp/go-cleanhttp"
)

// NewClient returns an HTTP client with optional HTTPS
func NewClient(target string) (*http.Client, error) {

	url, err := url.Parse(target)
	if err != nil {
		return nil, err
	}

	tr := cleanhttp.DefaultTransport()

	if url.Scheme == "https" {
		tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	client := &http.Client{
		Transport: tr,
	}

	return client, nil
}
