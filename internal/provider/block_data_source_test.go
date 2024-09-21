package provider

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

func TestAccBlockWithoutUtilizationAndVnetDataSource(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("GET", "https://mockedHost.azurewebsites.net/api/spaces/au/blocks/AustraliaEast?expand=false&utilization=false",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/datasource/blocks/block_without_utilization_and_vnet.json").String()), nil
		})

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccProviderConfig + `data "azureipam_block" "test" {
					space = "au"
					name  = "AustraliaEast"
					expand             = false
  					append_utilization = false
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					//Verify common attributes to ensure that all are set
					resource.TestCheckResourceAttr("data.azureipam_block.test", "space", "au"),
					resource.TestCheckResourceAttr("data.azureipam_block.test", "append_utilization", "false"),
					resource.TestCheckResourceAttr("data.azureipam_block.test", "expand", "false"),
					resource.TestCheckResourceAttr("data.azureipam_block.test", "name", "AustraliaEast"),
					resource.TestCheckResourceAttr("data.azureipam_block.test", "cidr", "10.82.0.0/16"),
					resource.TestCheckNoResourceAttr("data.azureipam_block.test", "size"),
					resource.TestCheckNoResourceAttr("data.azureipam_block.test", "used"),
					resource.TestCheckResourceAttr("data.azureipam_block.test", "externals.#", "0"),
					resource.TestCheckResourceAttr("data.azureipam_block.test", "vnets.#", "2"),
					resource.TestCheckResourceAttr("data.azureipam_block.test", "vnets.0.id", "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/RG-WE-ALL-COMMS-01/providers/Microsoft.Network/virtualNetworks/vnet-we-d-terratest-hub-01"),
					resource.TestCheckResourceAttr("data.azureipam_block.test", "vnets.1.id", "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/RG-WE-ALL-COMMS-01/providers/Microsoft.Network/virtualNetworks/vnet-we-a-testzavd-01"),
					resource.TestCheckResourceAttr("data.azureipam_block.test", "reservations.#", "2"),
					//first reservation
					resource.TestCheckResourceAttr("data.azureipam_block.test", "reservations.0.id", "YYtppsvYQsRSBpZLsioZSV"),
					resource.TestCheckResourceAttr("data.azureipam_block.test", "reservations.0.cidr", "10.82.6.0/23"),
					resource.TestCheckResourceAttr("data.azureipam_block.test", "reservations.0.description", "acceptance-test"),
					resource.TestCheckResourceAttrWith("data.azureipam_block.test", "reservations.0.created_on", func(value string) error {
						expected, _ := time.Parse(time.RFC3339, "2024-09-07T06:21:42+02:00")
						current, _ := time.Parse(time.RFC3339, value)
						if current.Equal(expected) {
							return nil
						}
						return errors.New("expected " + expected.String() + " got " + current.String())
					}),
					resource.TestCheckResourceAttr("data.azureipam_block.test", "reservations.0.created_by", "dummyemail@gmail.com"),
					resource.TestCheckNoResourceAttr("data.azureipam_block.test", "reservations.0.settled_on"),
					resource.TestCheckNoResourceAttr("data.azureipam_block.test", "reservations.0.settled_by"),
					resource.TestCheckResourceAttr("data.azureipam_block.test", "reservations.0.status", "wait"),
					//second reservation
					resource.TestCheckResourceAttr("data.azureipam_block.test", "reservations.1.id", "hi3fxt9PeSpxhykfSszVUb"),
					resource.TestCheckResourceAttr("data.azureipam_block.test", "reservations.1.cidr", "10.82.1.160/27"),
					resource.TestCheckResourceAttr("data.azureipam_block.test", "reservations.1.description", "vnet-we-c-arq3tier-01"),
					resource.TestCheckResourceAttrWith("data.azureipam_block.test", "reservations.1.created_on", func(value string) error {
						expected, _ := time.Parse(time.RFC3339, "2023-11-08T13:51:07+01:00")
						current, _ := time.Parse(time.RFC3339, value)
						if current.Equal(expected) {
							return nil
						}
						return errors.New("expected " + expected.String() + " got " + current.String())
					}),
					resource.TestCheckResourceAttr("data.azureipam_block.test", "reservations.1.created_by", "spn:9fc2493a-b515-49a6-9d73-93e1bac5f6cc"),
					resource.TestCheckResourceAttrWith("data.azureipam_block.test", "reservations.1.settled_on", func(value string) error {
						expected, _ := time.Parse(time.RFC3339, "2024-04-03T09:18:41+02:00")
						current, _ := time.Parse(time.RFC3339, value)
						if current.Equal(expected) {
							return nil
						}
						return errors.New("expected " + expected.String() + " got " + current.String())
					}),
					resource.TestCheckResourceAttr("data.azureipam_block.test", "reservations.1.settled_by", "dummyemail@gmail.com"),
					resource.TestCheckResourceAttr("data.azureipam_block.test", "reservations.1.status", "cancelledByUser"),
				),
			},
		},
	})
}

func TestAccBlocWithUtilizationAndVnetDataSource(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("GET", "https://mockedHost.azurewebsites.net/api/spaces/au/blocks/AustraliaEast?expand=true&utilization=true",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/datasource/blocks/block_with_utilization_and_vnet.json").String()), nil
		})

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccProviderConfig + `data "azureipam_block" "test" {
					space = "au"
					name  = "AustraliaEast"
					expand             = true
  					append_utilization = true
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					//Verify common attributes to ensure that all are set
					resource.TestCheckResourceAttr("data.azureipam_block.test", "space", "au"),
					resource.TestCheckResourceAttr("data.azureipam_block.test", "append_utilization", "true"),
					resource.TestCheckResourceAttr("data.azureipam_block.test", "expand", "true"),
					resource.TestCheckResourceAttr("data.azureipam_block.test", "name", "AustraliaEast"),
					resource.TestCheckResourceAttr("data.azureipam_block.test", "cidr", "10.82.0.0/16"),
					resource.TestCheckResourceAttr("data.azureipam_block.test", "size", "65536"),
					resource.TestCheckResourceAttr("data.azureipam_block.test", "used", "288"),
					resource.TestCheckResourceAttr("data.azureipam_block.test", "externals.#", "0"),
					resource.TestCheckResourceAttr("data.azureipam_block.test", "vnets.#", "2"),
					//first vnet
					resource.TestCheckResourceAttr("data.azureipam_block.test", "vnets.0.id", "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/RG-WE-ALL-COMMS-01/providers/Microsoft.Network/virtualNetworks/vnet-we-d-terratest-hub-01"),
					resource.TestCheckResourceAttr("data.azureipam_block.test", "vnets.0.name", "vnet-we-d-terratest-hub-01"),
					resource.TestCheckResourceAttr("data.azureipam_block.test", "vnets.0.prefixes.#", "1"),
					resource.TestCheckResourceAttr("data.azureipam_block.test", "vnets.0.prefixes.0", "10.82.0.0/24"),
					resource.TestCheckResourceAttr("data.azureipam_block.test", "vnets.0.resource_group", "rg-we-all-comms-01"),
					resource.TestCheckResourceAttr("data.azureipam_block.test", "vnets.0.subscription_id", "00000000-0000-0000-0000-000000000000"),
					resource.TestCheckResourceAttr("data.azureipam_block.test", "vnets.0.tenant_id", "11111111-1111-1111-1111-111111111111"),
					resource.TestCheckResourceAttr("data.azureipam_block.test", "vnets.0.size", "256"),
					resource.TestCheckResourceAttr("data.azureipam_block.test", "vnets.0.used", "144"),
					resource.TestCheckResourceAttr("data.azureipam_block.test", "vnets.0.subnets.#", "2"),
					//first vnet, first subnet
					resource.TestCheckResourceAttr("data.azureipam_block.test", "vnets.0.subnets.0.name", "main"),
					resource.TestCheckResourceAttr("data.azureipam_block.test", "vnets.0.subnets.0.prefix", "10.82.0.0/25"),
					resource.TestCheckResourceAttr("data.azureipam_block.test", "vnets.0.subnets.0.size", "128"),
					resource.TestCheckResourceAttr("data.azureipam_block.test", "vnets.0.subnets.0.used", "5"),
					//first vnet, second subnet
					resource.TestCheckResourceAttr("data.azureipam_block.test", "vnets.0.subnets.1.name", "GatewaySubnet"),
					resource.TestCheckResourceAttr("data.azureipam_block.test", "vnets.0.subnets.1.prefix", "10.82.0.128/28"),
					resource.TestCheckResourceAttr("data.azureipam_block.test", "vnets.0.subnets.1.size", "16"),
					resource.TestCheckResourceAttr("data.azureipam_block.test", "vnets.0.subnets.1.used", "6"),

					//second vnet
					resource.TestCheckResourceAttr("data.azureipam_block.test", "vnets.1.id", "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/RG-WE-ALL-COMMS-01/providers/Microsoft.Network/virtualNetworks/vnet-we-a-testzavd-01"),
					resource.TestCheckResourceAttr("data.azureipam_block.test", "vnets.1.name", "vnet-we-a-testzavd-01"),
					resource.TestCheckResourceAttr("data.azureipam_block.test", "vnets.1.prefixes.#", "1"),
					resource.TestCheckResourceAttr("data.azureipam_block.test", "vnets.1.prefixes.0", "10.82.1.224/27"),
					resource.TestCheckResourceAttr("data.azureipam_block.test", "vnets.1.resource_group", "rg-we-all-comms-01"),
					resource.TestCheckResourceAttr("data.azureipam_block.test", "vnets.1.subscription_id", "00000000-0000-0000-0000-000000000000"),
					resource.TestCheckResourceAttr("data.azureipam_block.test", "vnets.1.tenant_id", "11111111-1111-1111-1111-111111111111"),
					resource.TestCheckResourceAttr("data.azureipam_block.test", "vnets.1.size", "32"),
					resource.TestCheckResourceAttr("data.azureipam_block.test", "vnets.1.used", "24"),
					resource.TestCheckResourceAttr("data.azureipam_block.test", "vnets.1.subnets.#", "2"),
					//second vnet, first subnet
					resource.TestCheckResourceAttr("data.azureipam_block.test", "vnets.1.subnets.0.name", "snet-we-a-private-01"),
					resource.TestCheckResourceAttr("data.azureipam_block.test", "vnets.1.subnets.0.prefix", "10.82.1.224/29"),
					resource.TestCheckResourceAttr("data.azureipam_block.test", "vnets.1.subnets.0.size", "8"),
					resource.TestCheckResourceAttr("data.azureipam_block.test", "vnets.1.subnets.0.used", "5"),
					//second vnet, second subnet
					resource.TestCheckResourceAttr("data.azureipam_block.test", "vnets.1.subnets.1.name", "snet-we-a-sessionhost-01"),
					resource.TestCheckResourceAttr("data.azureipam_block.test", "vnets.1.subnets.1.prefix", "10.82.1.240/28"),
					resource.TestCheckResourceAttr("data.azureipam_block.test", "vnets.1.subnets.1.size", "16"),
					resource.TestCheckResourceAttr("data.azureipam_block.test", "vnets.1.subnets.1.used", "5"),
					//reservations
					resource.TestCheckResourceAttr("data.azureipam_block.test", "reservations.#", "2"),
					//first reservation
					resource.TestCheckResourceAttr("data.azureipam_block.test", "reservations.0.id", "YYtppsvYQsRSBpZLsioZSV"),
					resource.TestCheckResourceAttr("data.azureipam_block.test", "reservations.0.cidr", "10.82.6.0/23"),
					resource.TestCheckResourceAttr("data.azureipam_block.test", "reservations.0.description", "acceptance-test"),
					resource.TestCheckResourceAttrWith("data.azureipam_block.test", "reservations.0.created_on", func(value string) error {
						expected, _ := time.Parse(time.RFC3339, "2024-09-07T06:21:42+02:00")
						current, _ := time.Parse(time.RFC3339, value)
						if current.Equal(expected) {
							return nil
						}
						return errors.New("expected " + expected.String() + " got " + current.String())
					}),
					resource.TestCheckResourceAttr("data.azureipam_block.test", "reservations.0.created_by", "dummyemail@gmail.com"),
					resource.TestCheckNoResourceAttr("data.azureipam_block.test", "reservations.0.settled_on"),
					resource.TestCheckNoResourceAttr("data.azureipam_block.test", "reservations.0.settled_by"),
					resource.TestCheckResourceAttr("data.azureipam_block.test", "reservations.0.status", "wait"),
					//second reservation
					resource.TestCheckResourceAttr("data.azureipam_block.test", "reservations.1.id", "hi3fxt9PeSpxhykfSszVUb"),
					resource.TestCheckResourceAttr("data.azureipam_block.test", "reservations.1.cidr", "10.82.1.160/27"),
					resource.TestCheckResourceAttr("data.azureipam_block.test", "reservations.1.description", "vnet-we-c-arq3tier-01"),
					resource.TestCheckResourceAttrWith("data.azureipam_block.test", "reservations.1.created_on", func(value string) error {
						expected, _ := time.Parse(time.RFC3339, "2023-11-08T13:51:07+01:00")
						current, _ := time.Parse(time.RFC3339, value)
						if current.Equal(expected) {
							return nil
						}
						return errors.New("expected " + expected.String() + " got " + current.String())
					}),
					resource.TestCheckResourceAttr("data.azureipam_block.test", "reservations.1.created_by", "spn:9fc2493a-b515-49a6-9d73-93e1bac5f6cc"),
					resource.TestCheckResourceAttrWith("data.azureipam_block.test", "reservations.1.settled_on", func(value string) error {
						expected, _ := time.Parse(time.RFC3339, "2024-04-03T09:18:41+02:00")
						current, _ := time.Parse(time.RFC3339, value)
						if current.Equal(expected) {
							return nil
						}
						return errors.New("expected " + expected.String() + " got " + current.String())
					}),
					resource.TestCheckResourceAttr("data.azureipam_block.test", "reservations.1.settled_by", "dummyemail@gmail.com"),
					resource.TestCheckResourceAttr("data.azureipam_block.test", "reservations.1.status", "cancelledByUser"),
				),
			},
		},
	})
}
