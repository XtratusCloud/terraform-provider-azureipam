package azureipamclient

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strings"
)

// internal Models
type externalRequest struct {
	Name        string `json:"name"`
	Description string `json:"desc"`
	Cidr        string `json:"cidr"`
}

// GetExternalsInfo - Returns a list of all External Network within a specific Space and Block.
func (c *Client) GetExternalsInfo(space string, block string) (*[]ExternalInfo, error) {
	//prepare request
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/spaces/%s/blocks/%s/externals", c.HostURL, space, block), nil)
	if err != nil {
		return nil, err
	}
	response, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	//process response
	var externalsInfo []ExternalInfo
	err = json.Unmarshal(response, &externalsInfo)
	if err != nil {
		return nil, err
	}

	return &externalsInfo, nil
}

// GetExternal - Returns a specifc external network by space, block and name
func (c *Client) GetExternal(space string, block string, name string) (*External, error) {

	//prepare request
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/spaces/%s/blocks/%s/externals/%s", c.HostURL, space, block, name), nil)
	if err != nil {
		return nil, err
	}
	response, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	//process response
	var readed ExternalInfo
	err = json.Unmarshal(response, &readed)
	if err != nil {
		return nil, err
	}
	//add attributes not included in response
	ret := External{
		Space: space,
		Block: block,
		Name: readed.Name,
		Description: readed.Description,
		Cidr: readed.Cidr,
	} 

	return &ret, nil
}

// CreateExternal - Create new external network within a specific Space and Block.
func (c *Client) CreateExternal(space string, block string, name string, desc string, cidr string) (*External, error) {

	//construct body
	request := &externalRequest{
		Name:        name,
		Description: desc,
		Cidr:        cidr,
	}
	rb, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	//prepare request
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/spaces/%s/blocks/%s/externals", c.HostURL, space, block), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}
	response, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	//process response
	externals := []ExternalInfo{}
	err = json.Unmarshal(response, &externals)
	if err != nil {
		return nil, err
	}
	//Find the created external in the collection
	i := slices.IndexFunc(externals, func(e ExternalInfo) bool { return e.Name == name })
	created := externals[i]

	//Create return object
	ret := External{
		Space: space,
		Block: block,
		Name: created.Name,
		Description: created.Description,
		Cidr: created.Cidr,
	} 
	return &ret, nil
}

func (c *Client) UpdateExternal(space string, block string, name string, newName *string, newDescription *string, newCidr *string) (*External, error) {

	//Read all external network collection
	externals, err := c.GetExternalsInfo(space, block)
	if err != nil {
		return nil, err
	}
	//Iterate the readed collection and construct request body
	var request = []externalRequest{}
	for _, current := range *externals {
		if current.Name == name {
			request = append(request, externalRequest{
				Name: *newName,
				Description: *newDescription,
				Cidr: *newCidr,
			})
		} else {
			request = append(request, externalRequest{
				Name: current.Name,
				Description: current.Description,
				Cidr: current.Cidr,
			})
		}
	}
	rb, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	//prepare request
	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/api/spaces/%s/blocks/%s/externals", c.HostURL, space, block), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}
	response, err := c.doRequest(req)
	if err != nil {
		return nil, errors.New(string(response) + ";inner error: " + err.Error())
	}

	//read and return the updated external network
	retVal, err := c.GetExternal(space, block, *newName)
	if err != nil {
		return nil, err
	}
	return retVal, nil
}

// DeleteExternal- Deletes a external networl
func (c *Client) DeleteExternal(space string, block string, name string) error {

	//prepare request
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/api/spaces/%s/blocks/%s/externals/%s", c.HostURL, space, block, name), nil)
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
