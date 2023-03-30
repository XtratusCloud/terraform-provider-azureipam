package azureipam

import (
	"context"
	"fmt"
	"strings"
	"time"

	cli "terraform-provider-azureipam/ipamclient"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceReservation() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceReservationCreate,
		ReadContext:   resourceReservationRead,
		DeleteContext: resourceReservationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"space": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"block": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"size": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"reverse_search": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"smallest_cidr": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cidr": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_on": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceReservationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*cli.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	//the reservation Id is stored as space/block/id
	idParts := strings.Split(d.Id(), "/")
	space := idParts[0]
	block := idParts[1]
	id := idParts[2]
	//read reservation
	reservation, err := c.GetReservation(space, block, id)
	if err != nil {
		return diag.FromErr(err)
	}
	flattenReservation(reservation, space, block, d)

	return diags
}

func resourceReservationCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	space := d.Get("space").(string)
	block := d.Get("block").(string)
	description := d.Get("description").(string)
	size := d.Get("size").(int)
	reserveSearch := d.Get("reverse_search").(bool)
	smallestCidr := d.Get("smallest_cidr").(bool)
	c := m.(*cli.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	reservation, err := c.CreateReservation(space, block, description, size, reserveSearch, smallestCidr)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", space, block, reservation.Id))

	flattenReservation(reservation, space, block, d)

	return diags
}

func resourceReservationDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*cli.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	//the reservation Id is stored as space/block/id
	idParts := strings.Split(d.Id(), "/")
	space := idParts[0]
	block := idParts[1]
	id := idParts[2]

	//delete reservation
	err := c.DeleteReservation(space, block, id)
	if err != nil {
		return diag.FromErr(err)
	}

	// d.SetId("") is automatically called assuming delete returns no errors, but
	// it is added here for explicitness.
	d.SetId("")

	return diags
}

func flattenReservation(reservation *cli.Reservation, space, block string, d *schema.ResourceData) {
	d.Set("space", space)
	d.Set("block", block)
	d.Set("id", reservation.Id)
	d.Set("cidr", reservation.Cidr)
	d.Set("description", reservation.Description)
	d.Set("created_on", time.Unix(int64(reservation.CreatedOn), 0).Format(time.RFC1123))
	d.Set("created_by", reservation.CreatedBy)
	d.Set("settled_on", time.Unix(int64(reservation.SettledOn), 0).Format(time.RFC1123))
	d.Set("settled_by", reservation.SettledBy)
	d.Set("status", reservation.Status)
	d.Set("tags", reservation.Tags)
}
