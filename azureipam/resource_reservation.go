package azureipam

import (
	"context"
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
			StateContext: schema.ImportStatePassthroughContext,
		},
		Description: "The reservation resource allows you to create a IPAM reservation in the specific space and block.",
		Schema: map[string]*schema.Schema{
			"space": {
				Description: "Name of the existing space in the IPAM application.",
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"block": {
				Description: "Name of the existing block, related to the specified space, in which the reservation is to be made.",
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"size": {
				Description: "Integer value to indicate the subnet mask bits, which defines the size of the vnet to reserve (example 24 for a /24 subnet).",
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Description: "Description text that describe the reservation, that will be added as an additional tag.",
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"reverse_search": {
				Description: "New networks will be created as close to the end of the block as possible?. Defaults to `false`.",
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"smallest_cidr": {
				Description: "New networks will be created using the smallest possible available block? (e.g. it will not break up large CIDR blocks when possible) .Defaults to `false`.",
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"id": {
				Description: "The unique identifier of the generated reservation.",
				Type:     schema.TypeString,
				Computed: true,
			},
			"cidr": {
				Description: "The assigned and reserved range, in cidr notation.",
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_by": {
				Description: "Email or identification of user that created the reservation.",
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_on": {
				Description: "The date and time that the reservacion was created.",
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Description: "Status of the reservation, a 'wait' status indicates that is waiting for the related vnet creation",
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
	}
}

func resourceReservationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	id := d.Get("id").(string)
	space := d.Get("space").(string)
	block := d.Get("block").(string)
	c := m.(*cli.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

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

	//Create reservation
	reservation, err := c.CreateReservation(space, block, description, size, reserveSearch, smallestCidr)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(reservation.Id)

	flattenReservation(reservation, space, block, d)

	return diags
}

func resourceReservationDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	id := d.Get("id").(string)
	space := d.Get("space").(string)
	block := d.Get("block").(string)
	c := m.(*cli.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

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
