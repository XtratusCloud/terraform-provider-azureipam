package provider

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

func TestAccReservationsNotSettledDataSource(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("GET", "https://mockedHost.azurewebsites.net/api/spaces/au/blocks/AustraliaEast/reservations?settled=false",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/datasource/reservations/reservations_not_settled.json").String()), nil
		})

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccProviderConfig + `data "azureipam_reservations" "test" {
					space = "au"
					block = "AustraliaEast"
					include_settled = false
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					//Verify common attributes to ensure that all are set
					resource.TestCheckResourceAttr("data.azureipam_reservations.test", "space", "au"),
					resource.TestCheckResourceAttr("data.azureipam_reservations.test", "block", "AustraliaEast"),
					resource.TestCheckResourceAttr("data.azureipam_reservations.test", "include_settled", "false"),
					// Verify number of reservations returned
					resource.TestCheckResourceAttr("data.azureipam_reservations.test", "reservations.#", "1"),
					// Verify the first reservation to ensure all attributes are set
					resource.TestCheckResourceAttr("data.azureipam_reservations.test", "reservations.0.id", "3MFHm4s88SVrH8nQ4cK9Um"),
					resource.TestCheckResourceAttr("data.azureipam_reservations.test", "reservations.0.cidr", "10.82.3.0/24"),
					resource.TestCheckResourceAttr("data.azureipam_reservations.test", "reservations.0.description", "this is a test"),
					resource.TestCheckResourceAttrWith("data.azureipam_reservations.test", "reservations.0.created_on", func(value string) error {
						expected, _ := time.Parse(time.RFC3339, "2024-09-06T20:01:30+02:00")
						current, _ := time.Parse(time.RFC3339, value)
						if current.Equal(expected) {
							return nil
						}
						return errors.New("expected " + expected.String() + " got " + current.String())
					}),
					resource.TestCheckResourceAttr("data.azureipam_reservations.test", "reservations.0.created_by", "dummyemail@gmail.com"),
					resource.TestCheckNoResourceAttr("data.azureipam_reservations.test", "reservations.0.settled_on"),
					resource.TestCheckNoResourceAttr("data.azureipam_reservations.test", "reservations.0.settled_by"),
					resource.TestCheckResourceAttr("data.azureipam_reservations.test", "reservations.0.status", "wait"),
					// Verify the tag collection
					resource.TestCheckResourceAttr("data.azureipam_reservations.test", "reservations.0.tags.%", "1"),
					resource.TestCheckResourceAttr("data.azureipam_reservations.test", "reservations.0.tags.X-IPAM-RES-ID", "3MFHm4s88SVrH8nQ4cK9Um"),
				),
			},
		},
	})
}

func TestAccReservationsWithSettledDataSource(t *testing.T) {
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
				Config: testAccProviderConfig + `data "azureipam_reservations" "test" {
					space = "au"
					block = "AustraliaEast"
					include_settled = true
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					//Verify common attributes to ensure that all are set
					resource.TestCheckResourceAttr("data.azureipam_reservations.test", "space", "au"),
					resource.TestCheckResourceAttr("data.azureipam_reservations.test", "block", "AustraliaEast"),
					resource.TestCheckResourceAttr("data.azureipam_reservations.test", "include_settled", "true"),
					// Verify number of reservations returned
					resource.TestCheckResourceAttr("data.azureipam_reservations.test", "reservations.#", "4"),
					// Verify the first reservation to ensure all attributes are set
					resource.TestCheckResourceAttr("data.azureipam_reservations.test", "reservations.0.id", "3MFHm4s88SVrH8nQ4cK9Um"),
					resource.TestCheckResourceAttr("data.azureipam_reservations.test", "reservations.0.cidr", "10.82.3.0/24"),
					resource.TestCheckResourceAttr("data.azureipam_reservations.test", "reservations.0.description", "this is a test"),
					resource.TestCheckResourceAttrWith("data.azureipam_reservations.test", "reservations.0.created_on", func(value string) error {
						expected, _ := time.Parse(time.RFC3339, "2024-09-06T20:01:30+02:00")
						current, _ := time.Parse(time.RFC3339, value)
						if current.Equal(expected) {
							return nil
						}
						return errors.New("expected " + expected.String() + " got " + current.String())
					}),
					resource.TestCheckResourceAttr("data.azureipam_reservations.test", "reservations.0.created_by", "dummyemail@gmail.com"),
					resource.TestCheckNoResourceAttr("data.azureipam_reservations.test", "reservations.0.settled_on"),
					resource.TestCheckNoResourceAttr("data.azureipam_reservations.test", "reservations.0.settled_by"),
					resource.TestCheckResourceAttr("data.azureipam_reservations.test", "reservations.0.status", "wait"),
					resource.TestCheckResourceAttr("data.azureipam_reservations.test", "reservations.0.tags.%", "1"),
					resource.TestCheckResourceAttr("data.azureipam_reservations.test", "reservations.0.tags.X-IPAM-RES-ID", "3MFHm4s88SVrH8nQ4cK9Um"),
					// Verify the last reservation to ensure all attributes are set
					resource.TestCheckResourceAttr("data.azureipam_reservations.test", "reservations.3.id", "hi3fxt9PeSpxhykfSszVUb"),
					resource.TestCheckResourceAttr("data.azureipam_reservations.test", "reservations.3.cidr", "10.82.1.160/27"),
					resource.TestCheckResourceAttr("data.azureipam_reservations.test", "reservations.3.description", "vnet-we-c-arq3tier-01"),
					resource.TestCheckResourceAttrWith("data.azureipam_reservations.test", "reservations.3.created_on", func(value string) error {
						expected, _ := time.Parse(time.RFC3339, "2023-11-08T13:51:07+01:00")
						current, _ := time.Parse(time.RFC3339, value)
						if current.Equal(expected) {
							return nil
						}
						return errors.New("expected " + expected.String() + " got " + current.String())
					}),					 
					resource.TestCheckResourceAttr("data.azureipam_reservations.test", "reservations.3.created_by", "spn:9fc2493a-b515-49a6-9d73-93e1bac5f6cc"),
					resource.TestCheckResourceAttrWith("data.azureipam_reservations.test", "reservations.3.settled_on", func(value string) error {
						expected, _ := time.Parse(time.RFC3339, "2024-04-03T09:18:41+02:00")
						current, _ := time.Parse(time.RFC3339, value)
						if current.Equal(expected) {
							return nil
						}
						return errors.New("expected " + expected.String() + " got " + current.String())
					}),					 
					resource.TestCheckResourceAttr("data.azureipam_reservations.test", "reservations.3.settled_by", "dummyemail@gmail.com"),
					resource.TestCheckResourceAttr("data.azureipam_reservations.test", "reservations.3.status", "cancelledByUser"),
					resource.TestCheckResourceAttr("data.azureipam_reservations.test", "reservations.3.tags.%", "1"),
					resource.TestCheckResourceAttr("data.azureipam_reservations.test", "reservations.3.tags.X-IPAM-RES-ID", "hi3fxt9PeSpxhykfSszVUb"),
				),
			},
		},
	})
}
