package azureipamclient

// Shared models

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

//SpaceInfo
type SpaceInfo struct {
	Name        string      `json:"name,omitempty"`
	Description string      `json:"desc,omitempty"`
	Blocks      []BlockInfo `json:"blocks,omitempty"`
	Size        *float64    `json:"size,omitempty"`
	Used        *float64    `json:"used,omitempty"`
}

//BlockInfo
type BlockInfo struct {
	Name         string            `json:"name,omitempty"`
	Cidr         string            `json:"cidr,omitempty"`
	Vnets        []VnetInfo        `json:"vnets,omitempty"`
	Externals    []ExternalInfo    `json:"externals,omitempty"`
	Reservations []ReservationInfo `json:"resv,omitempty"`
	Size         *float64          `json:"size,omitempty"`
	Used         *float64          `json:"used,omitempty"`
}

//VnetInfo
type VnetInfo struct {
	Name           *string      `json:"name,omitempty"`
	Id             string       `json:"id,omitempty"`
	Prefixes       []string     `json:"prefixes,omitempty"`
	Subnets        []SubnetInfo `json:"subnets,omitempty"`
	ResourceGroup  *string      `json:"resource_group,omitempty"`
	SubscriptionId *string      `json:"subscription_id,omitempty"`
	TenantId       *string      `json:"tenant_id,omitempty"`
	Size           *float64     `json:"size,omitempty"`
	Used           *float64     `json:"used,omitempty"`
}

//SubnetInfo
type SubnetInfo struct {
	Name   string   `json:"name,omitempty"`
	Prefix string   `json:"prefix,omitempty"`
	Size   *float64 `json:"size,omitempty"`
	Used   *float64 `json:"used,omitempty"`
}

//ExternalInfo
type ExternalInfo struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"desc,omitempty"`
	Cidr        string `json:"cidr,omitempty"`
}

//ReservationInfo
type ReservationInfo struct {
	Id          string   `json:"id,omitempty"`
	Cidr        string   `json:"cidr,omitempty"`
	Description string   `json:"desc,omitempty"`
	CreatedOn   float64  `json:"createdOn,omitempty"`
	CreatedBy   string   `json:"createdBy,omitempty"`
	SettledOn   *float64 `json:"settledOn,omitempty"`
	SettledBy   *string  `json:"settledBy,omitempty"`
	Status      string   `json:"status,omitempty"`
}

//Block
type Block struct {
	Name  string `json:"name,omitempty"`
	Space string `json:"space,omitempty"`
	Cidr  string `json:"cidr,omitempty"`
}

//External
type External struct {
	Space       string `json:"space,omitempty"`
	Block       string `json:"block,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"desc,omitempty"`
	Cidr        string `json:"cidr,omitempty"`
}

//BlockNetworkInfo
type BlockNetworkInfo struct {
	Name           string   `json:"name,omitempty"`
	Id             string   `json:"id,omitempty"`
	Prefixes       []string `json:"prefixes,omitempty"`
	ResourceGroup  *string  `json:"resource_group,omitempty"`
	SubscriptionId *string  `json:"subscription_id,omitempty"`
	TenantId       *string  `json:"tenant_id,omitempty"`
}
