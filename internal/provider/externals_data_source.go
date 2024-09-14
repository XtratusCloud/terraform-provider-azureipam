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
	_ datasource.DataSource              = &externalsDataSource{}
	_ datasource.DataSourceWithConfigure = &externalsDataSource{}
)

// NewExternalsDataSource is a helper function to simplify the provider implementation.
func NewExternalsDataSource() datasource.DataSource {
	return &externalsDataSource{}
}

// externalsDataSource is the data source implementation.
type externalsDataSource struct {
	client *ipamclient.Client
}

// externalsDataSourceModel maps the data source schema data.
type externalsDataSourceModel struct {
	Space     types.String      `tfsdk:"space"`
	Block     types.String      `tfsdk:"block"`
	Externals []externalModel `tfsdk:"externals"`
}

// Metadata returns the data source type name.
func (d *externalsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_externals"
}

// Schema defines the schema for the data source.
func (d *externalsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "The externals data source allows you to retrieve information about all external networks associated with a space and block.",
		Attributes: map[string]schema.Attribute{
			"space": schema.StringAttribute{
				Description: "Name of the space for which to search the associated `externals`.",
				Required:    true,
			},
			"block": schema.StringAttribute{
				Description: "Name of the block for which to search the associated `externals`.",
				Required:    true,
			},
			 
			"externals": schema.ListNestedAttribute{
				Description: "List containing the `externals` found.",
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
							Description: "The IP range configured in the external network, in cidr notation.",
							Computed:    true,
						},
						 
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *externalsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state externalsDataSourceModel

	// Read Terraform configuration state into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	externals, err := d.client.GetExternalsInfo(
		state.Space.ValueString(),
		state.Block.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read AzureIpam External Networks",
			err.Error(),
		)
		return
	}

	// Map response body to model
	for _, external := range *externals {
		state.Externals = append(state.Externals, flattenExternalInfo(&external))
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *externalsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
