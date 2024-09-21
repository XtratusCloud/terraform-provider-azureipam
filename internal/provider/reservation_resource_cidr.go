package provider

import (
	"context"
	"fmt"
	"time"

	ipamclient "terraform-provider-azureipam/ipamclient"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &reservationResourceCidr{}
	_ resource.ResourceWithConfigure   = &reservationResourceCidr{}
	_ resource.ResourceWithImportState = &reservationResourceCidr{}
)

// NewReservationCidrResource is a helper function to simplify the provider implementation.
func NewReservationCidrResource() resource.Resource {
	return &reservationResourceCidr{}
}

// reservationResourceCidrModel maps the resource schema data.
type reservationResourceCidrModel struct {
	Space        types.String      `tfsdk:"space"`
	Block        types.String      `tfsdk:"block"`
	SpecificCidr types.String      `tfsdk:"specific_cidr"`
	Description  types.String      `tfsdk:"description"`
	Id           types.String      `tfsdk:"id"`
	Cidr         types.String      `tfsdk:"cidr"`
	CreatedBy    types.String      `tfsdk:"created_by"`
	CreatedOn    timetypes.RFC3339 `tfsdk:"created_on"`
	SettledBy    types.String      `tfsdk:"settled_by"`
	SettledOn    timetypes.RFC3339 `tfsdk:"settled_on"`
	Status       types.String      `tfsdk:"status"`
	Tags         types.Map         `tfsdk:"tags"`
}

// reservationResourceCidr is the resource implementation.
type reservationResourceCidr struct {
	client *ipamclient.Client
}

// Metadata returns the resource type name.
func (r *reservationResourceCidr) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_reservation_cidr"
}

// Schema defines the schema for the resource.
func (r *reservationResourceCidr) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "The reservation resource allows you to create a IPAM reservation in the specific space and block with a fixed cidr.",
		Attributes: map[string]schema.Attribute{
			"space": schema.StringAttribute{
				Description: "Name of the existing space in the IPAM application. Changing this forces a new resource to be created.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"block": schema.StringAttribute{
				Description: "List with the names of blocks in the specified space in which the reservation is to be create. The list is evaluated in the order provider. Changing this forces a new resource to be created.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"specific_cidr": schema.StringAttribute{
				Description: "The specific CIDR to reserve, in cidr notation. At least one of size or specific_cidr attribute must be specified. Not allowed if more than one block is specified.",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				Description: "Description text that describe the reservation, that will be added as an additional tag.",
				Optional:    true,
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
func (r *reservationResourceCidr) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan reservationResourceCidrModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	block := []string{plan.Block.ValueString()}
	reservation, err := r.client.CreateReservation(
		plan.Space.ValueString(),
		block,
		plan.Description.ValueStringPointer(),
		nil,
		plan.SpecificCidr.ValueStringPointer(),
		false,
		false,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating reservation",
			"Could not create reservation, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	flattenReservationCidr(reservation, &plan)
	plan.Tags, _ = types.MapValueFrom(ctx, types.StringType, reservation.Tags)
	// //Calculate requested size from assigned Cidr
	// size, err := strconv.Atoi(strings.Split(reservation.Cidr, "/")[1])
	// plan.Size = types.Int32Value(int32(size))
	// if err != nil {
	// 	resp.Diagnostics.AddError(
	// 		"Error Reading AzureIpam Reservation",
	// 		"Could not determinate requested size for Reservation with id "+plan.Id.ValueString()+": "+err.Error(),
	// 	)
	// 	return
	// }

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information.
func (r *reservationResourceCidr) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state reservationResourceCidrModel
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
	flattenReservationCidr(reservation, &state)
	state.Tags, _ = types.MapValueFrom(ctx, types.StringType, reservation.Tags)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update not allowed, returning readed plan as current state.
func (n *reservationResourceCidr) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var model reservationResourceCidr

	resp.Diagnostics.Append(req.Plan.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

func (r *reservationResourceCidr) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state reservationResourceCidrModel
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
func (r *reservationResourceCidr) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *reservationResourceCidr) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func flattenReservationCidr(reservation *ipamclient.Reservation, model *reservationResourceCidrModel) {
	model.Id = types.StringValue(reservation.Id)
	model.Space = types.StringValue(reservation.Space)
	model.Block = types.StringValue(reservation.Block)
	model.SpecificCidr = types.StringValue(reservation.Cidr) //the requestd cidr is also the assigned cidr
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
