package provider

import (
	"context"
	"os"

	ipamclient "terraform-provider-azureipam/ipamclient"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ provider.Provider              = &azureIpamProvider{}
	_ provider.ProviderWithFunctions = &azureIpamProvider{}
)

// NewAzureIpamProvider is a helper function to simplify provider server and testing implementation.
func NewAzureIpamProvider(version string) func() provider.Provider {
	return func() provider.Provider {
		return &azureIpamProvider{
			version: version,
		}
	}
}

// azureIpamProvider defines the provider implementation.
type azureIpamProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// azureIpamProviderModel describes the provider data model.
type azureIpamProviderModel struct {
	ApiUrl                      types.String `tfsdk:"api_url"`
	Token                       types.String `tfsdk:"token"`
	SkipCertificateVerification types.Bool   `tfsdk:"skip_cert_verification"`
}

// Metadata returns the provider type name.
func (p *azureIpamProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "azureipam"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *azureIpamProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Terraform provider to manage reservations in Azure IPAM solution through REST API.",
		Attributes: map[string]schema.Attribute{
			"api_url": schema.StringAttribute{
				MarkdownDescription: "The root url of the APIM REST API solution to be used, without the /api url suffix. Must be also assigned at AZUREIPAM_API_URL environment variable.",
				Optional:            true,
			},
			"token": schema.StringAttribute{
				MarkdownDescription: "The bearer token to be used when authenticating to the API. Must be also assigned at AZUREIPAM_TOKEN environment variable.",
				Optional:            true,
				Sensitive:           true,
			},
			"skip_cert_verification": schema.BoolAttribute{
				MarkdownDescription: "Specifies it the certificate chain validation must be skipped calling the API endpoint. Default to false.",
				Optional:            true,
			},
		},
	}
}

func (p *azureIpamProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring AzureIpam client")

	// Retrieve provider data from configuration
	var config azureIpamProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.
	if config.ApiUrl.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_url"),
			"Unknown AzureIpam API url",
			"The provider cannot create the AzureIpam API client as there is an unknown configuration value for the AzureIpam API host. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the AZUREIPAM_API_URL environment variable.",
		)
	}
	if config.Token.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("token"),
			"Unknown AzureIpam API access token",
			"The provider cannot create the AzureIpam API client as there is an unknown configuration value for the AzureIpam API access token. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the AZUREIPAM_TOKEN environment variable.",
		)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.
	apiUrl := os.Getenv("AZUREIPAM_API_URL")
	token := os.Getenv("AZUREIPAM_TOKEN")
	if !config.ApiUrl.IsNull() {
		apiUrl = config.ApiUrl.ValueString()
	}
	if !config.Token.IsNull() {
		token = config.Token.ValueString()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.
	if apiUrl == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_url"),
			"Missing AzureIpam API url",
			"The provider cannot create the AzureIpam API client as there is a missing or empty value for the AzureIpam API url. "+
				"Set the url value in the configuration or use the AZUREIPAM_API_URL environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}
	if token == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("token"),
			"Missing AzureIpam API access token",
			"The provider cannot create the AzureIpam API client as there is a missing or empty value for the AzureIpam API access token. "+
				"Set the access token value in the configuration or use the AZUREIPAM_TOKEN environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	var skipCertVerification bool
	if p.version == "test" {
		skipCertVerification = false //always false for acceptance tests, to enforce the http.DefaultTransport usage
	} else if config.SkipCertificateVerification.IsNull() {
		skipCertVerification = false
	} else {
		skipCertVerification = config.SkipCertificateVerification.ValueBool()
	}

	ctx = tflog.SetField(ctx, "azureipam_api_url", apiUrl)
	ctx = tflog.SetField(ctx, "azureipam_token", token)
	ctx = tflog.SetField(ctx, "azureipam_skip_cert_verification", skipCertVerification)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "azureipam_token")

	tflog.Debug(ctx, "Creating AzureIpam client")
	// Create a new AzureIpam client using the configuration values
	client, err := ipamclient.NewClient(&apiUrl, &token, skipCertVerification)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create AzureIpam API Client",
			"An unexpected error occurred when creating the AzureIpam API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"AzureIpam Client Error: "+err.Error(),
		)
		return
	}

	// Make the AzureIpam client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client

	tflog.Info(ctx, "Configured AzureIpam client", map[string]any{"success": true})
}

// DataSources defines the data sources implemented in the provider.
func (p *azureIpamProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewReservationsDataSource,
		NewReservationDataSource,
		NewSpacesDataSource,
		NewSpaceDataSource,
		NewBlocksDataSource,
		NewBlockDataSource,
		NewExternalsDataSource,
		NewExternalDataSource,
		NewBlockNetworksAvailablesDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *azureIpamProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewReservationResource,
		NewSpaceResource,
		NewBlockResource,
		NewExternalResource,
		NewReservationCidrResource,
		NewBlockNetworkResource,
	}
}

func (p *azureIpamProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{}
}
