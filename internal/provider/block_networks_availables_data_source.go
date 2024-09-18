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
	_ datasource.DataSource              = &blockNetworkAvailablesDataSource{}
	_ datasource.DataSourceWithConfigure = &blockNetworkAvailablesDataSource{}
)

// NewBlockNetworksAvailablesDataSource is a helper function to simplify the provider implementation.
func NewBlockNetworksAvailablesDataSource() datasource.DataSource {
	return &blockNetworkAvailablesDataSource{}
}

// blockNetworkAvailablesDataSource is the data source implementation.
type blockNetworkAvailablesDataSource struct {
	client *ipamclient.Client
}

// blockNetworkAvailablesDataSourceModel maps the data source schema data.
type blockNetworkAvailablesDataSourceModel struct {
	Space types.String   `tfsdk:"space"`
	Block types.String   `tfsdk:"block"`
	Ids   []types.String `tfsdk:"ids"`
}

// Metadata returns the data source type name.
func (d *blockNetworkAvailablesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_block_networks_availables"
}

// Schema defines the schema for the data source.
func (d *blockNetworkAvailablesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "The block network availables data source allows you to retrieve information of the azure virtual networks availables to be associated with a space and block.",
		Attributes: map[string]schema.Attribute{
			"space": schema.StringAttribute{
				Description: "Name of the space where to search for networks availables.",
				Required:    true,
			},
			"block": schema.StringAttribute{
				Description: "Name of the block where to search for networks availables.",
				Required:    true,
			},
			"ids": schema.ListAttribute{
				Description: "The list of of the Azure virtual networs resource Ids that can be associated to the block.",
				Computed:    true,
				ElementType: types.StringType,
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *blockNetworkAvailablesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state blockNetworkAvailablesDataSourceModel

	// Read Terraform configuration state into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	ids, err := d.client.GetBlockNetworksAvailables(
		state.Space.ValueString(),
		state.Block.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read AzureIpam Block Network Availables",
			err.Error(),
		)
		return
	}

	// Map response body to model
	for _, id := range *ids {
		state.Ids = append(state.Ids, types.StringValue(id))
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *blockNetworkAvailablesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
