package azureipam

import (
	"context"
	"strconv"
	"time"

	cli "terraform-provider-azureipam/ipamclient"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSpaces() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSpacesRead,
		Description: "The reservations data source allows you to retrieve information about all existing reservations in the specific space and block.",
		Schema: map[string]*schema.Schema{
			"expand": {
				Description: "Indicates if network references to full network objects must be included.",
				Type:     schema.TypeBool,
				Optional: true,
			},
			"append_utilization": {
				Description: "Indicates if utilization information for each network must be included.",
				Type:     schema.TypeBool,
				Optional: true,
			},
			"spaces": {
				Description: "List containing the `spaces` found.",
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Description: "Name of the space.",
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Description: "Text that describes the space.",
							Type:     schema.TypeString,
							Computed: true,
						},
						"blocks": {
							Description: "List containing the `blocks` included in this`space`.",
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Description: "Name of the block.",
										Type:     schema.TypeString,
										Computed: true,
									},
									"cidr": {
										Description: "The IPV4 range assigned to this block, in cidr notation.",
										Type:     schema.TypeString,
										Computed: true,
									},
									"vnets": {
										Description: "List containing the `vnet` included in this `block`.",
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Description: "Name of the virtual network.",
													Type:     schema.TypeString,
													Computed: true,
												},
												"id": {
													Description: "Resourece Id of the virtual network.",
													Type:     schema.TypeString,
													Computed: true,
												},
												"prefixes": {
													Description: "The list of IPV4 prefixes assigned to this vnet, in cidr notation.",
													Type:     schema.TypeList,
													Computed: true,
													Elem:     schema.TypeString,
												},
												"subnets": {
													Description: "List containing the `subnets` included in this `vnet`.",
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"name": {
																Description: "Name of the subnet-",
																Type:     schema.TypeString,
																Computed: true,
															},
															"prefix": {
																Description: "The IPV4 prefix assigned to this block, in cidr notation.",
																Type:     schema.TypeString,
																Computed: true,
															},
															"size": {
																Description: "Total IP's allowed in the `subnet` by its size.",
																Type:     schema.TypeFloat,
																Computed: true,
															},
															"used": {
																Description: "Assigned IP's in the `subnet`.",
																Type:     schema.TypeFloat,
																Computed: true,
															},
														},
													},
												},
												"resource_group": {
													Description: "Name of the resource group where the `vnet` is deployed.",
													Type:     schema.TypeString,
													Computed: true,
												},
												"subscription_id": {
													Description: "Id of the Azure subscription where the `vnet` is deployed.",
													Type:     schema.TypeString,
													Computed: true,
												},
												"tenant_id": {
													Description: "Id of the Azure tenant where the `vnet` is deployed.",
													Type:     schema.TypeString,
													Computed: true,
												},
												"size": {
													Description: "Total IP's allowed in the `vnet` by its size.",
													Type:     schema.TypeFloat,
													Computed: true,
												},
												"used": {
													Description: "Assigned IP's in the `vnet`.",
													Type:     schema.TypeFloat,
													Computed: true,
												},
											},
										},
									},
									"externals": {
										Description: "List containing the `external networks` included in this `block`.",
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Description: "Name of the external network.",
													Type:     schema.TypeString,
													Computed: true,
												},
												"description": {
													Description: "Text that describes the external network.",
													Type:     schema.TypeString,
													Computed: true,
												},
												"cidr": {
													Description: "The IPV4 range reserved for the external network, in cidr notation.",
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"reservations": {
										Description: "List containing the `reservations` included in this `block`.",
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
													Description: "The IPv4 range assigned to this reservation, in cidr notation.",
													Type:     schema.TypeString,
													Computed: true,
												},
												"description": {
													Description: "Text that describes the reservation.",
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
													Description: "Status of the reservation, a 'wait' status indicates that is waiting for the related vnet creation",
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"size": {
										Type:     schema.TypeFloat,
										Computed: true,
									},
									"used": {
										Type:     schema.TypeFloat,
										Computed: true,
									},
								},
							},
						},
						"size": {
							Type:     schema.TypeFloat,
							Computed: true,
							Optional: true,
						},
						"used": {
							Type:     schema.TypeFloat,
							Computed: true,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceSpacesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	expand := d.Get("expand").(bool)
	appendUtilization := d.Get("append_utilization").(bool)
	c := m.(*cli.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	//decode json response to struct defined in models
	spacesInfo, err := c.GetSpaces(expand, appendUtilization)
	if err != nil {
		return diag.FromErr(err)
	}

	//parse to schema
	spaceItems := parseSpaces(spacesInfo, expand, appendUtilization)
	if err := d.Set("spaces", &spaceItems); err != nil {
		return diag.FromErr(err)
	}

	// always run
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}

func parseSpaces(spaces *[]cli.Space, expand bool, appendUtilization bool) []interface{} {
	results := make([]interface{}, 0)
	if spaces != nil {
		for _, space := range *spaces {
			results = append(results, map[string]interface{}{
				"name":        space.Name,
				"description": space.Description,
				"blocks": func() []interface{} {
					if expand {
						return parseBlocks(space.Blocks, expand, appendUtilization)
					} else {
						return nil
					}
				}(),
				"size": func() *float64 {
					if appendUtilization {
						return space.Size
					} else {
						return nil
					}
				}(),
				"used": func() *float64 {
					if appendUtilization {
						return space.Used
					} else {
						return nil
					}
				}(),
			})
		}
	}
	return results
}

func parseBlocks(blocks *[]cli.Block, expand bool, appendUtilization bool) []interface{} {
	results := make([]interface{}, 0)
	if blocks != nil {
		for _, block := range *blocks {
			results = append(results, map[string]interface{}{
				"name": block.Name,
				"cidr": block.Cidr,
				"vnets": func() []interface{} {
					if expand {
						return parseVnets(block.Vnets)
					} else {
						return nil
					}
				}(),
				"externals": func() []interface{} {
					if expand {
						return parseExternals(block.Externals)
					} else {
						return nil
					}
				}(),
				"reservations": func() []interface{} {
					if expand {
						return parseReservationsLite(block.Reservations)
					} else {
						return nil
					}
				}(),
				"size": func() *float64 {
					if appendUtilization {
						return block.Size
					} else {
						return nil
					}
				}(),
				"used": func() *float64 {
					if appendUtilization {
						return block.Used
					} else {
						return nil
					}
				}(),
			})
		}
	}
	return results
}

func parseVnets(vnets *[]cli.Vnet) []interface{} {
	results := make([]interface{}, 0)
	if vnets != nil {
		for _, vnet := range *vnets {
			results = append(results, map[string]interface{}{
				"name":            vnet.Name,
				"id":              vnet.Id,
				"prefixes":        vnet.Prefixes,
				"subnets":         parseSubnets(vnet.Subnets),
				"resource_group":  vnet.ResourceGroup,
				"subscription_id": vnet.SubscriptionId,
				"tenant_id":       vnet.TenantId,
				"size":            vnet.Size,
				"used":            vnet.Used,
			})
		}
	}
	return results
}

func parseSubnets(subnets *[]cli.Subnet) []interface{} {
	results := make([]interface{}, 0)
	if subnets != nil {
		for _, subnet := range *subnets {
			results = append(results, map[string]interface{}{
				"name":   subnet.Name,
				"prefix": subnet.Prefix,
				"size":   subnet.Size,
				"used":   subnet.Used,
			})
		}
	}
	return results
}

func parseExternals(externals *[]cli.External) []interface{} {
	results := make([]interface{}, 0)
	if externals != nil {
		for _, external := range *externals {
			results = append(results, map[string]interface{}{
				"name":        external.Name,
				"description": external.Description,
				"cidr":        external.Cidr,
			})
		}
	}
	return results
}

func parseReservationsLite(reservations *[]cli.ReservationLite) []interface{} {
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
			})
		}
	}
	return results
}
