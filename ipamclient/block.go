package azureipamclient

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// internal Models
type blockRequest struct {
	Name string `json:"name"`
	Cidr string `json:"cidr"`
}
type blockUpdateRequest struct {
	Op    string `json:"op"`
	Path  string `json:"path"`
	Value string `json:"value"`
}

// GetBlocks - Returns a list of all Blocks within a specific Space.
func (c *Client) GetBlocks(space string, expand bool, appendUtilization bool) (*[]BlockInfo, error) {
	//prepare request
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/spaces/%s/blocks?expand=%t&utilization=%t", c.HostURL, space, expand, appendUtilization), nil)
	if err != nil {
		return nil, err
	}
	response, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	//process response
	var blocksInfo []BlockInfo
	err = json.Unmarshal(response, &blocksInfo)
	if err != nil {
		return nil, err
	}

	return &blocksInfo, nil
}

// GetBlock - Returns a specifc block by name
func (c *Client) GetBlock(space string, name string, expand bool, appendUtilization bool) (*Block, error) {

	//prepare request
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/spaces/%s/blocks/%s?expand=%t&utilization=%t", c.HostURL, space, name, expand, appendUtilization), nil)
	if err != nil {
		return nil, err
	}
	response, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	//process response
	var block Block	
	err = json.Unmarshal(response, &block)
	if err != nil {
		return nil, err
	}
	//add attributes not included in response
	block.Space = space

	return &block, nil
}

// CreateBlock - Create new block
func (c *Client) CreateBlock(space string, name string, cidr string) (*Block, error) {

	//construct body
	request := &blockRequest{
		Name: name,
		Cidr: cidr,
	}
	rb, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	//prepare request
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/spaces/%s/blocks", c.HostURL, space), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}
	response, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	//process response
	block := Block{}
	err = json.Unmarshal(response, &block)
	if err != nil {
		return nil, err
	}
	//add attributes not included in response
	block.Space = space

	return &block, nil
}

func (c *Client) UpdateBlock(space string, name string, newName *string, newCidr *string) (*Block, error) {

	//construct body
	var request = []blockUpdateRequest{}
	if newName != nil {
		request = append(request, blockUpdateRequest{
			Op:    "replace",
			Path:  "/name",
			Value: *newName,
		})
	}
	if newCidr != nil {
		request = append(request, blockUpdateRequest{
			Op:    "replace",
			Path:  "/cidr",
			Value: *newCidr,
		})
	}

	rb, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	//prepare request
	req, err := http.NewRequest("PATCH", fmt.Sprintf("%s/api/spaces/%s/blocks/%s", c.HostURL, space, name), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}
	response, err := c.doRequest(req)
	if err != nil {
		return nil, errors.New(string(response) + ";inner error: " + err.Error())
	}

	//read and return the updated block
	retVal, err := c.GetBlock(space, *newName, false, false)
	if err != nil {
		return nil, err
	}
	return retVal, nil
}

// DeleteBlock- Deletes a block
func (c *Client) DeleteBlock(space string, name string, force bool) error {

	//prepare request
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/api/spaces/%s/blocks/%s?force=%t", c.HostURL, space, name, force), nil)
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
