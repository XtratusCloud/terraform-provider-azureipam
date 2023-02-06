package azureipamclient

// Reservation -
type Reservation struct {
	Id        string            `json:"id,omitempty"`
	Cidr      string            `json:"cidr,omitempty"`
	UserId    string            `json:"userId,omitempty"`
	CreatedOn float64           `json:"createdOn,omitempty"`
	Status    string            `json:"status,omitempty"`
	Tags      map[string]string `json:"tag,omitempty"`
}

type reservationRequest struct {
	Size int               `json:"size"`
	Tags map[string]string `json:"tags"`
}
