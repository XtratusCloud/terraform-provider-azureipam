package provider

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

func TestAccReservationDataSource(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("GET", "https://mockedHost.azurewebsites.net/api/spaces/au/blocks/AustraliaEast/reservations?settled=true",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/datasource/reservations/reservations_with_settled.json").String()), nil
		})

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccProviderConfig + `data "azureipam_reservation" "test" {
					space = "au"
					block = "AustraliaEast"
					id = "3MFHm4s88SVrH8nQ4cK9Um"
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					//Verify common attributes to ensure that all are set
					resource.TestCheckResourceAttr("data.azureipam_reservation.test", "space", "au"),
					resource.TestCheckResourceAttr("data.azureipam_reservation.test", "block", "AustraliaEast"),
					resource.TestCheckResourceAttr("data.azureipam_reservation.test", "id", "3MFHm4s88SVrH8nQ4cK9Um"),
					resource.TestCheckResourceAttr("data.azureipam_reservation.test", "cidr", "10.82.3.0/24"),
					resource.TestCheckResourceAttr("data.azureipam_reservation.test", "description", "this is a test"),
					resource.TestCheckResourceAttrWith("data.azureipam_reservation.test", "created_on", func(value string) error {
						expected, _ := time.Parse(time.RFC3339, "2024-09-06T20:01:30+02:00")
						current, _ := time.Parse(time.RFC3339, value)
						if current.Equal(expected) {
							return nil
						}
						return errors.New("expected " + expected.String() + " got " + current.String())
					}),
					resource.TestCheckResourceAttr("data.azureipam_reservation.test", "created_by", "dummyemail@gmail.com"),
					resource.TestCheckNoResourceAttr("data.azureipam_reservation.test", "settled_on"),
					resource.TestCheckNoResourceAttr("data.azureipam_reservation.test", "settled_by"),
					resource.TestCheckResourceAttr("data.azureipam_reservation.test", "status", "wait"),
					// Verify the tag collection
					resource.TestCheckResourceAttr("data.azureipam_reservation.test", "tags.%", "1"),
					resource.TestCheckResourceAttr("data.azureipam_reservation.test", "tags.X-IPAM-RES-ID", "3MFHm4s88SVrH8nQ4cK9Um"),
				),
			},
		},
	})
}