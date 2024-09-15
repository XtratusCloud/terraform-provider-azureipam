package provider

import (
	"context"
	"fmt"

	ipamclient "terraform-provider-azureipam/ipamclient"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &blockDataSource{}
	_ datasource.DataSourceWithConfigure = &blockDataSource{}
)

// NewBlockDataSource is a helper function to simplify the provider implementation.
func NewBlockDataSource() datasource.DataSource {
	return &blockDataSource{}
}

// blockDataSource is the data source implementation.
type blockDataSource struct {
	client *ipamclient.Client
}

// blockDataSourceModel maps the data source schema data.
type blockDataSourceModel struct {
	Space             types.String       `tfsdk:"space"`
	Expand            types.Bool         `tfsdk:"expand"`
	AppendUtilization types.Bool         `tfsdk:"append_utilization"`
	Name              types.String       `tfsdk:"name"`
	Cidr              types.String       `tfsdk:"cidr"`
	Vnets             []vnetModel        `tfsdk:"vnets"`
	Externals         []externalModel    `tfsdk:"externals"`
	Reservations      []reservationModel `tfsdk:"reservations"`
	Size              types.Float64      `tfsdk:"size"`
	Used              types.Float64      `tfsdk:"used"`
}

// Metadata returns the data source type name.
func (d *blockDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_block"
}

// Schema defines the schema for the data source.
func (d *blockDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "The block data source allows you to retrieve one specific block and their related information.",
		Attributes: map[string]schema.Attribute{
			"space": schema.StringAttribute{
				Description: "Name of the `space` for which to read the related `blocks`.",
				Required:    true,
			},
			"expand": schema.BoolAttribute{
				Description: "Indicates if network references to full network objects must be included.",
				Optional:    true,
			},
			"append_utilization": schema.BoolAttribute{
				Description: "Indicates if utilization information for each network must be included.",
				Optional:    true,
			},
			"name": schema.StringAttribute{
				Description: "Name of the block.",
				Required:    true,
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
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *blockDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state blockDataSourceModel

	// Read Terraform configuration state into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	block, err := d.client.GetBlockInfo(
		state.Space.ValueString(),
		state.Name.ValueString(),
		state.Expand.ValueBool(),
		state.AppendUtilization.ValueBool(),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read AzureIpam Block named: "+ state.Name.ValueString() + " is space: " + state.Space.ValueString(),
			err.Error(),
		)
		return
	}

	// Map response body to state model
	model := flattenBlockInfo(block)
	state.Cidr = model.Cidr
	state.Vnets = model.Vnets
	state.Externals = model.Externals
	state.Reservations = model.Reservations
	state.Size = model.Size
	state.Used = model.Used

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *blockDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
