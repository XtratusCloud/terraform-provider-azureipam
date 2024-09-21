package provider

import (
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

func TestAccBlockResource(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("POST", "https://mockedHost.azurewebsites.net/api/spaces/au/blocks",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/block/new_block.json").String()), nil
		})
	httpmock.RegisterResponder("GET", "https://mockedHost.azurewebsites.net/api/spaces/au/blocks/AustraliaNorth?expand=false&utilization=false",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/block/new_block.json").String()), nil
		})
	httpmock.RegisterResponder("PATCH", "https://mockedHost.azurewebsites.net/api/spaces/au/blocks/AustraliaNorth",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/block/updated_block.json").String()), nil
		})
	httpmock.RegisterResponder("GET", "https://mockedHost.azurewebsites.net/api/spaces/au/blocks/AustraliaWest?expand=false&utilization=false",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/block/updated_block.json").String()), nil
		})
	httpmock.RegisterResponder("DELETE", "https://mockedHost.azurewebsites.net/api/spaces/au/blocks/AustraliaWest?force=true",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, ""), nil
		})

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccProviderConfig + `resource "azureipam_block" "test" {
					space = "au"
					name = "AustraliaNorth"
					cidr = "10.85.0.0/16"
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					//Verify common attributes to ensure that all are set
					resource.TestCheckResourceAttr("azureipam_block.test", "space", "au"),
					resource.TestCheckResourceAttr("azureipam_block.test", "name", "AustraliaNorth"),
					resource.TestCheckResourceAttr("azureipam_block.test", "cidr", "10.85.0.0/16"),
				),
			},
			// ImportState testing
			{
				ResourceName:                         "azureipam_block.test",
				ImportState:                          true,
				ImportStateId:                        "au/AustraliaNorth",
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "name",
			},
			// Update and Read testing
			{
				Config: testAccProviderConfig + `resource "azureipam_block" "test" {
					space = "au"
					name = "AustraliaWest"
					cidr = "10.86.0.0/16"
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					//Verify attributes after update to ensure that all are set
					resource.TestCheckResourceAttr("azureipam_block.test", "space", "au"),
					resource.TestCheckResourceAttr("azureipam_block.test", "name", "AustraliaWest"),
					resource.TestCheckResourceAttr("azureipam_block.test", "cidr", "10.86.0.0/16"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
