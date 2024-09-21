package provider

import (
	"context"
	"fmt"
	"regexp"

	ipamclient "terraform-provider-azureipam/ipamclient"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &blockResource{}
	_ resource.ResourceWithConfigure   = &blockResource{}
	_ resource.ResourceWithImportState = &blockResource{}
)

// NewBlockResource is a helper function to simplify the provider implementation.
func NewBlockResource() resource.Resource {
	return &blockResource{}
}

// blockResourceModel maps the resource schema data.
type blockResourceModel struct {
	Name  types.String `tfsdk:"name"`
	Space types.String `tfsdk:"space"`
	Cidr  types.String `tfsdk:"cidr"`
}

// blockResource is the resource implementation.
type blockResource struct {
	client *ipamclient.Client
}

// Metadata returns the resource type name.
func (r *blockResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_block"
}

// Schema defines the schema for the resource.
func (r *blockResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "The block resource allows you to create a IPAM block in a specific Space.",
		Attributes: map[string]schema.Attribute{
			"space": schema.StringAttribute{
				Description: "Name of the space where the block must be created. Changing this forces a new resource to be created.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Name of the block.",
				Required:    true,
			},
			"cidr": schema.StringAttribute{
				Description: "The IP range to configure to the block, in cidr notation.",
				Required:    true,
			},
		},
	}
}

// Create a new resource.
func (r *blockResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan blockResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	block, err := r.client.CreateBlock(
		plan.Space.ValueString(),
		plan.Name.ValueString(),
		plan.Cidr.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating block",
			"Could not create block, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	flattenBlock(block, &plan)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information.
func (r *blockResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state blockResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	//read block
	block, err := r.client.GetBlock(
		state.Space.ValueString(),
		state.Name.ValueString(),
		false,
		false,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading AzureIpam Block",
			"Could not read AzureIpam Block with name "+state.Name.ValueString()+": "+err.Error(),
		)
		return
	}

	// Overwrite items with refreshed state
	flattenBlock(block, &state)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update not allowed, returning readed plan as current state.
func (n *blockResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve current state of the resource
	var state blockResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Retrieve values from plan
	var plan blockResourceModel
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	//Modify the block resource
	block, err := n.client.UpdateBlock(
		state.Space.ValueString(),
		state.Name.ValueString(),
		plan.Name.ValueStringPointer(),
		plan.Cidr.ValueStringPointer(),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating block",
			"Could not update block, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	flattenBlock(block, &plan)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *blockResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state blockResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing block
	err := r.client.DeleteBlock(
		state.Space.ValueString(),
		state.Name.ValueString(),
		true,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting AzureIpam Block",
			"Could not delete block, unexpected error: "+err.Error(),
		)
		return
	}
}

// Configure adds the provider configured client to the resource.
func (r *blockResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *blockResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID, validate, split and save to id attribute
	re := regexp.MustCompile("^(?<space>[a-zA-Z0-9]*)/(?<block>[a-zA-Z0-9]*)$")
	//validate
	if !re.MatchString(req.ID) {
		resp.Diagnostics.AddError(
			"Error Importing AzureIpam Block",
			"The specified ID is not in the correct format {SpaceName}/{BlockName}.",
		)
		return
	}
	//extract values
	matches := re.FindStringSubmatch(req.ID)
	space := matches[re.SubexpIndex("space")]
	block := matches[re.SubexpIndex("block")]

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("space"), space)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), block)...)
	//resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

func flattenBlock(block *ipamclient.Block, model *blockResourceModel) {
	model.Name = types.StringValue(block.Name)
	model.Space = types.StringValue(block.Space)
	model.Cidr = types.StringValue(block.Cidr)
}
