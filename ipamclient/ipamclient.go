package azureipamclient

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client -
type Client struct {
	HostURL    string
	HTTPClient *http.Client
	Token      string
}

// NewClient - Construct a new HTTP Client to interact with the APIM REST API
func NewClient(host, authToken *string, defaultHttpTransport bool) (*Client, error) {
	var tr http.RoundTripper
	if defaultHttpTransport {
		tr = http.DefaultTransport //Allow to set http.DefaultTransport to allow Acceptance tests with [jarcoal/httpmock](https://github.com/jarcoal/httpmock)
	} else {
		tr = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // By default skip tls certificate verification
		}
	}
	c := Client{
		HTTPClient: &http.Client{Timeout: 10 * time.Second, Transport: tr},
	}

	// set client values, if provided
	if host != nil {
		c.HostURL = *host
	}
	if authToken != nil {
		c.Token = *authToken
	}

	return &c, nil
}

// doRequest -
func (c *Client) doRequest(req *http.Request) ([]byte, error) {
	//perform request
	req.Header.Add("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token))
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	//read response body
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	//write error not StatusOK
	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated && res.StatusCode != http.StatusAccepted && res.StatusCode != http.StatusNoContent {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	return body, err
}
