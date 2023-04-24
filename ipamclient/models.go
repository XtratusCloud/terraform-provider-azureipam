package azureipamclient

// Reservation -
type Reservation struct {
	Id          string            `json:"id,omitempty"`
	Space       string            `json:"space,omitempty"`
	Block       string            `json:"block,omitempty"`
	Cidr        string            `json:"cidr,omitempty"`
	Description string            `json:"desc,omitempty"`
	CreatedOn   float64           `json:"createdOn,omitempty"`
	CreatedBy   string            `json:"createdBy,omitempty"`
	SettledOn   float64           `json:"settledOn,omitempty"`
	SettledBy   string            `json:"settledBy,omitempty"`
	Status      string            `json:"status,omitempty"`
	Tags        map[string]string `json:"tag,omitempty"`
}

type reservationRequest struct {
	Size          int    `json:"size"`
	Description   string `json:"desc"`
	ReverseSearch bool   `json:"reverse_search"`
	SmallestCidr  bool   `json:"smallest_cidr"`
}
