package provider

import (
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

func TestAccExternalResource(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("POST", "https://mockedHost.azurewebsites.net/api/spaces/au/blocks/AustraliaSoutheast/externals",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/external/externals_with_new_external.json").String()), nil
		})
	httpmock.RegisterResponder("GET", "https://mockedHost.azurewebsites.net/api/spaces/au/blocks/AustraliaSoutheast/externals/acctest",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/external/new_external.json").String()), nil
		})
	httpmock.RegisterResponder("GET", "https://mockedHost.azurewebsites.net/api/spaces/au/blocks/AustraliaSoutheast/externals",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/external/externals_with_new_external.json").String()), nil
		})
	httpmock.RegisterResponder("PUT", "https://mockedHost.azurewebsites.net/api/spaces/au/blocks/AustraliaSoutheast/externals",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/external/externals_with_updated_external.json").String()), nil
		})
	httpmock.RegisterResponder("GET", "https://mockedHost.azurewebsites.net/api/spaces/au/blocks/AustraliaSoutheast/externals/acctestupdated",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/external/updated_external.json").String()), nil
		})
	httpmock.RegisterResponder("DELETE", "https://mockedHost.azurewebsites.net/api/spaces/au/blocks/AustraliaSoutheast/externals/acctestupdated",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, ""), nil
		})

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccProviderConfig + `resource "azureipam_external" "test" {
					space = "au"
					block = "AustraliaSoutheast"
					name = "acctest"
					description = "External Network for Acceptance Tests"
					cidr = "10.83.1.0/24"
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					//Verify common attributes to ensure that all are set
					resource.TestCheckResourceAttr("azureipam_external.test", "space", "au"),
					resource.TestCheckResourceAttr("azureipam_external.test", "block", "AustraliaSoutheast"),
					resource.TestCheckResourceAttr("azureipam_external.test", "name", "acctest"),
					resource.TestCheckResourceAttr("azureipam_external.test", "description", "External Network for Acceptance Tests"),
					resource.TestCheckResourceAttr("azureipam_external.test", "cidr", "10.83.1.0/24"),
				),
			},
			// ImportState testing
			{
				ResourceName:                         "azureipam_external.test",
				ImportState:                          true,
				ImportStateId:                        "au/AustraliaSoutheast/acctest",
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "name",
			},
			// Update and Read testing
			{
				Config: testAccProviderConfig + `resource "azureipam_external" "test" {
					space = "au"
					block = "AustraliaSoutheast"
					name = "acctestupdated"
					description = "External Network for Acceptance Tests Updated"
					cidr = "10.83.10.0/24"
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					//Verify attributes after update to ensure that all are set
					resource.TestCheckResourceAttr("azureipam_external.test", "space", "au"),
					resource.TestCheckResourceAttr("azureipam_external.test", "block", "AustraliaSoutheast"),
					resource.TestCheckResourceAttr("azureipam_external.test", "name", "acctestupdated"),
					resource.TestCheckResourceAttr("azureipam_external.test", "description", "External Network for Acceptance Tests Updated"),
					resource.TestCheckResourceAttr("azureipam_external.test", "cidr", "10.83.10.0/24"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
