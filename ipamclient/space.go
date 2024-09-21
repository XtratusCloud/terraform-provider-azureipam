package azureipamclient

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

//internal Models
type spaceRequest struct {
	Name        string `json:"name"`
	Description string `json:"desc"`
}
type spaceUpdateRequest struct {
	Op    string `json:"op"`
	Path  string `json:"path"`
	Value string `json:"value"`
}

// GetSpaces - Returns all existing spaces
func (c *Client) GetSpaces(expand bool, appendUtilization bool) (*[]SpaceInfo, error) {
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
	var spacesInfo []SpaceInfo
	err = json.Unmarshal(response, &spacesInfo)
	if err != nil {
		return nil, err
	}

	return &spacesInfo, nil
}

// GetSpace - Returns a specifc space by name
func (c *Client) GetSpace(name string, expand bool, appendUtilization bool) (*SpaceInfo, error) {

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
	var spaceInfo SpaceInfo
	err = json.Unmarshal(response, &spaceInfo)
	if err != nil {
		return nil, err
	}

	return &spaceInfo, nil
}

// CreateSpace - Create new space
func (c *Client) CreateSpace(name string, description string) (*SpaceInfo, error) {

	//construct body
	request := &spaceRequest{
		Name:        name,
		Description: description,
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
	space := SpaceInfo{}
	err = json.Unmarshal(response, &space)
	if err != nil {
		return nil, err
	}

	return &space, nil
}

func (c *Client) UpdateSpace(name string, newName *string, newDescription *string) (*SpaceInfo, error) {

	//construct body
	var request = []spaceUpdateRequest{}
	if newName != nil {
		request = append(request, spaceUpdateRequest{
			Op:    "replace",
			Path:  "/name",
			Value: *newName,
		})
	}
	if newDescription != nil {
		request = append(request, spaceUpdateRequest{
			Op:    "replace",
			Path:  "/desc",
			Value: *newDescription,
		})
	}

	rb, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	//prepare request
	req, err := http.NewRequest("PATCH", fmt.Sprintf("%s/api/spaces/%s", c.HostURL, name), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}
	response, err := c.doRequest(req)
	if err != nil {
		return nil, errors.New(string(response) + ";inner error: " + err.Error())
	}

	//read and return the updated space
	retVal, err := c.GetSpace(*newName, false, false)
	if err != nil {
		return nil, err
	}
	return retVal, nil
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
