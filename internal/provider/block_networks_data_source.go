package provider

import (
	"context"
	"fmt"

	ipamclient "terraform-provider-azureipam/ipamclient"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &blockNetworksDataSource{}
	_ datasource.DataSourceWithConfigure = &blockNetworksDataSource{}
)

// NewBlockNetworksDataSource is a helper function to simplify the provider implementation.
func NewBlockNetworksDataSource() datasource.DataSource {
	return &blockNetworksDataSource{}
}

// blockNetworksDataSource is the data source implementation.
type blockNetworksDataSource struct {
	client *ipamclient.Client
}

// blockNetworksDataSourceModel maps the data source schema data.
type blockNetworksDataSourceModel struct {
	Space    types.String        `tfsdk:"space"`
	Block    types.String        `tfsdk:"block"`
	Networks []blockNetworkModel `tfsdk:"networks"`
}

// Metadata returns the data source type name.
func (d *blockNetworksDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_block_networks"
}

// Schema defines the schema for the data source.
func (d *blockNetworksDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "The block networks data source allows you to retrieve information of the azure virtual networks actively associated to the specified space and block.",
		Attributes: map[string]schema.Attribute{
			"space": schema.StringAttribute{
				Description: "Name of the space where to search the associated block networks.",
				Required:    true,
			},
			"block": schema.StringAttribute{
				Description: "Name of the block where to search the associated block networks.",
				Required:    true,
			},
			"networks": schema.ListNestedAttribute{
				Description: "List containing the `vnet` included in this `block`.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "Azure Resource ID of the virtual network already associated.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "Name of the Azure virtual network.",
							Computed:    true,
						},
						"prefixes": schema.ListAttribute{
							Description: "The list of IPV4 prefixes assigned to this vnet, in cidr notation.",
							Computed:    true,
							ElementType: types.StringType,
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
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *blockNetworksDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state blockNetworksDataSourceModel

	// Read Terraform configuration state into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	blockNetworks, err := d.client.GetBlockNetworksInfo(
		state.Space.ValueString(),
		state.Block.ValueString(),
		true,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read AzureIpam Block Networks",
			err.Error(),
		)
		return
	}

	// Map response body to model
	for _, blockNetwork := range *blockNetworks {
		info, diagFlattten := flattenBlockNetworkInfo(ctx, &blockNetwork)
		//Append current block network to schema collection
		state.Networks = append(state.Networks, info)
		//Append diagnostics information
		resp.Diagnostics.Append(diagFlattten...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *blockNetworksDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
