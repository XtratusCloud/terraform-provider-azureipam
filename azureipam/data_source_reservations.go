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
		Schema: map[string]*schema.Schema{
			"space": {
				Type:     schema.TypeString,
				Required: true,
			},
			"block": {
				Type:     schema.TypeString,
				Required: true,
			},
			"include_settled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"reservations": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cidr": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_on": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_by": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"settled_on": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"settled_by": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tags": {
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
