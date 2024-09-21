package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

const (
	// providerConfig is a shared configuration to combine with the actual
	// test configuration so the AzureIpam client is properly configured.
	// It is also possible to use the AZUREIPAM_ environment variables instead,
	// such as updating the Makefile and running the testing through that tool.
	testAccProviderConfig = `
provider "azureipam" {
  api_url = "https://mockedHost.azurewebsites.net"
  token = "dummyForTesting"
}
`
)

var (
	// testAccProtoV6ProviderFactories are used to instantiate a provider during
	// acceptance testing. The factory function will be invoked for every Terraform
	// CLI command executed to create a provider server to which the CLI can
	// reattach.
	testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"azureipam": providerserver.NewProtocol6WithError(NewAzureIpamProvider("test")()),
	}
)