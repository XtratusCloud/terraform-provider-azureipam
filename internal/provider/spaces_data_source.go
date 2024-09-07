package provider

import (
	"context"
	"fmt"
	"time"

	ipamclient "terraform-provider-azureipam/ipamclient"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &spacesDataSource{}
	_ datasource.DataSourceWithConfigure = &spacesDataSource{}
)

// NewSpacesDataSource is a helper function to simplify the provider implementation.
func NewSpacesDataSource() datasource.DataSource {
	return &spacesDataSource{}
}

// spacesDataSource is the data source implementation.
type spacesDataSource struct {
	client *ipamclient.Client
}

// spacesDataSourceModel maps the data source schema data.
type spacesDataSourceModel struct {
	Expand            types.Bool   `tfsdk:"expand"`
	AppendUtilization types.Bool   `tfsdk:"append_utilization"`
	Spaces            []spaceModel `tfsdk:"spaces"`
}

// spaceModel maps spaces schema data.
type spaceModel struct {
	Name        types.String  `tfsdk:"name"`
	Description types.String  `tfsdk:"description"`
	Blocks      []blockModel  `tfsdk:"blocks"`
	Size        types.Float64 `tfsdk:"size"`
	Used        types.Float64 `tfsdk:"used"`
}

// blockModel maps blocks schema data.
type blockModel struct {
	Name         types.String            `tfsdk:"name"`
	Cidr         types.String            `tfsdk:"cidr"`
	Vnets        []vnetModel             `tfsdk:"vnets"`
	Externals    []externalModel         `tfsdk:"externals"`
	Reservations []reservationsLiteModel `tfsdk:"reservations"`
	Size         types.Float64           `tfsdk:"size"`
	Used         types.Float64           `tfsdk:"used"`
}

// vnetModel maps vnets schema data.
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

// subnetModel maps subnets  schema data.
type subnetModel struct {
	Name   types.String  `tfsdk:"name"`
	Prefix types.String  `tfsdk:"prefix"`
	Size   types.Float64 `tfsdk:"size"`
	Used   types.Float64 `tfsdk:"used"`
}

// externalModel maps externals schema data.
type externalModel struct {
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Cidr        types.String `tfsdk:"cidr"`
}

type reservationsLiteModel struct {
	Id          types.String      `tfsdk:"id"`
	Cidr        types.String      `tfsdk:"cidr"`
	Description types.String      `tfsdk:"description"`
	CreatedOn   timetypes.RFC3339 `tfsdk:"created_on"`
	CreatedBy   types.String      `tfsdk:"created_by"`
	SettledOn   timetypes.RFC3339 `tfsdk:"settled_on"`
	SettledBy   types.String      `tfsdk:"settled_by"`
	Status      types.String      `tfsdk:"status"`
}

// Metadata returns the data source type name.
func (d *spacesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_spaces"
}

// Schema defines the schema for the data source.
func (d *spacesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "The reservations data source allows you to retrieve information about all existing reservations in the specific space and block.",
		Attributes: map[string]schema.Attribute{
			"expand": schema.BoolAttribute{
				Description: "Indicates if network references to full network objects must be included.",
				Optional:    true,
			},
			"append_utilization": schema.BoolAttribute{
				Description: "Indicates if utilization information for each network must be included.",
				Optional:    true,
			},
			"spaces": schema.ListNestedAttribute{
				Description: "List containing the `spaces` found.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "Name of the space.",
							Computed:    true,
						},
						"description": schema.StringAttribute{
							Description: "Text that describes the space.",
							Computed:    true,
						},
						"blocks": schema.ListNestedAttribute{
							Description: "List containing the `blocks` included in this `space`.",
							Computed:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"name": schema.StringAttribute{
										Description: "Name of the block.",
										Computed:    true,
									},
									"cidr": schema.StringAttribute{
										Description: "The IPV4 range assigned to this block, in cidr notation.",
										Computed:    true,
									},
									"vnets": schema.ListNestedAttribute{
										Description: "List containing the `vnet` included in this `block`.",
										Computed:    true,
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"name": schema.StringAttribute{
													Description: "Name of the virtual network.",
													Computed:    true,
												},
												"id": schema.StringAttribute{
													Description: "Resourece Id of the virtual network.",
													Computed:    true,
												},
												"prefixes": schema.ListAttribute{
													Description: "The list of IPV4 prefixes assigned to this vnet, in cidr notation.",
													Computed:    true,
													ElementType: types.StringType,
												},
												"subnets": schema.ListNestedAttribute{
													Description: "List containing the `subnets` included in this `vnet`.",
													Computed:    true,
													NestedObject: schema.NestedAttributeObject{
														Attributes: map[string]schema.Attribute{
															"name": schema.StringAttribute{
																Description: "Name of the subnet.",
																Computed:    true,
															},
															"prefix": schema.StringAttribute{
																Description: "The IPV4 prefix assigned to this block, in cidr notation.",
																Computed:    true,
															},
															"size": schema.Float64Attribute{
																Description: "Total IP's allowed in the `subnet` by its size.",
																Computed:    true,
															},
															"used": schema.Float64Attribute{
																Description: "Assigned IP's in the `subnet`.",
																Computed:    true,
															},
														},
													},
												},
												"resource_group": schema.StringAttribute{
													Description: "Name of the resource group where the `vnet` is deployed.",
													Computed:    true,
												},
												"subscription_id": schema.StringAttribute{
													Description: "Id of the Azure subscription where the `vnet` is deployed.",
													Computed:    true,
												},
												"tenant_id": schema.StringAttribute{
													Description: "Id of the Azure tenant where the `vnet` is deployed.",
													Computed:    true,
												},
												"size": schema.Float64Attribute{
													Description: "Total IP's allowed in the `vnet` by its size.",
													Computed:    true,
												},
												"used": schema.Float64Attribute{
													Description: "Assigned IP's in the `vnet`.",
													Computed:    true,
												},
											},
										},
									},
									"externals": schema.ListNestedAttribute{
										Description: "List containing the `external networks` included in this `block`.",
										Computed:    true,
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"name": schema.StringAttribute{
													Description: "Name of the external network.",
													Computed:    true,
												},
												"description": schema.StringAttribute{
													Description: "Text that describes the external network.",
													Computed:    true,
												},
												"cidr": schema.StringAttribute{
													Description: "The IPV4 range reserved for the external network, in cidr notation.",
													Computed:    true,
												},
											},
										},
									},
									"reservations": schema.ListNestedAttribute{
										Description: "List containing the `reservations` included in this `block`.",
										Computed:    true,
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"id": schema.StringAttribute{
													Description: "The unique identifier of the reservation.",
													Computed:    true,
												},
												"cidr": schema.StringAttribute{
													Description: "The IPv4 range assigned to this reservation, in cidr notation.",
													Computed:    true,
												},
												"description": schema.StringAttribute{
													Description: "Text that describes the reservation.",
													Computed:    true,
												},
												"created_on": schema.StringAttribute{
													CustomType:  timetypes.RFC3339Type{},
													Description: "The date and time that the reservacion was created.",
													Computed:    true,
												},
												"created_by": schema.StringAttribute{
													Description: "Email or identification of user that created the reservation.",
													Computed:    true,
												},
												"settled_on": schema.StringAttribute{
													CustomType:  timetypes.RFC3339Type{},
													Description: "The date and time when the reservation was settled.",
													Computed:    true,
												},
												"settled_by": schema.StringAttribute{
													Description: "Email or identification of user that have settled the reservation.",
													Computed:    true,
												},
												"status": schema.StringAttribute{
													Description: "Status of the reservation, a 'wait' status indicates that is waiting for the related vnet creation",
													Computed:    true,
												},
											},
										},
									},
									"size": schema.Float64Attribute{
										Description: "Total IP's allowed in the `block` by its size.",
										Computed:    true,
									},
									"used": schema.Float64Attribute{
										Description: "Assigned IP's in the `block`.",
										Computed:    true,
									},
								},
							},
						},
						"size": schema.Float64Attribute{
							Description: "Total IP's allowed in the `space` by its size.",
							Computed:    true,
						},
						"used": schema.Float64Attribute{
							Description: "Assigned IP's in the `space`.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *spacesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state spacesDataSourceModel

	// Read Terraform configuration state into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	spaces, err := d.client.GetSpaces(
		state.Expand.ValueBool(),
		state.AppendUtilization.ValueBool(),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read AzureIpam Spaces",
			err.Error(),
		)
		return
	}

	// Map response body to model
	for _, space := range *spaces {
		state.Spaces = append(state.Spaces, flattenSpaceLite(&space))
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *spacesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*ipamclient.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *azureipam.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = client
}

func flattenSpaceLite(space *ipamclient.Space) spaceModel {
	var model spaceModel

	model.Name = types.StringValue(space.Name)
	model.Description = types.StringValue(space.Description)
	for _, block := range space.Blocks {
		model.Blocks = append(model.Blocks, flattenBlock(&block))
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

func flattenBlock(block *ipamclient.Block) blockModel {
	var model blockModel

	model.Name = types.StringValue(block.Name)
	model.Cidr = types.StringValue(block.Cidr)
	for _, vnet := range block.Vnets {
		model.Vnets = append(model.Vnets, flattenVnet(&vnet))
	}
	for _, external := range block.Externals {
		model.Externals = append(model.Externals, flattenExternal(&external))
	}
	for _, reservation := range block.Reservations {
		model.Reservations = append(model.Reservations, flattenReservationLite(&reservation))
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

func flattenVnet(vnet *ipamclient.Vnet) vnetModel {
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
		model.Subnets = append(model.Subnets, flattenSubnet(&subnet))
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

func flattenSubnet(subnet *ipamclient.Subnet) subnetModel {
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
func flattenExternal(external *ipamclient.External) externalModel {
	var model externalModel

	model.Name = types.StringValue(external.Name)
	model.Description = types.StringValue(external.Description)
	model.Cidr = types.StringValue(external.Cidr)

	return model
}

func flattenReservationLite(reservation *ipamclient.ReservationLite) reservationsLiteModel {
	var model reservationsLiteModel

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
