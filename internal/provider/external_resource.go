package provider

import (
	"context"
	"fmt"
	"regexp"

	ipamclient "terraform-provider-azureipam/ipamclient"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &externalResource{}
	_ resource.ResourceWithConfigure   = &externalResource{}
	_ resource.ResourceWithImportState = &externalResource{}
)

// NewExternalResource is a helper function to simplify the provider implementation.
func NewExternalResource() resource.Resource {
	return &externalResource{}
}

// externalResourceModel maps the resource schema data.
type externalResourceModel struct {
	Space       types.String `tfsdk:"space"`
	Block       types.String `tfsdk:"block"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Cidr        types.String `tfsdk:"cidr"`
}

// externalResource is the resource implementation.
type externalResource struct {
	client *ipamclient.Client
}

// Metadata returns the resource type name.
func (r *externalResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_external"
}

// Schema defines the schema for the resource.
func (r *externalResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "The external resource allows you to associate an external network to the target space and block.",
		Attributes: map[string]schema.Attribute{
			"space": schema.StringAttribute{
				Description: "Name of the space where the external must be associated..",
				Required:    true,
			},
			"block": schema.StringAttribute{
				Description: "Name of the block where the external must be associated.",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "Name of the external network.",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "Text that describes the external network.",
				Required:    true,
			},
			"cidr": schema.StringAttribute{
				Description: "The IP range to configure to the external network, in cidr notation.",
				Required:    true,
			},
		},
	}
}

// Create a new resource.
func (r *externalResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan externalResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	external, err := r.client.CreateExternal(
		plan.Space.ValueString(),
		plan.Block.ValueString(),
		plan.Name.ValueString(),
		plan.Description.ValueString(),
		plan.Cidr.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating external network",
			"Could not create external network, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	flattenExternal(external, &plan)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information.
func (r *externalResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state externalResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	//read external
	external, err := r.client.GetExternal(
		state.Space.ValueString(),
		state.Block.ValueString(),
		state.Name.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading AzureIpam external network",
			"Could not read AzureIpam external  network with name "+state.Name.ValueString()+": "+err.Error(),
		)
		return
	}

	// Overwrite items with refreshed state
	flattenExternal(external, &state)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update not allowed, returning readed plan as current state.
func (n *externalResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve current state of the resource
	var state externalResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Retrieve values from plan
	var plan externalResourceModel
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	//Modify the external resource
	external, err := n.client.UpdateExternal(
		state.Space.ValueString(),
		state.Block.ValueString(),
		state.Name.ValueString(),
		plan.Name.ValueStringPointer(),
		plan.Description.ValueStringPointer(),
		plan.Cidr.ValueStringPointer(),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating external network",
			"Could not update external network, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	flattenExternal(external, &plan)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *externalResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state externalResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing external network
	err := r.client.DeleteExternal(
		state.Space.ValueString(),
		state.Block.ValueString(),
		state.Name.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting AzureIpam external network",
			"Could not delete external network, unexpected error: "+err.Error(),
		)
		return
	}
}

// Configure adds the provider configured client to the resource.
func (r *externalResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *externalResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID, validate, split and save to id attribute
	re := regexp.MustCompile("^(?<space>[a-zA-Z0-9]*)/(?<block>[a-zA-Z0-9]*)/(?<name>[a-zA-Z0-9]*)$")
	//validate
	if !re.MatchString(req.ID) {
		resp.Diagnostics.AddError(
			"Error Importing AzureIpam external network",
			"The specified ID is not in the correct format {SpaceName}/{BlockName}/{ExternalNetworkName}.",
		)
		return
	}
	//extract values
	matches := re.FindStringSubmatch(req.ID)
	space := matches[re.SubexpIndex("space")]
	block := matches[re.SubexpIndex("block")]
	name := matches[re.SubexpIndex("name")]

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("space"), space)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("block"), block)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), name)...)
	//resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

func flattenExternal(external *ipamclient.External, model *externalResourceModel) {
	model.Space = types.StringValue(external.Space)
	model.Block = types.StringValue(external.Block)
	model.Name = types.StringValue(external.Name)
	model.Description = types.StringValue(external.Description)
	model.Cidr = types.StringValue(external.Cidr)
}
