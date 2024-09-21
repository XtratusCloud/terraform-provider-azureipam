package provider

import (
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

func TestAccSpaceResource(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("POST", "https://mockedHost.azurewebsites.net/api/spaces",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/space/new_space.json").String()), nil
		})
	httpmock.RegisterResponder("GET", "https://mockedHost.azurewebsites.net/api/spaces/as?expand=false&utilization=false",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/space/new_space.json").String()), nil
		})
	httpmock.RegisterResponder("PATCH", "https://mockedHost.azurewebsites.net/api/spaces/as",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/space/updated_space.json").String()), nil
		})
	httpmock.RegisterResponder("GET", "https://mockedHost.azurewebsites.net/api/spaces/asia?expand=false&utilization=false",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/space/updated_space.json").String()), nil
		})
	httpmock.RegisterResponder("DELETE", "https://mockedHost.azurewebsites.net/api/spaces/asia?force=true",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, ""), nil
		})

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccProviderConfig + `resource "azureipam_space" "test" {
					name = "as"
					description = "Asia"
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					//Verify common attributes to ensure that all are set
					resource.TestCheckResourceAttr("azureipam_space.test", "name", "as"),
					resource.TestCheckResourceAttr("azureipam_space.test", "description", "Asia"),
				),
			},
			// ImportState testing
			{
				ResourceName:                         "azureipam_space.test",
				ImportState:                          true,
				ImportStateId:                        "as",
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "name",
			},
			// Update and Read testing
			{
				Config: testAccProviderConfig + `resource "azureipam_space" "test" {
					name = "asia"
					description = "Asia Description"
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					//Verify attributes after update to ensure that all are set
					resource.TestCheckResourceAttr("azureipam_space.test", "name", "asia"),
					resource.TestCheckResourceAttr("azureipam_space.test", "description", "Asia Description"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
