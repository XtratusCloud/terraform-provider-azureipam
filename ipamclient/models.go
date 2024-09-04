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
	SettledOn   *float64          `json:"settledOn,omitempty"`
	SettledBy   *string           `json:"settledBy,omitempty"`
	Status      string            `json:"status,omitempty"`
	Tags        map[string]string `json:"tag,omitempty"`
}

type ReservationLite struct {
	Id          string  `json:"id,omitempty"`
	Cidr        string  `json:"cidr,omitempty"`
	Description string  `json:"desc,omitempty"`
	CreatedOn   float64 `json:"createdOn,omitempty"`
	CreatedBy   string  `json:"createdBy,omitempty"`
	SettledOn   float64 `json:"settledOn,omitempty"`
	SettledBy   string  `json:"settledBy,omitempty"`
	Status      string  `json:"status,omitempty"`
}

type reservationRequest struct {
	Size          int    `json:"size"`
	Description   string `json:"desc"`
	ReverseSearch bool   `json:"reverse_search"`
	SmallestCidr  bool   `json:"smallest_cidr"`
}

//Space
type Space struct {
	Name        string   `json:"name,omitempty"`
	Description string   `json:"desc,omitempty"`
	Blocks      []Block  `json:"blocks,omitempty"`
	Size        *float64 `json:"size,omitempty"`
	Used        *float64 `json:"used,omitempty"`
}

type spaceRequest struct {
	Name        string `json:"name"`
	Description string `json:"desc"`
}

//Block
type Block struct {
	Name         string            `json:"name,omitempty"`
	Cidr         string            `json:"cidr,omitempty"`
	Vnets        []Vnet            `json:"vnets,omitempty"`
	Externals    []External        `json:"externals,omitempty"`
	Reservations []ReservationLite `json:"resv,omitempty"`
	Size         *float64          `json:"size,omitempty"`
	Used         *float64          `json:"used,omitempty"`
}

//Vnet
type Vnet struct {
	Name           *string  `json:"name,omitempty"`
	Id             string   `json:"id,omitempty"`
	Prefixes       []string `json:"prefixes,omitempty"`
	Subnets        []Subnet `json:"subnets,omitempty"`
	ResourceGroup  *string  `json:"resource_group,omitempty"`
	SubscriptionId *string  `json:"subscription_id,omitempty"`
	TenantId       *string  `json:"tenant_id,omitempty"`
	Size           *float64 `json:"size,omitempty"`
	Used           *float64 `json:"used,omitempty"`
}

//Subnet
type Subnet struct {
	Name   string   `json:"name,omitempty"`
	Prefix string   `json:"prefix,omitempty"`
	Size   *float64 `json:"size,omitempty"`
	Used   *float64 `json:"used,omitempty"`
}

//External
type External struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"desc,omitempty"`
	Cidr        string `json:"cidr,omitempty"`
}
