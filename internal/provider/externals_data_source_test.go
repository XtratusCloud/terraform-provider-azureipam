package provider

import (
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

func TestAccRExternalsDataSource(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("GET", "https://mockedHost.azurewebsites.net/api/spaces/au/blocks/AustraliaSoutheast/externals",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/datasource/externals/externals_all.json").String()), nil
		})

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccProviderConfig + `data "azureipam_externals" "test" {
					space = "au"
					block = "AustraliaSoutheast"
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					//Verify common attributes to ensure that all are set
					resource.TestCheckResourceAttr("data.azureipam_externals.test", "space", "au"),
					resource.TestCheckResourceAttr("data.azureipam_externals.test", "block", "AustraliaSoutheast"),
					// Verify number of externals returned
					resource.TestCheckResourceAttr("data.azureipam_externals.test", "externals.#", "2"),
					// Verify the first external to ensure all attributes are set
					resource.TestCheckResourceAttr("data.azureipam_externals.test", "externals.0.name", "prueba"),
					resource.TestCheckResourceAttr("data.azureipam_externals.test", "externals.0.description", "descripcion prueba"),
					resource.TestCheckResourceAttr("data.azureipam_externals.test", "externals.0.cidr", "10.83.0.0/24"),

					// Verify the second external to ensure all attributes are set
					resource.TestCheckResourceAttr("data.azureipam_externals.test", "externals.1.name", "acctest"),
					resource.TestCheckResourceAttr("data.azureipam_externals.test", "externals.1.description", "External Network for Acceptance Tests"),
					resource.TestCheckResourceAttr("data.azureipam_externals.test", "externals.1.cidr", "10.83.1.0/24"),
				),
			},
		},
	})
}
