package azureipam

import (
	"context"
	"strconv"
	"time"

	cli "terraform-provider-azureipam/ipamclient"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceReservations() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceReservationsRead,
		Description: "The reservations data source allows you to retrieve information about all existing reservations in the specified space and block.",
		Schema: map[string]*schema.Schema{
			"space": {
				Description: "Name of the existing space in the IPAM application.",
				Type:     schema.TypeString,
				Required: true,
			},
			"block": {
				Description: "Name of the existing block, related to the specified space.",
				Type:     schema.TypeString,
				Required: true,
			},
			"include_settled": {
				Description: "Settled reservations must be also included? Defaults to `false`.",
				Type:     schema.TypeBool,
				Optional: true,
			},
			"reservations": {
				Description: "List containing the `reservations` found for the specified attributes.",
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Description: "The unique identifier of the reservation.",
							Type:     schema.TypeString,
							Computed: true,
						},
						"cidr": {
							Description: "The assigned and reserved range, in cidr notation.",
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Description: "Description text that describe the reservation.",
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_on": {
							Description: "The date and time that the reservacion was created.",
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_by": {
							Description: "Email or identification of user that created the reservation.",
							Type:     schema.TypeString,
							Computed: true,
						},
						"settled_on": {
							Description: "The date and time when the reservation was settled.",
							Type:     schema.TypeString,
							Computed: true,
						},
						"settled_by": {
							Description: "Email or identification of user that have settled the reservation.",
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Description: "Status of the reservation, a 'wait' status indicates that is waiting for the related vnet creation.",
							Type:     schema.TypeString,
							Computed: true,
						},
						"tags": {
							Description: "Auto-generated tags for the reservation. Particular relevance the 'X-IPAM-RES-ID' tag, since it must be included in the vnet creation in order that the IPAM solution automatically considers the reservation as completed.",
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceReservationsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	space := d.Get("space").(string)
	block := d.Get("block").(string)
	includeSettled := d.Get("include_settled").(bool)
	c := m.(*cli.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	//decode json response to struct defined in models
	reservationsInfo, err := c.GetReservations(space, block, includeSettled)
	if err != nil {
		return diag.FromErr(err)
	}

	//parse to schema
	reservationItems := parseReservations(&reservationsInfo)
	if err := d.Set("reservations", &reservationItems); err != nil {
		return diag.FromErr(err)
	}

	// always run
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}

func parseReservations(reservations *[]cli.Reservation) []interface{} {
	results := make([]interface{}, 0)
	if reservations != nil {
		for _, reservation := range *reservations {
			results = append(results, map[string]interface{}{
				"id":          reservation.Id,
				"cidr":        reservation.Cidr,
				"description": reservation.Description,
				"created_on":  time.Unix(int64(reservation.CreatedOn), 0).Format(time.RFC1123),
				"created_by":  reservation.CreatedBy,
				"settled_on":  time.Unix(int64(reservation.SettledOn), 0).Format(time.RFC1123),
				"settled_by":  reservation.SettledBy,
				"status":      reservation.Status,
				"tags":        reservation.Tags,
			})
		}
	}
	return results
}
