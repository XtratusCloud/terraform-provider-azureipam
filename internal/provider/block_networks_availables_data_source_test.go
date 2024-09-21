package provider

import (
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

func TestAccBlockNetworksAvailablesDataSource(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("GET", "https://mockedHost.azurewebsites.net/api/spaces/au/blocks/AustraliaEast/available",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/datasource/block_networks/availables.json").String()), nil
		})

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccProviderConfig + `data "azureipam_block_networks_availables" "test" {
					space  = "au"
  					block  = "AustraliaEast"
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					//Verify common attributes to ensure that all are set
					resource.TestCheckResourceAttr("data.azureipam_block_networks_availables.test", "space", "au"),
					resource.TestCheckResourceAttr("data.azureipam_block_networks_availables.test", "block", "AustraliaEast"),
					resource.TestCheckResourceAttr("data.azureipam_block_networks_availables.test", "ids.#", "2"),
					resource.TestCheckResourceAttr("data.azureipam_block_networks_availables.test", "ids.0", "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/RG-WE-ALL-COMMS-01/providers/Microsoft.Network/virtualNetworks/vnet-we-a-testzavd-01"),
					resource.TestCheckResourceAttr("data.azureipam_block_networks_availables.test", "ids.1", "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/RG-WE-ALL-COMMS-01/providers/Microsoft.Network/virtualNetworks/vnet-we-d-terratest-hub-01"),
				),
			},
		},
	})
}
