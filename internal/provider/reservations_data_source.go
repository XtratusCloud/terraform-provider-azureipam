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
	_ datasource.DataSource              = &reservationsDataSource{}
	_ datasource.DataSourceWithConfigure = &reservationsDataSource{}
)

// NewReservationsDataSource is a helper function to simplify the provider implementation.
func NewReservationsDataSource() datasource.DataSource {
	return &reservationsDataSource{}
}

// reservationsDataSource is the data source implementation.
type reservationsDataSource struct {
	client *ipamclient.Client
}

// reservationsDataSourceModel maps the data source schema data.
type reservationsDataSourceModel struct {
	Space          types.String        `tfsdk:"space"`
	Block          types.String        `tfsdk:"block"`
	IncludeSettled types.Bool          `tfsdk:"include_settled"`
	Reservations   []reservationsModel `tfsdk:"reservations"`
}

// reservationsModel maps reservations schema data.
type reservationsModel struct {
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
func (d *reservationsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_reservations"
}

// Schema defines the schema for the data source.
func (d *reservationsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "The reservations data source allows you to retrieve information about all existing reservations in the specified space and block.",
		Attributes: map[string]schema.Attribute{
			"space": schema.StringAttribute{
				Description: "Name of the existing space in the IPAM application.",
				Required:    true,
			},
			"block": schema.StringAttribute{
				Description: "Name of the existing block, related to the specified space.",
				Required:    true,
			},
			"include_settled": schema.BoolAttribute{
				Description: "Settled reservations must be also included? Defaults to `false`.",
				Optional:    true,
			},
			"reservations": schema.ListNestedAttribute{
				Description: "List containing the `reservations` found for the specified attributes.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "The unique identifier of the reservation.",
							Computed:    true,
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
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *reservationsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state reservationsDataSourceModel

	// Read Terraform configuration state into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	reservations, err := d.client.GetReservations(state.Space.ValueString(), state.Block.ValueString(), state.IncludeSettled.ValueBool())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read AzureIpam Reservations",
			err.Error(),
		)
		return
	}

	// Map response body to model
	for _, reservation := range *reservations {
		var reservationState reservationsModel
		reservationState.Id = types.StringValue(reservation.Id)
		reservationState.Cidr = types.StringValue(reservation.Cidr)
		reservationState.Description = types.StringValue(reservation.Description)
		reservationState.CreatedOn = timetypes.NewRFC3339TimeValue(time.Unix(int64(reservation.CreatedOn), 0))
		reservationState.CreatedBy = types.StringValue(reservation.CreatedBy)
		if reservation.SettledOn == nil {
			reservationState.SettledOn = timetypes.NewRFC3339Null()
		} else {
			reservationState.SettledOn = timetypes.NewRFC3339TimeValue(time.Unix(int64(*reservation.SettledOn), 0))
		}
		if reservation.SettledBy == nil {
			reservationState.SettledBy = types.StringNull()
		} else {
			reservationState.SettledBy = types.StringValue(*reservation.SettledBy)
		}
		reservationState.Status = types.StringValue(reservation.Status)
		reservationState.Tags, _ = types.MapValueFrom(ctx, types.StringType, reservation.Tags)

		state.Reservations = append(state.Reservations, reservationState)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *reservationsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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