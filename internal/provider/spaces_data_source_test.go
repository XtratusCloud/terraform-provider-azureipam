package provider

import (
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

func TestAccSpacesWithoutUtilizationAndVnetDataSource(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("GET", "https://mockedHost.azurewebsites.net/api/spaces?expand=false&utilization=false",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/datasource/spaces/spaces_without_utilization_and_vnet.json").String()), nil
		})

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccProviderConfig + `data "azureipam_spaces" "test" {
					expand             = false
  					append_utilization = false
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					//Verify common attributes to ensure that all are set
					resource.TestCheckResourceAttr("data.azureipam_spaces.test", "append_utilization", "false"),
					resource.TestCheckResourceAttr("data.azureipam_spaces.test", "expand", "false"),
					// Verify number of spaces returned
					resource.TestCheckResourceAttr("data.azureipam_spaces.test", "spaces.#", "1"),
					// Verify the first space to ensure all attributes are set
					resource.TestCheckResourceAttr("data.azureipam_spaces.test", "spaces.0.name", "au"),
					resource.TestCheckResourceAttr("data.azureipam_spaces.test", "spaces.0.description", "Australia"),
					resource.TestCheckNoResourceAttr("data.azureipam_spaces.test", "spaces.0.size"),
					resource.TestCheckNoResourceAttr("data.azureipam_spaces.test", "spaces.0.used"),

					// Verify number of blocks returned in first space
					resource.TestCheckResourceAttr("data.azureipam_spaces.test", "spaces.0.blocks.#", "2"),
					// Verify first block to ensure all attributes are set
					resource.TestCheckResourceAttr("data.azureipam_spaces.test", "spaces.0.blocks.0.name", "AustraliaSoutheast"),
					resource.TestCheckResourceAttr("data.azureipam_spaces.test", "spaces.0.blocks.0.cidr", "10.83.0.0/16"),
					resource.TestCheckNoResourceAttr("data.azureipam_spaces.test", "spaces.0.blocks.0.size"),
					resource.TestCheckNoResourceAttr("data.azureipam_spaces.test", "spaces.0.blocks.0.used"),
					resource.TestCheckResourceAttr("data.azureipam_spaces.test", "spaces.0.blocks.0.vnets.#", "0"),
					resource.TestCheckResourceAttr("data.azureipam_spaces.test", "spaces.0.blocks.0.reservations.#", "0"),
					resource.TestCheckResourceAttr("data.azureipam_spaces.test", "spaces.0.blocks.0.externals.#", "1"),
					resource.TestCheckResourceAttr("data.azureipam_spaces.test", "spaces.0.blocks.0.externals.0.name", "test"),
					resource.TestCheckResourceAttr("data.azureipam_spaces.test", "spaces.0.blocks.0.externals.0.description", "external for tests"),
					resource.TestCheckResourceAttr("data.azureipam_spaces.test", "spaces.0.blocks.0.externals.0.cidr", "10.83.0.0/24"),
					// Verify second block to ensure all attributes are set
					resource.TestCheckResourceAttr("data.azureipam_spaces.test", "spaces.0.blocks.1.name", "AustraliaEast"),
					resource.TestCheckResourceAttr("data.azureipam_spaces.test", "spaces.0.blocks.1.cidr", "10.82.0.0/16"),
					resource.TestCheckNoResourceAttr("data.azureipam_spaces.test", "spaces.0.blocks.1.size"),
					resource.TestCheckNoResourceAttr("data.azureipam_spaces.test", "spaces.0.blocks.1.used"),
					resource.TestCheckResourceAttr("data.azureipam_spaces.test", "spaces.0.blocks.1.externals.#", "0"),
					resource.TestCheckResourceAttr("data.azureipam_spaces.test", "spaces.0.blocks.1.vnets.#", "2"),
					resource.TestCheckResourceAttr("data.azureipam_spaces.test", "spaces.0.blocks.1.vnets.0.id", "/subscriptions/554a3781-68fa-4429-b46a-f86664d836b4/resourceGroups/RG-WE-ALL-COMMS-01/providers/Microsoft.Network/virtualNetworks/vnet-we-d-terratest-hub-01"),
					resource.TestCheckResourceAttr("data.azureipam_spaces.test", "spaces.0.blocks.1.vnets.1.id", "/subscriptions/49a789bb-9a8d-42be-8465-53aa77552366/resourceGroups/RG-WE-ALL-COMMS-01/providers/Microsoft.Network/virtualNetworks/vnet-we-a-testzavd-01"),
					resource.TestCheckResourceAttr("data.azureipam_spaces.test", "spaces.0.blocks.1.reservations.#", "2"),
					//second block, first reservation
					resource.TestCheckResourceAttr("data.azureipam_spaces.test", "spaces.0.blocks.1.reservations.0.id", "YYtppsvYQsRSBpZLsioZSV"),
					resource.TestCheckResourceAttr("data.azureipam_spaces.test", "spaces.0.blocks.1.reservations.0.cidr", "10.82.6.0/23"),
					resource.TestCheckResourceAttr("data.azureipam_spaces.test", "spaces.0.blocks.1.reservations.0.description", "acceptance-test"),
					resource.TestCheckResourceAttr("data.azureipam_spaces.test", "spaces.0.blocks.1.reservations.0.created_on", "2024-09-07T06:21:42+02:00"),
					resource.TestCheckResourceAttr("data.azureipam_spaces.test", "spaces.0.blocks.1.reservations.0.created_by", "dummyemail@gmail.com"),
					resource.TestCheckNoResourceAttr("data.azureipam_spaces.test", "spaces.0.blocks.1.reservations.0.settled_on"),
					resource.TestCheckNoResourceAttr("data.azureipam_spaces.test", "spaces.0.blocks.1.reservations.0.settled_by"),
					resource.TestCheckResourceAttr("data.azureipam_spaces.test", "spaces.0.blocks.1.reservations.0.status", "wait"),
					//second block, second reservation
					resource.TestCheckResourceAttr("data.azureipam_spaces.test", "spaces.0.blocks.1.reservations.1.id", "hi3fxt9PeSpxhykfSszVUb"),
					resource.TestCheckResourceAttr("data.azureipam_spaces.test", "spaces.0.blocks.1.reservations.1.cidr", "10.82.1.160/27"),
					resource.TestCheckResourceAttr("data.azureipam_spaces.test", "spaces.0.blocks.1.reservations.1.description", "vnet-we-c-arq3tier-01"),
					resource.TestCheckResourceAttr("data.azureipam_spaces.test", "spaces.0.blocks.1.reservations.1.created_on", "2023-11-08T13:51:07+01:00"),
					resource.TestCheckResourceAttr("data.azureipam_spaces.test", "spaces.0.blocks.1.reservations.1.created_by", "spn:9fc2493a-b515-49a6-9d73-93e1bac5f6cc"),
					resource.TestCheckResourceAttr("data.azureipam_spaces.test", "spaces.0.blocks.1.reservations.1.settled_on", "2024-04-03T09:18:41+02:00"),
					resource.TestCheckResourceAttr("data.azureipam_spaces.test", "spaces.0.blocks.1.reservations.1.settled_by", "dummyemail@gmail.com"),
					resource.TestCheckResourceAttr("data.azureipam_spaces.test", "spaces.0.blocks.1.reservations.1.status", "cancelledByUser"),
				),
			},
		},
	})
}

func TestAccSpacesWithUtilizationAndVnetDataSource(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("GET", "https://mockedHost.azurewebsites.net/api/spaces?expand=true&utilization=true",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/datasource/spaces/spaces_with_utilization_and_vnet.json").String()), nil
		})

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccProviderConfig + `data "azureipam_spaces" "test" {
					expand             = true
  					append_utilization = true
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					//Verify common attributes to ensure that all are set
					resource.TestCheckResourceAttr("data.azureipam_spaces.test", "append_utilization", "true"),
					resource.TestCheckResourceAttr("data.azureipam_spaces.test", "expand", "true"),
					// Verify number of spaces returned
					resource.TestCheckResourceAttr("data.azureipam_spaces.test", "spaces.#", "1"),
					// Verify the first space to ensure all attributes are set
					resource.TestCheckResourceAttr("data.azureipam_spaces.test", "spaces.0.name", "au"),
					resource.TestCheckResourceAttr("data.azureipam_spaces.test", "spaces.0.description", "Australia"),
					resource.TestCheckResourceAttr("data.azureipam_spaces.test", "spaces.0.size", "131072"),
					resource.TestCheckResourceAttr("data.azureipam_spaces.test", "spaces.0.used", "544"),

					// Verify number of blocks returned in first space
					resource.TestCheckResourceAttr("data.azureipam_spaces.test", "spaces.0.blocks.#", "2"),
					// Verify first block to ensure all attributes are set
					resource.TestCheckResourceAttr("data.azureipam_spaces.test", "spaces.0.blocks.0.name", "AustraliaSoutheast"),
					resource.TestCheckResourceAttr("data.azureipam_spaces.test", "spaces.0.blocks.0.cidr", "10.83.0.0/16"),
					resource.TestCheckResourceAttr("data.azureipam_spaces.test", "spaces.0.blocks.0.size", "65536"),
					resource.TestCheckResourceAttr("data.azureipam_spaces.test", "spaces.0.blocks.0.used", "256"),
					resource.TestCheckResourceAttr("data.azureipam_spaces.test", "spaces.0.blocks.0.vnets.#", "0"),
					resource.TestCheckResourceAttr("data.azureipam_spaces.test", "spaces.0.blocks.0.reservations.#", "0"),
					resource.TestCheckResourceAttr("data.azureipam_spaces.test", "spaces.0.blocks.0.externals.#", "1"),
					resource.TestCheckResourceAttr("data.azureipam_spaces.test", "spaces.0.blocks.0.externals.0.name", "test"),
					resource.TestCheckResourceAttr("data.azureipam_spaces.test", "spaces.0.blocks.0.externals.0.description", "external for tests"),
					resource.TestCheckResourceAttr("data.azureipam_spaces.test", "spaces.0.blocks.0.externals.0.cidr", "10.83.0.0/24"),
					
					// Verify second block to ensure all attributes are set
					resource.TestCheckResourceAttr("data.azureipam_spaces.test", "spaces.0.blocks.1.name", "AustraliaEast"),
					resource.TestCheckResourceAttr("data.azureipam_spaces.test", "spaces.0.blocks.1.cidr", "10.82.0.0/16"),
					resource.TestCheckResourceAttr("data.azureipam_spaces.test", "spaces.0.blocks.1.size", "65536"),
					resource.TestCheckResourceAttr("data.azureipam_spaces.test", "spaces.0.blocks.1.used", "288"),
					resource.TestCheckResourceAttr("data.azureipam_spaces.test", "spaces.0.blocks.1.externals.#", "0"),
					resource.TestCheckResourceAttr("data.azureipam_spaces.test", "spaces.0.blocks.1.vnets.#", "2"),
					//second block, first vnet
					resource.TestCheckResourceAttr("data.azureipam_spaces.test", "spaces.0.blocks.1.vnets.0.id", "/subscriptions/554a3781-68fa-4429-b46a-f86664d836b4/resourceGroups/RG-WE-ALL-COMMS-01/providers/Microsoft.Network/virtualNetworks/vnet-we-d-terratest-hub-01"),
					resource.TestCheckResourceAttr("data.azureipam_spaces.test", "spaces.0.blocks.1.vnets.0.name", "vnet-we-d-terratest-hub-01"),
					resource.TestCheckResourceAttr("data.azureipam_spaces.test", "spaces.0.blocks.1.vnets.0.prefixes.#",  "1"),
					resource.TestCheckResourceAttr("data.azureipam_spaces.test", "spaces.0.blocks.1.vnets.0.prefixes.0",  "10.82.0.0/24"),
					resource.TestCheckResourceAttr("data.azureipam_spaces.test", "spaces.0.blocks.1.vnets.0.resource_group", "rg-we-all-comms-01"),
					resource.TestCheckResourceAttr("data.azureipam_spaces.test", "spaces.0.blocks.1.vnets.0.subscription_id", "554a3781-68fa-4429-b46a-f86664d836b4"),
					resource.TestCheckResourceAttr("data.azureipam_spaces.test", "spaces.0.blocks.1.vnets.0.tenant_id", "562871ed-4b14-46b8-899a-0c5ffa5e4272"),
					resource.TestCheckResourceAttr("data.azureipam_spaces.test", "spaces.0.blocks.1.vnets.0.size", "256"),
					resource.TestCheckResourceAttr("data.azureipam_spaces.test", "spaces.0.blocks.1.vnets.0.used", "144"),
					resource.TestCheckResourceAttr("data.azureipam_spaces.test", "spaces.0.blocks.1.vnets.0.subnets.#", "2"),
					//second block, first vnet, first subnet
					resource.TestCheckResourceAttr("data.azureipam_spaces.test", "spaces.0.blocks.1.vnets.0.subnets.0.name", "main"),
					resource.TestCheckResourceAttr("data.azureipam_spaces.test", "spaces.0.blocks.1.vnets.0.subnets.0.prefix", "10.82.0.0/25"),
					resource.TestCheckResourceAttr("data.azureipam_spaces.test", "spaces.0.blocks.1.vnets.0.subnets.0.size", "128"),
					resource.TestCheckResourceAttr("data.azureipam_spaces.test", "spaces.0.blocks.1.vnets.0.subnets.0.used", "5"),
					//second block, first vnet, second subnet
					resource.TestCheckResourceAttr("data.azureipam_spaces.test", "spaces.0.blocks.1.vnets.0.subnets.1.name", "GatewaySubnet"),
					resource.TestCheckResourceAttr("data.azureipam_spaces.test", "spaces.0.blocks.1.vnets.0.subnets.1.prefix", "10.82.0.128/28"),
					resource.TestCheckResourceAttr("data.azureipam_spaces.test", "spaces.0.blocks.1.vnets.0.subnets.1.size", "16"),
					resource.TestCheckResourceAttr("data.azureipam_spaces.test", "spaces.0.blocks.1.vnets.0.subnets.1.used", "6"),

					//second block, second vnet
					resource.TestCheckResourceAttr("data.azureipam_spaces.test", "spaces.0.blocks.1.vnets.1.id", "/subscriptions/49a789bb-9a8d-42be-8465-53aa77552366/resourceGroups/RG-WE-ALL-COMMS-01/providers/Microsoft.Network/virtualNetworks/vnet-we-a-testzavd-01"),
					resource.TestCheckResourceAttr("data.azureipam_spaces.test", "spaces.0.blocks.1.vnets.1.name", "vnet-we-a-testzavd-01"),
					resource.TestCheckResourceAttr("data.azureipam_spaces.test", "spaces.0.blocks.1.vnets.1.prefixes.#",  "1"),
					resource.TestCheckResourceAttr("data.azureipam_spaces.test", "spaces.0.blocks.1.vnets.1.prefixes.0",  "10.82.1.224/27"),
					resource.TestCheckResourceAttr("data.azureipam_spaces.test", "spaces.0.blocks.1.vnets.1.resource_group", "rg-we-all-comms-01"),
					resource.TestCheckResourceAttr("data.azureipam_spaces.test", "spaces.0.blocks.1.vnets.1.subscription_id", "49a789bb-9a8d-42be-8465-53aa77552366"),
					resource.TestCheckResourceAttr("data.azureipam_spaces.test", "spaces.0.blocks.1.vnets.1.tenant_id", "562871ed-4b14-46b8-899a-0c5ffa5e4272"),
					resource.TestCheckResourceAttr("data.azureipam_spaces.test", "spaces.0.blocks.1.vnets.1.size", "32"),
					resource.TestCheckResourceAttr("data.azureipam_spaces.test", "spaces.0.blocks.1.vnets.1.used", "24"),
					resource.TestCheckResourceAttr("data.azureipam_spaces.test", "spaces.0.blocks.1.vnets.1.subnets.#", "2"),
					//second block, second vnet, first subnet
					resource.TestCheckResourceAttr("data.azureipam_spaces.test", "spaces.0.blocks.1.vnets.1.subnets.0.name", "snet-we-a-private-01"),
					resource.TestCheckResourceAttr("data.azureipam_spaces.test", "spaces.0.blocks.1.vnets.1.subnets.0.prefix", "10.82.1.224/29"),
					resource.TestCheckResourceAttr("data.azureipam_spaces.test", "spaces.0.blocks.1.vnets.1.subnets.0.size", "8"),
					resource.TestCheckResourceAttr("data.azureipam_spaces.test", "spaces.0.blocks.1.vnets.1.subnets.0.used", "5"),
					//second block, second vnet, second subnet
					resource.TestCheckResourceAttr("data.azureipam_spaces.test", "spaces.0.blocks.1.vnets.1.subnets.1.name", "snet-we-a-sessionhost-01"),
					resource.TestCheckResourceAttr("data.azureipam_spaces.test", "spaces.0.blocks.1.vnets.1.subnets.1.prefix", "10.82.1.240/28"),
					resource.TestCheckResourceAttr("data.azureipam_spaces.test", "spaces.0.blocks.1.vnets.1.subnets.1.size", "16"),
					resource.TestCheckResourceAttr("data.azureipam_spaces.test", "spaces.0.blocks.1.vnets.1.subnets.1.used", "5"),
					//second block reservations
					resource.TestCheckResourceAttr("data.azureipam_spaces.test", "spaces.0.blocks.1.reservations.#", "2"),
					//second block, first reservation
					resource.TestCheckResourceAttr("data.azureipam_spaces.test", "spaces.0.blocks.1.reservations.0.id", "YYtppsvYQsRSBpZLsioZSV"),
					resource.TestCheckResourceAttr("data.azureipam_spaces.test", "spaces.0.blocks.1.reservations.0.cidr", "10.82.6.0/23"),
					resource.TestCheckResourceAttr("data.azureipam_spaces.test", "spaces.0.blocks.1.reservations.0.description", "acceptance-test"),
					resource.TestCheckResourceAttr("data.azureipam_spaces.test", "spaces.0.blocks.1.reservations.0.created_on", "2024-09-07T06:21:42+02:00"),
					resource.TestCheckResourceAttr("data.azureipam_spaces.test", "spaces.0.blocks.1.reservations.0.created_by", "dummyemail@gmail.com"),
					resource.TestCheckNoResourceAttr("data.azureipam_spaces.test", "spaces.0.blocks.1.reservations.0.settled_on"),
					resource.TestCheckNoResourceAttr("data.azureipam_spaces.test", "spaces.0.blocks.1.reservations.0.settled_by"),
					resource.TestCheckResourceAttr("data.azureipam_spaces.test", "spaces.0.blocks.1.reservations.0.status", "wait"),
					//second block, second reservation
					resource.TestCheckResourceAttr("data.azureipam_spaces.test", "spaces.0.blocks.1.reservations.1.id", "hi3fxt9PeSpxhykfSszVUb"),
					resource.TestCheckResourceAttr("data.azureipam_spaces.test", "spaces.0.blocks.1.reservations.1.cidr", "10.82.1.160/27"),
					resource.TestCheckResourceAttr("data.azureipam_spaces.test", "spaces.0.blocks.1.reservations.1.description", "vnet-we-c-arq3tier-01"),
					resource.TestCheckResourceAttr("data.azureipam_spaces.test", "spaces.0.blocks.1.reservations.1.created_on", "2023-11-08T13:51:07+01:00"),
					resource.TestCheckResourceAttr("data.azureipam_spaces.test", "spaces.0.blocks.1.reservations.1.created_by", "spn:9fc2493a-b515-49a6-9d73-93e1bac5f6cc"),
					resource.TestCheckResourceAttr("data.azureipam_spaces.test", "spaces.0.blocks.1.reservations.1.settled_on", "2024-04-03T09:18:41+02:00"),
					resource.TestCheckResourceAttr("data.azureipam_spaces.test", "spaces.0.blocks.1.reservations.1.settled_by", "dummyemail@gmail.com"),
					resource.TestCheckResourceAttr("data.azureipam_spaces.test", "spaces.0.blocks.1.reservations.1.status", "cancelledByUser"),
				),
			},
		},
	})
}
