package provider

import (
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

func TestAccBlockNetworkResource(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("POST", "https://mockedHost.azurewebsites.net/api/spaces/au/blocks/AustraliaEast/networks",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/networks/space_with_new_network.json").String()), nil
		})
	httpmock.RegisterResponder("GET", "https://mockedHost.azurewebsites.net/api/spaces/au/blocks/AustraliaEast/networks?expand=true",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/networks/networks_with_new_network.json").String()), nil
		})
	httpmock.RegisterResponder("DELETE", "https://mockedHost.azurewebsites.net/api/spaces/au/blocks/AustraliaEast/networks",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, ""), nil
		})

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccProviderConfig + `resource "azureipam_block_network" "test" {
					space = "au"
  					block = "AustraliaEast"
  					id    = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/RG-WE-ALL-COMMS-01/providers/Microsoft.Network/virtualNetworks/vnet-we-d-terratest-hub-01"
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					//Verify common attributes to ensure that all are set
					resource.TestCheckResourceAttr("azureipam_block_network.test", "space", "au"),
					resource.TestCheckResourceAttr("azureipam_block_network.test", "block", "AustraliaEast"),
					resource.TestCheckResourceAttr("azureipam_block_network.test", "id", "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/RG-WE-ALL-COMMS-01/providers/Microsoft.Network/virtualNetworks/vnet-we-d-terratest-hub-01"),
					resource.TestCheckResourceAttr("azureipam_block_network.test", "name", "vnet-we-d-terratest-hub-01"),
					resource.TestCheckResourceAttr("azureipam_block_network.test", "resource_group", "rg-we-all-comms-01"),
					resource.TestCheckResourceAttr("azureipam_block_network.test", "subscription_id", "00000000-0000-0000-0000-000000000000"),
					resource.TestCheckResourceAttr("azureipam_block_network.test", "tenant_id", "11111111-1111-1111-1111-111111111111"),
					//Verify number of prefixes returned
					resource.TestCheckResourceAttr("azureipam_block_network.test", "prefixes.#", "1"),
					//Verify first prefix
					resource.TestCheckResourceAttr("azureipam_block_network.test", "prefixes.0", "10.82.0.0/24"),
				),
			},
			// ImportState testing
			{
				ResourceName:                         "azureipam_block_network.test",
				ImportState:                          true,
				ImportStateId:                        "au/AustraliaEast//subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/RG-WE-ALL-COMMS-01/providers/Microsoft.Network/virtualNetworks/vnet-we-d-terratest-hub-01",
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "name",
			},
			// Update  NOT ALLOWED by provider

			// Delete testing automatically occurs in TestCase
		},
	})
}
