package azureipamclient

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// internal Models
type blockNetworkRequest struct {
	Id     string `json:"id"`
	Active bool   `json:"active"`
}

//GetBlockNetworksAvailables - Return a list of the Azure resource ids virtual networks availables to be associated to the space and block specified
func (c *Client) GetBlockNetworksAvailables(space string, block string) (*[]string, error) {
	//prepare request
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/spaces/%s/blocks/%s/available", c.HostURL, space, block), nil)
	if err != nil {
		return nil, err
	}
	response, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	//process response
	var networkIds []string
	err = json.Unmarshal(response, &networkIds)
	if err != nil {
		return nil, err
	}

	return &networkIds, nil
}


// GetBlockNetworksInfo - Returns a list of all Block Networks within a specific Space and Block.
func (c *Client) GetBlockNetworksInfo(space string, block string, expand bool) (*[]BlockNetworkInfo, error) {
	//prepare request
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/spaces/%s/blocks/%s/networks?expand=%t", c.HostURL, space, block, expand), nil)
	if err != nil {
		return nil, err
	}
	response, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	//process response
	var externalsInfo []BlockNetworkInfo
	err = json.Unmarshal(response, &externalsInfo)
	if err != nil {
		return nil, err
	}

	return &externalsInfo, nil
}

// GetBlockNetworkInfo - Returns a specifc block network by space, block and id
func (c *Client) GetBlockNetworkInfo(space string, block string, id string, expand bool) (*BlockNetworkInfo, error) {

	//Read all external networks in a space and block
	networks, err := c.GetBlockNetworksInfo(space, block, expand)
	if err != nil {
		return nil, err
	}
	//find the block network by Id, and returns
	var networkInfo *BlockNetworkInfo
	for _, network := range *networks {
		if network.Id == id {
			networkInfo = &network
		}
	}
	if networkInfo == nil {
		return nil, errors.New("invalid block network id")
	}

	return networkInfo, nil
}

// CreateBlockNetwork - Create new block network within a specific Space and Block.
func (c *Client) CreateBlockNetwork(space string, block string, id string) (*BlockNetworkInfo, error) {

	//construct body
	request := &blockNetworkRequest{
		Id:     id,
		Active: true,
	}
	rb, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	//prepare request
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/spaces/%s/blocks/%s/networks", c.HostURL, space, block), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}
	response, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	//process response
	var returnedBlock BlockInfo
	err = json.Unmarshal(response, &returnedBlock)
	if err != nil {
		return nil, err
	}

	//Create return object
	ret,err := c.GetBlockNetworkInfo(space, block, id, true)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

 // DeleteBlockNetwork- Deletes a block network within a specific Space and Block.
func (c *Client) DeleteBlockNetwork(space string, block string, id string) error {
	
	//construct body
	request := []string {
		id,
	}	 
	rb, err := json.Marshal(request)
	if err != nil {
		return err
	}

	//prepare request
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/api/spaces/%s/blocks/%s/networks", c.HostURL, space, block), strings.NewReader(string(rb)))
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
