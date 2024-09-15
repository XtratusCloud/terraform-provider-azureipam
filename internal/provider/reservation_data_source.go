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
	_ datasource.DataSource              = &reservationDataSource{}
	_ datasource.DataSourceWithConfigure = &reservationDataSource{}
)

// NewReservationDataSource is a helper function to simplify the provider implementation.
func NewReservationDataSource() datasource.DataSource {
	return &reservationDataSource{}
}

// reservationDataSource is the data source implementation.
type reservationDataSource struct {
	client *ipamclient.Client
}

// reservationDataSourceModel maps the data source schema data.
type reservationDataSourceModel struct {
	Space       types.String      `tfsdk:"space"`
	Block       types.String      `tfsdk:"block"`
	Id          types.String      `tfsdk:"id"`
	Cidr        types.String      `tfsdk:"cidr"`
	Description types.String      `tfsdk:"description"`
	CreatedOn   timetypes.RFC3339 `tfsdk:"created_on"`
	CreatedBy   types.String      `tfsdk:"created_by"`
	SettledOn   timetypes.RFC3339 `tfsdk:"settled_on"`
	SettledBy   types.String      `tfsdk:"settled_by"`
	Status      types.String      `tfsdk:"status"`
	Tags        types.Map         `tfsdk:"tags"`
}

// Metadata returns the data source type name.
func (d *reservationDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_reservation"
}

// Schema defines the schema for the data source.
func (d *reservationDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "The reservation data source allows you to retrieve a specific reservation by id in the specified space and block.",
		Attributes: map[string]schema.Attribute{
			"space": schema.StringAttribute{
				Description: "Name of the space where the reservation is allocated.",
				Required:    true,
			},
			"block": schema.StringAttribute{
				Description: "Name of the  block where the reservation is allocated.",
				Required:    true,
			},
			"id": schema.StringAttribute{
				Description: "The unique identifier of the reservation.",
				Required:    true,
			},
			"cidr": schema.StringAttribute{
				Description: "The assigned and reserved range, in cidr notation.",
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "Description text that describe the reservation.",
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
				Description: "Status of the reservation, a 'wait' status indicates that is waiting for the related vnet creation.",
				Computed:    true,
			},
			"tags": schema.MapAttribute{
				Description: "Auto-generated tags for the reservation. Particular relevance the 'X-IPAM-RES-ID' tag, since it must be included in the vnet creation in order that the IPAM solution automatically considers the reservation as completed.",
				Computed:    true,
				ElementType: types.StringType,
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *reservationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state reservationDataSourceModel

	// Read Terraform configuration state into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	reservation, err := d.client.GetReservation(
		state.Space.ValueString(),
		state.Block.ValueString(),
		state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read AzureIpam Reservations",
			err.Error(),
		)
		return
	}

	//  Map response body to state model
	var model reservationResourceModel //to reuse existing flatten method
	flattenReservation(reservation,&model) 
	state.Id = model.Id
	state.Cidr = model.Cidr
	state.Description = model.Description
	state.CreatedOn = model.CreatedOn
	state.CreatedBy = model.CreatedBy
	state.SettledOn = model.SettledOn 
	state.SettledBy = model.SettledBy
	state.Status = model.Status 
	state.Tags, _ = types.MapValueFrom(ctx, types.StringType, reservation.Tags)

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *reservationDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
