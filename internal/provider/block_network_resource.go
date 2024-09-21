package provider

import (
	"context"
	"fmt"
	"regexp"

	ipamclient "terraform-provider-azureipam/ipamclient"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &blockNetworkResource{}
	_ resource.ResourceWithConfigure   = &blockNetworkResource{}
	_ resource.ResourceWithImportState = &blockNetworkResource{}
)

// NewBlockNetworkResource is a helper function to simplify the provider implementation.
func NewBlockNetworkResource() resource.Resource {
	return &blockNetworkResource{}
}

// blockNetworkResourceModel maps the resource schema data.
type blockNetworkResourceModel struct {
	Space          types.String `tfsdk:"space"`
	Block          types.String `tfsdk:"block"`
	Id             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	Prefixes       types.List   `tfsdk:"prefixes"`
	ResourceGroup  types.String `tfsdk:"resource_group"`
	SubscriptionId types.String `tfsdk:"subscription_id"`
	TenantId       types.String `tfsdk:"tenant_id"`
}

// blockNetworkResource is the resource implementation.
type blockNetworkResource struct {
	client *ipamclient.Client
}

// Metadata returns the resource type name.
func (r *blockNetworkResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_block_network"
}

// Schema defines the schema for the resource.
func (r *blockNetworkResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "The block_network resource allow to associate an existing azure network to the target block.",
		Attributes: map[string]schema.Attribute{
			"space": schema.StringAttribute{
				Description: "Name of the space where the external must be associated. Changing this forces a new resource to be created.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"block": schema.StringAttribute{
				Description: "Name of the block where the external must be associated. Changing this forces a new resource to be created.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"id": schema.StringAttribute{
				Description: "Azure Resource ID of the virtual network to associate.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
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
	}
}

// Create a new block network.
func (r *blockNetworkResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan blockNetworkResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	block, err := r.client.CreateBlockNetwork(
		plan.Space.ValueString(),
		plan.Block.ValueString(),
		plan.Id.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating block network",
			"Could not create block network, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	diags = flattenBlockNetwork(ctx, block, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information.
func (r *blockNetworkResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state blockNetworkResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	//read external
	block, err := r.client.GetBlockNetworkInfo(
		state.Space.ValueString(),
		state.Block.ValueString(),
		state.Id.ValueString(),
		true,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading AzureIpam block network",
			"Could not read AzureIpam block network with id "+state.Id.ValueString()+": "+err.Error(),
		)
		return
	}

	// Overwrite items with refreshed state
	diags = flattenBlockNetwork(ctx, block, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
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
func (n *blockNetworkResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var model blockNetworkResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

func (r *blockNetworkResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state blockNetworkResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing external network
	err := r.client.DeleteBlockNetwork(
		state.Space.ValueString(),
		state.Block.ValueString(),
		state.Id.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting AzureIpam block network",
			"Could not delete block network, unexpected error: "+err.Error(),
		)
		return
	}
}

// Configure adds the provider configured client to the resource.
func (r *blockNetworkResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *blockNetworkResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID, validate, split and save to id attribute
	re := regexp.MustCompile("^(?<space>[a-zA-Z0-9]*)/(?<block>[a-zA-Z0-9]*)/(?<id>/subscriptions/(?<SubscriptionId>.*)/resourceGroups/(?<ResourceGroupName>.*)/providers/(?<ResourceProviderNamespace>.*)/(?<ResourceType>.*)/(?<ResourceName>.*))$")
	
	//validate
	if !re.MatchString(req.ID) {
		resp.Diagnostics.AddError(
			"Error Importing AzureIpam block network",
			"The specified ID is not in the correct format {SpaceName}/{BlockName}/{AzureResourceIdOfNetwork}.",
		)
		return
	}
	//extract values
	matches := re.FindStringSubmatch(req.ID)
	space := matches[re.SubexpIndex("space")]
	block := matches[re.SubexpIndex("block")]
	id := matches[re.SubexpIndex("id")]

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("space"), space)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("block"), block)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

func flattenBlockNetwork(ctx context.Context, external *ipamclient.BlockNetworkInfo, model *blockNetworkResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics 

	model.Id = types.StringValue(external.Id)
	model.Name = types.StringValue(external.Name)
	model.Prefixes, diags = types.ListValueFrom(ctx, types.StringType, external.Prefixes)
	model.ResourceGroup = types.StringValue(*external.ResourceGroup)
	model.SubscriptionId = types.StringValue(*external.SubscriptionId)
	model.TenantId = types.StringValue(*external.TenantId)

	return diags
}
