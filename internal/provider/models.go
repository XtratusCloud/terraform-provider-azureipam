package provider

import (
	ipamclient "terraform-provider-azureipam/ipamclient"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

//shared models

// spaceModel maps space schema data.
type spaceModel struct {
	Name        types.String  `tfsdk:"name"`
	Description types.String  `tfsdk:"description"`
	Blocks      []blockModel  `tfsdk:"blocks"`
	Size        types.Float64 `tfsdk:"size"`
	Used        types.Float64 `tfsdk:"used"`
}

// blockModel maps block schema data.
type blockModel struct {
	Name         types.String            `tfsdk:"name"`
	Cidr         types.String            `tfsdk:"cidr"`
	Vnets        []vnetModel             `tfsdk:"vnets"`
	Externals    []externalModel         `tfsdk:"externals"`
	Reservations []reservationModel `tfsdk:"reservations"`
	Size         types.Float64           `tfsdk:"size"`
	Used         types.Float64           `tfsdk:"used"`
}

// vnetModel maps vnet schema data.
type vnetModel struct {
	Name           types.String   `tfsdk:"name"`
	Id             types.String   `tfsdk:"id"`
	Prefixes       []types.String `tfsdk:"prefixes"`
	Subnets        []subnetModel  `tfsdk:"subnets"`
	ResourceGroup  types.String   `tfsdk:"resource_group"`
	SubscriptionId types.String   `tfsdk:"subscription_id"`
	TenantId       types.String   `tfsdk:"tenant_id"`
	Size           types.Float64  `tfsdk:"size"`
	Used           types.Float64  `tfsdk:"used"`
}

// subnetModel maps subnet schema data.
type subnetModel struct {
	Name   types.String  `tfsdk:"name"`
	Prefix types.String  `tfsdk:"prefix"`
	Size   types.Float64 `tfsdk:"size"`
	Used   types.Float64 `tfsdk:"used"`
}

// externalModel maps external schema data.
type externalModel struct {
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Cidr        types.String `tfsdk:"cidr"`
}

// reservationModel maps Reservation schema data.
type reservationModel struct {
	Id          types.String      `tfsdk:"id"`
	Cidr        types.String      `tfsdk:"cidr"`
	Description types.String      `tfsdk:"description"`
	CreatedOn   timetypes.RFC3339 `tfsdk:"created_on"`
	CreatedBy   types.String      `tfsdk:"created_by"`
	SettledOn   timetypes.RFC3339 `tfsdk:"settled_on"`
	SettledBy   types.String      `tfsdk:"settled_by"`
	Status      types.String      `tfsdk:"status"`
}



//shared map functions 
func flattenSpaceInfo(space *ipamclient.SpaceInfo) spaceModel {
	var model spaceModel

	model.Name = types.StringValue(space.Name)
	model.Description = types.StringValue(space.Description)
	for _, block := range space.Blocks {
		model.Blocks = append(model.Blocks, flattenBlockInfo(&block))
	}
	if space.Size == nil {
		model.Size = types.Float64Null()
	} else {
		model.Size = types.Float64Value(*space.Size)
	}
	if space.Used == nil {
		model.Used = types.Float64Null()
	} else {
		model.Used = types.Float64Value(*space.Used)
	}

	return model
}

func flattenBlockInfo(block *ipamclient.BlockInfo) blockModel {
	var model blockModel

	model.Name = types.StringValue(block.Name)
	model.Cidr = types.StringValue(block.Cidr)
	for _, vnet := range block.Vnets {
		model.Vnets = append(model.Vnets, flattenVnetInfo(&vnet))
	}
	for _, external := range block.Externals {
		model.Externals = append(model.Externals, flattenExternalInfo(&external))
	}
	for _, reservation := range block.Reservations {
		model.Reservations = append(model.Reservations, flattenReservationInfo(&reservation))
	}
	if block.Size == nil {
		model.Size = types.Float64Null()
	} else {
		model.Size = types.Float64Value(*block.Size)
	}
	if block.Used == nil {
		model.Used = types.Float64Null()
	} else {
		model.Used = types.Float64Value(*block.Used)
	}

	return model
}

func flattenVnetInfo(vnet *ipamclient.VnetInfo) vnetModel {
	var model vnetModel

	model.Id = types.StringValue(vnet.Id)
	if vnet.Name == nil {
		model.Name = types.StringNull()
	} else {
		model.Name = types.StringValue(*vnet.Name)
	}
	for _, prefix := range vnet.Prefixes {
		model.Prefixes = append(model.Prefixes, types.StringValue(prefix))
	}
	for _, subnet := range vnet.Subnets {
		model.Subnets = append(model.Subnets, flattenSubnetInfo(&subnet))
	}
	if vnet.ResourceGroup == nil {
		model.ResourceGroup = types.StringNull()
	} else {
		model.ResourceGroup = types.StringValue(*vnet.ResourceGroup)
	}
	if vnet.SubscriptionId == nil {
		model.SubscriptionId = types.StringNull()
	} else {
		model.SubscriptionId = types.StringValue(*vnet.SubscriptionId)
	}
	if vnet.TenantId == nil {
		model.TenantId = types.StringNull()
	} else {
		model.TenantId = types.StringValue(*vnet.TenantId)
	}
	if vnet.Size == nil {
		model.Size = types.Float64Null()
	} else {
		model.Size = types.Float64Value(*vnet.Size)
	}
	if vnet.Used == nil {
		model.Used = types.Float64Null()
	} else {
		model.Used = types.Float64Value(*vnet.Used)
	}

	return model
}

func flattenSubnetInfo(subnet *ipamclient.SubnetInfo) subnetModel {
	var model subnetModel

	model.Name = types.StringValue(subnet.Name)
	model.Prefix = types.StringValue(subnet.Prefix)
	if subnet.Size == nil {
		model.Size = types.Float64Null()
	} else {
		model.Size = types.Float64Value(*subnet.Size)
	}
	if subnet.Used == nil {
		model.Used = types.Float64Null()
	} else {
		model.Used = types.Float64Value(*subnet.Used)
	}

	return model
}

func flattenExternalInfo(external *ipamclient.ExternalInfo) externalModel {
	var model externalModel

	model.Name = types.StringValue(external.Name)
	model.Description = types.StringValue(external.Description)
	model.Cidr = types.StringValue(external.Cidr)

	return model
}

func flattenReservationInfo(reservation *ipamclient.ReservationInfo) reservationModel {
	var model reservationModel

	model.Id = types.StringValue(reservation.Id)
	model.Cidr = types.StringValue(reservation.Cidr)
	model.Description = types.StringValue(reservation.Description)
	model.CreatedOn = timetypes.NewRFC3339TimeValue(time.Unix(int64(reservation.CreatedOn), 0))
	model.CreatedBy = types.StringValue(reservation.CreatedBy)
	if reservation.SettledOn == nil {
		model.SettledOn = timetypes.NewRFC3339Null()
	} else {
		model.SettledOn = timetypes.NewRFC3339TimeValue(time.Unix(int64(*reservation.SettledOn), 0))
	}
	if reservation.SettledBy == nil {
		model.SettledBy = types.StringNull()
	} else {
		model.SettledBy = types.StringValue(*reservation.SettledBy)

	}
	model.Status = types.StringValue(reservation.Status)

	return model
}
