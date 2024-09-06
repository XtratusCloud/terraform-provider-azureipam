package provider

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	ipamclient "terraform-provider-azureipam/ipamclient"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &reservationResource{}
	_ resource.ResourceWithConfigure = &reservationResource{}
)

// NewReservationResource is a helper function to simplify the provider implementation.
func NewReservationResource() resource.Resource {
	return &reservationResource{}
}

// reservationResourceModel maps the resource schema data.
type reservationResourceModel struct {
	Space         types.String      `tfsdk:"space"`
	Block         types.String      `tfsdk:"block"`
	Size          types.Int32       `tfsdk:"size"`
	Description   types.String      `tfsdk:"description"`
	ReverseSearch types.Bool        `tfsdk:"reverse_search"`
	SmallestCidr  types.Bool        `tfsdk:"smallest_cidr"`
	Id            types.String      `tfsdk:"id"`
	Cidr          types.String      `tfsdk:"cidr"`
	CreatedBy     types.String      `tfsdk:"created_by"`
	CreatedOn     timetypes.RFC3339 `tfsdk:"created_on"`
	SettledBy     types.String      `tfsdk:"settled_by"`
	SettledOn     timetypes.RFC3339 `tfsdk:"settled_on"`
	Status        types.String      `tfsdk:"status"`
	Tags          types.Map         `tfsdk:"tags"`
}

// reservationResource is the resource implementation.
type reservationResource struct {
	client *ipamclient.Client
}

// Metadata returns the resource type name.
func (r *reservationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_reservation"
}

// Schema defines the schema for the resource.
func (r *reservationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "The reservation resource allows you to create a IPAM reservation in the specific space and block.",
		Attributes: map[string]schema.Attribute{
			"space": schema.StringAttribute{
				Description: "Name of the existing space in the IPAM application.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"block": schema.StringAttribute{
				Description: "Name of the existing block, related to the specified space, in which the reservation is to be made.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"size": schema.Int32Attribute{
				Description: "Integer value to indicate the subnet mask bits, which defines the size of the vnet to reserve (example 24 for a /24 subnet).",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "Description text that describe the reservation, that will be added as an additional tag.",
				Optional:    true,
			},
			"reverse_search": schema.BoolAttribute{
				Description: "New networks will be created as close to the end of the block as possible?. Defaults to `false`.",
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"smallest_cidr": schema.BoolAttribute{
				Description: "New networks will be created using the smallest possible available block? (e.g. it will not break up large CIDR blocks when possible) .Defaults to `false`.",
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"id": schema.StringAttribute{
				Description: "The unique identifier of the generated reservation.",
				Computed:    true,
			},
			"cidr": schema.StringAttribute{
				Description: "The assigned and reserved range, in cidr notation.",
				Computed:    true,
			},
			"created_by": schema.StringAttribute{
				Description: "Email or identification of user that created the reservation.",
				Computed:    true,
			},
			"created_on": schema.StringAttribute{
				CustomType:  timetypes.RFC3339Type{},
				Description: "The date and time that the reservacion was created.",
				Computed:    true,
			},
			"settled_by": schema.StringAttribute{
				Description: "Email or identification of user that settled the reservation.",
				Computed:    true,
			},
			"settled_on": schema.StringAttribute{
				CustomType:  timetypes.RFC3339Type{},
				Description: "The date and time that the reservacion was settled.",
				Computed:    true,
			},
			"status": schema.StringAttribute{
				Description: "Status of the reservation, a 'wait' status indicates that is waiting for the related vnet creation",
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

// Create a new resource.
func (r *reservationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan reservationResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	reservation, err := r.client.CreateReservation(
		plan.Space.ValueString(),
		plan.Block.ValueString(),
		plan.Description.ValueString(),
		int(plan.Size.ValueInt32()),
		plan.ReverseSearch.ValueBool(),
		plan.SmallestCidr.ValueBool(),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating reservation",
			"Could not create reservation, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	flattenReservation(reservation, &plan)
	plan.Tags, _ = types.MapValueFrom(ctx, types.StringType, reservation.Tags)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information.
func (r *reservationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state reservationResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
 
	//read reservation
	reservation, err := r.client.FindReservationById(
		state.Id.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading AzureIpam Reservation",
			"Could not read AzureIpam Reservation with id "+state.Id.ValueString()+": "+err.Error(),
		)
		return
	}

	// Overwrite items with refreshed state
	flattenReservation(reservation, &state)
	// state.ReverseSearch = types.BoolValue(reverse_search)
	// state.SmallestCidr = types.BoolValue(smallest_cidr)
	state.Tags, _ = types.MapValueFrom(ctx, types.StringType, reservation.Tags)
	//Calculate requested size from assigned Cidr
	size, err := strconv.Atoi(strings.Split(reservation.Cidr, "/")[1])
	state.Size = types.Int32Value(int32(size))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading AzureIpam Reservation",
			"Could not determinate requested size for Reservation with id "+state.Id.ValueString()+": "+err.Error(),
		)
		return
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update not allowed, returning readed plan as current state.
func (n *reservationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var model reservationResource

	resp.Diagnostics.Append(req.Plan.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

func (r *reservationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state reservationResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing reservation
	err := r.client.DeleteReservation(
		state.Space.ValueString(),
		state.Block.ValueString(),
		state.Id.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting AzureIpam Reservation",
			"Could not delete reservation, unexpected error: "+err.Error(),
		)
		return
	}
}

// Configure adds the provider configured client to the resource.
func (r *reservationResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = client
}

func (r *reservationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func flattenReservation(reservation *ipamclient.Reservation, model *reservationResourceModel) {
	model.Id = types.StringValue(reservation.Id)
	model.Space = types.StringValue(reservation.Space)
	model.Block = types.StringValue(reservation.Block)
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
}
