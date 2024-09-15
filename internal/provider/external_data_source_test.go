package provider

import (
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

func TestAccRExternalDataSource(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("GET", "https://mockedHost.azurewebsites.net/api/spaces/au/blocks/AustraliaSoutheast/externals/acctest",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/datasource/externals/external_one.json").String()), nil
		})

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccProviderConfig + `data "azureipam_external" "test" {
					space = "au"
					block = "AustraliaSoutheast"
					name = "acctest"
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					//Verify common attributes to ensure that all are set
					resource.TestCheckResourceAttr("data.azureipam_external.test", "space", "au"),
					resource.TestCheckResourceAttr("data.azureipam_external.test", "block", "AustraliaSoutheast"),
					resource.TestCheckResourceAttr("data.azureipam_external.test", "name", "acctest"),
					resource.TestCheckResourceAttr("data.azureipam_external.test", "description", "External Network for Acceptance Tests"),
					resource.TestCheckResourceAttr("data.azureipam_external.test", "cidr", "10.83.1.0/24"),
				),
			},
		},
	})
}
