package provider

import (
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

func TestAccBlockNetworksDataSource(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("GET", "https://mockedHost.azurewebsites.net/api/spaces/au/blocks/AustraliaEast/networks?expand=true",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/datasource/block_networks/block_networks_all.json").String()), nil
		})

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccProviderConfig + `data "azureipam_block_networks" "test" {
					space  = "au"
  					block  = "AustraliaEast"
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					//Verify common attributes to ensure that all are set
					resource.TestCheckResourceAttr("data.azureipam_block_networks.test", "space", "au"),
					resource.TestCheckResourceAttr("data.azureipam_block_networks.test", "block", "AustraliaEast"),
					resource.TestCheckResourceAttr("data.azureipam_block_networks.test", "networks.#", "2"),
					// Verify first network to ensure all attributes are set
					resource.TestCheckResourceAttr("data.azureipam_block_networks.test", "networks.0.id", "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/RG-WE-ALL-COMMS-01/providers/Microsoft.Network/virtualNetworks/vnet-we-a-testzavd-01"),
					resource.TestCheckResourceAttr("data.azureipam_block_networks.test", "networks.0.name", "vnet-we-a-testzavd-01"),
					resource.TestCheckResourceAttr("data.azureipam_block_networks.test", "networks.0.prefixes.#", "1"),
					resource.TestCheckResourceAttr("data.azureipam_block_networks.test", "networks.0.prefixes.0", "10.82.1.224/27"),
					resource.TestCheckResourceAttr("data.azureipam_block_networks.test", "networks.0.resource_group", "rg-we-all-comms-01"),
					resource.TestCheckResourceAttr("data.azureipam_block_networks.test", "networks.0.subscription_id", "00000000-0000-0000-0000-000000000000"),
					resource.TestCheckResourceAttr("data.azureipam_block_networks.test", "networks.0.tenant_id", "11111111-1111-1111-1111-111111111111"),
					// Verify second network to ensure all attributes are set
					resource.TestCheckResourceAttr("data.azureipam_block_networks.test", "networks.1.id", "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/RG-WE-ALL-COMMS-01/providers/Microsoft.Network/virtualNetworks/vnet-we-d-terratest-hub-01"),
					resource.TestCheckResourceAttr("data.azureipam_block_networks.test", "networks.1.name", "vnet-we-d-terratest-hub-01"),
					resource.TestCheckResourceAttr("data.azureipam_block_networks.test", "networks.0.prefixes.#", "1"),
					resource.TestCheckResourceAttr("data.azureipam_block_networks.test", "networks.1.prefixes.0", "10.82.0.0/24"),
					resource.TestCheckResourceAttr("data.azureipam_block_networks.test", "networks.1.resource_group", "rg-we-all-comms-01"),
					resource.TestCheckResourceAttr("data.azureipam_block_networks.test", "networks.1.subscription_id", "00000000-0000-0000-0000-000000000000"),
					resource.TestCheckResourceAttr("data.azureipam_block_networks.test", "networks.1.tenant_id", "11111111-1111-1111-1111-111111111111"),
				),
			},
		},
	})
}
