package azureipamclient

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// GetReservations - Returns all existing reservations by space and block
func (c *Client) GetReservations(space, block string, includeSettled bool) ([]Reservation, error) {
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

	return reservationsInfo, nil
}

// GetReservation - Returns a specifc reservation by space and block
func (c *Client) GetReservation(space, block, id string) (*Reservation, error) {

	//read all reservations
	reservationsInfo, err := c.GetReservations(space, block, true)
	if err != nil {
		return nil, err
	}

	//find the reservation by Id, and returns
	for _, reservation := range reservationsInfo {
		if reservation.Id == id {
			return &reservation, nil
		}
	}

	//not found -> Error
	return nil, fmt.Errorf("Reservation not found: %s", id)
}

// CreateReservation - Create new reservation
func (c *Client) CreateReservation(space, block, description string, size int, reverseSearch bool, smallestCidr bool) (*Reservation, error) {
	//construct body
	request := &reservationRequest{
		Size:          size,
		ReverseSearch: reverseSearch,
		SmallestCidr:  smallestCidr,
	}
	if description != "" {
		request.Description = description
	}
	rb, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	//prepare request
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/spaces/%s/blocks/%s/reservations", c.HostURL, space, block), strings.NewReader(string(rb)))
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
