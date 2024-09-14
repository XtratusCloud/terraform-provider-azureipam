package azureipamclient

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// internal Models
type reservationRequest struct {
	Blocks        []string `json:"blocks"`
	Size          int      `json:"size"`
	Description   *string  `json:"desc"`
	ReverseSearch bool     `json:"reverse_search"`
	SmallestCidr  bool     `json:"smallest_cidr"`
}

// GetReservations - Returns all existing reservations by space and block
func (c *Client) GetReservations(space, block string, includeSettled bool) (*[]Reservation, error) {
	//prepare request
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/spaces/%s/blocks/%s/reservations?settled=%t", c.HostURL, space, block, includeSettled), nil)
	if err != nil {
		return nil, err
	}
	response, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	//process response
	var reservationsInfo []Reservation
	err = json.Unmarshal(response, &reservationsInfo)
	if err != nil {
		return nil, err
	}

	return &reservationsInfo, nil
}

// GetReservation - Search for a specifc reservation ID iterating spaces and blocks
func (c *Client) FindReservationById(id string) (*Reservation, error) {

	//read all reservations
	spaces, err := c.GetSpaces(false, false)
	if err != nil {
		return nil, err
	}

	//find the reservation by Id, and returns
	for _, space := range *spaces {
		for _, block := range space.Blocks {
			for _, reservation := range block.Reservations {
				if reservation.Id == id {
					return c.GetReservation(space.Name, block.Name, reservation.Id)
				}
			}
		}
	}

	//not found -> Error
	return nil, fmt.Errorf("Reservation not found: %s", id)
}

// GetReservation - Search for a specifc reservation ID iterating spaces and blocks
func (c *Client) GetReservation(space string, block string, id string) (*Reservation, error) {

	//read all reservations
	reservationsInfo, err := c.GetReservations(space, block, true)
	if err != nil {
		return nil, err
	}

	//find the reservation by Id, and returns
	for _, reservation := range *reservationsInfo {
		if reservation.Id == id {
			return &reservation, nil
		}
	}

	//not found -> Error
	return nil, fmt.Errorf("Reservation not found: %s", id)
}

// CreateReservation - Create new reservation
func (c *Client) CreateReservation(space string, blocks []string, description *string, size int, reverseSearch bool, smallestCidr bool) (*Reservation, error) {
	//construct body
	request := &reservationRequest{
		Blocks:        blocks,
		Size:          size,
		ReverseSearch: reverseSearch,
		SmallestCidr:  smallestCidr,
		Description:   description,
	}
	rb, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	//prepare request
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/spaces/%s/reservations", c.HostURL, space), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}
	response, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	//process response
	reservation := Reservation{}
	err = json.Unmarshal(response, &reservation)
	if err != nil {
		return nil, err
	}

	return &reservation, nil
}

// DeleteReservation - Deletes a reservation
func (c *Client) DeleteReservation(space, block, id string) error {
	//construct body
	request := [1]string{id}
	rb, err := json.Marshal(request)
	if err != nil {
		return err
	}

	//prepare request
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/api/spaces/%s/blocks/%s/reservations", c.HostURL, space, block), strings.NewReader(string(rb)))
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
