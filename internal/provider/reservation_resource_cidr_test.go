package provider

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

func TestAccReservationCidrResource(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("POST", "https://mockedHost.azurewebsites.net/api/spaces/au/blocks/AustraliaSoutheast/reservations",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/reservation_cidr/new_reservation.json").String()), nil
		})
	httpmock.RegisterResponder("DELETE", "https://mockedHost.azurewebsites.net/api/spaces/au/blocks/AustraliaSoutheast/reservations",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, ""), nil
		})
	httpmock.RegisterResponder("GET", "https://mockedHost.azurewebsites.net/api/spaces?expand=false&utilization=false",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/reservation_cidr/spaces_with_new_reservation_info.json").String()), nil
		})
	httpmock.RegisterResponder("GET", "https://mockedHost.azurewebsites.net/api/spaces/au/blocks/AustraliaSoutheast/reservations?settled=true",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/reservation_cidr/reservations_with_new_reservation.json").String()), nil
		})

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccProviderConfig + `resource "azureipam_reservation_cidr" "test" {
					space          = "au"
					block          = "AustraliaSoutheast"
					specific_cidr  = "10.82.4.0/24"
					description    = "acceptance-test"
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					//Verify common attributes to ensure that all are set
					resource.TestCheckResourceAttr("azureipam_reservation_cidr.test", "space", "au"),
					resource.TestCheckResourceAttr("azureipam_reservation_cidr.test", "block", "AustraliaSoutheast"),
					resource.TestCheckResourceAttr("azureipam_reservation_cidr.test", "id", "Etc4svKttPXMQyvCb9sjy2"),
					resource.TestCheckResourceAttr("azureipam_reservation_cidr.test", "cidr", "10.82.4.0/24"),
					resource.TestCheckResourceAttr("azureipam_reservation_cidr.test", "description", "acceptance-test"),
					resource.TestCheckResourceAttrWith("azureipam_reservation_cidr.test", "created_on", func(value string) error {
						expected, _ := time.Parse(time.RFC3339, "2024-09-07T06:21:42+02:00")
						current, _ := time.Parse(time.RFC3339, value)
						if current.Equal(expected) {
							return nil
						}
						return errors.New("expected " + expected.String() + " got " + current.String())
					}),
					resource.TestCheckResourceAttr("azureipam_reservation_cidr.test", "created_by", "dummyemail@gmail.com"),
					resource.TestCheckNoResourceAttr("azureipam_reservation_cidr.test", "settled_on"),
					resource.TestCheckNoResourceAttr("azureipam_reservation_cidr.test", "settled_by"),
					resource.TestCheckResourceAttr("azureipam_reservation_cidr.test", "status", "wait"),
					resource.TestCheckResourceAttr("azureipam_reservation_cidr.test", "tags.%", "1"),
					resource.TestCheckResourceAttr("azureipam_reservation_cidr.test", "tags.X-IPAM-RES-ID", "Etc4svKttPXMQyvCb9sjy2"),
				),
			},
			// ImportState testing
			{
				ResourceName:            "azureipam_reservation_cidr.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
			// Update  NOT ALLOWED by provider

			// Delete testing automatically occurs in TestCase
		},
	})
}
