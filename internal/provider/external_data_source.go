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
	_ datasource.DataSource              = &externalDataSource{}
	_ datasource.DataSourceWithConfigure = &externalDataSource{}
)

// NewExternalDataSource is a helper function to simplify the provider implementation.
func NewExternalDataSource() datasource.DataSource {
	return &externalDataSource{}
}

// externalDataSource is the data source implementation.
type externalDataSource struct {
	client *ipamclient.Client
}

// externalDataSourceModel maps the data source schema data.
type externalDataSourceModel struct {
	Space       types.String `tfsdk:"space"`
	Block       types.String `tfsdk:"block"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Cidr        types.String `tfsdk:"cidr"`
}

// Metadata returns the data source type name.
func (d *externalDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_external"
}

// Schema defines the schema for the data source.
func (d *externalDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "The externals data source allows you to retrieve information about one specific external network associated with a space and block.",
		Attributes: map[string]schema.Attribute{
			"space": schema.StringAttribute{
				Description: "Name of the space of the external network to search.",
				Required:    true,
			},
			"block": schema.StringAttribute{
				Description: "Name of the block of the external network to search.",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "Name of the external network to search.",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "Text that describes the external network.",
				Computed:    true,
			},
			"cidr": schema.StringAttribute{
				Description: "The IP range configured in the external network, in cidr notation.",
				Computed:    true,
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *externalDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state externalDataSourceModel

	// Read Terraform configuration state into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	external, err := d.client.GetExternalInfo(
		state.Space.ValueString(),
		state.Block.ValueString(),
		state.Name.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read AzureIpam External Network",
			err.Error(),
		)
		return
	}

	// Map response body to model
	state.Description = types.StringValue(external.Description)
	state.Cidr = types.StringValue(external.Cidr)

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *externalDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
