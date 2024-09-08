package provider

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

func TestAccReservationResource(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("POST", "https://mockedHost.azurewebsites.net/api/spaces/au/blocks/AustraliaEast/reservations",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/reservation/new_reservation.json").String()), nil
		})
	httpmock.RegisterResponder("DELETE", "https://mockedHost.azurewebsites.net/api/spaces/au/blocks/AustraliaEast/reservations",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, ""), nil
		})
	httpmock.RegisterResponder("GET", "https://mockedHost.azurewebsites.net/api/spaces?expand=false&utilization=false",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/reservation/spaces_with_new_reservation_info.json").String()), nil
		})
	httpmock.RegisterResponder("GET", "https://mockedHost.azurewebsites.net/api/spaces/au/blocks/AustraliaEast/reservations?settled=true",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/reservation/reservations_with_new_reservation.json").String()), nil
		})

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccProviderConfig + `resource "azureipam_reservation" "test" {
					space          = "au"
					block          = "AustraliaEast"
					size           = 23
					description    = "acceptance-test"
					reverse_search = true
					smallest_cidr  = true
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					//Verify common attributes to ensure that all are set
					resource.TestCheckResourceAttr("azureipam_reservation.test", "space", "au"),
					resource.TestCheckResourceAttr("azureipam_reservation.test", "block", "AustraliaEast"),
					resource.TestCheckResourceAttr("azureipam_reservation.test", "reverse_search", "true"),
					resource.TestCheckResourceAttr("azureipam_reservation.test", "smallest_cidr", "true"),
					resource.TestCheckResourceAttr("azureipam_reservation.test", "id", "YYtppsvYQsRSBpZLsioZSV"),
					resource.TestCheckResourceAttr("azureipam_reservation.test", "cidr", "10.82.6.0/23"),
					resource.TestCheckResourceAttr("azureipam_reservation.test", "description", "acceptance-test"),
					resource.TestCheckResourceAttrWith("azureipam_reservation.test", "created_on", func(value string) error {
						expected, _ := time.Parse(time.RFC3339, "2024-09-07T06:21:42+02:00")
						current, _ := time.Parse(time.RFC3339, value)
						if current.Equal(expected) {
							return nil
						}
						return errors.New("expected " + expected.String() + " got " + current.String())
					}),
					resource.TestCheckResourceAttr("azureipam_reservation.test", "created_by", "dummyemail@gmail.com"),
					resource.TestCheckNoResourceAttr("azureipam_reservation.test", "settled_on"),
					resource.TestCheckNoResourceAttr("azureipam_reservation.test", "settled_by"),
					resource.TestCheckResourceAttr("azureipam_reservation.test", "status", "wait"),
					resource.TestCheckResourceAttr("azureipam_reservation.test", "tags.%", "1"),
					resource.TestCheckResourceAttr("azureipam_reservation.test", "tags.X-IPAM-RES-ID", "YYtppsvYQsRSBpZLsioZSV"),
				),
			},
			// ImportState testing
			{
				ResourceName:            "azureipam_reservation.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"reverse_search", "smallest_cidr"},
			},
			// Update  NOT ALLOWED by provider

			// Delete testing automatically occurs in TestCase
		},
	})
}
