package azureipamclient

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// GetSpaces - Returns all existing spaces
func (c *Client) GetSpaces(expand bool, appendUtilization bool) (*[]Space, error) {
	//prepare request
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/spaces?expand=%t&utilization=%t", c.HostURL, expand, appendUtilization), nil)
	if err != nil {
		return nil, err
	}
	response, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	//process response
	var spacesInfo []Space
	err = json.Unmarshal(response, &spacesInfo)
	if err != nil {
		return nil, err
	}

	return &spacesInfo, nil
}

// GetSpace - Returns a specifc space by name
func (c *Client) GetSpace(name string, expand bool, appendUtilization bool) (*Space, error) {

	//prepare request
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/spaces/%s?expand=%t&utilization=%t", c.HostURL, name, expand, appendUtilization), nil)
	if err != nil {
		return nil, err
	}
	response, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	//process response
	var spaceInfo Space
	err = json.Unmarshal(response, &spaceInfo)
	if err != nil {
		return nil, err
	}

	return &spaceInfo, nil
}

// CreateSpace - Create new space
func (c *Client) CreateSpace(name string, description string) (*Space, error) {

	//construct body
	request := &spaceRequest{
		Name: name,
	}
	if description != "" {
		request.Description = description
	}
	rb, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	//prepare request
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/spaces", c.HostURL), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}
	response, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	//process response
	space := Space{}
	err = json.Unmarshal(response, &space)
	if err != nil {
		return nil, err
	}

	return &space, nil
}

// DeleteSpace- Deletes a space
func (c *Client) DeleteSpace(name string, force bool) error {

	//prepare request
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/api/spaces/%s?force=%t", c.HostURL, name, force), nil)
	if err != nil {
		return err
	}
	response, err := c.doRequest(req)
	if err != nil {
		return err
	}

	//process response
	if string(response) != "" {
		return errors.New(string(response))
	}

	return nil
}
